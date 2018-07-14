package actions

import "github.com/gobuffalo/buffalo"

// ImportsIndex default implementation.
func ImportsIndex(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]interface{}{
		"hello": "world",
	}))
}
