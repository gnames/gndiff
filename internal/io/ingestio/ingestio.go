package ingestio

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gndiff/pkg/ent/ingester"
	"github.com/gnames/gndiff/pkg/ent/record"
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

	ext := strings.ToLower(filepath.Ext(f.Name()))

	records, sep, err := ing.RecordsFromText(f, fileName, ext)
	if err != nil {
		return nil, err
	}

	if records != nil {
		return records, nil
	}

	if sep == rune(0) {
		return nil, errors.New("cannot determine field separator")
	}

	// rewind file back
	f.Seek(0, io.SeekStart)
	return ing.RecordsFromCSV(f, sep, fileName)
}

func (ing ingestio) RecordsFromText(
	r io.Reader,
	fileName, ext string,
) ([]record.Record, rune, error) {
	var res []record.Record
	scanner := bufio.NewScanner(r)
	var count int
	for scanner.Scan() {
		line := scanner.Text()
		if count == 0 {
			sep := fileSep(line, ext)
			if sep != rune(0) {
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

func (ing ingestio) RecordsFromCSV(ior io.Reader, sep rune, fileName string) ([]record.Record, error) {
	r := csv.NewReader(ior)
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
			ID:      getField(row, id),
			Name:    getField(row, name),
			Family:  getField(row, family),
		}
		res = append(res, rec)
	}
	return parse(res), nil
}

func getField(row []string, idx int) string {
	if idx == -1 {
		return ""
	}
	return row[idx]
}

func readHeader(s []string) (int, int, int, bool) {
	name, id, family := -1, -1, -1
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
			res = append(res, addParsed(recs[i]))
		}
	}
	return res
}

func addParsed(rec record.Record) record.Record {
	p := rec.Parsed
	rec.ParsingQuality = p.ParseQuality
	rec.Cardinality = p.Cardinality
	rec.CanonicalSimple = p.Canonical.Simple
	rec.CanonicalFull = p.Canonical.Full
	if p.Authorship != nil {
		rec.Authors = p.Authorship.Authors
		yrStr := strings.Trim(p.Authorship.Year, "()")
		yr, err := strconv.Atoi(yrStr)
		if err == nil {
			rec.Year = yr
		}
	}
	return rec
}

func fileSep(s, ext string) rune {
	if ext == ".csv" {
		return ','
	} else if strings.Contains(s, "\t") || ext == ".tsv" {
		return '\t'
	} else if !strings.Contains(s, ",") || strings.Contains(s, " ") {
		return rune(0)
	} else {
		return ','
	}
}
