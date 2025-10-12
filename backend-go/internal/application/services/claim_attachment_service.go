package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/logger"

	"github.com/google/uuid"
)

type CreateAttachmentCommand struct {
	Type string
	URL  string
}

type ClaimAttachmentService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ClaimAttachment, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimAttachment, error)

	Create(tx application.Transaction, claimID uuid.UUID, cmd *CreateAttachmentCommand) (*entities.ClaimAttachment, error)
	HardDelete(tx application.Transaction, claimID, attachmentID uuid.UUID) error
}

type claimAttachmentService struct {
	log        logger.Logger
	claimRepo  repositories.ClaimRepository
	attachRepo repositories.ClaimAttachmentRepository
}

func NewClaimAttachmentService(log logger.Logger, claimRepo repositories.ClaimRepository, attachRepo repositories.ClaimAttachmentRepository) ClaimAttachmentService {
	return &claimAttachmentService{
		log:        log,
		claimRepo:  claimRepo,
		attachRepo: attachRepo,
	}
}

func (s *claimAttachmentService) GetByID(ctx context.Context, id uuid.UUID) (*entities.ClaimAttachment, error) {
	claimAttachment, err := s.attachRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return claimAttachment, nil
}

func (s *claimAttachmentService) GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimAttachment, error) {
	claimAttachments, err := s.attachRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return claimAttachments, nil
}

func (s *claimAttachmentService) Create(tx application.Transaction, claimID uuid.UUID,
	cmd *CreateAttachmentCommand) (*entities.ClaimAttachment, error) {

	_, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if !entities.IsValidAttachmentType(cmd.Type) {
		_ = tx.Rollback()
		return nil, apperrors.NewInvalidCredentials()
	}
	if !IsValidURL(cmd.URL) {
		_ = tx.Rollback()
		return nil, apperrors.NewInvalidCredentials()
	}

	attachment := entities.NewClaimAttachment(claimID, cmd.Type, cmd.URL)
	err = s.attachRepo.Create(tx, attachment)
	if err != nil {
		return nil, rollbackOnErr(tx, err)
	}

	return attachment, commitOrLog(tx)
}

func (s *claimAttachmentService) HardDelete(tx application.Transaction, claimID, attachmentID uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if claim.Status != entities.ClaimStatusDraft {
		_ = tx.Rollback()
		return apperrors.NewNotAllowDeleteClaim()
	}
	err = s.attachRepo.HardDelete(tx, attachmentID)
	if err != nil {
		return rollbackOnErr(tx, err)
	}

	return commitOrLog(tx)
}
