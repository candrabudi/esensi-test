package models

import "time"

type Invoice struct {
	ID           int       `json:"id"`
	InvoiceNo    string    `json:"invoice_no" gorm:"size:15;not_null"`
	CustomerID   int       `json:"customer_id" gorm:"not_null"`
	CustomerName string    `json:"customer_name" gorm:"not_null"`
	IssueDate    time.Time `json:"issue_date" gorm:"not_null"`
	Subject      string    `json:"subject" gorm:"size:191;not_null"`
	TotalItem    int       `json:"total_item" gorm:"not_null"`
	GrandTotal   float64   `json:"grandtotal" gorm:"not_null"`
	SubTotal     float64   `json:"subtotal" gorm:"not_null"`
	DueDate      time.Time `json:"due_date" gorm:"not_null"`
	Status       string    `json:"status" gorm:"type:enum('Paid', 'Unpaid', 'Cancel', 'Expired');not_null"`
	Common
}
