package actions

import "github.com/gobuffalo/buffalo"

type votesCreateRequest struct {
	PictureId uuid.UUID `json:"uuid"`
}

// VotesCreate creates a vote for a picture
func VotesCreate(c buffalo.Context) error {
	return c.Render(204, nil)
}
