package repository

import (
	"context"
	"esensi-test/internal/models"

	"gorm.io/gorm"
)

type User interface {
	Insert(ctx context.Context, input *models.User) (err error)
}

type user struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *user {
	return &user{
		db,
	}
}

func (u *user) Insert(ctx context.Context, input *models.User) (err error) {
	err = u.Db.WithContext(ctx).Create(&input).Error
	if err != nil {
		return err
	}

	return nil
}
