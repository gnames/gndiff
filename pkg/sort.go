package gndiff

import (
	"strconv"

	"github.com/gnames/gnames/pkg/ent/score"
	"github.com/gnames/gnames/pkg/ent/verifier"
	"github.com/gnames/gndiff/pkg/ent/output"
	"github.com/gnames/gndiff/pkg/ent/record"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

func sortByScore(res []output.Match) {
	for i := range res {
		matchRes := getMatchRes(res[i])
		s := score.New()
		s.SortResults(matchRes)
		res[i].ReferenceRecords = toRecords(matchRes.MatchResults)
	}
}

func toRecords(rd []*vlib.ResultData) []record.Record {
	res := make([]record.Record, len(rd))
	for i := range rd {
		idx, _ := strconv.Atoi(rd[i].LocalID)
		rec := record.Record{
			ID:              rd[i].RecordID,
			DataSet:         rd[i].DataSourceTitleShort,
			Index:           idx,
			EditDistance:    rd[i].EditDistance,
			Name:            rd[i].MatchedName,
			ParsingQuality:  rd[i].ParsingQuality,
			Cardinality:     rd[i].MatchedCardinality,
			CanonicalFull:   rd[i].MatchedCanonicalFull,
			CanonicalSimple: rd[i].MatchedCanonicalSimple,
			Authors:         rd[i].MatchedAuthors,
			Year:            rd[i].MatchedYear,
			Family:          rd[i].ClassificationPath,
			MatchType:       rd[i].MatchType,
			Score:           rd[i].SortScore,
			ScoreDetails:    &rd[i].ScoreDetails,
		}
		res[i] = rec
	}
	return res
}

func getMatchRes(m output.Match) *verifier.MatchRecord {
	res := &verifier.MatchRecord{
		ID:              m.QueryRecord.ID,
		Name:            m.QueryRecord.Name,
		Cardinality:     m.QueryRecord.Cardinality,
		CanonicalFull:   m.QueryRecord.CanonicalFull,
		CanonicalSimple: m.QueryRecord.Canonical.Simple,
		Authors:         m.QueryRecord.Authors,
		Year:            m.QueryRecord.Year,
	}

	var rds []*vlib.ResultData
	for _, v := range m.ReferenceRecords {
		rd := &vlib.ResultData{
			DataSourceTitleShort:   v.DataSet,
			RecordID:               v.ID,
			LocalID:                strconv.Itoa(v.Index),
			ParsingQuality:         v.ParseQuality,
			MatchedName:            v.Name,
			MatchedCardinality:     v.Cardinality,
			MatchedCanonicalFull:   v.CanonicalFull,
			MatchedCanonicalSimple: v.CanonicalSimple,
			MatchedAuthors:         v.Authors,
			MatchedYear:            v.Year,
			EditDistance:           v.EditDistance,
			MatchType:              v.MatchType,
			ClassificationPath:     v.Family,
		}
		rds = append(rds, rd)
	}
	res.MatchResults = rds
	return res
}
