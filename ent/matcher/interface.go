package matcher

import "github.com/gnames/gndiff/ent/record"

type Matcher interface {
	Init([]record.Record) error
	Match(record.Record) ([]record.Record, error)
	MatchExact(string) ([]record.Record, error)
	MatchFuzzy(string, string) ([]record.Record, error)
}
