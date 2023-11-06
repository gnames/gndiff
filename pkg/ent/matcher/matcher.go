package matcher

import (
	"strings"

	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/dbase"
	"github.com/gnames/gndiff/pkg/ent/fuzzy"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser/ent/stemmer"
)

type matcher struct {
	cfg config.Config
	db  dbase.DBase
	f   fuzzy.Fuzzy
}

func New(cfg config.Config) Matcher {
	res := matcher{
		cfg: cfg,
		db:  dbase.New(),
		f:   fuzzy.New(),
	}
	return &res
}

func (m *matcher) Init(recs []record.Record) error {
	err := m.db.Init(recs)
	if err == nil {
		err = m.f.Init(recs)
	}
	return err
}

func (m *matcher) MatchExact(can, stem string) ([]record.Record, error) {
	var err error
	var res []record.Record
	if len(m.f.FindExact(stem)) > 0 {
		res, err = m.db.Select(stem)
	}

	for i := range res {
		var ed int
		canRes := res[i].CanonicalSimple
		matchType := verifier.Exact

		// check if match was by species group.
		matchType, ed = checkSpGrMatch(
			can, stem, canRes,
			res[i].EditDistance,
			matchType,
		)

		// if did not match to species group, see if exact match by stem
		// is actually a fuzzy match by simple canonical.
		if matchType == verifier.Exact {
			matchType, ed = checkFuzzyMatch(
				can, canRes, matchType,
			)
		}

		res[i].MatchType = matchType
		res[i].EditDistance = ed
	}

	return res, err
}

func (m *matcher) MatchFuzzy(can, stem string) ([]record.Record, error) {
	var err error
	var res []record.Record
	var resStems []string
	if resStems = m.f.FindFuzzy(stem); len(resStems) > 0 {
		res, err = m.fetchFuzzyCanonicals(can, stem, resStems, false)
		if err != nil {
			return res, err
		}
		for i := range res {
			matchType := verifier.Fuzzy
			canRes := res[i].CanonicalSimple
			matchType, ed := checkSpGrMatch(
				can, stem, canRes,
				res[i].EditDistance,
				res[i].MatchType,
			)
			if matchType == verifier.Fuzzy {
				ed = fuzzy.EditDistance(can, canRes, true)
			}
			res[i].EditDistance = ed
			res[i].MatchType = matchType
		}

		return res, err
	}

	return res, nil
}

func (m *matcher) fetchFuzzyCanonicals(
	can, stem string,
	resStems []string,
	noCheck bool,
) ([]record.Record, error) {
	var err error
	var recs, res []record.Record
	for i := range resStems {
		recs, err = m.db.Select(resStems[i])
		if err != nil {
			return res, err
		}
		for ii := range recs {
			res = append(res, recs[ii])
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
	can, stem, canRes string,
	ed int,
	mt verifier.MatchTypeValue,
) (verifier.MatchTypeValue, int) {

	// stem is shortened for species group trinomials. For example
	// for `Bubo bubo bubo` stem would be `Bubo bub`, while
	// stemmer.StemCanonical would return `Bubo bub bub`.
	if stemmer.StemCanonical(can) == stem {
		return mt, ed
	}

	elCan := strings.Split(can, " ")
	elCanRes := strings.Split(canRes, " ")

	l := min(len(elCan), len(elCanRes))

	can = strings.Join(elCan[:l], " ")
	canRes = strings.Join(elCanRes[:l], " ")

	if canRes == can {
		return verifier.ExactSpeciesGroup, ed
	}

	ed = fuzzy.EditDistance(can, canRes, true)
	return verifier.FuzzySpeciesGroup, ed
}
