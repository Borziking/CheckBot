package datadir

import (
	"os"
	"path/filepath"
)

func Dir() string {
	if d := os.Getenv("DATA_DIR"); d != "" {
		return d
	}
	return "."
}

func Path(name string) string {
	dir := Dir()
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, name)
}
