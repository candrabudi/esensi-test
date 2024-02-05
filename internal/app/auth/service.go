package auth

import (
	"context"
	"errors"
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/internal/models"
	"esensi-test/internal/repository"
	"esensi-test/pkg/util"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	UserRepository        repository.User
	UserSessionRepository repository.UserSession
}

type Service interface {
	Login(ctx context.Context, payload dto.PayloadLogin) (dto.ResponseLoginUser, error)
	Logout(ctx context.Context, bearer string) (err error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository:        f.UserRepository,
		UserSessionRepository: f.UserSessionRepository,
	}
}

func (s *service) Login(ctx context.Context, payload dto.PayloadLogin) (dto.ResponseLoginUser, error) {
	email := payload.Email
	password := payload.Password
	queries := []string{"email = ?"}
	argsSlice := [][]interface{}{
		{email},
	}
	data, err := s.UserRepository.FindUser(ctx, queries, argsSlice...)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		return dto.ResponseLoginUser{}, err
	}
	err = ComparePasswords(data.Password, password)
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}

	userID := strconv.Itoa(data.ID)
	secretKey := []byte(util.GetEnv("JWT_KEY", "fallback"))
	jwt, err := GenerateToken(secretKey, userID, data.Email)
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}

	inserSession := models.UserSession{
		UserID: data.ID,
		Token:  jwt,
	}

	err = s.UserSessionRepository.Insert(ctx, &inserSession)
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}

	response := dto.ResponseLoginUser{
		Name:        data.Name,
		Email:       data.Email,
		AccessToken: jwt,
	}
	return response, nil
}

func (s *service) Logout(ctx context.Context, bearer string) (err error) {

	err = s.UserSessionRepository.Logout(ctx, bearer)
	if err != nil {
		return errors.New("Sorry can`t logout, check your bearer")
	}

	return nil
}

func ComparePasswords(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func GenerateToken(secretKey []byte, userID string, email string) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", err
	}

	var expiredTime int64
	jwtMode := util.GetEnv("JWT_MODE", "fallback")
	if jwtMode == "release" {
		expiredTime = time.Now().In(loc).Add(time.Hour * 72).Unix()
	} else {
		expiredTime = time.Now().In(loc).Add(time.Hour * 72).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     expiredTime,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
