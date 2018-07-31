package grifts

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/happenslol/picwiz/models"
	"github.com/markbates/grift/grift"
	"github.com/nfnt/resize"
	"github.com/spaolacci/murmur3"
	"github.com/spf13/afero"
)

var _ = grift.Namespace("imports", func() {
	grift.Desc("process", "Processes all pending imports")
	grift.Add("process", func(c *grift.Context) error {
		// get all unprocessed imports
		pending := []models.Import{}
		if err := models.DB.RawQuery(
			"SELECT * FROM imports WHERE processed=false",
		).All(&pending); err != nil {
			return err
		}

		if len(pending) == 0 {
			fmt.Printf("no pending imports!\n")
			return nil
		}

		for _, p := range pending {
			fmt.Printf("found pending import dir %s\n", p.Source)
			if err := processPendingImport(p); err != nil {
				fmt.Printf("error while importing: %v", err)
			}
		}

		return nil
	})

	grift.Desc("rehash", "Calculates new hashes for all processed imports")
	grift.Add("rehash", func(c *grift.Context) error {
		// get all processed imports
		processed := models.Imports{}
		if err := models.DB.RawQuery(
			"SELECT * FROM imports WHERE processed=true",
		).All(&processed); err != nil {
			return err
		}

		for _, i := range processed {
			if err := hashImport(i); err != nil {
				fmt.Printf("error while hashing import: %v\n", err)
			}
		}

		return nil
	})

	grift.Desc("dedupe", "Dedupes using hash")
	grift.Add("dedupe", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			// dedupe after rehashing
			duplicateHashes := []string{}
			if err := tx.RawQuery(
				"SELECT hash FROM pictures GROUP BY " +
					"hash HAVING COUNT(*) > 1",
			).All(&duplicateHashes); err != nil {
				return err
			}

			for _, h := range duplicateHashes {
				query := fmt.Sprintf(
					"SELECT * FROM pictures WHERE hash='%s'", h,
				)

				pics := models.Pictures{}
				if err := tx.RawQuery(query).All(&pics); err != nil {
					return err
				}

				if len(pics) == 0 {
					continue
				}

				first := pics[0]

				ids := []string{}
				for i, pic := range pics {
					if i == 0 {
						continue
					}

					quoted := fmt.Sprintf("'%s'", pic.ID.String())
					ids = append(ids, quoted)
				}

				arg := strings.Join(ids, ",")

				voteQuery := fmt.Sprintf(
					"SELECT * FROM votes WHERE picture_id IN (%s)",
					arg,
				)

				votes := models.Votes{}
				if err := tx.RawQuery(voteQuery).All(&votes); err != nil {
					return err
				}

				for _, v := range votes {
					v.PictureID = first.ID
					if err := tx.Save(&v); err != nil {
						return err
					}
				}

				otherPics := models.Pictures{}
				picsQuery := fmt.Sprintf("SELECT * FROM pictures WHERE id IN (%s)", arg)
				if err := tx.RawQuery(picsQuery).All(&otherPics); err != nil {
					return err
				}

				for _, p := range otherPics {
					first.Upvotes += p.Upvotes
					first.Downvotes += p.Downvotes
				}

				if err := tx.Save(&first); err != nil {
					return err
				}

				deleteQuery := fmt.Sprintf(
					"DELETE FROM pictures WHERE hash='%s' AND id != '%s'",
					first.Hash,
					first.ID.String(),
				)

				if err := tx.RawQuery(deleteQuery).Exec(); err != nil {
					return err
				}
			}

			return nil
		})
	})
})

func hashImport(i models.Import) error {
	fmt.Printf("rehashing import %s\n", i.ID.String())

	loc := fmt.Sprintf(
		"%s%simports%s%s",
		storagePath,
		afero.FilePathSeparator,
		afero.FilePathSeparator,
		i.ID,
	)

	isDir, _ := afero.IsDir(fs, loc)
	if !isDir {
		return errors.New(
			fmt.Sprintf("%s was not a directory\n", loc),
		)
	}

	query := fmt.Sprintf(
		"SELECT * FROM pictures WHERE import_id='%s'",
		i.ID.String(),
	)

	pics := models.Pictures{}
	if err := models.DB.RawQuery(query).All(&pics); err != nil {
		return err
	}

	concurrency := 3
	sem := make(chan bool, concurrency)

	all := len(pics)
	var wg sync.WaitGroup
	wg.Add(all)
	for _, p := range pics {
		sem <- true
		go func() {
			o := fmt.Sprintf(
				"%s%s%s",
				loc,
				afero.FilePathSeparator,
				p.Filename,
			)

			if err := hashPicture(o, p, &wg, sem); err != nil {
				fmt.Printf("\terror hashing picture: %v\n", err)
			}
		}()
	}

	wg.Wait()

	return nil
}

