package dbase

import "github.com/gnames/gndiff/ent/record"

type DBase interface {
	Init([]record.Record) error
	Select(string) ([]record.Record, error)
}
