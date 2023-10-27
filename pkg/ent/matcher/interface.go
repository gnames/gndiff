package matcher

import "github.com/gnames/gndiff/pkg/ent/record"

// Matcher creates a queryable reference dataset and facilitates
// exact and fuzzy matching.
type Matcher interface {
	Init([]record.Record) error
	Match(record.Record) ([]record.Record, error)
	MatchExact(can, stem string) ([]record.Record, error)
	MatchFuzzy(can, stem string) ([]record.Record, error)
}
