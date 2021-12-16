package dbase

import "github.com/gnames/gndiff/ent/record"

type DBase interface {
	Init() error
	Create([]record.Record) error
	Select(string) ([]record.Record, error)
}