func hashPicture(
	loc string,
	p models.Picture,
	wg *sync.WaitGroup,
	sem chan bool,
) error {
	defer func() {
		<-sem
		wg.Done()
	}()

	fmt.Printf("\trehashing picture %s: %s\n", loc, p.ID.String())

	file, err := os.Open(loc)
	defer file.Close()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	hasher := murmur3.New128()
	hasher.Write(bytes)
	p.Hash = hex.EncodeToString(hasher.Sum(nil))

	if err := models.DB.Save(&p); err != nil {
		return err
	}

	return nil
}

func processPendingImport(p models.Import) error {
	loc := fmt.Sprintf(
		"%s%simports%s%s",
		storagePath,
		afero.FilePathSeparator,
		afero.FilePathSeparator,
		p.ID,
	)

	isDir, _ := afero.IsDir(fs, loc)
	if !isDir {
		return errors.New(fmt.Sprintf("%s was not a directory\n", loc))
	}

	files, err := afero.ReadDir(fs, loc)
	if err != nil {
		return err
	}

	concurrency := 3
	sem := make(chan bool, concurrency)

	all := len(files)
	var wg sync.WaitGroup
	wg.Add(all)
	for i, f := range files {
		sem <- true
		go processImage(f, loc, p.ID, &wg, sem, i, all)
	}

	wg.Wait()

	p.Processed = true
	models.DB.Save(&p)

	fmt.Printf("import %s successful!\n", p.ID.String())

	return nil
}

func processImage(
	f os.FileInfo,
	loc string,
	importID uuid.UUID,
	wg *sync.WaitGroup,
	sem chan bool,
	i int,
	all int,
) {
	defer func() {
		<-sem
		wg.Done()
	}()

	if f.IsDir() {
		wg.Done()
		return
	}

	filePath := fmt.Sprintf(
		"%s%s%s",
		loc,
		afero.FilePathSeparator,
		f.Name(),
	)

	fmt.Printf("\t(%d/%d) importing file %s\n", i+1, all, filePath)

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Printf("\terror opening image: %v\n", err)
		wg.Done()
		return
	}

	decoded, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("\terror decoding image: %v\n", err)
		wg.Done()
		return
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("\terror reading image bytes: %v\n", err)
		wg.Done()
		return
	}

	hasher := murmur3.New128()
	hasher.Write(bytes)
	hash := hex.EncodeToString(hasher.Sum(nil))

	picId := uuid.Must(uuid.NewV4())
	picture := models.Picture{
		ID:              picId,
		Filename:        f.Name(),
		Hash:            hash,
		ImportID:        importID,
		ConfidenceLevel: 0.0,
		Sorting:         0.0,
	}

	count, err := models.DB.
		Where("hash = ?", hash).
		Count(&models.Pictures{})

	if err != nil {
		fmt.Printf("\terror checking dupes: %v\n", err)
		wg.Done()
		return
	}

	if count > 0 {
		fmt.Printf("\tskipping duplicate file: %v\n", f.Name())
		wg.Done()
		return
	}

	b := decoded.Bounds()

	var outX, outY uint
	if b.Dx() > b.Dy() {
		outX = 1920
		outY = 0
	} else {
		outY = 1920
		outX = 0
	}

	resized := resize.Resize(outX, outY, decoded, resize.Lanczos3)

	saveLoc := fmt.Sprintf(
		"%s%sstatic%s%s.jpg",
		storagePath,
		afero.FilePathSeparator,
		afero.FilePathSeparator,
		picId,
	)

	out, err := os.Create(saveLoc)
	defer out.Close()
	if err != nil {
		fmt.Printf("\terror creating file: %v\n", err)
		wg.Done()
		return
	}

	if err := models.DB.Create(&picture); err != nil {
		fmt.Printf("error saving image: %v\n", err)
	}

	jpeg.Encode(out, resized, nil)
}
