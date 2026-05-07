package main

import (
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mu         sync.Mutex
	processing = make(map[int64]bool)
)

func canUse(userID int64) bool {
	mu.Lock()
	defer mu.Unlock()

	if processing[userID] {
		return false
	}

	processing[userID] = true
	return true
}

func doneProcessing(userID int64) {
	mu.Lock()
	defer mu.Unlock()
	processing[userID] = false
}

func startBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	log.Println("bot enable", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		log.Printf(">>> message: '%s' от userID: %d", text, userID)

		switch {
		case text == "/duty":
			go handleWithCooldown(bot, chatID, userID, func() {
				handleDuty(bot, chatID)
			})

		case text == "/time_schedule":
			go handleWithCooldown(bot, chatID, userID, func() {
				handleTimesheet(bot, chatID)
			})

		case text == "/settings":
			handleSettings(bot, chatID)

		case text == "/monitor":
			go handleWithCooldown(bot, chatID, userID, func() {
				handleMonitor(bot, chatID)
			})

		case strings.HasPrefix(text, "/setsheet "):
			handleSetSheet(bot, chatID, text)
		}
	}
}

func handleWithCooldown(bot *tgbotapi.BotAPI, chatID int64, userID int64, handler func()) {
	if !canUse(userID) {
		bot.Send(tgbotapi.NewMessage(chatID, "⏳ Подожди, запрос уже выполняется..."))
		return
	}
	defer doneProcessing(userID)
	handler()
}

func handleDuty(bot *tgbotapi.BotAPI, chatID int64) {
	url, err := getSheetURL("duty")
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Лист не настроен. Используй /settings"))
		return
	}

	rows, err := fetchCSVFromURL(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки данных"))
		return
	}

	drawDutyTable(rows)

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("duty.png"))
	photo.Caption = "📅 График дежурств"
	bot.Send(photo)
}

func handleTimesheet(bot *tgbotapi.BotAPI, chatID int64) {
	url, err := getSheetURL("timesheet")
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Лист не настроен. Используй /settings"))
		return
	}

	rows, err := fetchCSVFromURL(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки данных"))
		return
	}

	table := parseCSV(rows)
	drawTable(table)

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("table.png"))
	photo.Caption = "⏱ График учёта времени"
	bot.Send(photo)
}

func handleSettings(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID,
		"⚙️ *Настройки:*\n\n"+
			"Изменить лист дежурства:\n`/setsheet duty ГИД`\n\n"+
			"Изменить лист учёта времени:\n`/setsheet time_schedule ГИД`\n\n"+
			"ГИД — число в конце ссылки после `gid=`")
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func handleSetSheet(bot *tgbotapi.BotAPI, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) != 3 {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Формат: /setsheet duty 123456"))
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки конфига"))
		return
	}

	cfg.Sheets[parts[1]] = parts[2]
	err = saveConfig(cfg)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения конфига"))
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, "✅ Лист '"+parts[1]+"' обновлён!"))
}

func handleMonitor(bot *tgbotapi.BotAPI, chatID int64) {
	url, err := getURL("monitor")
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Лист не настроен"))
		return
	}

	rows, err := fetchCSVFromURL(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки данных"))
		return
	}

	drawMonitorTable(rows)

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("monitor.png"))
	photo.Caption = "📊 Мониторинг"
	bot.Send(photo)
}
