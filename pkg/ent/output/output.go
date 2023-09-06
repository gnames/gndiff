package output

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gnames/gndiff/pkg/ent/record"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/verifier"
)

type Output struct {
	Matches []Match
}

type Match struct {
	SourceRecord     record.Record   `json:"sourceRecord"`
	ReferenceRecords []record.Record `json:"referenceRecords"`
}

// NameOutput takes result of verification for one string and converts it into
// required format (CSV or JSON).
func MatchOutput(o Output, f gnfmt.Format) string {
	sort.Slice(o.Matches, func(i, j int) bool {
		return o.Matches[i].SourceRecord.Index < o.Matches[j].SourceRecord.Index
	})

	switch f {
	case gnfmt.CSV:
		return csvOutput(o, ',')
	case gnfmt.TSV:
		return csvOutput(o, '\t')
	case gnfmt.CompactJSON:
		return jsonOutput(o, false)
	case gnfmt.PrettyJSON:
		return jsonOutput(o, true)
	}
	return "N/A"
}

// CSVHeader returns the header string for CSV output format.
func CSVHeader(f gnfmt.Format) string {
	header := []string{
		"SourceFile", "SrcRowNum", "SrcTaxonId", "SrcName", "SrcCanonical",
		"ReferenceFile", "RefRowNum", "RefTaxonId", "RefName", "RefCanonical",
		"Score", "MatchType", "EditDistance",
	}
	switch f {
	case gnfmt.CSV:
		return gnfmt.ToCSV(header, ',')
	case gnfmt.TSV:
		return gnfmt.ToCSV(header, '\t')
	default:
		return ""
	}
}

func csvOutput(o Output, sep rune) string {
	var res []string
	for i := range o.Matches {
		rows := csvRow(o.Matches[i], sep)
		for i := range rows {
			res = append(res, rows[i])
		}
	}
	return strings.Join(res, "\n")
}

func csvRow(m Match, sep rune) []string {
	var res []string
	s := m.SourceRecord
	r := m.ReferenceRecords
	for i := range r {
		row := []string{
			s.DataSet,
			strconv.Itoa(s.Index),
			s.ID,
			s.Name,
			s.CanonicalFull,
			r[i].DataSet,
			strconv.Itoa(r[i].Index),
			r[i].ID,
			r[i].Name,
			r[i].CanonicalFull,
			strconv.FormatFloat(r[i].Score, 'f', 5, 64),
			r[i].MatchType.String(),
			strconv.Itoa(r[i].EditDistance),
		}
		res = append(res, gnfmt.ToCSV(row, sep))
	}
	if len(r) == 0 {
		row := []string{
			s.DataSet,
			strconv.Itoa(s.Index),
			s.ID,
			s.Name,
			s.CanonicalFull,
			"",
			"",
			"",
			"",
			"",
			"",
			verifier.NoMatch.String(),
			"",
		}
		res = append(res, gnfmt.ToCSV(row, sep))
	}
	return res
}

func jsonOutput(o Output, pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
