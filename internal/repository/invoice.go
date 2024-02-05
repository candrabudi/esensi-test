package repository

import (
	"context"
	"errors"
	"esensi-test/internal/dto"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Invoice interface {
	FindAll(ctx context.Context, limit int, offset int, selectedFields string, query string, args ...interface{}) (dto.ResultInvoice, error)
	FindByID(ctx context.Context, invoiceID int) (dto.DetailInvoice, error)
	Insert(ctx context.Context, input *dto.InsertInvoice) error
	Update(ctx context.Context, invoiceID int, input *dto.InsertInvoice) error
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
			GrandTotal:   invoice.GrandTotal.InexactFloat64(),
			SubTotal:     invoice.SubTotal.InexactFloat64(),
			TotalItem:    invoice.TotalItem,
			DueDate:      formattedDueDate,
			Status:       string(invoice.Status),
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

func (i *invoice) FindByID(ctx context.Context, invoiceID int) (dto.DetailInvoice, error) {
	var invoice models.Invoice
	var invoiceDetails []models.InvoiceDetail

	if err := i.Db.WithContext(ctx).First(&invoice, invoiceID).Error; err != nil {
		return dto.DetailInvoice{}, err
	}

	if err := i.Db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&invoiceDetails).Error; err != nil {
		return dto.DetailInvoice{}, err
	}

	formattedIssueDate := invoice.IssueDate.Format("2006-01-02")
	formattedDueDate := invoice.DueDate.Format("2006-01-02")

	detailInvoice := dto.DetailInvoice{
		ID:              invoice.ID,
		InvoiceNo:       invoice.InvoiceNo,
		CustomerName:    invoice.CustomerName,
		CustomerAddress: invoice.CustomerAddress,
		Subject:         invoice.Subject,
		IssueDate:       formattedIssueDate,
		GrandTotal:      invoice.GrandTotal.InexactFloat64(),
		SubTotal:        invoice.SubTotal.InexactFloat64(),
		TotalItem:       invoice.TotalItem,
		DueDate:         formattedDueDate,
		Status:          string(invoice.Status),
		Items:           make([]dto.InvoiceDetail, len(invoiceDetails)),
	}

	for idx, detail := range invoiceDetails {
		detailInvoice.Items[idx] = dto.InvoiceDetail{
			ItemID:    detail.ItemID,
			ItemName:  detail.ItemName,
			Quantity:  detail.Qty,
			UnitPrice: detail.ItemPrice.InexactFloat64(),
			Amount:    detail.Amount.InexactFloat64(),
		}
	}

	return detailInvoice, nil
}

func (i *invoice) Insert(ctx context.Context, input *dto.InsertInvoice) error {
	nextInvoiceNo, err := i.getNextInvoiceNo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next invoice number: %w", err)
	}

	customerID, err := i.getCustomerIDByName(ctx, input.CustomerName)
	if err != nil {
		return fmt.Errorf("failed to get customer ID: %w", err)
	}

	issueDate, err := time.Parse("2006-01-02", input.IssueDate)
	if err != nil {
		return fmt.Errorf("failed to parse IssueDate: %w", err)
	}

	dueDate := time.Now().Add(3 * 24 * time.Hour)

	tx := i.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
			tx.Rollback()
		}
	}()
	grandtotal := decimal.NewFromFloat(input.GrandTotal)
	subtotal := decimal.NewFromFloat(input.SubTotal)
	insertInvoice := models.Invoice{
		InvoiceNo:       nextInvoiceNo,
		CustomerID:      customerID,
		CustomerName:    input.CustomerName,
		CustomerAddress: input.CustomerAddress,
		IssueDate:       issueDate,
		Subject:         input.Subject,
		TotalItem:       len(input.Items),
		GrandTotal:      grandtotal,
		SubTotal:        subtotal,
		DueDate:         dueDate,
		Status:          "Unpaid",
	}

	if err := tx.WithContext(ctx).Create(&insertInvoice).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to store invoice: %w", err)
	}

	for _, item := range input.Items {
		unitPrice := decimal.NewFromFloat(item.UnitPrice)
		amount := decimal.NewFromFloat(item.Amount)
		insertInvoiceDetail := models.InvoiceDetail{
			InvoiceID: uint(insertInvoice.ID),
			ItemID:    item.ItemID,
			ItemName:  item.ItemName,
			ItemPrice: unitPrice,
			Qty:       item.Quantity,
			Amount:    amount,
		}

		if err := tx.WithContext(ctx).Create(&insertInvoiceDetail).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to store invoice detail: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (i *invoice) Update(ctx context.Context, invoiceID int, input *dto.InsertInvoice) error {
	customerID, err := i.getCustomerIDByName(ctx, input.CustomerName)
	if err != nil {
		return fmt.Errorf("failed to get customer ID: %w", err)
	}

	issueDate, err := time.Parse("2006-01-02", input.IssueDate)
	if err != nil {
		return fmt.Errorf("failed to parse IssueDate: %w", err)
	}

	tx := i.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
			tx.Rollback()
		}
	}()

	existingInvoice := models.Invoice{}
	if err := tx.WithContext(ctx).First(&existingInvoice, invoiceID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find existing invoice: %w", err)
	}
	grandtotal := decimal.NewFromFloat(input.GrandTotal)
	subtotal := decimal.NewFromFloat(input.SubTotal)

	existingInvoice.CustomerID = customerID
	existingInvoice.CustomerName = input.CustomerName
	existingInvoice.CustomerAddress = input.CustomerAddress
	existingInvoice.IssueDate = issueDate
	existingInvoice.Subject = input.Subject
	existingInvoice.TotalItem = len(input.Items)
	existingInvoice.GrandTotal = grandtotal
	existingInvoice.SubTotal = subtotal
	existingInvoice.DueDate = time.Now().Add(3 * 24 * time.Hour)
	existingInvoice.Status = "Unpaid"

	if err := tx.WithContext(ctx).Save(&existingInvoice).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	for _, item := range input.Items {
		unitPrice := decimal.NewFromFloat(item.UnitPrice)
		amount := decimal.NewFromFloat(item.Amount)
		insertInvoiceDetail := models.InvoiceDetail{
			InvoiceID: uint(invoiceID),
			ItemID:    item.ItemID,
			ItemName:  item.ItemName,
			ItemPrice: unitPrice,
			Qty:       item.Quantity,
			Amount:    amount,
		}

		if err := tx.WithContext(ctx).Where(models.InvoiceDetail{InvoiceID: uint(invoiceID), ItemID: item.ItemID}).Assign(insertInvoiceDetail).FirstOrCreate(&models.InvoiceDetail{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to store invoice detail: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

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
