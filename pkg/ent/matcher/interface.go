package matcher

import "github.com/gnames/gndiff/pkg/ent/record"

type Matcher interface {
	Init([]record.Record) error
	Match(record.Record) ([]record.Record, error)
	MatchExact(can string, spGr bool) ([]record.Record, error)
	MatchFuzzy(string, string) ([]record.Record, error)
}
