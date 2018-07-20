package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/happenslol/picwiz/actions"
	"github.com/spf13/afero"
	bimg "gopkg.in/h2non/bimg.v1"
)

var fs = afero.NewOsFs()
var storagePath = envy.Get("STORAGE_LOCATION", "")
var resizeOptsPortrait = bimg.Options{
	Height: 1080,
}

var resizeOptsLandscape = bimg.Options{
	Width: 1920,
}

func init() {
	buffalo.Grifts(actions.App())
}
