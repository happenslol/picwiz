package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/happenslol/picwiz/models"
)

type votesCreateRequest struct {
	PictureID uuid.UUID `json:"uuid"`
}

// VotesCreate creates a vote for a picture
func VotesCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	pictureID, _ := uuid.FromString(c.Param("pictureId"))

	picture := models.Picture{}
	if err := tx.Find(&picture, pictureID); err != nil {
		return c.Render(404, nil)
	}

	toSave := models.Vote{
		ID:        uuid.Must(uuid.NewV1()),
		PictureID: picture.ID,
	}

	if err := tx.Create(&toSave); err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	// TODO: Update picture score here
	// picture.ConfidenceLevel = ?
	if err := tx.Save(&picture); err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return c.Render(204, nil)
}
