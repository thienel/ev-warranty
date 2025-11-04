package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type claimHistoryRepository struct {
	db *gorm.DB
}

func NewClaimHistoryRepository(db *gorm.DB) repository.ClaimHistoryRepository {
	return &claimHistoryRepository{db: db}
}

func (c *claimHistoryRepository) Create(tx application.Tx, history *entity.ClaimHistory) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Create(history).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperror.ErrDuplicateKey.WithMessage(dup + " already existed").WithError(err)
		}
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimHistoryRepository) SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Delete(&entity.ClaimHistory{}, "claim_id = ?", claimID).Error; err != nil {
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimHistoryRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimHistory, error) {
	var histories []*entity.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("changed_at DESC").
		Find(&histories).Error; err != nil {
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return histories, nil
}

func (c *claimHistoryRepository) FindLatestByClaimID(ctx context.Context, claimID uuid.UUID) (*entity.ClaimHistory, error) {
	var history entity.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("changed_at DESC").
		First(&history).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFoundError.WithMessage("claim history not found").WithError(err)
		}
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return &history, nil
}

func (c *claimHistoryRepository) FindByDateRange(ctx context.Context, claimID uuid.UUID, startDate, endDate time.Time) ([]*entity.ClaimHistory, error) {
	var histories []*entity.ClaimHistory
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND changed_at BETWEEN ? AND ?", claimID, startDate, endDate).
		Order("changed_at DESC").
		Find(&histories).Error; err != nil {
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return histories, nil
}
