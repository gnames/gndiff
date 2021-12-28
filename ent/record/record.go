package record

import (
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnparser/ent/parsed"
)

type Record struct {
	DataSet       string                  `json:"dataSet"`
	Index         int                     `json:"index"`
	EditDistance  int                     `json:"editDistance,omitempty"`
	ID            string                  `json:"id,omitempty"`
	Name          string                  `json:"name"`
	Family        string                  `json:"family,omitempty"`
	MatchType     verifier.MatchTypeValue `json:"matchType,omitempty"`
	parsed.Parsed `json:"-"`
}
