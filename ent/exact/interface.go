package exact

import "github.com/gnames/gndiff/ent/record"

type Exact interface {
	Init([]record.Record)
	Find(string) bool
}
