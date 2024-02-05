package models

type Customer struct {
	ID              int    `json:"id"`
	CustomerName    string `json:"customer_name" gorm:"size:100;not_null"`
	CustomerAddress string `json:"customer_address" gorm:"not_null"`
	CustomerPhone   string `json:"customer_phone" gorm:"not_null"`
	Common
}
