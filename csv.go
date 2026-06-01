package main

import (
	"encoding/csv"
	"net/http"
)

func fetchCSVFromURL(url string) ([][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}
