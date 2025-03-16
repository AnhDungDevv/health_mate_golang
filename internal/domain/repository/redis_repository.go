package repository

import (
	"context"
	"health_backend/internal/domain/entity"
)

type ChatRedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*entity.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *entity.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
