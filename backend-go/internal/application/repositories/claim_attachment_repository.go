package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entity"

	"github.com/google/uuid"
)

type ClaimAttachmentRepository interface {
	Create(tx application.Tx, attachment *entity.ClaimAttachment) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDeleteByClaimID(tx application.Tx, id uuid.UUID) error

	FindByID(ctx context.Context, id uuid.UUID) (*entity.ClaimAttachment, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimAttachment, error)
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	FindByType(ctx context.Context, claimID uuid.UUID, attachmentType string) ([]*entity.ClaimAttachment, error)
}
