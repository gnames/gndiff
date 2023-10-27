package web_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gndiff/internal/io/web"
	"github.com/gnames/gndiff/pkg/ent/output"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/gnvers"
	"github.com/gnames/gnlib/ent/verifier"
	"github.com/stretchr/testify/assert"
)

var restURL = "http://0.0.0.0:8080/api/v0/"

func TestPing(t *testing.T) {
	resp, err := http.Get(restURL + "ping")
	assert.Nil(t, err)

	response, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	assert.Equal(t, "pong", string(response))
}

func TestVersion(t *testing.T) {
	resp, err := http.Get(restURL + "version")
	assert.Nil(t, err)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	enc := gnfmt.GNjson{}
	var response gnvers.Version
	err = enc.Decode(respBytes, &response)
	assert.Nil(t, err)
	assert.Regexp(t, `^v\d+\.\d+\.\d+`, response.Version)
}

func TestDiff(t *testing.T) {
	assert := assert.New(t)
	var response output.Output
	qryFile := "../../../pkg/testdata/ebird-100.csv"
	refFile := "../../../pkg/testdata/ioc-bird-100.csv"
	qryFileName := filepath.Base(qryFile)
	refFileName := filepath.Base(refFile)
	qryTxt, err := os.ReadFile(qryFile)
	assert.Nil(err)
	refTxt, err := os.ReadFile(refFile)

	request := web.Input{
		QueryFileName:     qryFileName,
		QueryText:         string(qryTxt),
		ReferenceFileName: refFileName,
		ReferenceText:     string(refTxt),
	}

	req, err := gnfmt.GNjson{}.Encode(request)
	assert.Nil(err)
	r := bytes.NewReader(req)
	resp, err := http.Post(restURL+"diff", "application/json", r)
	assert.Nil(err)
	respBytes, err := io.ReadAll(resp.Body)

	assert.Nil(err)
	err = gnfmt.GNjson{}.Decode(respBytes, &response)
	assert.Nil(err)

	assert.Greater(response.Metadata.TimeTotalSec, float64(0))
	assert.Greater(len(response.Matches), 50)
	var count int
	for _, v := range response.Matches {
		fmt.Println(v.ReferenceRecords[0].MatchType.String())
		if v.ReferenceRecords[0].MatchType == verifier.PartialExact {
			count++
		}
	}
	assert.Greater(count, 50)
}
