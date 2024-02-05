package item

import (
	"context"
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/internal/models"
	"esensi-test/internal/repository"
)

type service struct {
	ItemRepository repository.Item
}

type Service interface {
	FindAll(ctx context.Context, search string) ([]dto.FindAllItem, error)
	Store(ctx context.Context, input dto.InsertItem) (err error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		ItemRepository: f.ItemRepository,
	}
}

func (s *service) FindAll(ctx context.Context, search string) ([]dto.FindAllItem, error) {
	fields := "id, item_name, item_description, item_type, item_price, item_stock"
	condition := "deleted_at IS NULL"
	condition = "item_name LIKE ? AND deleted_at IS NULL"
	query := "%" + search + "%"

	items, err := s.ItemRepository.FindAll(ctx, fields, condition, query)
	if err != nil {
		return []dto.FindAllItem{}, err
	}

	var results []dto.FindAllItem
	for _, item := range items {
		ditem := dto.FindAllItem{
			ItemName:        item.ItemName,
			ItemDescription: item.ItemDescription,
			ItemType:        item.ItemType,
			ItemPrice:       item.ItemPrice,
			ItemStock:       item.ItemStock,
		}
		results = append(results, ditem)
	}

	return results, nil
}

func (s *service) Store(ctx context.Context, input dto.InsertItem) (err error) {
	inputItem := models.Item{
		ItemName:        input.ItemName,
		ItemDescription: input.ItemDescription,
		ItemType:        input.ItemType,
		ItemPrice:       input.ItemPrice,
		ItemStock:       input.ItemStock,
	}

	err = s.ItemRepository.Insert(ctx, &inputItem)
	if err != nil {
		return err
	}

	return nil
}
