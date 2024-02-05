package repository

import (
	"context"
	"esensi-test/internal/models"

	"gorm.io/gorm"
)

type User interface {
	Insert(ctx context.Context, input *models.User) (err error)
	FindUser(ctx context.Context, queries []string, argsSlice ...[]interface{}) (models.User, error)
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

func (u *user) FindUser(ctx context.Context, queries []string, argsSlice ...[]interface{}) (models.User, error) {
	var res models.User

	db := u.Db.WithContext(ctx).Model(models.User{})

	for idx, query := range queries {
		if idx < len(argsSlice) {
			db = db.Or(query, argsSlice[idx]...)
		}
	}

	if err := db.Find(&res).Error; err != nil {
		return models.User{}, err
	}

	return res, nil
}
