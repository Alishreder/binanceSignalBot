package telegram

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleError(message string) {
	b.notifyAdmin(message)
	log.Printf(message)
}

func (b *Bot) notifyAdmin(message string) {
	adminChatID, err := strconv.Atoi(os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Printf("can't convert string to int while trying to send message to admin: %s", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(int64(adminChatID), message)
	_, err = b.bot.Send(msg)
	if err != nil {
		log.Printf("error while trying to send message to admin: %s", err.Error())
		return
	}
}
