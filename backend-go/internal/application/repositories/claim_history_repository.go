package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type ClaimHistoryRepository interface {
	Create(ctx context.Context, history *entities.ClaimHistory) error
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error)
	FindLatestByClaimID(ctx context.Context, claimID uuid.UUID) (*entities.ClaimHistory, error)
	FindByDateRange(ctx context.Context, claimID uuid.UUID, startDate, endDate time.Time) ([]*entities.ClaimHistory, error)
	SoftDeleteByClaimID(ctx context.Context, claimID uuid.UUID) error
}
