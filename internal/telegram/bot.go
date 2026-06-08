package telegram

import (
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"duty-bot/internal/config"
	"duty-bot/internal/store"
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

func Start(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	log.Println("bot enable", bot.Self.UserName)

	setupCommands(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			handleCallback(bot, update.CallbackQuery)
		case update.Message != nil:
			handleMessage(bot, update.Message)
		}
	}
}

func setupCommands(bot *tgbotapi.BotAPI) {
	cmds := []tgbotapi.BotCommand{
		{Command: "menu", Description: "Открыть меню"},
		{Command: "duty", Description: "График дежурств"},
		{Command: "time_schedule", Description: "График учёта времени"},
		{Command: "monitor", Description: "Мониторинг"},
		{Command: "settings", Description: "Настройки таблиц (админ)"},
	}
	if _, err := bot.Request(tgbotapi.NewSetMyCommands(cmds...)); err != nil {
		log.Println("не удалось установить команды:", err)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	userID := msg.From.ID
	chatID := msg.Chat.ID
	text := msg.Text

	log.Printf(">>> message: '%s' от userID: %d", text, userID)

	rememberUser(msg.From, msg.Chat)

	if strings.HasPrefix(text, "/") || isMenuButton(text) {
		clearPending(userID)
	} else if consumePending(bot, chatID, userID, text) {
		return
	}

	switch {
	case text == "/start" || text == "/menu":
		sendMenu(bot, chatID, userID)

	case text == "/duty" || text == btnDuty:
		go handleWithCooldown(bot, chatID, userID, func() { sendSource(bot, chatID, config.SourceDuty) })

	case text == "/time_schedule" || text == btnTimesheet:
		go handleWithCooldown(bot, chatID, userID, func() { sendSource(bot, chatID, config.SourceTimesheet) })

	case text == "/monitor" || text == btnMonitor:
		go handleWithCooldown(bot, chatID, userID, func() { sendSource(bot, chatID, config.SourceMonitor) })

	case text == "/settings" || text == btnSettings:
		sendSettings(bot, chatID, userID)
	}
}

func handleWithCooldown(bot *tgbotapi.BotAPI, chatID, userID int64, handler func()) {
	if !canUse(userID) {
		bot.Send(tgbotapi.NewMessage(chatID, "⏳ Подожди, запрос уже выполняется..."))
		return
	}
	defer doneProcessing(userID)
	handler()
}

func rememberUser(from *tgbotapi.User, chat *tgbotapi.Chat) {
	if from == nil || chat == nil || !chat.IsPrivate() {
		return
	}
	name := strings.TrimSpace(from.FirstName + " " + from.LastName)
	store.Remember(store.User{
		ID:       from.ID,
		ChatID:   chat.ID,
		Username: from.UserName,
		Name:     name,
	})
}
