package repository

import (
	"context"
	"errors"
	"esensi-test/internal/dto"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Invoice interface {
	FindAll(ctx context.Context, limit int, offset int, selectedFields string, query string, args ...interface{}) (dto.ResultInvoice, error)
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

func (i *invoice) FindAll(ctx context.Context, limit int, offset int, selectedFields string, query string, args ...interface{}) (dto.ResultInvoice, error) {
	var res []models.Invoice
	var totalData int64
	var count int64

	db := i.Db.WithContext(ctx).Model(models.Invoice{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Count(&totalData).Error; err != nil {
		return dto.ResultInvoice{}, err
	}

	if limit > 0 {
		db = db.Limit(limit)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.Where(query, args...).Find(&res).Error; err != nil {
		return dto.ResultInvoice{}, err
	}

	if err := i.Db.WithContext(ctx).Model(models.Invoice{}).Where(query, args...).Count(&count).Error; err != nil {
		return dto.ResultInvoice{}, err
	}

	var invoices []dto.FindAllInvoice
	for _, invoice := range res {
		formattedIssueDate := invoice.IssueDate.Format("2006-01-02")
		formattedDueDate := invoice.DueDate.Format("2006-01-02")

		dinvoice := dto.FindAllInvoice{
			ID:           invoice.ID,
			InvoiceNo:    invoice.InvoiceNo,
			CustomerName: invoice.CustomerName,
			Subject:      invoice.Subject,
			IssueDate:    formattedIssueDate,
			GrandTotal:   invoice.GrandTotal,
			SubTotal:     invoice.SubTotal,
			TotalItem:    invoice.TotalItem,
			DueDate:      formattedDueDate,
			Status:       invoice.Status,
		}
		invoices = append(invoices, dinvoice)
	}

	result := dto.ResultInvoice{
		Items: invoices,
		Metadata: dto.MetaData{
			Limit:        limit,
			Offset:       offset,
			TotalResults: int(totalData),
			Count:        len(invoices),
		},
	}
	if len(invoices) == 0 {
		result.Items = []dto.FindAllInvoice{}
	}
	return result, nil
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
