package actions

import "github.com/gobuffalo/buffalo"

// AdminIndex default implementation.
func AdminIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("admin.html"))
}
