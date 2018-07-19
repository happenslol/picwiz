package grifts

import (
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("imports", func() {
	grift.Desc("process", "Processes all pending imports")
	grift.Add("process", func(c *grift.Context) error {
		return nil
	})
})
