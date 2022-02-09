package ingestio

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/ingester"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnsys"
)

type ingestio struct {
	config.Config
}

func New(cfg config.Config) ingester.Ingester {
	return ingestio{Config: cfg}
}

func (ing ingestio) Records(path string) ([]record.Record, error) {
	exists, err := gnsys.FileExists(path)
	if !exists {
		return nil, fmt.Errorf("file '%s' does not exist", path)
	}
	if err != nil {
		return nil, err
	}

	fileName := strings.Split(filepath.Base(path), ".")[0]
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, sep, err := tryNamesOnly(f, fileName)
	if err != nil {
		return nil, err
	}

	if records != nil {
		return records, nil
	}

	if sep == rune(0) {
		return nil, errors.New("cannot determine field separator")
	}

	r := csv.NewReader(f)
	r.Comma = sep

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	name, id, family, valid := readHeader(header)
	if !valid {
		return nil, errors.New("the CSV file needs `scientificName` field")
	}
	var count int
	var res []record.Record
	for {
		count++
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		rec := record.Record{
			DataSet: fileName,
			Index:   count,
			ID:      row[id],
			Name:    row[name],
			Family:  row[family],
		}
		res = append(res, rec)
	}
	return parse(res), nil
}

func readHeader(s []string) (int, int, int, bool) {
	var name, id, family int
	var valid bool
	// remove BOM character if exists
	if len(s) > 0 && len(s[0]) > 3 && s[0][0:3] == "\xef\xbb\xbf" {
		s[0] = s[0][3:]
	}
	for i := range s {
		field := strings.ToLower(s[i])
		switch field {
		case "scientificname":
			name = i
			valid = true
		case "taxonid":
			id = i
		case "family":
			family = i
		}
	}
	return name, id, family, valid
}

func parse(recs []record.Record) []record.Record {
	res := make([]record.Record, 0, len(recs))
	names := make([]string, len(recs))
	for i := range recs {
		names[i] = recs[i].Name
	}
	cfg := gnparser.NewConfig(gnparser.OptJobsNum(100))
	gnp := gnparser.New(cfg)
	parsed := gnp.ParseNames(names)
	for i := range recs {
		recs[i].Parsed = parsed[i]
		if recs[i].Parsed.Parsed {
			res = append(res, recs[i])
		}
	}
	return res
}

func tryNamesOnly(f *os.File, fileName string) ([]record.Record, rune, error) {
	var res []record.Record
	scanner := bufio.NewScanner(f)

	var count int

	for scanner.Scan() {
		line := scanner.Text()
		if count == 0 {
			sep := fileSep(line)
			if sep != rune(0) {
				f.Seek(0, io.SeekStart)
				return nil, sep, nil
			}
		}

		res = append(res, record.Record{DataSet: fileName, Index: count, Name: line})
		count++
	}

	if err := scanner.Err(); err != nil {
		return res, rune(0), err
	}
	return parse(res), rune(0), nil
}

func fileSep(s string) rune {
	if strings.Contains(s, "\t") {
		return '\t'
	} else if !strings.Contains(s, ",") || strings.Contains(s, " ") {
		return rune(0)
	} else {
		return ','
	}
}
