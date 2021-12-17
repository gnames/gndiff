package fuzzy

import (
	"sort"

	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gndiff/ent/record"
)

type fuzzy struct {
	trie       *levenshtein.MinTree
	canonicals map[string][]string
}

func New() Fuzzy {
	res := fuzzy{canonicals: make(map[string][]string)}
	return &res
}

func (f *fuzzy) Init(recs []record.Record) error {
	var err error
	stems := make([]string, len(recs))
	for i := range recs {
		stem := recs[i].Canonical.Stemmed
		stems[i] = stem
		f.canonicals[stem] = append(f.canonicals[stem], recs[i].Canonical.Simple)
	}
	sort.Strings(stems)
	f.trie, err = levenshtein.NewMinTree(stems)
	if err != nil {
		return err
	}
	return nil
}

func (f *fuzzy) FindExact(stem string) []string {
	return f.find(stem, 0)
}

func (f *fuzzy) FindFuzzy(stem string) []string {
	return f.find(stem, 1)
}

func (f *fuzzy) find(stem string, maxDist int) []string {
	stems := f.trie.FuzzyMatches(stem, maxDist)
	res := make([]string, 0, len(stems)*2)
	for i := range stems {
		res = append(res, f.canonicals[stems[i]]...)
	}
	return res
}
