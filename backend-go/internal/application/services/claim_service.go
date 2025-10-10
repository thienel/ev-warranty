package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/logger"
	"log"
	"time"

	"github.com/google/uuid"
)

type ClaimService interface {
	Create(tx application.Transaction, cmd *CreateClaimCommand) (*entities.Claim, error)
	Update(tx application.Transaction, id uuid.UUID, cmd *UpdateClaimCommand) error
	Delete(tx application.Transaction, id uuid.UUID) error
	UpdateStatus(tx application.Transaction, id uuid.UUID, cmd *UpdateClaimStatusCommand) error

	GetByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error)
	GetAll(ctx context.Context, filters ClaimFilters, pagination Pagination) (*ClaimListResult, error)
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

func (s *claimService) Create(tx application.Transaction, cmd *CreateClaimCommand) (*entities.Claim, error) {
	claim := entities.NewClaim(
		cmd.VehicleID,
		cmd.CustomerID,
		cmd.Description,
		entities.ClaimStatusDraft,
		uuid.Nil,
	)

	if err := s.claimRepo.Create(tx, claim); err != nil {
		return nil, rollbackOnErr(tx, err)
	}

	history := entities.NewClaimHistory(claim.ID, entities.ClaimStatusDraft, cmd.CreatorID)
	if err := s.historyRepo.Create(tx, history); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return claim, nil
}

func (s *claimService) Update(tx application.Transaction, id uuid.UUID, cmd *UpdateClaimCommand) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft && claim.Status != entities.ClaimStatusRequestInfo {
		return apperrors.ErrInvalidCredentials("can only update draft or request info claims")
	}

	claim.Description = cmd.Description

	if err = s.claimRepo.Update(tx, claim); err != nil {
		return rollbackOnErr(tx, err)
	}

	return tx.Commit()
}

func (s *claimService) Delete(tx application.Transaction, id uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft && claim.Status != entities.ClaimStatusCancelled {
		return apperrors.ErrInvalidCredentials("can only delete draft or cancelled claims")
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
				return rollbackOnErr(tx, err)
			}
		}
		return tx.Commit()
	}

	if err = s.claimRepo.HardDelete(tx, id); err != nil {
		return rollbackOnErr(tx, err)
	}

	return tx.Commit()
}

func (s *claimService) UpdateStatus(tx application.Transaction,
	id uuid.UUID, cmd *UpdateClaimStatusCommand) error {

	claim, err := s.claimRepo.FindByID(tx.GetCtx(), id)
	if err != nil {
		return err
	}

	if !entities.IsValidClaimStatus(cmd.Status) {
		return apperrors.ErrInvalidCredentials("invalid claim status")
	}

	if !s.isValidStatusTransition(claim.Status, cmd.Status) {
		return apperrors.ErrInvalidCredentials("invalid status transition")
	}

	if err = s.claimRepo.UpdateStatus(tx, id, cmd.Status); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	history := &entities.ClaimHistory{
		ID:        uuid.New(),
		ClaimID:   id,
		Status:    cmd.Status,
		ChangedBy: cmd.ChangedBy,
		ChangedAt: time.Now(),
	}
	if err = s.historyRepo.Create(tx, history); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if cmd.Status == entities.ClaimStatusApproved || cmd.Status == entities.ClaimStatusPartiallyApproved {
		claim.ApprovedBy = cmd.ChangedBy
		claim.Status = cmd.Status
		if err = s.claimRepo.Update(tx, claim); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
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

func (s *claimService) isValidStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		entities.ClaimStatusDraft: {
			entities.ClaimStatusSubmitted,
			entities.ClaimStatusCancelled,
		},
		entities.ClaimStatusSubmitted: {
			entities.ClaimStatusReviewing,
			entities.ClaimStatusCancelled,
		},
		entities.ClaimStatusReviewing: {
			entities.ClaimStatusRequestInfo,
			entities.ClaimStatusApproved,
			entities.ClaimStatusPartiallyApproved,
			entities.ClaimStatusRejected,
		},
		entities.ClaimStatusRequestInfo: {
			entities.ClaimStatusSubmitted,
			entities.ClaimStatusCancelled,
		},
		entities.ClaimStatusApproved:          {},
		entities.ClaimStatusPartiallyApproved: {},
		entities.ClaimStatusRejected:          {},
		entities.ClaimStatusCancelled:         {},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

func rollbackOnErr(tx application.Transaction, originalErr error) error {
	if err := tx.Rollback(); err != nil {
		log.Printf("[TX ROLLBACK FAILED] original error: %v, rollback error: %v", originalErr, err)
	}
	return originalErr
}

func commitLog(tx application.Transaction) error {
	if err := tx.Rollback(); err != nil {
		log.Printf("[TX COMMIT FAILED] commit error: %v", err)
		return apperrors.NewInternal("transaction commit error", err)
	}
	return nil
}
