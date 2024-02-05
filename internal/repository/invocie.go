package repository

import (
	"context"
	"errors"
	"esensi-test/internal/dto"
	"esensi-test/internal/models"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Invoice interface {
	Insert(ctx context.Context, input *dto.InsertInvoice) error
}

type invoice struct {
	Db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *invoice {
	return &invoice{
		db,
	}
}

func (i *invoice) Insert(ctx context.Context, input *dto.InsertInvoice) error {
	nextInvoiceNo, err := i.getNextInvoiceNo(ctx)
	if err != nil {
		return err
	}

	customerID, err := i.getCustomerIDByName(ctx, input.CustomerName)
	if err != nil {
		return err
	}

	issueDate, err := time.Parse("2006-01-02", input.IssueDate)
	if err != nil {
		return errors.New("Failed to parse IssueDate")
	}

	dueDate := time.Now().Add(3 * 24 * time.Hour)

	tx := i.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	insertInvoice := models.Invoice{
		InvoiceNo:    nextInvoiceNo,
		CustomerID:   customerID,
		CustomerName: input.CustomerName,
		IssueDate:    issueDate,
		Subject:      input.Subject,
		TotalItem:    len(input.Items),
		GrandTotal:   input.GrandTotal,
		SubTotal:     input.SubTotal,
		DueDate:      dueDate,
		Status:       "Unpaid",
	}

	if err := tx.WithContext(ctx).Create(&insertInvoice).Error; err != nil {
		tx.Rollback()
		return errors.New("Failed to store invoice")
	}

	for _, item := range input.Items {
		insertInvoiceDetail := models.InvoiceDetail{
			InvoiceID: insertInvoice.ID,
			ItemID:    item.ItemID,
			ItemName:  item.ItemName,
			ItemPrice: item.UnitPrice,
			Qty:       item.Quantity,
			Amount:    item.Amount,
		}

		if err := tx.WithContext(ctx).Create(&insertInvoiceDetail).Error; err != nil {
			tx.Rollback()
			return errors.New("Failed to store invoice detail")
		}
	}

	tx.Commit()
	return nil
}

func (i *invoice) getNextInvoiceNo(ctx context.Context) (string, error) {
	var lastInvoice models.Invoice
	if err := i.Db.WithContext(ctx).Order("invoice_no desc").First(&lastInvoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "0001", nil
		}
		return "", errors.New("Failed to get the last invoice number")
	}

	lastInvoiceNo, err := strconv.Atoi(lastInvoice.InvoiceNo)
	if err != nil {
		return "", errors.New("Failed to convert last invoice number")
	}

	nextInvoiceNo := strconv.Itoa(lastInvoiceNo + 1)
	formattedInvoiceNo := fmt.Sprintf("%04s", nextInvoiceNo)

	return formattedInvoiceNo, nil
}

func (i *invoice) getCustomerIDByName(ctx context.Context, customerName string) (int, error) {
	var customer models.Customer
	if err := i.Db.WithContext(ctx).Where("customer_name = ?", customerName).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("Customer not found")
		}
		return 0, errors.New("Failed to get the customer ID")
	}

	return customer.ID, nil
}
