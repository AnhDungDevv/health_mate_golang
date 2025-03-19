package repository

import (
	"context"
	"database/sql"
	"health_backend/internal/auth"
	"health_backend/internal/models"
	"health_backend/pkg/utils"

	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

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

// Register implements auth.Repository.
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Register")
	defer span.Finish()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.Create")
	}
	return user, nil
}

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

func (r *authRepo) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.FinByEmail")
	defer span.Finish()

	foundUser := &models.User{}
	if err := r.db.WithContext(ctx).Where("email = ?", user.Email).First(foundUser).Error; err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail")
	}
	return foundUser, nil
}

// TODO: FindByName implements auth.Repository.
func (r *authRepo) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	return nil, errors.Errorf("Not implement")

}

// TODO : GetByID implements auth.Repository.
func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return nil, errors.Errorf("Not implement")

}

// TODO: GetUsers implements auth.Repository.
func (r *authRepo) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	return nil, errors.Errorf("Not implement")

}

// Update implements auth.Repository.
func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {

	return nil, errors.Errorf("Not implement")
}
