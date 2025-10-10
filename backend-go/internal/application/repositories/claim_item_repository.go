package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type ClaimItemRepository interface {
	Create(ctx context.Context, item *entities.ClaimItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimItem, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimItem, error)
	Update(ctx context.Context, item *entities.ClaimItem) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	SoftDeleteByClaimID(ctx context.Context, claimID uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	SumCostByClaimID(ctx context.Context, claimID uuid.UUID) (float64, error)
	FindByStatus(ctx context.Context, claimID uuid.UUID, status string) ([]*entities.ClaimItem, error)
}
