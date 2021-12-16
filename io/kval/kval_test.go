package kval_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/gnames/gndiff/io/kval"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var recs []record.Record
	cfg := config.New()
	db := kval.New(cfg)
	err := db.Init()
	assert.Nil(t, err)
	ing := ingestio.New(cfg)
	p := filepath.Join("../../testdata/", "ebird.csv")
	recs, err = ing.Records(p)
	assert.Nil(t, err)

	err = db.Create(recs)
	assert.Nil(t, err)
	recs, err = db.Select("Rhea american nobil")
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, "Rhea americana nobilis", recs[0].Name)
}
