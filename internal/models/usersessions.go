package models

type UserSession struct {
	ID     int    `json:"id"`
	UserID int    `json:"name" gorm:"not_null"`
	Token  string `json:"token" gorm:"size:255;not_null"`
	Common
}
