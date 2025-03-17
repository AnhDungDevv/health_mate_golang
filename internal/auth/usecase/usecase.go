package usecase

import (
	"context"
	"health_backend/config"
	"health_backend/internal/auth"
	"health_backend/internal/models"
	httpErrors "health_backend/pkg/httpError"
	"health_backend/pkg/logger"
	"health_backend/pkg/utils"
	"net/http"

	"github.com/go-playground/validator"
	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"github.com/pkg/errors"

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

// // Register implements auth.UseCase.
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	existUser, _ := u.authRepo.FindByEmail(ctx, user)

	if existUser != nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.EmailAlreadyExistsError, nil)
	}

	validate := validator.New()
	if err := validate.StructCtx(ctx, user); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.ValidateUser"))
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.Register"))
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(user, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil

}

// Login implements auth.UseCase.
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	defer span.Finish()

	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)

	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.GetUser.GenerateJWTToken"))

	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil

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

// Update implements auth.UseCase.
func (a *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}
