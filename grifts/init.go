package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/happenslol/picwiz/actions"
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()
var storagePath = envy.Get("STORAGE_LOCATION", "")

func init() {
	buffalo.Grifts(actions.App())
}
