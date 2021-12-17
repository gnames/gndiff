package kval

import (
	badger "github.com/dgraph-io/badger/v2"
	"github.com/gnames/gndiff/ent/record"
)

func (kv *kval) Create(recs []record.Record) error {
	data := prepareData(recs)
	kvTxn := kv.bdb.NewTransaction(true)

	for k, v := range data {
		val, err := kv.enc.Encode(v)
		if err != nil {
			return err
		}
		if err = kvTxn.Set([]byte(k), val); err == badger.ErrTxnTooBig {
			err = kvTxn.Commit()
			if err != nil {
				return err
			}
			kvTxn = kv.bdb.NewTransaction(true)
			err = kvTxn.Set([]byte(k), val)
			if err != nil {
				return err
			}
		}
	}
	return kvTxn.Commit()
}

func prepareData(recs []record.Record) map[string][]record.Record {
	res := make(map[string][]record.Record)

	for i := range recs {
		stem := recs[i].Canonical.Simple
		res[stem] = append(res[stem], recs[i])
	}
	return res
}
