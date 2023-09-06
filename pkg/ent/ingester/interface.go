package ingester

import "github.com/gnames/gndiff/pkg/ent/record"

type Ingester interface {
	Records(path string) ([]record.Record, error)
}
