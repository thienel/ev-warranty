package repositories

import (
	"auth-service/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	FindByOAuth(ctx context.Context, provider, oauthID string) (*entities.User, error)
}
