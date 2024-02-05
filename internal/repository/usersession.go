package repository

import (
	"context"
	"esensi-test/internal/models"
	"esensi-test/pkg/util"
	"time"

	"gorm.io/gorm"
)

type UserSession interface {
	Insert(ctx context.Context, input *models.UserSession) (err error)
	FindOneByFields(ctx context.Context, selectedFields string, query string, args ...any) (models.UserSession, error)
	Logout(ctx context.Context, bearer string) error
}

type usersession struct {
	Db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *usersession {
	return &usersession{
		db,
	}
}

func (us *usersession) Insert(ctx context.Context, input *models.UserSession) (err error) {
	err = us.Db.WithContext(ctx).Create(&input).Error
	if err != nil {
		return err
	}

	return nil
}

func (us *usersession) FindOneByFields(ctx context.Context, selectedFields string, query string, args ...any) (models.UserSession, error) {
	var res models.UserSession

	db := us.Db.WithContext(ctx).Model(models.UserSession{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Find(&res).Error; err != nil {
		return models.UserSession{}, err
	}

	return res, nil
}

func (us *usersession) Logout(ctx context.Context, bearer string) error {
	return us.Db.WithContext(ctx).Model(&models.UserSession{}).
		Where("token = ?", bearer).
		Update("deleted_at", time.Now()).Error
}
