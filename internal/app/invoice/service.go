package invoice

import (
	"context"
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/internal/repository"
)

type service struct {
	InvoiceRepository repository.Invoice
}

type Service interface {
	Store(ctx context.Context, input dto.InsertInvoice) (err error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		InvoiceRepository: f.InvoiceRepository,
	}
}

func (s *service) Store(ctx context.Context, input dto.InsertInvoice) (err error) {

	err = s.InvoiceRepository.Insert(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}
