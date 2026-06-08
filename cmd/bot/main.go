package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/Borziking/CheckBot/internal/config"
	"github.com/Borziking/CheckBot/internal/telegram"
)

func main() {
	log.Println("bot run")

	loadEnv()

	config.Ensure()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	log.Printf("admin id: %d", config.AdminID())

	telegram.Start(token)
}

// loadEnv loads the nearest .env found by walking up from the working directory.
func loadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			godotenv.Load(path)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}
