package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) repositories.RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (t *refreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	if err := t.db.WithContext(ctx).Create(token).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.NewDBDuplicateKeyError(dup)
		}
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (t *refreshTokenRepository) Update(ctx context.Context, token *entities.RefreshToken) error {
	if err := t.db.WithContext(ctx).Model(token).Select("is_revoked").Updates(token).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (t *refreshTokenRepository) Find(ctx context.Context, tokenStr string) (*entities.RefreshToken, error) {
	var token entities.RefreshToken
	if err := t.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewRefreshTokenNotFound()
		}
		return nil, apperrors.NewDBOperationError(err)
	}
	return &token, nil
}

func (t *refreshTokenRepository) Revoke(ctx context.Context, tokenStr string) error {
	if err := t.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("is_revoked", true).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}
