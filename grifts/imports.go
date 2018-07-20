package grifts

import (
	"fmt"

	"github.com/happenslol/picwiz/models"
	"github.com/markbates/grift/grift"
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
	return nil
}
