package models

type User struct {
	ID     int `gorm:"primaryKey"`
	ChatID int64
}
