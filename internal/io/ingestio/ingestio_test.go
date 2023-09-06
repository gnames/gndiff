package ingestio_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/internal/io/ingestio"
	"github.com/gnames/gndiff/pkg/config"
	"github.com/stretchr/testify/assert"
)

var path = "../../../pkg/testdata/"

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
	assert.Contains(t, err.Error(), "needs `scientificName` field")
	assert.Nil(t, rec)
}

func TestRecords(t *testing.T) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	p := filepath.Join(path, "ebird.csv")
	rec, err := ing.Records(p)
	assert.Nil(t, err)
	assert.True(t, len(rec) > 1000)
	assert.Equal(t, "Rhea americana nobilis", rec[9].Name)
	assert.Equal(t, "Rheidae (Rheas)", rec[9].Family)
	assert.Equal(t, "gn_11", rec[9].ID)
	assert.True(t, rec[9].Parsed.Parsed)
	assert.Equal(t, "Rhea american nobil", rec[9].Canonical.Stemmed)
}

func TestTSV(t *testing.T) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	p := filepath.Join(path, "ioc-bird.tsv")
	rec, err := ing.Records(p)
	assert.Nil(t, err)
	assert.True(t, len(rec) > 2)
	assert.Equal(t, "Rhea americana (Linnaeus, 1758)", rec[2].Name)
}

func TestNamesList(t *testing.T) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	p := filepath.Join(path, "names.txt")
	rec, err := ing.Records(p)
	assert.Nil(t, err)
	assert.True(t, len(rec) > 2)
	assert.Equal(t, "Rhea americana (Linnaeus, 1758)", rec[2].Name)
}
