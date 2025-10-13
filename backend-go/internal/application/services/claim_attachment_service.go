package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/logger"
	"net/url"

	"github.com/google/uuid"
)

type CreateAttachmentCommand struct {
	Type string
	URL  string
}

type ClaimAttachmentService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ClaimAttachment, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimAttachment, error)

	Create(tx application.Tx, claimID uuid.UUID, cmd *CreateAttachmentCommand) (*entities.ClaimAttachment, error)
	HardDelete(tx application.Tx, claimID, attachmentID uuid.UUID) error
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

func (s *claimAttachmentService) Create(tx application.Tx, claimID uuid.UUID,
	cmd *CreateAttachmentCommand) (*entities.ClaimAttachment, error) {
	_, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return nil, err
	}

	if !entities.IsValidAttachmentType(cmd.Type) {
		return nil, apperrors.NewInvalidCredentials()
	}
	if !IsValidURL(cmd.URL) {
		return nil, apperrors.NewInvalidCredentials()
	}

	attachment := entities.NewClaimAttachment(claimID, cmd.Type, cmd.URL)
	err = s.attachRepo.Create(tx, attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *claimAttachmentService) HardDelete(tx application.Tx, claimID, attachmentID uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft {
		return apperrors.NewNotAllowDeleteClaim()
	}
	err = s.attachRepo.HardDelete(tx, attachmentID)
	if err != nil {
		return err
	}

	return nil
}

func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
