package repository

import (
	"context"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"

	"gorm.io/gorm"
)

type Item interface {
	FindByFields(ctx context.Context, selectedFields string, query string, args ...any) ([]models.Item, error)
}

type item struct {
	Db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *item {
	return &item{
		db,
	}
}

func (i *item) FindByFields(ctx context.Context, selectedFields string, query string, args ...any) ([]models.Item, error) {
	var res []models.Item

	db := i.Db.WithContext(ctx).Model(models.Item{})
	db = util.SetSelectFields(db, selectedFields)

	if args != nil {
		if err := db.Where(query, args...).Find(&res).Error; err != nil {
			return []models.Item{}, err
		}
	} else {
		if err := db.Where(query).Find(&res).Error; err != nil {
			return []models.Item{}, err
		}
	}

	return res, nil
}
