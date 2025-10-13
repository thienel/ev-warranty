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

type CreateClaimItemCommand struct {
	PartCategoryID    int
	FaultyPartID      uuid.UUID
	ReplacementPartID *uuid.UUID
	IssueDescription  string
	Status            string
	Type              string
	Cost              float64
}

type UpdateClaimItemCommand struct {
	IssueDescription string
	Type             string
	Cost             float64
}

type UpdateClaimItemStatusCommand struct {
}

type ClaimItemService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ClaimItem, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimItem, error)

	Create(tx application.Tx, claimID uuid.UUID, cmd *CreateClaimItemCommand) (*entities.ClaimItem, error)
	Update(tx application.Tx, claimID, itemID uuid.UUID, cmd *UpdateClaimItemCommand) error
	HardDelete(tx application.Tx, claimID, itemID uuid.UUID) error

	Approve(tx application.Tx, claimID, itemID uuid.UUID) error
	Reject(tx application.Tx, claimID, itemID uuid.UUID) error
}

type claimItemService struct {
	log       logger.Logger
	claimRepo repositories.ClaimRepository
	itemRepo  repositories.ClaimItemRepository
}

func NewClaimItemService(log logger.Logger, claimRepo repositories.ClaimRepository, itemRepo repositories.ClaimItemRepository) ClaimItemService {
	return &claimItemService{
		log:       log,
		claimRepo: claimRepo,
		itemRepo:  itemRepo,
	}
}

func (s *claimItemService) GetByID(ctx context.Context, id uuid.UUID) (*entities.ClaimItem, error) {
	item, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *claimItemService) GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entities.ClaimItem, error) {
	items, err := s.itemRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *claimItemService) Create(tx application.Tx, claimID uuid.UUID,
	cmd *CreateClaimItemCommand) (*entities.ClaimItem, error) {
	_, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return nil, err
	}

	if !entities.IsValidClaimItemStatus(cmd.Status) {
		return nil, apperrors.NewInvalidCredentials()
	}
	if !entities.IsValidClaimItemType(cmd.Type) {
		return nil, apperrors.NewInvalidCredentials()
	}

	item := entities.NewClaimItem(claimID, cmd.PartCategoryID, cmd.FaultyPartID, cmd.ReplacementPartID,
		cmd.IssueDescription, cmd.Status, cmd.Type, cmd.Cost)
	err = s.itemRepo.Create(tx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *claimItemService) Update(tx application.Tx, claimID, itemID uuid.UUID, cmd *UpdateClaimItemCommand) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return err
	}

	switch claim.Status {
	case entities.ClaimStatusDraft, entities.ClaimStatusRequestInfo:
	default:
		return apperrors.NewNotAllowUpdateClaim()
	}

	item, err := s.itemRepo.FindByID(tx.GetCtx(), itemID)
	if err != nil {
		return err
	}

	switch item.Status {
	case entities.ClaimItemStatusApproved, entities.ClaimItemStatusRejected:
	default:
		return apperrors.NewNotAllowUpdateClaim()
	}

	if !entities.IsValidClaimItemType(cmd.Type) {
		return apperrors.NewInvalidCredentials()
	}
	item.IssueDescription = cmd.IssueDescription
	item.Type = cmd.Type
	item.Cost = cmd.Cost

	err = s.itemRepo.Update(tx, item)
	if err != nil {
		return err
	}

	claim.TotalCost, err = s.itemRepo.SumCostByClaimID(tx, claimID)
	if err != nil {
		return err
	}

	err = s.claimRepo.Update(tx, claim)
	if err != nil {
		return err
	}

	return nil
}

func (s *claimItemService) HardDelete(tx application.Tx, claimID, itemID uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusDraft {
		return apperrors.NewNotAllowDeleteClaim()
	}

	err = s.itemRepo.HardDelete(tx, itemID)
	if err != nil {
		return err
	}

	claim.TotalCost, err = s.itemRepo.SumCostByClaimID(tx, claimID)
	if err != nil {
		return err
	}

	err = s.claimRepo.Update(tx, claim)
	if err != nil {
		return err
	}

	return nil
}

func (s *claimItemService) Approve(tx application.Tx, claimID, itemID uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusReviewing {
		return apperrors.NewNotAllowUpdateClaim()
	}

	err = s.itemRepo.UpdateStatus(tx, itemID, entities.ClaimStatusApproved)
	if err != nil {
		return err
	}

	claim.TotalCost, err = s.itemRepo.SumCostByClaimID(tx, claimID)
	if err != nil {
		return err
	}

	err = s.claimRepo.Update(tx, claim)
	if err != nil {
		return err
	}

	return nil
}

func (s *claimItemService) Reject(tx application.Tx, claimID, itemID uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return err
	}

	if claim.Status != entities.ClaimStatusReviewing {
		return apperrors.NewNotAllowUpdateClaim()
	}

	err = s.itemRepo.UpdateStatus(tx, itemID, entities.ClaimStatusRejected)
	if err != nil {
		return err
	}

	claim.TotalCost, err = s.itemRepo.SumCostByClaimID(tx, claimID)
	if err != nil {
		return err
	}

	err = s.claimRepo.Update(tx, claim)
	if err != nil {
		return err
	}

	return nil
}
