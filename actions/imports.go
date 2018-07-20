package actions

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
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

	srcStr := fmt.Sprintf(
		"%s%s%s",
		req.Location,
		afero.FilePathSeparator,
		req.Directory,
	)

	c.Logger().Infof("built src string %s", srcStr)

	exists, _ := afero.DirExists(fs, srcStr)
	if !exists {
		return c.Error(500, errors.New("target directory not found"))
	}

	c.Logger().Infof("importing all files from %s", srcStr)
	uuid, _ := uuid.NewV1()

	toCreate := &models.Import{
		ID:        uuid,
		Author:    req.Author,
		Source:    srcStr,
		Processed: false,
	}

	go func() {
		fmt.Printf("storage path: %s\n", storagePath)
		destStr := fmt.Sprintf(
			"%s%simports%s%s",
			storagePath,
			afero.FilePathSeparator,
			afero.FilePathSeparator,
			toCreate.ID,
		)
		fmt.Printf("target path: %s\n", destStr)

		if runtime.GOOS == "windows" {
			xcopyCmd := fmt.Sprintf(
				"xcopy \"%s\" \"%s\" /c /i /q /s /y",
				srcStr,
				destStr,
			)

			fmt.Printf("running xcopy with opts: %s\n", xcopyCmd)

			cmd := exec.Command("powershell", xcopyCmd)

			if err := cmd.Run(); err != nil {
				fmt.Printf("import error: %v\n", err)
			} else {
				fmt.Printf("import %s successful\n", toCreate.ID)
			}
		} else {
			cpOpts := fmt.Sprintf(
				"cp -rf '%s' '%s'",
				srcStr,
				destStr,
			)

			fmt.Printf("running cp with opts %s\n", cpOpts)

			cmd := exec.Command("bash", "-c", cpOpts)

			if err := cmd.Run(); err != nil {
				fmt.Printf("import error: %v\n", err)
			} else {
				fmt.Printf("import %s successful\n", toCreate.ID)
			}
		}
	}()

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
