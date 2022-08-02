package repository

import "github.com/Alishreder/binanceSignalBot/pkg/models"

type UsersInterface interface {
	Add(user models.User) error
	GetAll() ([]models.User, error)
}
