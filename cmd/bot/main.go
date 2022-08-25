package main

import (
	"log"

	"github.com/Alishreder/binanceSignalBot/pkg/crypto"
	"github.com/Alishreder/binanceSignalBot/pkg/models"
	"github.com/Alishreder/binanceSignalBot/pkg/repository/postgresDB"
	"github.com/Alishreder/binanceSignalBot/pkg/telegram"
	"github.com/Alishreder/binanceSignalBot/pkg/utilits"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	telegramToken := utilits.GetEnvVariable("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	db, err := initDB()
	if err != nil {
		log.Fatalln(err)
	}

	usersRepository := postgresDB.NewUsersRepository(db)
	priceSender := crypto.NewPriceSender()

	telegramBot := telegram.NewBot(bot, priceSender, usersRepository)
	telegramBot.Start()

}

func initDB() (*gorm.DB, error) {
	dbURL := utilits.GetEnvVariable("CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err

	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
