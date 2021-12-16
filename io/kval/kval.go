package kval

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/dbase"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gnfmt"
)

type kval struct {
	cfg config.Config
	bdb *badger.DB
	enc gnfmt.GNgob
}

func New(cfg config.Config) dbase.DBase {
	res := kval{cfg: cfg}
	return &res
}

func (kv *kval) Init() error {
	var err error
	options := badger.DefaultOptions("")
	options.Logger = nil
	options.InMemory = true
	kv.bdb, err = badger.Open(options)
	return err
}

func (kv *kval) Select(stem string) ([]record.Record, error) {
	txn := kv.bdb.NewTransaction(false)
	defer func() {
		err := txn.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}()
	val, err := txn.Get([]byte(stem))
	if err == badger.ErrKeyNotFound {
		err = fmt.Errorf("%s not found", stem)
		return nil, err
	} else if err != nil {
		return nil, err
	}
	var gob []byte
	gob, err = val.ValueCopy(gob)
	if err != nil {
		log.Fatal(err)
	}
	var res []record.Record
	err = kv.enc.Decode(gob, &res)
	return res, err
}
