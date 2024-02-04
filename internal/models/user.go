package models

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name" gorm:"size:100;not_null"`
	Email       string `json:"email" gorm:"size:100;not_null;unique"`
	Password    string `json:"password" gorm:"size:255;not_null"`
	PhoneNumber string `json:"phone_number" gorm:"size:20;not_null;unique"`
	Common
}
