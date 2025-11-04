package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type claimItemRepository struct {
	db *gorm.DB
}

func NewClaimItemRepository(db *gorm.DB) repository.ClaimItemRepository {
	return &claimItemRepository{db: db}
}

func (c *claimItemRepository) Create(tx application.Tx, item *entity.ClaimItem) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Create(item).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperror.ErrDuplicateKey.WithMessage(dup + " already existed")
		}
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimItemRepository) Update(tx application.Tx, item *entity.ClaimItem) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Model(item).
		Select("part_category_id", "faulty_part_id", "replacement_part_id",
			"issue_description", "status", "type", "cost").
		Updates(item).Error; err != nil {
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimItemRepository) HardDelete(tx application.Tx, id uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Unscoped().Delete(&entity.ClaimItem{}, "id = ?", id).Error; err != nil {
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimItemRepository) SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Delete(&entity.ClaimItem{}, "claim_id = ?", claimID).Error; err != nil {
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimItemRepository) UpdateStatus(tx application.Tx, id uuid.UUID, status string) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Model(&entity.ClaimItem{}).Where("id = ?", id).
		Update("status", status).Error; err != nil {

		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (c *claimItemRepository) SumCostByClaimID(tx application.Tx, claimID uuid.UUID) (float64, error) {
	db := tx.GetTx().(*gorm.DB)
	var totalCost float64
	if err := db.Model(&entity.ClaimItem{}).Where("claim_id = ? AND status = 'APPROVED'", claimID).
		Select("COALESCE(SUM(cost), 0)").Scan(&totalCost).Error; err != nil {

		return 0, apperror.ErrDBOperation.WithError(err)
	}
	return totalCost, nil
}

func (c *claimItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ClaimItem, error) {
	var item entity.ClaimItem
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFoundError.WithMessage("Claim item not found").WithError(err)
		}
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return &item, nil
}

func (c *claimItemRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimItem, error) {
	var items []*entity.ClaimItem
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return items, nil
}

func (c *claimItemRepository) CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error) {
	var count int64
	if err := c.db.WithContext(ctx).
		Model(&entity.ClaimItem{}).
		Where("claim_id = ?", claimID).
		Count(&count).Error; err != nil {
		return 0, apperror.ErrDBOperation.WithError(err)
	}
	return count, nil
}

func (c *claimItemRepository) FindByStatus(ctx context.Context, claimID uuid.UUID, status string) ([]*entity.ClaimItem, error) {
	var items []*entity.ClaimItem
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND status = ?", claimID, status).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return items, nil
}
