package telegram

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"duty-bot/internal/config"
)

var sourceTitles = map[string]string{
	config.SourceDuty:      "📅 Дежурства",
	config.SourceTimesheet: "⏱ График времени",
	config.SourceMonitor:   "📊 Мониторинг",
}

var (
	editMu      sync.Mutex
	pendingEdit = make(map[int64]string)
)

func sendSettings(bot *tgbotapi.BotAPI, chatID, userID int64) {
	if !config.IsAdmin(userID) {
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
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[config.SourceDuty], "edit:"+config.SourceDuty)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[config.SourceTimesheet], "edit:"+config.SourceTimesheet)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[config.SourceMonitor], "edit:"+config.SourceMonitor)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu")},
	)
}

func startEdit(bot *tgbotapi.BotAPI, chatID, userID int64, key string) {
	if !config.IsAdmin(userID) {
		bot.Send(tgbotapi.NewMessage(chatID, "⛔ Доступ только для администратора"))
		return
	}
	title, ok := sourceTitles[key]
	if !ok {
		return
	}

	editMu.Lock()
	pendingEdit[userID] = key
	editMu.Unlock()

	bot.Send(tgbotapi.NewMessage(chatID,
		"✏️ Меняем таблицу «"+title+"»\n\n"+
			"Открой нужный лист в Google Таблицах и пришли ссылку из адресной строки.\n\n"+
			"Отмена — /menu"))
}

func clearPendingEdit(userID int64) {
	editMu.Lock()
	delete(pendingEdit, userID)
	editMu.Unlock()
}

func consumePendingEdit(bot *tgbotapi.BotAPI, chatID, userID int64, text string) bool {
	editMu.Lock()
	key, waiting := pendingEdit[userID]
	if waiting {
		delete(pendingEdit, userID)
	}
	editMu.Unlock()

	if !waiting {
		return false
	}

	if err := config.SetSource(key, text); err != nil {

		editMu.Lock()
		pendingEdit[userID] = key
		editMu.Unlock()

		bot.Send(tgbotapi.NewMessage(chatID, "❌ "+err.Error()+"\nПопробуй ещё раз или /menu для отмены."))
		return true
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Таблица «"+sourceTitles[key]+"» обновлена!")
	msg.ReplyMarkup = settingsKeyboard()
	bot.Send(msg)
	return true
}
