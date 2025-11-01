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

type claimRepository struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) repositories.ClaimRepository {
	return &claimRepository{db: db}
}

func (c *claimRepository) Create(tx application.Tx, claim *entities.Claim) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Create(claim).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.NewDBDuplicateKeyError(dup)
		}
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimRepository) Update(tx application.Tx, claim *entities.Claim) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Model(claim).Select("vehicle_id",
		"customer_id", "description", "status", "total_cost", "approved_by").
		Updates(claim).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimRepository) HardDelete(tx application.Tx, id uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Unscoped().Delete(&entities.Claim{}, "id = ?", id).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimRepository) SoftDelete(tx application.Tx, id uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Delete(&entities.Claim{}, "id = ?", id).Error; err != nil {
		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimRepository) UpdateStatus(tx application.Tx, id uuid.UUID, status string) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Model(&entities.Claim{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {

		return apperrors.NewDBOperationError(err)
	}
	return nil
}

func (c *claimRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error) {
	var claim entities.Claim
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(&claim).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewClaimNotFound()
		}
		return nil, apperrors.NewDBOperationError(err)
	}
	return &claim, nil
}

func (c *claimRepository) FindAll(ctx context.Context) ([]*entities.Claim, error) {
	var claims []*entities.Claim
	if err := c.db.WithContext(ctx).Find(&claims).Error; err != nil {
		return nil, apperrors.NewDBOperationError(err)
	}

	return claims, nil
}

func (c *claimRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Claim, error) {
	var claims []*entities.Claim
	if err := c.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&claims).Error; err != nil {
		return nil, apperrors.NewDBOperationError(err)
	}
	return claims, nil
}

func (c *claimRepository) FindByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]*entities.Claim, error) {
	var claims []*entities.Claim
	if err := c.db.WithContext(ctx).
		Where("vehicle_id = ?", vehicleID).
		Order("created_at DESC").
		Find(&claims).Error; err != nil {
		return nil, apperrors.NewDBOperationError(err)
	}
	return claims, nil
}
