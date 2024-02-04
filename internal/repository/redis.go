package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type redisRepository struct {
	Rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *redisRepository {
	return &redisRepository{
		rdb,
	}
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.Rdb.Set(ctx, key, jsonData, 23*time.Hour).Err()
	return err
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	err := r.Rdb.Del(ctx, key)
	if err != nil {
		return errors.New("failed delete")
	}
	return nil
}
