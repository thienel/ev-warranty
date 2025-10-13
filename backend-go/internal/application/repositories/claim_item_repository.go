package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type ClaimItemRepository interface {
	Create(tx application.Tx, item *entities.ClaimItem) error
	Update(tx application.Tx, item *entities.ClaimItem) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error
	UpdateStatus(tx application.Tx, id uuid.UUID, status string) error
	SumCostByClaimID(tx application.Tx, claimID uuid.UUID) (float64, error)

	FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimItem, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimItem, error)
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	FindByStatus(ctx context.Context, claimID uuid.UUID, status string) ([]*entities.ClaimItem, error)
}
