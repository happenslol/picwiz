package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/happenslol/picwiz/models"
	"github.com/spf13/afero"
)

// ImportsCandidates lists all external directories that
// could be imported
func ImportsCandidates(c buffalo.Context) error {
	res := make(M)
	var locations []string

	for _, loc := range scanLocations {
		exists, _ := afero.DirExists(fs, loc)
		if exists {
			locations = append(locations, loc)
		}
	}

	c.Logger().Infof(
		"looking for directories in %s",
		strings.Join(locations, ", "),
	)

	for _, loc := range locations {
		files, err := afero.ReadDir(fs, loc)
		if err != nil {
			return c.Error(500, errors.New(err.Error()))
		}

		var dirs []string
		for _, f := range files {
			if f.IsDir() {
				dirs = append(dirs, f.Name())
			}
		}

		if len(dirs) > 0 {
			res[loc] = dirs
		}
	}

	return c.Render(200, r.JSON(res))
}

// ImportsCreateRequest ...
type importsCreateRequest struct {
	Author    string `json:"author"`
	Location  string `json:"location"`
	Directory string `json:"directory"`
}

// ImportsCreate creates a new import and imports all the files
// from the given directory
func ImportsCreate(c buffalo.Context) error {
	var req importsCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	srcString := fmt.Sprintf(
		"%s%s%s",
		req.Location,
		afero.FilePathSeparator,
		req.Directory,
	)

	c.Logger().Infof("built src string %s", srcString)

	exists, _ := afero.DirExists(fs, srcString)
	if !exists {
		return c.Error(500, errors.New("target directory not found"))
	}

	c.Logger().Infof("importing all files from %s", srcString)

	toCreate := &models.Import{
		Author:    req.Author,
		Source:    srcString,
		Processed: false,
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return c.Error(500, errors.New("no tx found"))
	}

	err := tx.Create(toCreate)
	if err != nil {
		return err
	}

	return c.Render(201, r.JSON(toCreate))
}
