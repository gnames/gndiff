package exact

import (
	"sync"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gndiff/pkg/ent/record"
)

type exact struct {
	canonical     *baseBloomfilter.Bloomfilter
	canonicalSize uint
	mux           sync.Mutex
}

func New() Exact {
	res := exact{}
	return &res
}

func (e *exact) Init(recs []record.Record) {
	e.canonicalSize = uint(len(recs))
	cfg := bloomfilter.Config{
		N:        e.canonicalSize,
		P:        0.00001,
		HashName: bloomfilter.HASHER_OPTIMAL,
	}
	bf := baseBloomfilter.New(cfg)

	for i := range recs {
		bf.Add([]byte(recs[i].Canonical.Simple))
	}
	e.canonical = bf
}

func (e *exact) Find(s string) bool {
	e.mux.Lock()
	isIn := e.canonical.Check([]byte(s))
	e.mux.Unlock()
	return isIn
}
