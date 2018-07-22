package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/happenslol/picwiz/models"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	picBuffer := models.Pictures{}
	if err := tx.
		RawQuery("SELECT * FROM pictures ORDER BY sorting LIMIT 5").
		All(&picBuffer); err != nil {
		return err
	}

	preloadImages := []string{}
	for _, p := range picBuffer {
		preloadImages = append(preloadImages, p.ID.String())
	}

	c.Set("preloadImages", preloadImages)

	return c.Render(200, r.HTML("voter.html"))
}
