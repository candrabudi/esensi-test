package customer

import (
	"context"
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/internal/repository"
)

type service struct {
	CustomerRepository repository.Customer
}

type Service interface {
	FindAll(ctx context.Context, search string) ([]dto.FindAllCustomer, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		CustomerRepository: f.CustomerRepository,
	}
}

func (s *service) FindAll(ctx context.Context, search string) ([]dto.FindAllCustomer, error) {
	fields := "id, customer_name, customer_address, customer_phone"
	condition := "deleted_at IS NULL"
	condition = "customer_name LIKE ? AND deleted_at IS NULL"
	query := "%" + search + "%"

	items, err := s.CustomerRepository.FindAll(ctx, fields, condition, query)
	if err != nil {
		return []dto.FindAllCustomer{}, err
	}

	var results []dto.FindAllCustomer
	for _, item := range items {
		ditem := dto.FindAllCustomer{
			ID:              item.ID,
			CustomerName:    item.CustomerName,
			CustomerAddress: item.CustomerAddress,
			CustomerPhone:   item.CustomerPhone,
		}
		results = append(results, ditem)
	}

	return results, nil
}
