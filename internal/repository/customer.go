package repository

import (
	"context"
	"errors"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"
	"fmt"

	"gorm.io/gorm"
)

type Customer interface {
	FindAll(ctx context.Context, selectedFields string, query string, args ...interface{}) ([]models.Customer, error)
	Insert(ctx context.Context, input *models.Customer) error
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

func (c *customer) Insert(ctx context.Context, input *models.Customer) error {
	var existingCustomer models.Customer
	err := c.Db.WithContext(ctx).Model(&models.Customer{}).
		Where("customer_name = ?", input.CustomerName).
		First(&existingCustomer).Error

	if err == nil {
		return fmt.Errorf("customer with name '%s' already exists", input.CustomerName)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Failed store customer")
	}

	err = c.Db.WithContext(ctx).Create(&input).Error
	if err != nil {
		return errors.New("Failed store customer")
	}

	return nil
}
