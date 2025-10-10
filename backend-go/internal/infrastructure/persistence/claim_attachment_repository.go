package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/domain/repositories"
	"ev-warranty-go/internal/errors/apperrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const entityClaimAttachment = "claim attachment"

type claimAttachmentRepository struct {
	db *gorm.DB
}

func NewClaimAttachmentRepository(db *gorm.DB) repositories.ClaimAttachmentRepository {
	return &claimAttachmentRepository{db: db}
}

func (c *claimAttachmentRepository) Create(ctx context.Context, attachment *entities.ClaimAttachment) error {
	if err := c.db.WithContext(ctx).Create(attachment).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.ErrDuplicateKey(dup)
		}
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimAttachmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimAttachment, error) {
	var attachment entities.ClaimAttachment
	if err := c.db.WithContext(ctx).Where("id = ?", id).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound(entityClaimAttachment)
		}
		return nil, apperrors.ErrDBOperation(err)
	}
	return &attachment, nil
}

func (c *claimAttachmentRepository) FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimAttachment, error) {
	var attachments []*entities.ClaimAttachment
	if err := c.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		return nil, apperrors.ErrDBOperation(err)
	}
	return attachments, nil
}

func (c *claimAttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.db.WithContext(ctx).Delete(&entities.ClaimAttachment{}, "id = ?", id).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (c *claimAttachmentRepository) CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error) {
	var count int64
	if err := c.db.WithContext(ctx).
		Model(&entities.ClaimAttachment{}).
		Where("claim_id = ?", claimID).
		Count(&count).Error; err != nil {
		return 0, apperrors.ErrDBOperation(err)
	}
	return count, nil
}

func (c *claimAttachmentRepository) FindByType(ctx context.Context, claimID uuid.UUID, attachmentType string) ([]*entities.ClaimAttachment, error) {
	var attachments []*entities.ClaimAttachment
	if err := c.db.WithContext(ctx).
		Where("claim_id = ? AND attachment_type = ?", claimID, attachmentType).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		return nil, apperrors.ErrDBOperation(err)
	}
	return attachments, nil
}
