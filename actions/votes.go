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
	pictureID, _ := uuid.FromString(c.Param("pictureId"))
	picture := models.Picture{
		ID: pictureID,
	}

	uuid, _ := uuid.NewV1()
	toSave := models.Vote{
		ID:      uuid,
		Picture: picture,
	}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Create(&toSave); err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return c.Render(204, nil)
}
