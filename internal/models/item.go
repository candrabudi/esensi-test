package models

type Item struct {
	ID              int    `json:"id"`
	ItemName        string `json:"item__name" gorm:"size:100;not_null"`
	ItemDescription string `json:"item_description" gorm:"not_null"`
	ItemType        string `json:"item_type" gorm:"type:enum('Hardware', 'Software','Service');not_null"`
	ItemPrice       int    `json:"item_price" gorm:"not_null"`
	ItemStock       int    `json:"item_stock" gorm:"not_null"`
	Common
}
