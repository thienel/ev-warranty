package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
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
	Claims     []*entities.Claim
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type ClaimService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error)
	GetAll(ctx context.Context, filters ClaimFilters, pagination Pagination) (*ClaimListResult, error)

	Create(tx application.Transaction, cmd *CreateClaimCommand) (*entities.Claim, error)
	Update(tx application.Transaction, id uuid.UUID, cmd *UpdateClaimCommand) error
	Delete(tx application.Transaction, id uuid.UUID) error

	UpdateStatus(tx application.Transaction, id uuid.UUID, status string, changedBy uuid.UUID) error
	Submit(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error
	Approve(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error
	Reject(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error
	PartiallyApprove(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error

	GetHistory(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error)
}

type claimService struct {
	log            logger.Logger
	claimRepo      repositories.ClaimRepository
	itemRepo       repositories.ClaimItemRepository
	attachmentRepo repositories.ClaimAttachmentRepository
	historyRepo    repositories.ClaimHistoryRepository
}

func NewClaimService(
	log logger.Logger,
	claimRepo repositories.ClaimRepository,
	itemRepo repositories.ClaimItemRepository,
	attachmentRepo repositories.ClaimAttachmentRepository,
	historyRepo repositories.ClaimHistoryRepository,
) ClaimService {
	return &claimService{
		log:            log,
		claimRepo:      claimRepo,
		itemRepo:       itemRepo,
		attachmentRepo: attachmentRepo,
		historyRepo:    historyRepo,
	}
}

func (s *claimService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error) {
	return s.claimRepo.FindByID(ctx, id)
}

func (s *claimService) GetAll(ctx context.Context, filters ClaimFilters, pagination Pagination) (*ClaimListResult, error) {
	repoFilters := repositories.ClaimFilters{
		CustomerID: filters.CustomerID,
		VehicleID:  filters.VehicleID,
		Status:     filters.Status,
		FromDate:   filters.FromDate,
		ToDate:     filters.ToDate,
	}

	repoPagination := repositories.Pagination{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	claims, total, err := s.claimRepo.FindAll(ctx, repoFilters, repoPagination)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if pagination.PageSize > 0 {
		totalPages = int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))
	}

	return &ClaimListResult{
		Claims:     claims,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *claimService) Create(tx application.Transaction, cmd *CreateClaimCommand) (*entities.Claim, error) {
	claim := entities.NewClaim(cmd.VehicleID, cmd.CustomerID, cmd.Description, entities.ClaimStatusDraft, uuid.Nil)
	defer rollbackOrLog(tx)

	if err := s.claimRepo.Create(tx, claim); err != nil {
		return nil, err
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusDraft, cmd.CreatorID)
	if err := s.historyRepo.Create(tx, history); err != nil {
		return nil, err
	}

	return claim, commitOrLog(tx)
}

func (s *claimService) Update(tx application.Transaction, id uuid.UUID, cmd *UpdateClaimCommand) error {
	defer rollbackOrLog(tx)

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

	return commitOrLog(tx)
}

func (s *claimService) Delete(tx application.Transaction, id uuid.UUID) error {
	defer rollbackOrLog(tx)

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft && claim.Status != entities.ClaimStatusCancelled {
		return apperrors.NewNotAllowDeleteClaim()
	}

	if claim.Status == entities.ClaimStatusCancelled {
		softDeleters := []func(application.Transaction, uuid.UUID) error{
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
	} else if claim.Status == entities.ClaimStatusDraft {
		if err = s.claimRepo.HardDelete(tx, id); err != nil {
			return err
		}
	}

	return commitOrLog(tx)
}

func (s *claimService) UpdateStatus(tx application.Transaction, id uuid.UUID, status string, changedBy uuid.UUID) error {
	defer rollbackOrLog(tx)

	if !entities.IsValidClaimStatus(status) {
		return apperrors.NewInvalidCredentials()
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

	return commitOrLog(tx)
}

func (s *claimService) Submit(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error {
	defer rollbackOrLog(tx)

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

	return commitOrLog(tx)
}

func (s *claimService) Approve(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error {
	defer rollbackOrLog(tx)

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entities.IsValidClaimStatusTransition(claim.Status, entities.ClaimStatusApproved) {
		return apperrors.NewInvalidClaimAction()
	}

	items, err := s.itemRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	for _, item := range items {
		err = s.itemRepo.UpdateStatus(tx, item.ID, entities.ClaimItemStatusApproved)
		if err != nil {
			return err
		}
	}

	err = s.claimRepo.UpdateStatus(tx, id, entities.ClaimStatusApproved)
	if err != nil {
		return err
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusApproved, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return commitOrLog(tx)
}

func (s *claimService) Reject(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error {
	defer rollbackOrLog(tx)

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entities.IsValidClaimStatusTransition(claim.Status, entities.ClaimStatusRejected) {
		return apperrors.NewInvalidClaimAction()
	}

	items, err := s.itemRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	for _, item := range items {
		err = s.itemRepo.UpdateStatus(tx, item.ID, entities.ClaimItemStatusRejected)
		if err != nil {
			return err
		}
	}

	err = s.claimRepo.UpdateStatus(tx, id, entities.ClaimStatusRejected)
	if err != nil {
		return err
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusRejected, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return commitOrLog(tx)
}

func (s *claimService) PartiallyApprove(tx application.Transaction, id uuid.UUID, changedBy uuid.UUID) error {
	defer rollbackOrLog(tx)

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entities.IsValidClaimStatusTransition(claim.Status, entities.ClaimStatusPartiallyApproved) {
		return apperrors.NewInvalidClaimAction()
	}

	items, err := s.itemRepo.FindByClaimID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Status == entities.ClaimItemStatusPending {
			return apperrors.NewInvalidClaimAction()
		}
	}

	err = s.claimRepo.UpdateStatus(tx, id, entities.ClaimStatusPartiallyApproved)
	if err != nil {
		return err
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusPartiallyApproved, changedBy)
	if err = s.historyRepo.Create(tx, history); err != nil {
		return err
	}

	return commitOrLog(tx)
}

func (s *claimService) GetHistory(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimHistory, error) {
	histories, err := s.historyRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return histories, nil
}
