package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type ClaimAttachmentRepository interface {
	Create(tx application.Transaction, attachment *entities.ClaimAttachment) error
	HardDelete(tx application.Transaction, id uuid.UUID) error
	SoftDeleteByClaimID(tx application.Transaction, id uuid.UUID) error

	FindByID(ctx context.Context, id uuid.UUID) (*entities.ClaimAttachment, error)
	FindByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimAttachment, error)
	CountByClaimID(ctx context.Context, claimID uuid.UUID) (int64, error)
	FindByType(ctx context.Context, claimID uuid.UUID, attachmentType string) ([]*entities.ClaimAttachment, error)
}
