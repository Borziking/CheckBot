package telegram

import (
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Borziking/CheckBot/internal/config"
	"github.com/Borziking/CheckBot/internal/render"
	"github.com/Borziking/CheckBot/internal/sheets"
)

const (
	btnDuty      = "📅 Дежурства"
	btnTimesheet = "⏱ График времени"
	btnMonitor   = "📊 Мониторинг"
	btnSettings  = "⚙️ Настройки"
)

func isMenuButton(text string) bool {
	switch text {
	case btnDuty, btnTimesheet, btnMonitor, btnSettings:
		return true
	}
	return false
}

func sendRendered(bot *tgbotapi.BotAPI, chatID int64, prefix, caption string, draw func(path string) error) {
	f, err := os.CreateTemp("", prefix+"-*.png")
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка отрисовки"))
		return
	}
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	if err := draw(path); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка отрисовки"))
		return
	}

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(path))
	photo.Caption = caption
	bot.Send(photo)
}

func mainKeyboard(userID int64) tgbotapi.ReplyKeyboardMarkup {
	rows := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton(btnDuty), tgbotapi.NewKeyboardButton(btnTimesheet)},
		{tgbotapi.NewKeyboardButton(btnMonitor)},
	}
	if config.IsAdmin(userID) {
		rows = append(rows, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(btnSettings)})
	}
	kb := tgbotapi.NewReplyKeyboard(rows...)
	kb.ResizeKeyboard = true
	return kb
}

func sendMenu(bot *tgbotapi.BotAPI, chatID, userID int64) {
	msg := tgbotapi.NewMessage(chatID, "📋 Главное меню — выбери действие на клавиатуре ниже.")
	msg.ReplyMarkup = mainKeyboard(userID)
	bot.Send(msg)
}

func handleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	userID := cq.From.ID
	data := cq.Data

	rememberUser(cq.From, cq.Message.Chat)

	bot.Request(tgbotapi.NewCallback(cq.ID, ""))

	switch {
	case data == "menu":
		sendMenu(bot, chatID, userID)

	case data == "broadcast":
		startBroadcast(bot, chatID, userID)

	case strings.HasPrefix(data, "bcast:"):
		handleBroadcastCallback(bot, chatID, userID, strings.TrimPrefix(data, "bcast:"))

	case strings.HasPrefix(data, "edit:"):
		key := strings.TrimPrefix(data, "edit:")
		startEdit(bot, chatID, userID, key)
	}
}

func sendSource(bot *tgbotapi.BotAPI, chatID int64, key string) {
	switch key {
	case config.SourceDuty:
		sendDuty(bot, chatID)
	case config.SourceTimesheet:
		sendTimesheet(bot, chatID)
	case config.SourceMonitor:
		sendMonitor(bot, chatID)
	}
}

func loadSource(bot *tgbotapi.BotAPI, chatID int64, key string) ([][]string, bool) {
	url, err := config.SourceURL(key)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Таблица не настроена. Открой ⚙️ Настройки"))
		return nil, false
	}
	rows, err := sheets.FetchCSV(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки данных"))
		return nil, false
	}
	return rows, true
}

func sendDuty(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, config.SourceDuty)
	if !ok {
		return
	}
	sendRendered(bot, chatID, "duty", "📅 График дежурств", func(path string) error {
		return render.Duty(rows, path)
	})
}

func sendTimesheet(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, config.SourceTimesheet)
	if !ok {
		return
	}
	sendRendered(bot, chatID, "timesheet", "⏱ График учёта времени", func(path string) error {
		return render.Timesheet(sheets.Parse(rows), path)
	})
}

func sendMonitor(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, config.SourceMonitor)
	if !ok {
		return
	}
	sendRendered(bot, chatID, "monitor", "📊 Мониторинг", func(path string) error {
		return render.Monitor(rows, path)
	})
}
