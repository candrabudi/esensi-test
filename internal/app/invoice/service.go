package invoice

import (
	"context"
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/internal/repository"
	"fmt"
)

type service struct {
	InvoiceRepository repository.Invoice
}

type Service interface {
	FindAll(ctx context.Context, limit, offset int, filterFields map[string]interface{}) (dto.ResultInvoice, error)
	Store(ctx context.Context, input dto.InsertInvoice) (err error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		InvoiceRepository: f.InvoiceRepository,
	}
}

func (s *service) FindAll(ctx context.Context, limit, offset int, filterFields map[string]interface{}) (dto.ResultInvoice, error) {
	fields := "id, invoice_no, customer_name, issue_date, due_date, total_item, grand_total, sub_total, subject, status"
	condition := "deleted_at IS NULL"
	args := []interface{}{}

	for field, value := range filterFields {
		if value != "" {
			switch field {
			case "invoice_no":
				condition += " AND invoice_no = ?"
				args = append(args, value)
			case "issue_date":
				condition += " AND DATE(issue_date) = ?"
				args = append(args, value)
			case "subject":
				subject, ok := value.(string)
				if !ok {
					return dto.ResultInvoice{}, fmt.Errorf("invalid type for subject filter")
				}
				condition += " AND subject LIKE ?"
				args = append(args, "%"+subject+"%")
			case "total_item":
				condition += " AND total_item = ?"
				args = append(args, value)
			case "customer_name":
				condition += " AND customer_name LIKE ?"
				args = append(args, "%"+value.(string)+"%")
			case "due_date":
				condition += " AND DATE(due_date) = ?"
				args = append(args, value)
			case "status":
				status, ok := value.(string)
				if !ok {
					return dto.ResultInvoice{}, fmt.Errorf("invalid type for status filter")
				}
				condition += " AND status = ?"
				args = append(args, status)
			}
		}
	}

	results, err := s.InvoiceRepository.FindAll(ctx, limit, offset, fields, condition, args...)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *service) Store(ctx context.Context, input dto.InsertInvoice) (err error) {

	err = s.InvoiceRepository.Insert(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}