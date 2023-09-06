package gndiff_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/internal/io/ingestio"
	gndiff "github.com/gnames/gndiff/pkg"
	"github.com/gnames/gndiff/pkg/config"
	"github.com/stretchr/testify/assert"
)

var path = "testdata/"

// Issue #17: sorting data according to the score.
func TestScore(t *testing.T) {
	assert := assert.New(t)
	cfg := config.New()
	ing := ingestio.New(cfg)

	src := filepath.Join(path, "issue-17-src.txt")
	recSrc, err := ing.Records(src)
	assert.Nil(err)

	ref := filepath.Join(path, "issue-17-ref.txt")
	recRef, err := ing.Records(ref)
	assert.Nil(err)

	gnd := gndiff.New(cfg)
	res, err := gnd.Compare(recSrc, recRef)
	assert.Nil(err)
	assert.Equal(len(recSrc), len(res.Matches))
	obione := res.Matches[0]
	abelia := res.Matches[1]
	bubo := res.Matches[2]

	assert.Equal("Obione maritima var. maritimaa", obione.ReferenceRecords[0].CanonicalFull)
	assert.NotNil(obione.ReferenceRecords[0].ScoreDetails)
	assert.Equal("Bubo bubo Linn. 1758", bubo.ReferenceRecords[0].Name)
	assert.Equal("Abelia forrestii var. gracilenta (W.W.Sm.) Landrein", abelia.ReferenceRecords[0].Name)
}

// Issue #19: duplicated results for similar names
func TestNoDuplicates(t *testing.T) {
	assert := assert.New(t)
	cfg := config.New()
	ing := ingestio.New(cfg)

	src := filepath.Join(path, "issue-19-src.csv")
	recSrc, err := ing.Records(src)
	assert.Nil(err)

	ref := filepath.Join(path, "issue-19-ref.csv")
	recRef, err := ing.Records(ref)
	assert.Nil(err)

	gnd := gndiff.New(cfg)
	res, err := gnd.Compare(recSrc, recRef)
	assert.Nil(err)
	assert.Equal(len(recSrc), len(res.Matches))

	rrs := res.Matches[0].ReferenceRecords
	assert.Equal(2, len(rrs))
	assert.Equal("Obione maritima (Alfredo) Pacino var. maritima", rrs[0].Name)
	assert.Equal("Obione maritima (Alfredo) Pacino subsp. maritima", rrs[1].Name)

}

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

func TestNoFamily(t *testing.T) {
	assert := assert.New(t)
	cfg := config.New()
	ing := ingestio.New(cfg)

	ref := filepath.Join(path, "issue-16.csv")
	recRef, err := ing.Records(ref)
	assert.Nil(err)

	src := filepath.Join(path, "issue-16.csv")
	recSrc, err := ing.Records(src)
	assert.Nil(err)

	gnd := gndiff.New(cfg)
	res, err := gnd.Compare(recSrc, recRef)
	assert.Nil(err)
	assert.Equal(len(recSrc), len(res.Matches))
	srcRes := res.Matches[0].SourceRecord
	refRes := res.Matches[0].ReferenceRecords[0]
	assert.Equal(srcRes.ID, "")
	assert.Equal(srcRes.Family, "")
	assert.Equal(refRes.ID, "")
	assert.Equal(refRes.Family, "")
}
