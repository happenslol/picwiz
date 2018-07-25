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
})

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

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filePath := fmt.Sprintf(
			"%s%s%s",
			loc,
			afero.FilePathSeparator,
			f.Name(),
		)

		fmt.Printf("\timporting file %s\n", filePath)

		file, err := os.Open(filePath)
		defer file.Close()
		if err != nil {
			fmt.Printf("\terror opening image: %v\n", err)
			continue
		}

		decoded, _, err := image.Decode(file)
		if err != nil {
			fmt.Printf("\terror decoding image: %v\n", err)
			continue
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("\terror reading image bytes: %v\n", err)
			continue
		}

		hasher := murmur3.New128()
		hasher.Write(bytes)
		hash := hex.EncodeToString(hasher.Sum(nil))

		picId := uuid.Must(uuid.NewV1())
		picture := models.Picture{
			ID:              picId,
			Filename:        f.Name(),
			Hash:            hash,
			ImportID:        p.ID,
			ConfidenceLevel: 0.5,
			Sorting:         0.5,
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
			continue
		}

		if err := models.DB.Create(&picture); err != nil {
			fmt.Printf("error saving image: %v\n", err)
			continue
		}

		jpeg.Encode(out, resized, nil)
	}

	// p.Processed = true
	// models.DB.Save(&p)

	return nil
}
