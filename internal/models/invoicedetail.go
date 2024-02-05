package models

type InvoiceDetail struct {
	ID        int     `json:"id"`
	InvoiceID int     `json:"invoice_id" gorm:"size:11;not_null"`
	ItemID    int     `json:"item_id" gorm:"size:11;not_null"`
	ItemName  string  `json:"item_name" gorm:"size:191;not_null"`
	ItemPrice float64 `json:"item_price" gorm:"not_null"`
	Qty       int     `json:"qty" gorm:"not_null"`
	Amount    float64 `json:"amount" gorm:"not_null"`
	Common
}
