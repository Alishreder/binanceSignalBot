package telegram

import (
	"log"
	"time"

	"github.com/Alishreder/binanceSignalBot/pkg/crypto"
	"github.com/Alishreder/binanceSignalBot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	sender          crypto.PriceSender
	usersRepository repository.UsersInterface
}

func NewBot(bot *tgbotapi.BotAPI, sender crypto.PriceSender, userRepository repository.UsersInterface) *Bot {
	return &Bot{
		bot:             bot,
		sender:          sender,
		usersRepository: userRepository,
	}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChanel()

	go b.sender.TrackPriceChange("30m", 30*time.Minute)
	go b.sender.TrackPriceChange("1h", time.Hour)

	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {

	go b.notifyUsers()

	for update := range updates {
		if update.Message == nil { // ignore any not-Message updates
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.notifyAdmin(err.Error())
			}
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
