package record

import "github.com/gnames/gnparser/ent/parsed"

type Record struct {
	DataSet          string
	Index            int
	ID, Name, Family string
	parsed.Parsed
}
