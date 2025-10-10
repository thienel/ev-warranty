package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type ClaimAttachmentRepository interface {
	Create(ctx context.Context, attachment *entities.ClaimAttachment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimAttachment, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimAttachment, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
	SoftDeleteByClaimID(ctx context.Context, id uuid.UUID) error
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	FindByType(ctx context.Context, claimID uuid.UUID, attachmentType string) ([]*entities.ClaimAttachment, error)
}
