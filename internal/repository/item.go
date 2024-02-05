package repository

import (
	"context"
	"errors"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"
	"fmt"

	"gorm.io/gorm"
)

type Item interface {
	FindAll(ctx context.Context, selectedFields string, query string, args ...any) ([]models.Item, error)
	Insert(ctx context.Context, input *models.Item) (err error)
}

type item struct {
	Db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *item {
	return &item{
		db,
	}
}

func (i *item) FindAll(ctx context.Context, selectedFields string, query string, args ...interface{}) ([]models.Item, error) {
	var res []models.Item

	db := i.Db.WithContext(ctx).Model(models.Item{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Limit(5).Find(&res).Error; err != nil {
		return []models.Item{}, err
	}

	return res, nil
}

func (u *item) Insert(ctx context.Context, input *models.Item) error {
	var existingItem models.Item
	err := u.Db.WithContext(ctx).Model(&models.Item{}).
		Where("item_name = ?", input.ItemName).
		First(&existingItem).Error

	if err == nil {
		return fmt.Errorf("item with name '%s' already exists", input.ItemName)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Failed store item")
	}

	err = u.Db.WithContext(ctx).Create(&input).Error
	if err != nil {
		return errors.New("Failed store item")
	}

	return nil
}
