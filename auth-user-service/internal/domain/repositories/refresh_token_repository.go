package repositories

import (
	"auth-service/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	Update(ctx context.Context, token *entities.RefreshToken) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.RefreshToken, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.RefreshToken, error)
	FindByToken(ctx context.Context, tokenStr string) (*entities.RefreshToken, error)
	RevokeByToken(ctx context.Context, tokenStr string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}
