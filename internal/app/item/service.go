package item

import (
	"context"
	"esensi-test/internal/factory"
	"esensi-test/internal/models"
	"esensi-test/internal/repository"
)

type service struct {
	ItemRepository repository.Item
}

type Service interface {
	FindAll(ctx context.Context, search string) ([]models.Item, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		ItemRepository: f.ItemRepository,
	}
}

func (s *service) FindAll(ctx context.Context, search string) ([]models.Item, error) {
	fields := "id, item_name, item_description, item_type"
	condition := "deleted_at IS NULL"
	query := ""

	if search != "" {
		condition = "item_name LIKE ? AND deleted_at IS NULL"
		query = "%" + search + "%"
	}

	items, err := s.ItemRepository.FindByFields(ctx, fields, condition, query)
	if err != nil {
		return []models.Item{}, err
	}

	return items, nil
}
