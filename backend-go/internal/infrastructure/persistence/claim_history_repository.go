package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/domain/repositories"
	"ev-warranty-go/internal/errors/apperrors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const entityClaimHistory = "claim history"

type claimHistoryRepository struct {
	db *gorm.DB
}

func NewClaimHistoryRepository(db *gorm.DB) repositories.ClaimHistoryRepository {
	return &claimHistoryRepository{db: db}
}

func (c *claimHistoryRepository) Create(ctx context.Context, history *entities.ClaimHistory) error {
	if err := c.db.WithContext(ctx).Create(history).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.ErrDuplicateKey(dup)
		}
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimHistoryRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error) {
	var histories []*entities.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("changed_at DESC").
		Find(&histories).Error; err != nil {
		return nil, apperrors.ErrDBOperation(err)
	}
	return histories, nil
}

func (c *claimHistoryRepository) FindLatestByClaimID(ctx context.Context, claimID uuid.UUID) (*entities.ClaimHistory, error) {
	var history entities.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("changed_at DESC").
		First(&history).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound(entityClaimHistory)
		}
		return nil, apperrors.ErrDBOperation(err)
	}
	return &history, nil
}

func (c *claimHistoryRepository) FindByDateRange(ctx context.Context, claimID uuid.UUID, startDate, endDate time.Time) ([]*entities.ClaimHistory, error) {
	var histories []*entities.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND changed_at BETWEEN ? AND ?", claimID, startDate, endDate).
		Order("changed_at DESC").
		Find(&histories).Error; err != nil {
		return nil, apperrors.ErrDBOperation(err)
	}
	return histories, nil
}
