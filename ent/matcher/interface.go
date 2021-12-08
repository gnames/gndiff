package Matcher

import "github.com/gnames/gndiff/ent/record"

type Matcher interface {
	Init([]record.Record) error
  MatchExact(record.Record) []record.Record
  MatchFuzzy(record.Record) []record.Record
}
