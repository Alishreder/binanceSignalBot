package telegram

import (
	"fmt"
	"log"

	"github.com/Alishreder/binanceSignalBot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"

	replyStart = "Hi there!\nFrom now I will send you alerts if price of any token from Binance(with market cap more then 100 millions) increased on 5% or more per 1 hour or per 30 minutes."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) (err error) {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	_, err := b.bot.Send(msg)
	if err != nil {
		b.notifyAdmin("error while handleMessage(): " + err.Error())
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, replyStart)

	if b.usersRepository.IsUserRegistered(message.Chat.ID) {
		msg.Text = "Sorry, you have been already registered :("
		_, err := b.bot.Send(msg)
		return fmt.Errorf("error while trying to send message to user in handleStartCommand(): %w", err)
	}

	user := models.User{
		ChatID: message.Chat.ID,
	}

	if err := b.usersRepository.Add(user); err != nil {
		return fmt.Errorf("trying to handleStartCommand: %w", err)
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know this command")

	_, err := b.bot.Send(msg)
	return fmt.Errorf("error while trying to send message to user %d in handleUnknownCommand(): %w", message.Chat.ID, err)
}
