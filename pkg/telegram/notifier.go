package telegram

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) notifyUsers() {
	for {

		select {
		case message := <-b.sender.PriceChanges:

			users, err := b.usersRepository.GetAll()
			if err != nil {
				b.notifyAdmin("error while getting users from bd")
				return
			}

			for _, v := range users {
				msg := tgbotapi.NewMessage(v.ChatID, message)
				_, err := b.bot.Send(msg)
				if err != nil {
					if err := b.usersRepository.Delete(v.ChatID); err != nil {
						b.notifyAdmin(err.Error())
						return
					}
				}
			}

		default:
			time.Sleep(time.Second * 5)
		}

	}
}
