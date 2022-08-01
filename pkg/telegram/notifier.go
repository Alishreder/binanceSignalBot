package telegram

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) notifyUsers() {
	for {

		select {
		case message := <-b.sender.PriceChanges:

			b.users.m.Lock()
			for chatID := range b.users.users {
				msg := tgbotapi.NewMessage(chatID, message)
				b.bot.Send(msg)
			}
			b.users.m.Unlock()
		default:
			time.Sleep(time.Minute)
		}

	}
}
