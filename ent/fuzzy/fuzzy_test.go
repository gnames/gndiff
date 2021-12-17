package fuzzy_test

import (
	"fmt"
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

// EditDist without constraints
func TestDist(t *testing.T) {
	// // to hide warnings
	// log.SetLevel(log.FatalLevel)

	testData := []struct {
		str1, str2 string
		dist       int
	}{
		{"Hello", "Hello", 0},
		{"Pomatomus", "Pom-tomus", 1},
		{"Pomatomus", "Pomщtomus", 1},
		// ed = 3, too big
		{"sitting", "kitten", -1},
		// words are too small
		{"Pom atomus", "Poma tomus", -1},
		{"Acacia mal", "Acacia may", -1},
		// differnt number of words is not covered yet
		{"Pomatomus", "Poma  tomus", 2},
		// edge cases that should not happen
		// more than one empty space
		{"Pomatomus saltator", "Pomatomus  saltator", 1},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist := fuzzy.EditDistance(v.str1, v.str2, false)
		assert.Equal(t, dist, v.dist, msg)
	}
}
