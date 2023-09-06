package dbase

import (
	"log"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnfmt"
)

type dbase struct {
	bdb *badger.DB
	enc gnfmt.GNgob
}

func New() DBase {
	var err error
	res := dbase{}
	options := badger.DefaultOptions("")
	options.Logger = nil
	options.InMemory = true
	res.bdb, err = badger.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	return &res
}

func (db *dbase) Select(stem string) ([]record.Record, error) {
	txn := db.bdb.NewTransaction(false)
	defer func() {
		err := txn.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}()
	val, err := txn.Get([]byte(stem))
	if err == badger.ErrKeyNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var gob []byte
	gob, err = val.ValueCopy(gob)
	if err != nil {
		log.Fatal(err)
	}
	var res []record.Record
	err = db.enc.Decode(gob, &res)
	return res, err
}
