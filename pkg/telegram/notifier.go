package telegram

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) notifyUsers() {
	for {

		select {
		case message := <-b.sender.PriceChanges:

			users, err := b.usersRepository.GetAll()
			if err != nil {
				fmt.Println(err)
			}

			for _, v := range users {
				msg := tgbotapi.NewMessage(v.ChatID, message)
				_, err := b.bot.Send(msg)
				if err != nil {
					// TODO handle this error, delete from db and cache
				}
			}

		default:
			time.Sleep(time.Second * 5)
		}

	}
}
