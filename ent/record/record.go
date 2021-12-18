package record

import (
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser/ent/parsed"
)

type Record struct {
	DataSet             string
	Index, EditDistance int
	ID, Name, Family    string
	MatchType           verifier.MatchTypeValue
	parsed.Parsed
}
