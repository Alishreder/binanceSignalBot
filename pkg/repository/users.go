package repository

import "github.com/Alishreder/binanceSignalBot/pkg/models"

type UsersInterface interface {
	Add(user models.User) error
	GetAll() ([]models.User, error)
	IsUserRegistered(chatID int64) bool
	Delete(chatID int64) error
}
