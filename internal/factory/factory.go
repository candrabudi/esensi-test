package factory

import (
	"esensi-test/database"
	"esensi-test/internal/repository"
)

type Factory struct {
	RedisRepository repository.RedisRepository
}

func NewFactory() *Factory {
	// Check db connection
	rdb := database.GetRedisClient()
	return &Factory{
		// Pass the db connection to repository package for database query calling
		repository.NewRedisRepository(rdb),
	}
}
