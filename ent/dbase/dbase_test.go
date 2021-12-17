package dbase_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/dbase"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var err error
	var recs []record.Record
	db := dbase.New()
	cfg := config.New()
	ing := ingestio.New(cfg)
	p := filepath.Join("../../testdata/", "ebird.csv")
	recs, err = ing.Records(p)
	assert.Nil(t, err)

	err = db.Init(recs)
	assert.Nil(t, err)
	recs, err = db.Select("Rhea americana nobilis")
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, "85c94df5-33a0-5550-89d2-216a7e75e564", recs[0].VerbatimID)
}
