package persistence

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"context"
	"errors"

	"gorm.io/gorm"
)

const entityTokenName = "refresh token"

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) repositories.RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (t *refreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	if err := t.db.WithContext(ctx).Create(token).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.ErrDuplicateKey(dup)
		}
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (t *refreshTokenRepository) Update(ctx context.Context, token *entities.RefreshToken) error {
	if err := t.db.WithContext(ctx).Model(token).Select("is_revoked").Updates(token).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (t *refreshTokenRepository) Find(ctx context.Context, tokenStr string) (*entities.RefreshToken, error) {
	var token entities.RefreshToken
	if err := t.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound(entityUser)
		}
		return nil, apperrors.ErrDBOperation(err)
	}
	return &token, nil
}

func (t *refreshTokenRepository) Revoke(ctx context.Context, tokenStr string) error {
	if err := t.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("is_revoked", true).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}
