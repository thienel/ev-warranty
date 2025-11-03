package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entity"

	"github.com/google/uuid"
)

type ClaimItemRepository interface {
	Create(tx application.Tx, item *entity.ClaimItem) error
	Update(tx application.Tx, item *entity.ClaimItem) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error
	UpdateStatus(tx application.Tx, id uuid.UUID, status string) error
	SumCostByClaimID(tx application.Tx, claimID uuid.UUID) (float64, error)

	FindByID(ctx context.Context, id uuid.UUID) (*entity.ClaimItem, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimItem, error)
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	FindByStatus(ctx context.Context, claimID uuid.UUID, status string) ([]*entity.ClaimItem, error)
}
