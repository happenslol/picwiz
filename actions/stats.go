package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/happenslol/picwiz/models"
)

type Stats struct {
	Votes      int32
	Pictures   int32
	VoteCounts []voteCount
	AllVotes   []models.Vote
}

type voteCount struct {
	Votes int32 `db:"votes"`
	Count int32 `db:"count"`
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

func getPicturesPerVotes(tx *pop.Connection) []voteCount {

	counts := []voteCount{}

	err := tx.
		RawQuery("SELECT sub.votes, count(*) FROM (SELECT (upvotes + downvotes) as votes FROM pictures) sub GROUP BY votes ORDER BY sub.votes").
		All(&counts)

	if err != nil {
		fmt.Printf("Orror du nap!! (picture vote count)%v\n", err)
	}
	return counts
}

func getAllVotes(tx *pop.Connection) []models.Vote {
	votes := []models.Vote{}
	err := tx.
		RawQuery("SELECT * FROM votes").
		All(&votes)

	if err != nil {
		fmt.Printf("Orror du nap!! (picture vote count)%v\n", err)
	}
	return votes
}

// GetStats gets all stats
func getStats(tx *pop.Connection) Stats {
	return Stats{
		Votes:      GetVoteCount(tx),
		Pictures:   GetPictureCount(tx),
		VoteCounts: getPicturesPerVotes(tx),
		AllVotes:   getAllVotes(tx),
	}
}

// RenderStatsPage renders a page contains all stats
func RenderStatsPage(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	c.Set("stats", getStats(tx))

	return c.Render(200, r.HTML("stats.html"))
}
