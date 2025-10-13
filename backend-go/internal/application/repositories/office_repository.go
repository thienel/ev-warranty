package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type OfficeRepository interface {
	Create(ctx context.Context, office *entities.Office) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Office, error)
	FindAll(ctx context.Context) ([]*entities.Office, error)
	Update(ctx context.Context, office *entities.Office) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
