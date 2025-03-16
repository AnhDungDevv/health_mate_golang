package usecase

import (
	"context"
	"health_backend/internal/domain/entity"
	"health_backend/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	Login(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*entity.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*entity.UsersList, error)
}
