package matcher

import (
	"github.com/gnames/gndiff/ent/dbase"
	"github.com/gnames/gndiff/ent/exact"
	"github.com/gnames/gndiff/ent/fuzzy"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gnlib/ent/verifier"
)

type matcher struct {
	db dbase.DBase
	e  exact.Exact
	f  fuzzy.Fuzzy
}

func New() Matcher {
	res := matcher{
		db: dbase.New(),
		e:  exact.New(),
		f:  fuzzy.New(),
	}
	return &res
}

func (m *matcher) Init(recs []record.Record) error {
	err := m.db.Init(recs)
	if err == nil {
		m.e.Init(recs)
		err = m.f.Init(recs)
	}
	return err
}

func (m *matcher) MatchExact(canonical string) ([]record.Record, error) {
	var err error
	var res []record.Record
	if m.e.Find(canonical) {
		res, err = m.db.Select(canonical)
	}
	for i := range res {
		res[i].MatchType = verifier.Exact
	}
	return res, err
}

func (m *matcher) MatchFuzzy(can, stem string) ([]record.Record, error) {
	var res []record.Record
	var canonicals []string
	if canonicals = m.f.FindExact(stem); len(canonicals) > 0 {
		return m.fetchCanonicals(can, canonicals, true)
	}
	if canonicals = m.f.FindFuzzy(stem); len(canonicals) > 0 {
		return m.fetchCanonicals(can, canonicals, false)
	}
	return res, nil
}

func (m *matcher) fetchCanonicals(can string, cans []string, noCheck bool) ([]record.Record, error) {
	var err error
	var recs, res []record.Record
	for i := range cans {
		recs, err = m.db.Select(cans[i])
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
