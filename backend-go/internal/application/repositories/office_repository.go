package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entity"

	"github.com/google/uuid"
)

type OfficeRepository interface {
	Create(ctx context.Context, office *entity.Office) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Office, error)
	FindAll(ctx context.Context) ([]*entity.Office, error)
	Update(ctx context.Context, office *entity.Office) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
