package record

import (
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser/ent/parsed"
)

type Record struct {
	DataSet       string                  `json:"dataSet"`
	Index         int                     `json:"index"`
	EditDistance  int                     `json:"editDistance"`
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Family        string                  `json:"family"`
	MatchType     verifier.MatchTypeValue `json:"matchType"`
	parsed.Parsed `json:"-"`
}
