package ingester

import (
	"io"

	"github.com/gnames/gndiff/pkg/ent/record"
)

// Ingester converts plain text, CSV, or TSV file into
// name records.
type Ingester interface {
	// Records takes a filepath, reads the file and converts its
	// data into Records, or returns error if read was not
	// successful.
	Records(path string) ([]record.Record, error)

	// RecordsFromText takes a io.Reader and converts its data
	// into records, or return back an error.
	// In this case the file assumed to be in a text format with one
	// name per line.
	RecordsFromText(
		r io.Reader,
		fileName, ext string,
	) (rec []record.Record, sep rune, err error)

	// RecordsFromCSV takes io.Reader, CSV separator, fileName
	// and converts data into records
	// In this case file suppose to be in CSV or TSV format and contain
	// ScientificName, TaxonID, Family fields (2 last ones are optional).
	RecordsFromCSV(
		r io.Reader,
		sep rune,
		fileName string,
	) ([]record.Record, error)
}
