package gndiff

import (
	"github.com/gnames/gndiff/ent/output"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gnlib/ent/gnvers"
)

// GNdiff is an interface that enables comparison of two files that contain
// scientific names. Comparison tries to find matched names, and indicates
// type of match.
type GNdiff interface {
	Differ
	GetVersion() gnvers.Version
}

type Differ interface {
	Compare([]record.Record, []record.Record) []output.Output
}
