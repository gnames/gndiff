package matcher

import (
	"github.com/gnames/gndiff/ent/dbase"
	"github.com/gnames/gndiff/ent/record"
)

type matcher struct {
	db dbase.DBase
}

func New(db dbase.DBase) Matcher {
	res := matcher{
		db: db,
	}
	return &res
}

func (m *matcher) Init(recs []record.Record) error {
	err := m.db.Init()
	if err != nil {
		return err
	}
	return m.db.Create(recs)
}

func (m *matcher) MatchExact(stem string) ([]record.Record, error) {
	var res []record.Record
	return res, nil
}

func (m *matcher) MatchFuzzy(stem string) ([]record.Record, error) {
	var res []record.Record
	return res, nil
}
