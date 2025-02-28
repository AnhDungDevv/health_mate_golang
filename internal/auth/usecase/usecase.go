package usecase

import (
	"context"
	"health_backend/config"
	auth "health_backend/internal/auth/interfaces"
	"health_backend/internal/models"
	"health_backend/pkg/logger"
	"health_backend/pkg/utils"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"

	"github.com/opentracing/opentracing-go"
)

type authUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

// Auth usecase contructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, logger: log}
}

// Delete implements auth.UseCase.
func (a *authUC) Delete(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

// FindByName implements auth.UseCase.
func (a *authUC) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	panic("unimplemented")
}

// GetByID implements auth.UseCase.
func (a *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	panic("unimplemented")
}

// GetUsers implements auth.UseCase.
func (a *authUC) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	panic("unimplemented")
}

// Login implements auth.UseCase.
func (a *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	panic("unimplemented")
}

// // Register implements auth.UseCase.
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	return &models.UserWithToken{
		User:  createdUser,
		Token: "jaskldfklas",
	}, nil

}

// Update implements auth.UseCase.
func (a *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}
