package main

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sourceTitles = map[string]string{
	SourceDuty:      "📅 Дежурства",
	SourceTimesheet: "⏱ График времени",
	SourceMonitor:   "📊 Мониторинг",
}

var (
	editMu      sync.Mutex
	pendingEdit = make(map[int64]string)
)

func startEdit(bot *tgbotapi.BotAPI, chatID, userID int64, key string) {
	if !isAdminID(userID) {
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

	if err := setSource(key, text); err != nil {

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
