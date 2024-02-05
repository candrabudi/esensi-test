package models

import "github.com/shopspring/decimal"

type InvoiceDetail struct {
	ID        int             `json:"id" gorm:"primary_key"`
	InvoiceID uint            `json:"invoice_id" gorm:"size:11;not_null"`
	ItemID    int             `json:"item_id" gorm:"size:11;not_null"`
	ItemName  string          `json:"item_name" gorm:"size:191;not_null"`
	ItemPrice decimal.Decimal `json:"item_price" gorm:"type:decimal(10,2);not_null"`
	Qty       int             `json:"qty" gorm:"not_null"`
	Amount    decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not_null"`
	Common
}
