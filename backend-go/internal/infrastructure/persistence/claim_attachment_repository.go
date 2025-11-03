package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type claimAttachmentRepository struct {
	db *gorm.DB
}

func NewClaimAttachmentRepository(db *gorm.DB) repositories.ClaimAttachmentRepository {
	return &claimAttachmentRepository{db: db}
}

func (c *claimAttachmentRepository) Create(tx application.Tx, attachment *entity.ClaimAttachment) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Create(attachment).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperror.NewDBDuplicateKeyError(dup)
		}
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (c *claimAttachmentRepository) HardDelete(tx application.Tx, id uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Unscoped().Delete(&entity.ClaimAttachment{}, "id = ?", id).Error; err != nil {
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (c *claimAttachmentRepository) SoftDeleteByClaimID(tx application.Tx, claimID uuid.UUID) error {
	db := tx.GetTx().(*gorm.DB)
	if err := db.Delete(&entity.ClaimAttachment{}, "claim_id = ?", claimID).Error; err != nil {
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (c *claimAttachmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ClaimAttachment, error) {
	var attachment entity.ClaimAttachment
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewClaimAttachmentNotFound()
		}
		return nil, apperror.NewDBOperationError(err)
	}
	return &attachment, nil
}

func (c *claimAttachmentRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimAttachment, error) {
	var attachments []*entity.ClaimAttachment
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		return nil, apperror.NewDBOperationError(err)
	}
	return attachments, nil
}

func (c *claimAttachmentRepository) CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error) {
	var count int64
	if err := c.db.WithContext(ctx).
		Model(&entity.ClaimAttachment{}).
		Where("claim_id = ?", claimID).
		Count(&count).Error; err != nil {
		return 0, apperror.NewDBOperationError(err)
	}
	return count, nil
}

func (c *claimAttachmentRepository) FindByType(ctx context.Context, claimID uuid.UUID, attachmentType string) ([]*entity.ClaimAttachment, error) {
	var attachments []*entity.ClaimAttachment
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND attachment_type = ?", claimID, attachmentType).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		return nil, apperror.NewDBOperationError(err)
	}
	return attachments, nil
}
