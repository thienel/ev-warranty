package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db}
}

func (u *userRepository) Create(ctx context.Context, user *entities.User) error {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.NewDBDuplicateKeyError(dup)
		}
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (u *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFound()
		}
		return nil, apperrors.NewDBOperationError(err)
	}
	return &user, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFound()
		}
		return nil, apperrors.NewDBOperationError(err)
	}
	return &user, nil
}

func (u *userRepository) FindAll(ctx context.Context) ([]*entities.User, error) {
	var users []*entities.User
	if err := u.db.WithContext(ctx).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFound()
		}
		return nil, apperrors.NewDBOperationError(err)
	}
	return users, nil
}

func (u *userRepository) Update(ctx context.Context, user *entities.User) error {
	if err := u.db.WithContext(ctx).Model(user).
		Select("name", "email", "role",
			"password_hash", "is_active", "office_id", "oauth_provider", "oauth_id").
		Updates(user).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (u *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := u.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (u *userRepository) FindByOAuth(ctx context.Context, provider, oauthID string) (*entities.User, error) {
	var user entities.User
	if err := u.db.WithContext(ctx).Where("oauth_provider = ? AND oauth_id = ?", provider, oauthID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFound()
		}
	}
	return &user, nil
}

func getDuplicateKeyConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName
	}
	return ""
}
