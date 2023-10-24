package gndiff

import (
	"github.com/gnames/gndiff/pkg/ent/output"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnlib/ent/gnvers"
)

// GNdiff is an interface that enables comparison of two files that contain
// scientific names. Comparison tries to find matched names, and indicates
// type of match.
type GNdiff interface {
	Differ
	GetVersion() gnvers.Version
}

// Differ contains method to match query and reference. `Query`
// contains scientific name-string to match, while `reference`
// contains scientific name-string to match against.
type Differ interface {
	// Compare takes a query
	// should be in a csv, tsv, or plain text format.
	Compare(
		query []record.Record,
		reference []record.Record,
	) (output.Output, error)
}
