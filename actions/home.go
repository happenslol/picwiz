package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	preloadImages := []string{}
	for i := 1; i < 5; i++ {
		pic, err := getNextVotingPicture(tx)

		if err != nil {
			return c.Render(500, r.JSON(M{"error": err.Error()}))
		}
		preloadImages = append(preloadImages, pic.ID.String())
	}

	c.Set("preloadImages", preloadImages)

	return c.Render(200, r.HTML("voter.html"))
}
