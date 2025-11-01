package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/infrastructure/cloudinary"
	"ev-warranty-go/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type CreateClaimCommand struct {
	VehicleID   uuid.UUID
	CustomerID  uuid.UUID
	CreatorID   uuid.UUID
	Description string
}

type UpdateClaimCommand struct {
	Description string
}

type ClaimFilters struct {
	CustomerID *uuid.UUID
	VehicleID  *uuid.UUID
	Status     *string
	FromDate   *time.Time
	ToDate     *time.Time
}

type Pagination struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string
}

type ClaimListResult struct {
	Claims     []*entities.Claim `json:"claims"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

type ClaimService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error)
	GetAll(ctx context.Context) ([]*entities.Claim, error)

	Create(tx application.Tx, cmd *CreateClaimCommand) (*entities.Claim, error)
	Update(tx application.Tx, id uuid.UUID, cmd *UpdateClaimCommand) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDelete(tx application.Tx, id uuid.UUID) error

	UpdateStatus(tx application.Tx, id uuid.UUID, status string, changedBy uuid.UUID) error
	Submit(tx application.Tx, id uuid.UUID, changedBy uuid.UUID) error
	Complete(tx application.Tx, id uuid.UUID, changedBy uuid.UUID) error

	GetHistory(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error)
}

type claimService struct {
	log            logger.Logger
	claimRepo      repositories.ClaimRepository
	itemRepo       repositories.ClaimItemRepository
	attachmentRepo repositories.ClaimAttachmentRepository
	historyRepo    repositories.ClaimHistoryRepository
	cloudService   cloudinary.CloudinaryService
}

func NewClaimService(
	log logger.Logger,
	claimRepo repositories.ClaimRepository,
	itemRepo repositories.ClaimItemRepository,
	attachmentRepo repositories.ClaimAttachmentRepository,
	historyRepo repositories.ClaimHistoryRepository,
	cloudService cloudinary.CloudinaryService,
) ClaimService {
	return &claimService{
		log:            log,
		claimRepo:      claimRepo,
		itemRepo:       itemRepo,
		attachmentRepo: attachmentRepo,
		historyRepo:    historyRepo,
		cloudService:   cloudService,
	}
}

func (s *claimService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error) {
	return s.claimRepo.FindByID(ctx, id)
}

func (s *claimService) GetAll(ctx context.Context) ([]*entities.Claim, error) {
	claims, err := s.claimRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return claims, err
}

func (s *claimService) Create(tx application.Tx, cmd *CreateClaimCommand) (*entities.Claim, error) {
	claim := entities.NewClaim(cmd.VehicleID, cmd.CustomerID, cmd.Description, entities.ClaimStatusDraft, nil)

	if err := s.claimRepo.Create(tx, claim); err != nil {
		return nil, err
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusDraft, cmd.CreatorID)
	if err := s.historyRepo.Create(tx, history); err != nil {
		return nil, err
	}

	return claim, nil
}

func (s *claimService) Update(tx application.Tx, id uuid.UUID, cmd *UpdateClaimCommand) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft && claim.Status != entities.ClaimStatusRequestInfo {
		return apperrors.NewNotAllowUpdateClaim()
	}

	claim.Description = cmd.Description

	if err = s.claimRepo.Update(tx, claim); err != nil {
		return err
	}

	return nil
}

func (s *claimService) HardDelete(tx application.Tx, id uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft {
		return apperrors.NewNotAllowDeleteClaim()
	}

	attachments, err := s.attachmentRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	err = s.claimRepo.HardDelete(tx, id)
	if err == nil {
		for _, attach := range attachments {
			err := s.cloudService.DeleteFileByURL(context.Background(), attach.URL)
			if err != nil {
				s.log.Error("[Cloudinary] Failed to delete file in delete claim use case", "error", err)
			}
		}
	}
	return err
}

func (s *claimService) SoftDelete(tx application.Tx, id uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusCancelled {
		return apperrors.NewNotAllowDeleteClaim()
	}

	softDeleters := []func(application.Tx, uuid.UUID) error{
		s.claimRepo.SoftDelete,
		s.itemRepo.SoftDeleteByClaimID,
		s.attachmentRepo.SoftDeleteByClaimID,
		s.historyRepo.SoftDeleteByClaimID,
	}

	for _, deleteFn := range softDeleters {
		if err = deleteFn(tx, id); err != nil {
			return err
		}
	}

	return nil
}

func (s *claimService) UpdateStatus(tx application.Tx, id uuid.UUID, status string, changedBy uuid.UUID) error {
	if !entities.IsValidClaimStatus(status) {
		return apperrors.NewInvalidClaimStatus()
	}

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entities.IsValidClaimStatusTransition(claim.Status, status) {
		return apperrors.NewInvalidClaimAction()
	}

	err = s.claimRepo.UpdateStatus(tx, id, status)
	if err != nil {
		return err
	}

	history := entities.NewClaimHistory(claim.ID, status, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return nil
}

func (s *claimService) Submit(tx application.Tx, id uuid.UUID, changedBy uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entities.IsValidClaimStatusTransition(claim.Status, entities.ClaimStatusSubmitted) {
		return apperrors.NewInvalidClaimAction()
	}

	items, err := s.itemRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}
	attachments, err := s.attachmentRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if len(items) < entities.ClaimItemRequirePerClaim || len(attachments) < entities.AttachmentRequirePerClaim {
		return apperrors.NewMissingInformationClaim()
	}

	err = s.claimRepo.UpdateStatus(tx, id, entities.ClaimStatusSubmitted)
	if err != nil {
		return err
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusSubmitted, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return nil
}

func (s *claimService) Complete(tx application.Tx, id uuid.UUID, changedBy uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	items, err := s.itemRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	approvedCount := 0
	for _, item := range items {
		switch item.Status {
		case entities.ClaimItemStatusApproved:
			approvedCount++
		case entities.ClaimItemStatusRejected:
		default:
			return apperrors.NewInvalidClaimAction()
		}
	}

	newStatus := entities.ClaimStatusPartiallyApproved
	if approvedCount == len(items) {
		newStatus = entities.ClaimStatusApproved
	} else if approvedCount == 0 {
		newStatus = entities.ClaimStatusRejected
	}

	if !entities.IsValidClaimStatusTransition(claim.Status, newStatus) {
		return apperrors.NewInvalidClaimAction()
	}

	err = s.claimRepo.UpdateStatus(tx, id, newStatus)
	if err != nil {
		return err
	}

	history := entities.NewClaimHistory(claim.ID, newStatus, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return nil
}

func (s *claimService) GetHistory(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error) {
	histories, err := s.historyRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return histories, nil
}
