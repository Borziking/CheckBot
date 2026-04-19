package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SpreadsheetID string            `json:"spreadsheet_id"`
	Sheets        map[string]string `json:"sheets"`
	URLs          map[string]string `json:"urls"`
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile("config.json")
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func saveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}

func getSheetURL(sheetKey string) (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	gid, exists := cfg.Sheets[sheetKey]
	if !exists {
		return "", fmt.Errorf("лист '%s' не настроен", sheetKey)
	}
	return "https://docs.google.com/spreadsheets/d/" + cfg.SpreadsheetID + "/export?format=csv&gid=" + gid, nil
}

func getURL(key string) (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	url, exists := cfg.URLs[key]
	if !exists {
		return "", fmt.Errorf("url '%s' не настроен", key)
	}
	return url, nil
}
