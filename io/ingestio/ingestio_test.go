package ingestio_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/stretchr/testify/assert"
)

var path = "../../testdata/"

func TestRecordsBad(t *testing.T) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	p := filepath.Join(path, "nofile")
	rec, err := ing.Records(p)
	assert.Contains(t, err.Error(), "does not exist")
	assert.NotNil(t, err)
	assert.Nil(t, rec)

	p = filepath.Join(path, "bad-header.csv")
	rec, err = ing.Records(p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "needs `scientifiName` field")
	assert.Nil(t, rec)
}

func TestRecords(t *testing.T) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	p := filepath.Join(path, "ebird.csv")
	rec, err := ing.Records(p)
	assert.Nil(t, err)
	assert.True(t, len(rec) > 1000)
	assert.Equal(t, "Rhea americana nobilis", rec[10].Name)
	assert.Equal(t, "Rheidae (Rheas)", rec[10].Family)
	assert.Equal(t, "gn_11", rec[10].ID)
	assert.True(t, rec[10].Parsed.Parsed)
	assert.Equal(t, "Rhea american nobil", rec[10].Canonical.Stemmed)
}
