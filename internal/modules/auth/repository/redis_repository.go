package repository

import (
	"context"
	auth "health_backend/internal/auth/interfaces"
	"health_backend/internal/models"

	"github.com/redis/go-redis/v9"
)

// Auth redis repository

type authRedisRepo struct {
	redisClient *redis.Client
}

// DeleteUserCtx implements auth.RedisRepository.
func (a *authRedisRepo) DeleteUserCtx(ctx context.Context, key string) error {
	panic("unimplemented")
}

// GetByIDCtx implements auth.RedisRepository.
func (a *authRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	panic("unimplemented")
}

// SetUserCtx implements auth.RedisRepository.
func (a *authRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error {
	panic("unimplemented")
}

func NewAuthRedisRepo(redisClient *redis.Client) auth.RedisRepository {
	return &authRedisRepo{
		redisClient: redisClient,
	}
}
