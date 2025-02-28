package repository

import (
	"context"
	"database/sql"
	auth "health_backend/internal/auth/interfaces"
	"health_backend/internal/models"

	"health_backend/pkg/utils"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	opentracing "github.com/opentracing/opentracing-go"
)

// Auth Repository
type authRepo struct {
	db *gorm.DB
}

// Auth Repository contructor
func NewAuthRepository(db *gorm.DB) auth.Repository {
	return &authRepo{
		db: db,
	}
}

// Delete implements auth.Repository.
func (r *authRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.DeleteUser")
	defer span.Finish()

	result := r.db.WithContext(ctx).Where("id = ? ", userID).Delete(&models.User{})
	if result.Error != nil {
		return errors.Wrap(result.Error, "authRepo.DeleteUser")
	}
	if result.RowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authRepo.DeleteUser")
	}
	return nil
}

// TODO :FinadByEmail implementaion
func (r *authRepo) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.FinByEmail")
	defer span.Finish()

	result := r.db.WithContext(ctx).Where("email= ?", user.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(gorm.ErrRecordNotFound, "user not found")
		}
		return nil, errors.Wrap(result.Error, "authRepo.FindByEmail")
	}

	return user, nil

}

// TODO: FindByName implements auth.Repository.
func (r *authRepo) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	panic("unimplemented")
}

// TODO : GetByID implements auth.Repository.
func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	panic("unimplemented")
}

// TODO: GetUsers implements auth.Repository.
func (r *authRepo) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	panic("unimplemented")
}

// Register implements auth.Repository.
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Register")
	defer span.Finish()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.Create")
	}

	return user, nil
}

// Update implements auth.Repository.
func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {

	panic("unimplemented")

}
