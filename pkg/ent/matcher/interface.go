package matcher

import "github.com/gnames/gndiff/pkg/ent/record"

// Matcher creates a queryable reference dataset and facilitates
// exact and fuzzy matching.
type Matcher interface {
	// Init creates lookup key-value store for reference, as well as builds
	// a Trie out of stemmed canonicals from the rereference.
	// The Trie is used for exact and fuzzy matching.
	Init([]record.Record) error

	// Match takes a record and returns back matched records from the reference.
	Match(record.Record) ([]record.Record, error)

	// MatchExact returns records that match **stemmed canonical** exactly.
	// The final match type would depend how stemmed and canonical versions
	// differ.
	MatchExact(can, stem string) ([]record.Record, error)

	// Match Fuzzy returns records that match **stemmed canonical** with some
	// differences between reference stem and given stem.
	// The final match type might differ depending on how stemmed and canonical
	// versions differ.
	MatchFuzzy(can, stem string) ([]record.Record, error)
}
