package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type ClaimHistoryRepository interface {
	Create(tx application.Tx, history *entity.ClaimHistory) error
	SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error

	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimHistory, error)
	FindLatestByClaimID(ctx context.Context, claimID uuid.UUID) (*entity.ClaimHistory, error)
	FindByDateRange(ctx context.Context, claimID uuid.UUID, startDate, endDate time.Time) ([]*entity.ClaimHistory, error)
}
