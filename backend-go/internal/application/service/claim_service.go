package service

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/infrastructure/cloudinary"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/logger"
	"fmt"

	"github.com/google/uuid"
)

type CreateClaimCommand struct {
	VehicleID    uuid.UUID
	CustomerID   uuid.UUID
	StaffID      uuid.UUID
	TechnicianID uuid.UUID
	OfficeID     uuid.UUID
	Description  string
}

type UpdateClaimCommand struct {
	Description string
}

type ClaimService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
	GetAll(ctx context.Context) ([]*entity.Claim, error)

	Create(tx application.Tx, cmd *CreateClaimCommand) (*entity.Claim, error)
	Update(tx application.Tx, id uuid.UUID, cmd *UpdateClaimCommand) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDelete(tx application.Tx, id uuid.UUID) error

	UpdateStatus(tx application.Tx, id uuid.UUID, status string, changedBy uuid.UUID) error
	Submit(tx application.Tx, id uuid.UUID, changedBy uuid.UUID) error
	Complete(tx application.Tx, id uuid.UUID, changedBy uuid.UUID) error

	GetHistory(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimHistory, error)
}

type claimService struct {
	log            logger.Logger
	claimRepo      repository.ClaimRepository
	userRepo       repository.UserRepository
	itemRepo       repository.ClaimItemRepository
	attachmentRepo repository.ClaimAttachmentRepository
	historyRepo    repository.ClaimHistoryRepository
	cloudService   cloudinary.CloudinaryService
}

func NewClaimService(
	log logger.Logger,
	claimRepo repository.ClaimRepository,
	userRepo repository.UserRepository,
	itemRepo repository.ClaimItemRepository,
	attachmentRepo repository.ClaimAttachmentRepository,
	historyRepo repository.ClaimHistoryRepository,
	cloudService cloudinary.CloudinaryService,
) ClaimService {
	return &claimService{
		log:            log,
		claimRepo:      claimRepo,
		userRepo:       userRepo,
		itemRepo:       itemRepo,
		attachmentRepo: attachmentRepo,
		historyRepo:    historyRepo,
		cloudService:   cloudService,
	}
}

func (s *claimService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Claim, error) {
	return s.claimRepo.FindByID(ctx, id)
}

func (s *claimService) GetAll(ctx context.Context) ([]*entity.Claim, error) {
	claims, err := s.claimRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return claims, err
}

func (s *claimService) Create(tx application.Tx, cmd *CreateClaimCommand) (*entity.Claim, error) {
	staff, err := s.userRepo.FindByID(tx.GetCtx(), cmd.StaffID)
	if err != nil {
		return nil, err
	}

	technician, err := s.userRepo.FindByID(tx.GetCtx(), cmd.TechnicianID)
	if err != nil {
		return nil, err
	}

	if staff.OfficeID != technician.OfficeID {
		return nil, apperror.ErrInvalidTechnician
	}

	count, err := s.claimRepo.CountPendingByTechnician(tx.GetCtx(), cmd.TechnicianID)
	if err != nil {
		return nil, err
	}
	if count >= entity.MaxClaimsPerTechnician {
		return nil, apperror.ErrTechnicianWorkloadExceed
	}

	claim := entity.NewClaim(cmd.VehicleID, cmd.CustomerID, cmd.Description, entity.ClaimStatusDraft,
		cmd.StaffID, cmd.TechnicianID)

	if err := s.claimRepo.Create(tx, claim); err != nil {
		return nil, err
	}

	history := entity.NewClaimHistory(claim.ID, entity.ClaimStatusDraft, cmd.StaffID)

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

	if claim.Status != entity.ClaimStatusDraft {
		return apperror.ErrInvalidClaimAction.WithMessage("Can only update when status if draft")
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

	if claim.Status != entity.ClaimStatusDraft {
		return apperror.ErrInvalidClaimAction.WithMessage("Can only hard delete when status if draft")
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

	if claim.Status != entity.ClaimStatusCancelled {
		return apperror.ErrInvalidClaimAction.WithMessage("Can only soft delete when status is cancelled")
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

func (s *claimService) UpdateStatus(tx application.Tx, id uuid.UUID, status string, changedBy uuid.UUID,
) error {
	if !entity.IsValidClaimStatus(status) {
		return apperror.ErrInvalidInput.WithMessage("Invalid claim status")
	}

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entity.IsValidClaimStatusTransition(claim.Status, status) {
		return apperror.ErrInvalidClaimAction.WithMessage("This action are not allowed")
	}

	err = s.claimRepo.UpdateStatus(tx, id, status)
	if err != nil {
		return err
	}

	history := entity.NewClaimHistory(claim.ID, status, changedBy)
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

	if !entity.IsValidClaimStatusTransition(claim.Status, entity.ClaimStatusSubmitted) {
		return apperror.ErrInvalidClaimAction.WithMessage("Invalid claim action")
	}

	items, err := s.itemRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}
	attachments, err := s.attachmentRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if len(items) < entity.MinItemPerClaim {
		return apperror.ErrMissingInformationClaim.
			WithMessage(fmt.Sprintf("Claim item must be atleast %d", entity.MinItemPerClaim))
	}
	if len(attachments) < entity.MinAttachmentPerClaim {
		return apperror.ErrMissingInformationClaim.
			WithMessage(fmt.Sprintf("Claim attachment must be atleast %d", entity.MinAttachmentPerClaim))
	}

	err = s.claimRepo.UpdateStatus(tx, id, entity.ClaimStatusSubmitted)
	if err != nil {
		return err
	}

	history := entity.NewClaimHistory(claim.ID, entity.ClaimStatusSubmitted, changedBy)
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
		case entity.ClaimItemStatusApproved:
			approvedCount++
		case entity.ClaimItemStatusRejected:
		default:
			return apperror.ErrInvalidInput.WithMessage("Can only complete when all item are approved or rejected")
		}
	}

	newStatus := entity.ClaimStatusPartiallyApproved
	if approvedCount == len(items) {
		newStatus = entity.ClaimStatusApproved
	} else if approvedCount == 0 {
		newStatus = entity.ClaimStatusRejected
	}

	if !entity.IsValidClaimStatusTransition(claim.Status, newStatus) {
		return apperror.ErrInvalidInput.WithMessage("This action are not allowed")
	}

	err = s.claimRepo.UpdateStatus(tx, id, newStatus)
	if err != nil {
		return err
	}

	history := entity.NewClaimHistory(claim.ID, newStatus, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return nil
}

func (s *claimService) GetHistory(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimHistory, error) {
	histories, err := s.historyRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return histories, nil
}
