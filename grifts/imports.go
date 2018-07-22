package grifts

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gobuffalo/uuid"
	"github.com/happenslol/picwiz/models"
	"github.com/markbates/grift/grift"
	"github.com/spaolacci/murmur3"
	"github.com/spf13/afero"
	bimg "gopkg.in/h2non/bimg.v1"
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

		buffer, err := bimg.Read(filePath)
		if err != nil {
			fmt.Printf(
				"\terror importing image: %s - %v\n",
				filePath,
				err,
			)
		}

		err = nil

		hasher := murmur3.New128()
		hasher.Write(buffer)
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

		var resized []byte
		img := bimg.NewImage(buffer)
		dims, _ := img.Size()

		if dims.Height > dims.Width {
			resized, err = bimg.Resize(buffer, resizeOptsLandscape)
		} else {
			resized, err = bimg.Resize(buffer, resizeOptsPortrait)
		}

		if err != nil {
			fmt.Printf(
				"\terror resizing image: %s - %v",
				filePath,
				err,
			)
		}

		saveLoc := fmt.Sprintf(
			"%s%sstatic%s%s.jpg",
			storagePath,
			afero.FilePathSeparator,
			afero.FilePathSeparator,
			picId,
		)

		if err := models.DB.Create(&picture); err != nil {
			fmt.Printf("error saving image: %v\n", err)
		}

		bimg.Write(saveLoc, resized)
	}

	// p.Processed = true
	// models.DB.Save(&p)

	return nil
}
