package actions

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/happenslol/picwiz/models"
)

// PicturesHot ...
func PicturesHot(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hotPics := models.Pictures{}
	limit, err := strconv.Atoi(c.Param("limit"))

	if err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	if limit == 0 {
		limit = 200
	}

	query := fmt.Sprintf(
		"SELECT * FROM pictures ORDER BY sorting DESC LIMIT %d", //todo(Hannes): censoring (but only when param all not set)
		limit,
	)

	err = tx.
		RawQuery(query).
		All(&hotPics)

	if err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return c.Render(200, r.JSON(hotPics))
}

func biasedRandom(max int32) int32 {
	rand.Seed(time.Now().UnixNano())

	n := 4
	unif := rand.Float64()

	oneOver2N := 1 / math.Pow(2.0, float64(n))
	oneOverXPlus1N := 1 / math.Pow(unif+1.0, float64(n))

	random := (oneOverXPlus1N - oneOver2N) / (1 - oneOver2N)
	return int32(random * float64(max))
}

// PicturesNext ...
func PicturesNext(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	nextPic, err := getNextVotingPicture(tx)

	if err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return c.Render(200, r.JSON(nextPic.ID))
}

func getNextVotingPicture(tx *pop.Connection) (*models.Picture, error) {
	nextPic := models.Picture{}

	var count int
	err := tx.
		RawQuery("SELECT count(*) FROM pictures WHERE sorting >= -3").
		First(&count)

	if err != nil {
		return nil, err
	}

	skip := biasedRandom(int32(count))

	query := fmt.Sprintf(
		"SELECT * FROM pictures WHERE sorting >= -3 ORDER BY confidence_level ASC LIMIT 1 OFFSET %d",
		skip,
	)

	err = tx.RawQuery(query).First(&nextPic)

	if err != nil {
		return nil, err
	}

	return &nextPic, nil
}

func SetCensorStatus(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	pictureID, _ := uuid.FromString(c.Param("pictureId"))

	censoringText := c.Param("censored")

	isCensored := false
	if censoringText == "true" {
		isCensored = true
	}

	picture := models.Picture{}
	if err := tx.Find(&picture, pictureID); err != nil {
		return c.Render(404, nil)
	}

	picture.Censored = isCensored

	if err := tx.Save(&picture); err != nil {
		return c.Render(500, r.JSON(M{"error": err.Error()}))
	}

	return nil
}
