package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

const configPath = "config.json"

const (
	SourceDuty      = "duty"
	SourceTimesheet = "timesheet"
	SourceMonitor   = "monitor"
)

type Config struct {
	AdminID int64             `json:"admin_id"`
	Sources map[string]string `json:"sources"`
}

var cfgMu sync.Mutex

func loadConfig() (Config, error) {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Sources == nil {
		cfg.Sources = map[string]string{}
	}
	return cfg, nil
}

func saveConfig(cfg Config) error {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func isAdminID(userID int64) bool {
	cfg, err := loadConfig()
	if err != nil {
		return false
	}
	return cfg.AdminID == userID
}

func sourceURL(key string) (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	url, ok := cfg.Sources[key]
	if !ok || url == "" {
		return "", fmt.Errorf("источник '%s' не настроен", key)
	}
	return url, nil
}

func setSource(key, rawURL string) error {
	url, err := normalizeSheetURL(rawURL)
	if err != nil {
		return err
	}
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.Sources[key] = url
	return saveConfig(cfg)
}

var (
	reSpreadsheetID = regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	reGID           = regexp.MustCompile(`[?&#]gid=(\d+)`)
)

func normalizeSheetURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	idMatch := reSpreadsheetID.FindStringSubmatch(raw)
	if idMatch == nil {
		return "", fmt.Errorf("не похоже на ссылку Google Таблицы")
	}
	id := idMatch[1]

	gid := "0"
	if g := reGID.FindStringSubmatch(raw); g != nil {
		gid = g[1]
	}

	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", id, gid), nil
}
