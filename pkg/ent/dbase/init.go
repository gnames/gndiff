package dbase

import (
	badger "github.com/dgraph-io/badger/v2"
	"github.com/gnames/gndiff/pkg/ent/record"
)

func (db *dbase) Init(recs []record.Record) error {
	data := prepareData(recs)
	kvTxn := db.bdb.NewTransaction(true)

	for k, v := range data {
		val, err := db.enc.Encode(v)
		if err != nil {
			return err
		}
		if err = kvTxn.Set([]byte(k), val); err == badger.ErrTxnTooBig {
			err = kvTxn.Commit()
			if err != nil {
				return err
			}
			kvTxn = db.bdb.NewTransaction(true)
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
		stem := recs[i].Canonical.Stemmed
		res[stem] = append(res[stem], recs[i])
	}
	return res
}
