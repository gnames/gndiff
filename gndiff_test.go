package gndiff_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff"
	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/stretchr/testify/assert"
)

var path = "testdata/"

func TestGNdiff(t *testing.T) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	ref := filepath.Join(path, "ebird.csv")
	recRef, err := ing.Records(ref)
	assert.Nil(t, err)

	src := filepath.Join(path, "ioc-bird.csv")
	recSrc, err := ing.Records(src)
	assert.Nil(t, err)

	gnd := gndiff.New(cfg)
	res, err := gnd.Compare(recSrc, recRef)
	assert.Nil(t, err)
	assert.Equal(t, len(recSrc), len(res.Matches))
}
