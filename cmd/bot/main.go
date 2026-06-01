package main

import (
	"log"
	"os"

	"duty-bot/internal/telegram"
)

const defaultToken = "token_bot"

func main() {
	log.Println("bot run")

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		token = defaultToken
	}

	telegram.Start(token)
}
