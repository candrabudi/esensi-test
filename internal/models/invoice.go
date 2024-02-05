package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Invoice struct {
	ID              int             `json:"id" gorm:"primary_key"`
	InvoiceNo       string          `json:"invoice_no" gorm:"size:15;not_null"`
	CustomerID      int             `json:"customer_id" gorm:"not_null"`
	CustomerName    string          `json:"customer_name" gorm:"not_null"`
	CustomerAddress string          `json:"customer_address" gorm:"not_null"`
	IssueDate       time.Time       `json:"issue_date" gorm:"not_null"`
	Subject         string          `json:"subject" gorm:"size:191;not_null"`
	TotalItem       int             `json:"total_item" gorm:"not_null"`
	GrandTotal      decimal.Decimal `json:"grandtotal" gorm:"type:decimal(10,2);not_null"`
	SubTotal        decimal.Decimal `json:"subtotal" gorm:"type:decimal(10,2);not_null"`
	DueDate         time.Time       `json:"due_date" gorm:"not_null"`
	Status          InvoiceStatus   `json:"status" gorm:"type:enum('Paid', 'Unpaid', 'Cancel', 'Expired');not_null"`
	Common
}

type InvoiceStatus string

const (
	Paid    InvoiceStatus = "Paid"
	Unpaid  InvoiceStatus = "Unpaid"
	Cancel  InvoiceStatus = "Cancel"
	Expired InvoiceStatus = "Expired"
)
