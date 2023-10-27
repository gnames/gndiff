package matcher

import (
	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/dbase"
	"github.com/gnames/gndiff/pkg/ent/fuzzy"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnlib/ent/verifier"
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
		matchType := verifier.Exact
		if can != res[i].CanonicalSimple {
			matchType = verifier.Fuzzy
			ed = fuzzy.EditDistance(
				can,
				res[i].CanonicalSimple,
				true,
			)
		}
		matchType = spGroupMatchType(can, stem, matchType)
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
			matchType := spGroupMatchType(can, stem, res[i].MatchType)
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
			ed := fuzzy.EditDistance(can, recs[ii].Canonical.Simple, noCheck)
			if ed < 0 {
				continue
			}

			recs[ii].EditDistance = ed
			recs[ii].MatchType = verifier.Fuzzy
			res = append(res, recs[ii])
		}
	}
	return res, err
}
