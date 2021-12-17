package fuzzy

import "github.com/gnames/gndiff/ent/record"

type Fuzzy interface {
	Init([]record.Record) error
	FindExact(string) []string
	FindFuzzy(string) []string
}
