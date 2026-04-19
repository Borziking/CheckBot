package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

const csvURL = "https://docs.google.com/spreadsheets/d/SHEET_ID/export?format=csv&gid=GID"

func fetchCSVFromURL(url string) ([][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	return reader.ReadAll()
}

func main() {
	fmt.Println("Запускаем бота...")
	startBot("token_bot")
}
