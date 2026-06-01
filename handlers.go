package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMenu(bot *tgbotapi.BotAPI, chatID, userID int64) {
	rows := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("📅 Дежурства", "show:"+SourceDuty)},
		{tgbotapi.NewInlineKeyboardButtonData("⏱ График времени", "show:"+SourceTimesheet)},
		{tgbotapi.NewInlineKeyboardButtonData("📊 Мониторинг", "show:"+SourceMonitor)},
	}
	if isAdminID(userID) {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		})
	}

	msg := tgbotapi.NewMessage(chatID, "📋 *Главное меню*\nВыбери действие:")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	bot.Send(msg)
}

func sendSettings(bot *tgbotapi.BotAPI, chatID, userID int64) {
	if !isAdminID(userID) {
		bot.Send(tgbotapi.NewMessage(chatID, "⛔ Доступ только для администратора"))
		return
	}

	msg := tgbotapi.NewMessage(chatID, "⚙️ *Настройки таблиц*\nВыбери, какую таблицу изменить:")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = settingsKeyboard()
	bot.Send(msg)
}

func settingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[SourceDuty], "edit:"+SourceDuty)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[SourceTimesheet], "edit:"+SourceTimesheet)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[SourceMonitor], "edit:"+SourceMonitor)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu")},
	)
}

func handleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	userID := cq.From.ID
	data := cq.Data

	bot.Request(tgbotapi.NewCallback(cq.ID, ""))

	switch {
	case data == "menu":
		sendMenu(bot, chatID, userID)

	case data == "settings":
		sendSettings(bot, chatID, userID)

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
	case SourceDuty:
		sendDuty(bot, chatID)
	case SourceTimesheet:
		sendTimesheet(bot, chatID)
	case SourceMonitor:
		sendMonitor(bot, chatID)
	}
}

func loadSource(bot *tgbotapi.BotAPI, chatID int64, key string) ([][]string, bool) {
	url, err := sourceURL(key)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Таблица не настроена. Открой ⚙️ Настройки"))
		return nil, false
	}
	rows, err := fetchCSVFromURL(url)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка загрузки данных"))
		return nil, false
	}
	return rows, true
}

func sendDuty(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, SourceDuty)
	if !ok {
		return
	}
	drawDutyTable(rows)

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("duty.png"))
	photo.Caption = "📅 График дежурств"
	bot.Send(photo)
}

func sendTimesheet(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, SourceTimesheet)
	if !ok {
		return
	}
	drawTable(parseCSV(rows))

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("table.png"))
	photo.Caption = "⏱ График учёта времени"
	bot.Send(photo)
}

func sendMonitor(bot *tgbotapi.BotAPI, chatID int64) {
	rows, ok := loadSource(bot, chatID, SourceMonitor)
	if !ok {
		return
	}
	drawMonitorTable(rows)

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("monitor.png"))
	photo.Caption = "📊 Мониторинг"
	bot.Send(photo)
}
