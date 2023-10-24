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
	RecordsFromText(
		r io.Reader,
		fileName, ext string,
	) (rec []record.Record, sep rune, err error)

	// RecordsFromCSV takes io.Reader, CSV separator, fileName
	//  and converts data into records
	RecordsFromCSV(
		r io.Reader,
		sep rune,
		fileName string,
	) ([]record.Record, error)
}
