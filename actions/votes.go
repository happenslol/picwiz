package actions

import (
	"math"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/happenslol/picwiz/models"
)

type votesCreateRequest struct {
	IsUpvote bool `json:"isUpvote"`
}

// VotesCreate creates a vote for a picture
func VotesCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	pictureID, _ := uuid.FromString(c.Param("pictureId"))

	req := votesCreateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

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

	if req.IsUpvote {
		picture.Upvotes = picture.Upvotes + 1
	} else {
		picture.Downvotes = picture.Downvotes + 1
	}

	picture.Sorting = float32(
		confidenceLevel(picture.Upvotes, picture.Downvotes),
	)

	picture.ConfidenceLevel = float32(
		math.Abs(float64(0.5-picture.Sorting)) *
			float64(picture.Upvotes+picture.Downvotes) / 10.0,
	)

	if err := tx.Save(&picture); err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return c.Render(204, nil)
}

func confidenceLevel(ups uint32, downs uint32) float64 {
	if ups == 0 {
		if downs == 0 {
			return 0.5
		}

		return float64(-(int32(downs)))
	}

	n := float64(ups + downs)
	z := float64(1.64485) //1.0 = 85%, 1.6 = 95%
	phat := float64(ups) / n

	// TODO: Clean this clusterfuck up hannes
	return (phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
}
