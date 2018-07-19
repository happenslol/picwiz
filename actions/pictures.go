package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/happenslol/picwiz/models"
)

// PicturesHot ...
func PicturesHot(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	nextPic := models.Picture{}

	err := tx.
		RawQuery("SELECT * FROM pictures ORDER BY sorting DESC LIMIT 1").
		First(&nextPic)

	if err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return c.Render(200, r.JSON(nextPic.ID))
}
