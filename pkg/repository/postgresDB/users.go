package postgresDB

import (
	"fmt"

	"github.com/Alishreder/binanceSignalBot/pkg/models"
	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u *UsersRepository) Add(user models.User) error {
	if result := u.db.Create(&user); result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}
func (u *UsersRepository) GetAll() ([]models.User, error) {
	var users []models.User

	if result := u.db.Find(&users); result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return users, nil
}
