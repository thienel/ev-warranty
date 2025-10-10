package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const entityClaimItem = "claim item"

type claimItemRepository struct {
	db *gorm.DB
}

func NewClaimItemRepository(db *gorm.DB) repositories.ClaimItemRepository {
	return &claimItemRepository{db: db}
}

func (c *claimItemRepository) Create(tx application.Transaction, item *entities.ClaimItem) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Create(item).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.ErrDuplicateKey(dup)
		}
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimItemRepository) Update(tx application.Transaction, item *entities.ClaimItem) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Model(item).
		Select("part_category_id", "faulty_part_id", "replacement_part_id",
			"issue_description", "line_status", "type", "cost").
		Updates(item).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimItemRepository) HardDelete(tx application.Transaction, id uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Unscoped().Delete(&entities.ClaimItem{}, "id = ?", id).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimItemRepository) SoftDeleteByClaimID(tx application.Transaction, claimID uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Delete(&entities.ClaimItem{}, "claim_id = ?", claimID).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimItemRepository) UpdateStatus(tx application.Transaction, id uuid.UUID, status string) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Model(&entities.ClaimItem{}).Where("id = ?", id).
		Update("line_status", status).Error; err != nil {

		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimItem, error) {
	var item entities.ClaimItem
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound(entityClaimItem)
		}
		return nil, apperrors.ErrDBOperation(err)
	}
	return &item, nil
}

func (c *claimItemRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimItem, error) {
	var items []*entities.ClaimItem
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, apperrors.ErrDBOperation(err)
	}
	return items, nil
}

func (c *claimItemRepository) CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error) {
	var count int64
	if err := c.db.WithContext(ctx).
		Model(&entities.ClaimItem{}).
		Where("claim_id = ?", claimID).
		Count(&count).Error; err != nil {
		return 0, apperrors.ErrDBOperation(err)
	}
	return count, nil
}

func (c *claimItemRepository) SumCostByClaimID(ctx context.Context, claimID uuid.UUID) (float64, error) {
	var totalCost float64
	if err := c.db.WithContext(ctx).
		Model(&entities.ClaimItem{}).
		Where("claim_id = ?", claimID).
		Select("COALESCE(SUM(cost), 0)").
		Scan(&totalCost).Error; err != nil {
		return 0, apperrors.ErrDBOperation(err)
	}
	return totalCost, nil
}

func (c *claimItemRepository) FindByStatus(ctx context.Context, claimID uuid.UUID, status string) ([]*entities.ClaimItem, error) {
	var items []*entities.ClaimItem
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND line_status = ?", claimID, status).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, apperrors.ErrDBOperation(err)
	}
	return items, nil
}
