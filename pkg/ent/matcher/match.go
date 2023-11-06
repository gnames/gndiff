package matcher

import (
	"strings"

	"github.com/gnames/gndiff/pkg/ent/fuzzy"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser/ent/stemmer"
)

func (m *matcher) Match(rec record.Record) ([]record.Record, error) {
	res, err := m.MatchExact(rec.Canonical.Simple, rec.Canonical.Stemmed)
	if len(res) > 0 || err != nil {
		return res, err
	}

	res, err = m.MatchFuzzy(rec.Canonical.Simple, rec.Canonical.Stemmed)
	if len(res) > 0 || err != nil {
		return res, err
	}

	if rec.Cardinality > 1 || m.cfg.WithUninomialFuzzyMatch {
		res, err = m.partialMatch(rec.Canonical.Simple, rec.Canonical.Stemmed)
	}

	return res, err
}

type canPair struct {
	can, stem   string
	cardinality int
}

func (m *matcher) partialMatch(can, stem string) ([]record.Record, error) {
	var res []record.Record
	var err error
	cans := partialCombos(can, stem)
	for i := range cans {
		res, err = m.MatchExact(cans[i].can, cans[i].stem)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			matchType := verifier.Exact
			for ii := range res {
				if cans[i].can != res[ii].CanonicalSimple {
					matchType = verifier.PartialFuzzy
				}
				res[ii].MatchType = matchType
			}
			return res, err
		}

		if cans[i].cardinality == 1 && !m.cfg.WithUninomialFuzzyMatch {
			break
		}

		res, err = m.MatchFuzzy(cans[i].can, cans[i].stem)
		if len(res) > 0 || err != nil {
			for i := range res {
				res[i].MatchType = verifier.PartialFuzzy
			}
			return res, err
		}
	}

	return res, err
}

func checkFuzzyMatch(
	can, canRes string,
	mt verifier.MatchTypeValue,
) (verifier.MatchTypeValue, int) {
	var ed int
	if can != canRes {
		mt = verifier.Fuzzy
		ed = fuzzy.EditDistance(can, canRes, true)
	}
	return mt, ed
}

func checkSpGrMatch(
	can, stem, catRes string,
	mt verifier.MatchTypeValue,
) verifier.MatchTypeValue {
	if stemmer.StemCanonical(can) != stem {
		switch mt {
		case verifier.Exact:
			mt = verifier.ExactSpeciesGroup
		case verifier.Fuzzy:
			mt = verifier.FuzzySpeciesGroup
		}
	}
	return mt
}

func partialCombos(can, stem string) []canPair {
	canWs := strings.Split(can, " ")
	stemWs := strings.Split(stem, " ")
	switch len(canWs) {
	case 2:
		return []canPair{
			{
				can:         canWs[0],
				stem:        stemWs[0],
				cardinality: 1,
			},
		}
	case 3:
		return []canPair{
			{
				can:         canWs[0] + " " + canWs[2],
				stem:        stemWs[0] + " " + stemWs[2],
				cardinality: 2,
			},
			{
				can:         canWs[0] + " " + canWs[1],
				stem:        stemWs[0] + " " + stemWs[1],
				cardinality: 2,
			},
			{
				can:         canWs[0],
				stem:        canWs[0],
				cardinality: 1,
			},
		}
	case 4:
		return []canPair{
			{
				can:         canWs[0] + " " + canWs[3],
				stem:        stemWs[0] + " " + stemWs[3],
				cardinality: 2,
			},
			{
				can:         canWs[0] + " " + canWs[2],
				stem:        stemWs[0] + " " + stemWs[2],
				cardinality: 2,
			},
			{
				can:         canWs[0] + " " + canWs[1] + " " + canWs[2],
				stem:        stemWs[0] + " " + stemWs[1] + " " + stemWs[2],
				cardinality: 2,
			},
			{
				can:         canWs[0] + " " + canWs[1],
				stem:        stemWs[0] + " " + stemWs[1],
				cardinality: 2,
			},
			{
				can:         canWs[0],
				stem:        canWs[0],
				cardinality: 1,
			},
		}
	default:
		return nil
	}
}
