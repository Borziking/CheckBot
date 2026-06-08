package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"duty-bot/internal/datadir"
)

func file() string { return datadir.Path("config.json") }

func Ensure() {
	dst := file()
	if _, err := os.Stat(dst); err == nil {
		return
	}
	if data, err := os.ReadFile("config.json"); err == nil {
		os.WriteFile(dst, data, 0644)
	}
}

const (
	SourceDuty      = "duty"
	SourceTimesheet = "timesheet"
	SourceMonitor   = "monitor"
)

type Config struct {
	Sources map[string]string `json:"sources"`
}

var mu sync.Mutex

func Load() (Config, error) {
	mu.Lock()
	defer mu.Unlock()

	var cfg Config
	data, err := os.ReadFile(file())
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

func Save(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(file(), data, 0644)
}

func AdminID() int64 {
	if id, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64); err == nil {
		return id
	}
	return 0
}

func IsAdmin(userID int64) bool {
	return userID != 0 && userID == AdminID()
}

func SourceURL(key string) (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	url, ok := cfg.Sources[key]
	if !ok || url == "" {
		return "", fmt.Errorf("источник '%s' не настроен", key)
	}
	return url, nil
}

func SetSource(key, rawURL string) error {
	url, err := normalizeSheetURL(rawURL)
	if err != nil {
		return err
	}
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.Sources[key] = url
	return Save(cfg)
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
