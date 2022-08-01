package telegram

import (
	"log"
	"sync"

	"github.com/Alishreder/binanceSignalBot/pkg/crypto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	sender crypto.PriceSender
	users
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot:    bot,
		sender: crypto.NewPriceSender(),
		users: users{
			users: make(map[int64]bool),
			m:     &sync.Mutex{},
		},
	}
}

type users struct {
	users map[int64]bool
	m     *sync.Mutex
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChanel()

	go b.sender.TrackPriceChange()

	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {

	go b.notifyUsers()

	for update := range updates {
		if update.Message == nil { // ignore any not-Message updates
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)

	}
}

func (b *Bot) initUpdatesChanel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) addNewUser(chatID int64) {
	b.users.m.Lock()
	b.users.users[chatID] = true
	b.users.m.Unlock()
}
