package main

type TableData struct {
	Headers []string
	Rows    [][]string
}

func parseCSV(rows [][]string) TableData {
	return TableData{
		Headers: rows[0],
		Rows:    rows[1:],
	}
}
