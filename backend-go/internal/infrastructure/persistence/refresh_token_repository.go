package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) repositories.RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (t *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	if err := t.db.WithContext(ctx).Create(token).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperror.NewDBDuplicateKeyError(dup)
		}
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (t *refreshTokenRepository) Update(ctx context.Context, token *entity.RefreshToken) error {
	if err := t.db.WithContext(ctx).Model(token).Select("is_revoked").Updates(token).Error; err != nil {
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (t *refreshTokenRepository) Find(ctx context.Context, tokenStr string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	if err := t.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewRefreshTokenNotFound()
		}
		return nil, apperror.NewDBOperationError(err)
	}
	return &token, nil
}

func (t *refreshTokenRepository) Revoke(ctx context.Context, tokenStr string) error {
	if err := t.db.WithContext(ctx).Model(&entity.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("is_revoked", true).Error; err != nil {
		return apperror.NewDBOperationError(err)
	}
	return nil
}
