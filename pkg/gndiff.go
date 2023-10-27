package gndiff

import (
	"time"

	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/matcher"
	"github.com/gnames/gndiff/pkg/ent/output"
	"github.com/gnames/gndiff/pkg/ent/record"
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

func (gnd *gndiff) Compare(query, reference []record.Record) (output.Output, error) {
	startTime := time.Now()
	var err error
	var recs []record.Record
	res := make([]output.Match, len(query))
	m := matcher.New(gnd.cfg)
	err = m.Init(reference)
	if err != nil {
		return output.Output{}, err
	}
	initTime := time.Since(startTime)
	for i := range query {
		recs, err = m.Match(query[i])
		if err != nil {
			return output.Output{}, err
		}
		res[i].QueryRecord = query[i]
		res[i].ReferenceRecords = recs
	}
	sortByScore(res)
	totalTime := time.Since(startTime)

	out := output.Output{
		Metadata: output.Metadata{
			TimeTotalSec:  totalTime.Seconds(),
			TimeIngestSec: initTime.Seconds(),
			TimeQuerySec:  (totalTime - initTime).Seconds(),
		},
		Matches: res,
	}

	return out, nil
}

// Version function returns version number of `gnparser` and the timestamp
// of its build.
func (gnd *gndiff) GetVersion() gnvers.Version {
	version := Version
	build := Build
	return gnvers.Version{Version: version, Build: build}
}
