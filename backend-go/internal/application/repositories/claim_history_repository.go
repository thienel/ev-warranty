package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type ClaimHistoryRepository interface {
	Create(tx application.Tx, history *entities.ClaimHistory) error
	SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error

	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error)
	FindLatestByClaimID(ctx context.Context, claimID uuid.UUID) (*entities.ClaimHistory, error)
	FindByDateRange(ctx context.Context, claimID uuid.UUID, startDate, endDate time.Time) ([]*entities.ClaimHistory, error)
}
