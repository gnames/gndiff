package record

import "github.com/gnames/gnparser/ent/parsed"

type Record struct {
	DataSet             string
	Index, EditDistance int
	ID, Name, Family    string
	parsed.Parsed
}
