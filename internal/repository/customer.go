package repository

import (
	"context"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"

	"gorm.io/gorm"
)

type Customer interface {
	FindAll(ctx context.Context, selectedFields string, query string, args ...interface{}) ([]models.Customer, error)
}

type customer struct {
	Db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *customer {
	return &customer{
		db,
	}
}

func (c *customer) FindAll(ctx context.Context, selectedFields string, query string, args ...interface{}) ([]models.Customer, error) {
	var res []models.Customer

	db := c.Db.WithContext(ctx).Model(models.Customer{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Limit(5).Find(&res).Error; err != nil {
		return []models.Customer{}, err
	}

	return res, nil
}
