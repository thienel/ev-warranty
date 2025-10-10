package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
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

func (c *claimHistoryRepository) Create(tx application.Transaction, history *entities.ClaimHistory) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Create(history).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.NewDBDuplicateKeyError(dup)
		}
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimHistoryRepository) SoftDeleteByClaimID(tx application.Transaction, claimID uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Delete(&entities.ClaimHistory{}, "claim_id = ?", claimID).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimHistoryRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error) {
	var histories []*entities.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("changed_at DESC").
		Find(&histories).Error; err != nil {
		return nil, apperrors.NewDBOperationError(err)
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
			return nil, apperrors.NewClaimHistoryNotFound()
		}
		return nil, apperrors.NewDBOperationError(err)
	}
	return &history, nil
}

func (c *claimHistoryRepository) FindByDateRange(ctx context.Context, claimID uuid.UUID, startDate, endDate time.Time) ([]*entities.ClaimHistory, error) {
	var histories []*entities.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND changed_at BETWEEN ? AND ?", claimID, startDate, endDate).
		Order("changed_at DESC").
		Find(&histories).Error; err != nil {
		return nil, apperrors.NewDBOperationError(err)
	}
	return histories, nil
}
