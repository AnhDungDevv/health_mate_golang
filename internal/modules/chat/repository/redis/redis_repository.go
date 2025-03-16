package redis

import (
	"context"
	"health_backend/internal/domain/models"
	"health_backend/internal/domain/repository"
	baseredis "health_backend/internal/infrastructure/redis"
)

type chatRedisRepository struct {
	*baseredis.BaseRedisRepository
}

func NewChatRedisRepository(base *baseredis.BaseRedisRepository) repository.ChatRedisRepository {
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
