package matcher_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/matcher"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser"
	"github.com/stretchr/testify/assert"
)

var (
	bootstrapped bool
	m            matcher.Matcher
)

func TestMatch(t *testing.T) {
	var err error
	var recs []record.Record
	if !bootstrapped {
		initMatcher(t)
	}
	can := "Rhea americana nobilis"
	recs, err = m.MatchExact(can)
	assert.Nil(t, err)
	rec := recs[0]

	recs, err = m.Match(rec)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, verifier.Exact, recs[0].MatchType)

	gnp := gnparser.New(gnparser.NewConfig())
	rec.Name = "Rhea americanus nobilis vulgaris"
	rec.Parsed = gnp.ParseName(rec.Name)
	recs, err = m.Match(rec)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, verifier.PartialFuzzy, recs[0].MatchType)
}

func TestMatchExactFuzzy(t *testing.T) {
	if !bootstrapped {
		initMatcher(t)
	}
	var err error
	var recs []record.Record

	can := "Rhea americana nobilis"
	canSuffix := "Rhea americanus nobilis"
	canFuz := "Rhea ameriana nobilis"
	stemEx := "Rhea american nobil"
	stemFuz := "Rhea amerian nobil"
	badStr := "Something irrelevant"

	recs, err = m.MatchExact(can)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, can, recs[0].Canonical.Simple)
	assert.Equal(t, 0, recs[0].EditDistance)

	recs, err = m.MatchFuzzy(canSuffix, stemEx)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, can, recs[0].Canonical.Simple)
	assert.Equal(t, 2, recs[0].EditDistance)

	recs, err = m.MatchFuzzy(canFuz, stemFuz)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, can, recs[0].Canonical.Simple)
	assert.Equal(t, 1, recs[0].EditDistance)

	recs, err = m.MatchFuzzy(canSuffix, stemFuz)
	assert.Nil(t, err)
	assert.True(t, len(recs) == 0)

	recs, err = m.MatchFuzzy(canSuffix, badStr)
	assert.Nil(t, err)
	assert.True(t, len(recs) == 0)
}

func initMatcher(t *testing.T) {
	var err error
	var recs []record.Record
	cfg := config.New()
	ing := ingestio.New(cfg)
	p := filepath.Join("../../testdata/", "ebird.csv")
	recs, err = ing.Records(p)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)

	m = matcher.New()
	err = m.Init(recs)
	assert.Nil(t, err)
}
