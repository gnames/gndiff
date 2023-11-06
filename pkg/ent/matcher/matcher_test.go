package matcher_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/internal/io/ingestio"
	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/matcher"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser"
	"github.com/stretchr/testify/assert"
)

var (
	bootstrapped bool
	m            matcher.Matcher
	gnp          gnparser.GNparser
)

func TestMatch(t *testing.T) {
	var err error
	var recs []record.Record
	var rec record.Record
	if !bootstrapped {
		initMatcher(t)
	}
	can := "Rhea americana nobilis"
	stem := "Rhea american nobil"

	recs, err = m.MatchExact(can, stem)
	assert.Nil(t, err)
	rec = recs[0]

	recs, err = m.Match(rec)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, verifier.Exact, recs[0].MatchType)

	rec.Name = "Rhea americanus nobilis vulgaris"
	rec.Parsed = gnp.ParseName(rec.Name)
	recs, err = m.Match(rec)
	assert.Nil(t, err)
	assert.True(t, len(recs) > 0)
	assert.Equal(t, verifier.PartialFuzzy, recs[0].MatchType)
}

func TestSpGroup(t *testing.T) {
	assert := assert.New(t)
	if !bootstrapped {
		initMatcher(t)
	}
	var rec record.Record

	rec.Name = "Apteryx mantelli A. D. Bartlett, 1852 mantelli"
	rec.Parsed = gnp.ParseName(rec.Name)
	recs, err := m.Match(rec)
	assert.Nil(err)
	res := recs[0]
	assert.Equal(verifier.ExactSpeciesGroup.String(), res.MatchType.String())
	assert.Equal(0, res.EditDistance)

	rec = record.Record{}
	rec.Name = "Apteryx mantellus A. D. Bartlett, 1852 mantellus"
	rec.Parsed = gnp.ParseName(rec.Name)
	recs, err = m.Match(rec)
	assert.Nil(err)
	res = recs[0]
	assert.Equal(verifier.FuzzySpeciesGroup.String(), res.MatchType.String())
	assert.Equal(2, res.EditDistance)

	rec = record.Record{}
	rec.Name = "Apteryx maantelli A. D. Bartlett, 1852 maantelli"
	rec.Parsed = gnp.ParseName(rec.Name)
	recs, err = m.Match(rec)
	assert.Nil(err)
	res = recs[0]
	assert.Equal(verifier.FuzzySpeciesGroup.String(), res.MatchType.String())
	assert.Equal(1, res.EditDistance)
}

func TestMatchExactFuzzy(t *testing.T) {
	assert := assert.New(t)
	if !bootstrapped {
		initMatcher(t)
	}
	var err error
	var recs []record.Record

	can := "Rhea americana nobilis"
	canSuffix := "Rhea americanus nobilis"
	stemEx := "Rhea american nobil"

	recs, err = m.MatchExact(can, stemEx)
	assert.Nil(err)
	assert.True(len(recs) > 0)
	assert.Equal(can, recs[0].Canonical.Simple)
	assert.Equal(0, recs[0].EditDistance)

	recs, err = m.MatchExact(canSuffix, stemEx)
	assert.Nil(err)
	assert.True(len(recs) > 0)
	assert.Equal(can, recs[0].Canonical.Simple)
	assert.Equal(verifier.Fuzzy, recs[0].MatchType)
	assert.Equal(2, recs[0].EditDistance)

	canFuz := "Rhea ameriana nobilis"
	stemFuz := "Rhea amerian nobil"
	recs, err = m.MatchFuzzy(canFuz, stemFuz)
	assert.Nil(err)
	assert.True(len(recs) > 0)
	assert.Equal(can, recs[0].Canonical.Simple)
	assert.Equal(1, recs[0].EditDistance)

	canFuzzy := "Rhea amerianus nobilis"
	recs, err = m.MatchFuzzy(canFuzzy, stemFuz)
	assert.Nil(err)
	assert.True(len(recs) > 0)
	assert.Equal(3, recs[0].EditDistance)

	badStr := "Something irrelevant"
	recs, err = m.MatchFuzzy(canSuffix, badStr)
	assert.Nil(err)
	assert.True(len(recs) == 0)
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
	gnp = gnparser.New(
		gnparser.NewConfig(gnparser.OptWithSpeciesGroupCut(true)),
	)

	m = matcher.New(cfg)
	err = m.Init(recs)
	assert.Nil(t, err)
}
