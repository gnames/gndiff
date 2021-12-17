package fuzzy_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/fuzzy"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/stretchr/testify/assert"
)

func TestFuzzy(t *testing.T) {
	var err error
	var recs []record.Record
	cfg := config.New()
	ing := ingestio.New(cfg)
	p := filepath.Join("../../testdata/", "ebird.csv")
	recs, err = ing.Records(p)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)

	ex := fuzzy.New()
	err = ex.Init(recs)
	assert.Nil(t, err)
	cans := ex.FindExact("Rhea american nobil")
	assert.True(t, len(cans) > 0)
	noCans := ex.FindExact("Not a name")
	assert.True(t, len(noCans) == 0)
	cans = ex.FindFuzzy("Rhea ameican nobil")
	assert.True(t, len(cans) > 0)
}
