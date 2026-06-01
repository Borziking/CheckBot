package sheets

import (
	"encoding/csv"
	"net/http"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func FetchCSV(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func Parse(rows [][]string) Table {
	if len(rows) == 0 {
		return Table{}
	}
	return Table{Headers: rows[0], Rows: rows[1:]}
}
