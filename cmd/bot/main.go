package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"duty-bot/internal/config"
	"duty-bot/internal/telegram"
)

func main() {
	log.Println("bot run")

	_ = godotenv.Load()

	config.Ensure()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	log.Printf("admin id: %d", config.AdminID())

	telegram.Start(token)
}
