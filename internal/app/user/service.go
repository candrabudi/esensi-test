package user

import (
	"context"
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/internal/models"
	"esensi-test/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	UserRepository repository.User
}

type Service interface {
	CreateUser(ctx context.Context, input dto.InsertUserRequest) (err error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
	}
}

func (s *service) CreateUser(ctx context.Context, input dto.InsertUserRequest) (err error) {
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	inputUser := models.User{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    string(password),
	}

	err = s.UserRepository.Insert(ctx, &inputUser)
	if err != nil {
		return err
	}

	return nil
}
