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

func (r *BaseRedisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key không tồn tại
	} else if err != nil {
		return "", err // Lỗi Redis
	}
	return val, nil
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

// LPush - Thêm phần tử vào đầu danh sách trong Redis
func (r *BaseRedisRepository) LPush(ctx context.Context, key string, value interface{}) error {
	return r.client.LPush(ctx, key, value).Err()
}

// LRange - Lấy các phần tử trong một dãy từ Redis
func (r *BaseRedisRepository) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// Expire - Thiết lập thời gian hết hạn cho khóa trong Redis
func (r *BaseRedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}
