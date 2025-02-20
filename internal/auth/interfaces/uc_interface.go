package auth

import (
	"context"
	"health_backend/internal/models"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	// Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	// Update(ctx context.Context, user *models.User) (*models.User, error)
	// Delete(ctx context.Context, userID uuid.UUID) error
	// GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	// GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}
