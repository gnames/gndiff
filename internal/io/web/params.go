package web

// Input is how JSON for POST method is organized.
type Input struct {
	// QueryFileName can be `csv`, `tsv`, `text`.
	// This field is optional, but it can help to avoid problems.
	QueryFileName string `json:"queryFileName"`

	// QueryText is the content of CSV, TSV, or plain text file.
	// It contains names that need to be matched to Reference.
	QueryText string `json:"queryText"`

	// ReferenceFileName can be `csv`, `tsv`, or `text`.
	// This field is optional, but it can help to avoid problems.
	ReferenceFileName string `json:"referenceFileName"`

	// ReferenceText is the content of CSV, TSV, or plain text file.
	// It contains names that are matched agains by the Query.
	ReferenceText string `json:"referenceText"`
}
