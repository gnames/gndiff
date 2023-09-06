package exact

import "github.com/gnames/gndiff/pkg/ent/record"

type Exact interface {
	Init([]record.Record)
	Find(string) bool
}
