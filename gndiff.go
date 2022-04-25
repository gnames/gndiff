package gndiff

import (
	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/matcher"
	"github.com/gnames/gndiff/ent/output"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gnlib/ent/gnvers"
)

type gndiff struct {
	cfg config.Config
}

func New(cfg config.Config) GNdiff {
	res := gndiff{
		cfg: cfg,
	}
	return &res
}

func (gnd *gndiff) Compare(source, reference []record.Record) (output.Output, error) {
	var err error
	var recs []record.Record
	res := make([]output.Match, len(source))
	m := matcher.New()
	err = m.Init(reference)
	if err != nil {
		return output.Output{}, err
	}
	for i := range source {
		recs, err = m.Match(source[i])
		if err != nil {
			return output.Output{}, err
		}
		res[i].SourceRecord = source[i]
		res[i].ReferenceRecords = recs
	}
	sortByScore(res)

	return output.Output{Matches: res}, nil
}

// Version function returns version number of `gnparser` and the timestamp
// of its build.
func (gnd *gndiff) GetVersion() gnvers.Version {
	version := Version
	build := Build
	return gnvers.Version{Version: version, Build: build}
}
