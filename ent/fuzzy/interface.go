package fuzzy

import "github.com/gnames/gndiff/ent/record"

type Fuzzy interface {
	Init([]record.Record)
	Find(string) []string
}
