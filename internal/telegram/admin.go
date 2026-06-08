package telegram

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Borziking/CheckBot/internal/config"
	"github.com/Borziking/CheckBot/internal/store"
)

var sourceTitles = map[string]string{
	config.SourceDuty:      "📅 Дежурства",
	config.SourceTimesheet: "⏱ График времени",
	config.SourceMonitor:   "📊 Мониторинг",
}

type pendingKind int

const (
	pendingEditSource pendingKind = iota
	pendingBroadcast
)

type pendingState struct {
	kind pendingKind
	key  string
}

var (
	pendingMu sync.Mutex
	pending   = map[int64]pendingState{}
)

func setPending(userID int64, st pendingState) {
	pendingMu.Lock()
	pending[userID] = st
	pendingMu.Unlock()
}

func clearPending(userID int64) {
	pendingMu.Lock()
	delete(pending, userID)
	pendingMu.Unlock()
}

func takePending(userID int64) (pendingState, bool) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	st, ok := pending[userID]
	if ok {
		delete(pending, userID)
	}
	return st, ok
}

var (
	draftMu sync.Mutex
	drafts  = map[int64]string{}
)

func sendSettings(bot *tgbotapi.BotAPI, chatID, userID int64) {
	if !config.IsAdmin(userID) {
		bot.Send(tgbotapi.NewMessage(chatID, "⛔ Доступ только для администратора"))
		return
	}

	msg := tgbotapi.NewMessage(chatID, "⚙️ *Настройки*\nВыбери действие:")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = settingsKeyboard()
	bot.Send(msg)
}

func settingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[config.SourceDuty], "edit:"+config.SourceDuty)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[config.SourceTimesheet], "edit:"+config.SourceTimesheet)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(sourceTitles[config.SourceMonitor], "edit:"+config.SourceMonitor)},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("📣 Рассылка всем", "broadcast")},
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

	setPending(userID, pendingState{kind: pendingEditSource, key: key})
	bot.Send(tgbotapi.NewMessage(chatID,
		"✏️ Меняем таблицу «"+title+"»\n\n"+
			"Открой нужный лист в Google Таблицах и пришли ссылку из адресной строки.\n\n"+
			"Отмена — /menu"))
}

func startBroadcast(bot *tgbotapi.BotAPI, chatID, userID int64) {
	if !config.IsAdmin(userID) {
		bot.Send(tgbotapi.NewMessage(chatID, "⛔ Доступ только для администратора"))
		return
	}

	setPending(userID, pendingState{kind: pendingBroadcast})
	bot.Send(tgbotapi.NewMessage(chatID,
		fmt.Sprintf("📣 Пришли текст рассылки — он уйдёт всем пользователям бота (%d чел.).\n\nОтмена — /menu", store.Count())))
}

func askBroadcastConfirm(bot *tgbotapi.BotAPI, chatID, userID int64, text string) {
	draftMu.Lock()
	drafts[userID] = text
	draftMu.Unlock()

	preview := tgbotapi.NewMessage(chatID,
		"📣 Предпросмотр рассылки:\n\n"+text+
			fmt.Sprintf("\n\n— Отправить %d получателям?", store.Count()))
	preview.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("✅ Отправить", "bcast:send"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "bcast:cancel"),
		},
	)
	bot.Send(preview)
}

func handleBroadcastCallback(bot *tgbotapi.BotAPI, chatID, userID int64, action string) {
	if !config.IsAdmin(userID) {
		return
	}

	switch action {
	case "send":
		draftMu.Lock()
		text, ok := drafts[userID]
		delete(drafts, userID)
		draftMu.Unlock()
		if !ok {
			bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Черновик рассылки не найден. Начни заново через ⚙️ Настройки."))
			return
		}
		runBroadcast(bot, chatID, text)

	case "cancel":
		draftMu.Lock()
		delete(drafts, userID)
		draftMu.Unlock()
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Рассылка отменена"))
	}
}

func runBroadcast(bot *tgbotapi.BotAPI, adminChatID int64, text string) {
	users := store.All()
	bot.Send(tgbotapi.NewMessage(adminChatID,
		fmt.Sprintf("📣 Рассылка запущена для %d пользователей…", len(users))))

	go func() {
		var ok, fail int
		for _, u := range users {
			if _, err := bot.Send(tgbotapi.NewMessage(u.ChatID, text)); err != nil {
				fail++
				log.Printf("broadcast: failed to deliver to %d: %v", u.ChatID, err)
				if blockedByUser(err) {
					store.Remove(u.ID)
				}
			} else {
				ok++
			}
			time.Sleep(40 * time.Millisecond)
		}
		bot.Send(tgbotapi.NewMessage(adminChatID,
			fmt.Sprintf("✅ Рассылка завершена.\nДоставлено: %d\nНе доставлено: %d", ok, fail)))
	}()
}

func blockedByUser(err error) bool {
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "blocked") ||
		strings.Contains(s, "deactivated") ||
		strings.Contains(s, "chat not found") ||
		strings.Contains(s, "user is deactivated")
}

func consumePending(bot *tgbotapi.BotAPI, chatID, userID int64, text string) bool {
	st, ok := takePending(userID)
	if !ok {
		return false
	}

	switch st.kind {
	case pendingEditSource:
		if err := config.SetSource(st.key, text); err != nil {
			setPending(userID, st)
			bot.Send(tgbotapi.NewMessage(chatID, "❌ "+err.Error()+"\nПопробуй ещё раз или /menu для отмены."))
			return true
		}
		msg := tgbotapi.NewMessage(chatID, "✅ Таблица «"+sourceTitles[st.key]+"» обновлена!")
		msg.ReplyMarkup = settingsKeyboard()
		bot.Send(msg)

	case pendingBroadcast:
		askBroadcastConfirm(bot, chatID, userID, text)
	}

	return true
}
