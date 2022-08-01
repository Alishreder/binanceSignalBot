package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const commandStart = "start"

func (b *Bot) handleCommand(message *tgbotapi.Message) (err error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know this command")

	switch message.Command() {
	case commandStart:

		b.addNewUser(message.Chat.ID)

		msg.Text = "You started bot"
		_, err = b.bot.Send(msg)
	default:
		_, err = b.bot.Send(msg)
	}

	return
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}
