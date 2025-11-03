package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entity"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	Update(ctx context.Context, token *entity.RefreshToken) error
	Find(ctx context.Context, tokenStr string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, tokenStr string) error
}
