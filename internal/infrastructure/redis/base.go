package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type BaseRedisRepository struct {
	client *redis.Client
}

func NewBaseRedisRepository(client *redis.Client) *BaseRedisRepository {
	return &BaseRedisRepository{
		client: client,
	}
}

// Commnon method

func (r *BaseRedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *BaseRedisRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *BaseRedisRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *BaseRedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *BaseRedisRepository) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(ctx, key, data, expiration).Result()
}

func (r *BaseRedisRepository) HSet(ctx context.Context, key string, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, key, field, data).Err()
}

func (r *BaseRedisRepository) HGet(ctx context.Context, key string, field string, dest interface{}) error {
	data, err := r.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
