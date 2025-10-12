package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type ClaimItemRepository interface {
	Create(tx application.Transaction, item *entities.ClaimItem) error
	Update(tx application.Transaction, item *entities.ClaimItem) error
	HardDelete(tx application.Transaction, id uuid.UUID) error
	SoftDeleteByClaimID(tx application.Transaction, claimID uuid.UUID) error
	UpdateStatus(tx application.Transaction, id uuid.UUID, status string) error
	SumCostByClaimID(tx application.Transaction, claimID uuid.UUID) (float64, error)

	FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimItem, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimItem, error)
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	FindByStatus(ctx context.Context, claimID uuid.UUID, status string) ([]*entities.ClaimItem, error)
}
