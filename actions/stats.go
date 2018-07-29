package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/happenslol/picwiz/models"
)

type Stats struct {
	Votes    int32
	Pictures int32
}

// GetVoteCount gets total vote count
func GetVoteCount(tx *pop.Connection) int32 {
	count, err := tx.
		Count(&models.Vote{})

	if err != nil {
		fmt.Printf("Orror du nap!! (picture count)%v\n", err)
		return 0
	}

	return int32(count)
}

// GetPictureCount gets the total picture count
func GetPictureCount(tx *pop.Connection) int32 {
	count, err := tx.
		Count(&models.Picture{})

	if err != nil {
		fmt.Printf("Orror du nap!! (picture count)%v\n", err)
		return 0
	}

	return int32(count)
}

// GetStats gets all stats
func getStats(tx *pop.Connection) Stats {
	return Stats{
		Votes:    GetVoteCount(tx),
		Pictures: GetPictureCount(tx),
	}
}

// RenderStatsPage renders a page contains all stats
func RenderStatsPage(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	c.Set("stats", getStats(tx))

	return c.Render(200, r.HTML("stats.html"))
}
