package telegram

import (
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"duty-bot/internal/config"
	"duty-bot/internal/render"
	"duty-bot/internal/sheets"
)

func outPath(name string) string {
	return filepath.Join(os.TempDir(), name)
}

func sendMenu(bot *tgbotapi.BotAPI, chatID, userID int64) {
	rows := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("📅 Дежурства", "show:"+config.SourceDuty)},
		{tgbotapi.NewInlineKeyboardButtonData("⏱ График времени", "show:"+config.SourceTimesheet)},
		{tgbotapi.NewInlineKeyboardButtonData("📊 Мониторинг", "show:"+config.SourceMonitor)},
	}
	if config.IsAdmin(userID) {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		})
	}

	msg := tgbotapi.NewMessage(chatID, "📋 *Главное меню*\nВыбери действие:")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
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

	case data == "settings":
		sendSettings(bot, chatID, userID)

	case data == "broadcast":
		startBroadcast(bot, chatID, userID)

	case strings.HasPrefix(data, "bcast:"):
		handleBroadcastCallback(bot, chatID, userID, strings.TrimPrefix(data, "bcast:"))

	case strings.HasPrefix(data, "show:"):
		key := strings.TrimPrefix(data, "show:")
		go handleWithCooldown(bot, chatID, userID, func() { sendSource(bot, chatID, key) })

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
	path := outPath("duty.png")
	if err := render.Duty(rows, path); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка отрисовки"))
		return
	}
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(path))
	photo.Caption = "📅 График дежурств"
	bot.Send(photo)
}

func sendTimesheet(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, config.SourceTimesheet)
	if !ok {
		return
	}
	path := outPath("table.png")
	if err := render.Timesheet(sheets.Parse(rows), path); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка отрисовки"))
		return
	}
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(path))
	photo.Caption = "⏱ График учёта времени"
	bot.Send(photo)
}

func sendMonitor(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, config.SourceMonitor)
	if !ok {
		return
	}
	path := outPath("monitor.png")
	if err := render.Monitor(rows, path); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка отрисовки"))
		return
	}
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(path))
	photo.Caption = "📊 Мониторинг"
	bot.Send(photo)
}
