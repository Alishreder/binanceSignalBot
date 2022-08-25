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
		return fmt.Errorf("error while adding user %v to db, error: %w", user, result.Error)
	}

	return nil
}
func (u *UsersRepository) GetAll() ([]models.User, error) {
	var users []models.User

	if result := u.db.Find(&users); result.Error != nil {
		return nil, fmt.Errorf("error while getting all users from db, error: %w", result.Error)
	}

	return users, nil
}

func (u *UsersRepository) IsUserRegistered(chatID int64) bool {
	var user models.User
	var count int64

	u.db.Find(&user, "chat_id = ?", chatID).Count(&count)

	if count == 0 {
		return false
	}

	return true
}

func (u *UsersRepository) Delete(chatID int64) error {

	if err := u.db.Delete("chat_id = ?", chatID).Error; err != nil {
		return fmt.Errorf("error while trying to delete user %d from db: %w", chatID, err)
	}

	return nil
}
