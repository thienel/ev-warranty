package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entities"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	Update(ctx context.Context, token *entities.RefreshToken) error
	Find(ctx context.Context, tokenStr string) (*entities.RefreshToken, error)
	Revoke(ctx context.Context, tokenStr string) error
}
