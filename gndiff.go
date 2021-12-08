package gndiff

import (
	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/ingester"
	"github.com/gnames/gndiff/ent/output"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gnlib/ent/gnvers"
)

type gndiff struct {
	Differ
	ingester.Ingester
	config.Config
}

func New() GNdiff {
	res := gndiff{}
	return &res
}

func (gnd *gndiff) Compare(rec1, rec2 []record.Record) []output.Output {
	return nil
}

// Version function returns version number of `gnparser` and the timestamp
// of its build.
func (gnd *gndiff) GetVersion() gnvers.Version {
	version := Version
	build := Build
	return gnvers.Version{Version: version, Build: build}
}
