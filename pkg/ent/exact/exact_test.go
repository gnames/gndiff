package exact_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/internal/io/ingestio"
	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/exact"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	var err error
	var recs []record.Record
	cfg := config.New()
	ing := ingestio.New(cfg)
	p := filepath.Join("../../testdata/", "ebird.csv")
	recs, err = ing.Records(p)
	assert.Nil(t, err)

	ex := exact.New()
	ex.Init(recs)
	isIn := ex.Find("Rhea americana nobilis")
	assert.True(t, isIn)
	notIn := ex.Find("Not a name")
	assert.False(t, notIn)
}
