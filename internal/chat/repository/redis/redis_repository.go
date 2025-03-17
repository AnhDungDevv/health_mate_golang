package redis

import (
	"context"

	"health_backend/internal/auth"
	baseredis "health_backend/internal/infrastructure/redis"
	"health_backend/internal/models"
)

type chatRedisRepository struct {
	*baseredis.BaseRedisRepository
}

func NewChatRedisRepository(base *baseredis.BaseRedisRepository) auth.RedisRepository {
	return &chatRedisRepository{
		BaseRedisRepository: base,
	}
}

// DeleteUserCtx implements repository.ChatRedisRepository.
func (c *chatRedisRepository) DeleteUserCtx(ctx context.Context, key string) error {
	panic("unimplemented")
}

// GetByIDCtx implements repository.ChatRedisRepository.
func (c *chatRedisRepository) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	panic("unimplemented")
}

// SetUserCtx implements repository.ChatRedisRepository.
func (c *chatRedisRepository) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error {
	panic("unimplemented")
}
