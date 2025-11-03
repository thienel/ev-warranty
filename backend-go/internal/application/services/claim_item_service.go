package services

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
)

type CreateClaimItemCommand struct {
	PartCategoryID    uuid.UUID
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
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimItem, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimItem, error)

	Create(tx application.Tx, claimID uuid.UUID, cmd *CreateClaimItemCommand) (*entity.ClaimItem, error)
	Update(tx application.Tx, claimID, itemID uuid.UUID, cmd *UpdateClaimItemCommand) error
	HardDelete(tx application.Tx, claimID, itemID uuid.UUID) error

	Approve(tx application.Tx, claimID, itemID uuid.UUID) error
	Reject(tx application.Tx, claimID, itemID uuid.UUID) error
}

type claimItemService struct {
	claimRepo repositories.ClaimRepository
	itemRepo  repositories.ClaimItemRepository
}

func NewClaimItemService(claimRepo repositories.ClaimRepository, itemRepo repositories.ClaimItemRepository) ClaimItemService {
	return &claimItemService{
		claimRepo: claimRepo,
		itemRepo:  itemRepo,
	}
}

func (s *claimItemService) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimItem, error) {
	item, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *claimItemService) GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimItem, error) {
	items, err := s.itemRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *claimItemService) Create(tx application.Tx, claimID uuid.UUID,
	cmd *CreateClaimItemCommand) (*entity.ClaimItem, error) {
	_, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return nil, err
	}

	if !entity.IsValidClaimItemStatus(cmd.Status) {
		return nil, apperror.NewInvalidClaimItemStatus()
	}
	if !entity.IsValidClaimItemType(cmd.Type) {
		return nil, apperror.NewInvalidClaimItemType()
	}

	item := entity.NewClaimItem(claimID, cmd.PartCategoryID, cmd.FaultyPartID, cmd.ReplacementPartID,
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
	case entity.ClaimStatusDraft, entity.ClaimStatusRequestInfo:
	default:
		return apperror.NewNotAllowUpdateClaim()
	}

	item, err := s.itemRepo.FindByID(tx.GetCtx(), itemID)
	if err != nil {
		return err
	}

	switch item.Status {
	case entity.ClaimItemStatusPending:
	default:
		return apperror.NewNotAllowUpdateClaim()
	}

	if !entity.IsValidClaimItemType(cmd.Type) {
		return apperror.NewInvalidClaimItemType()
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

	if claim.Status != entity.ClaimStatusDraft {
		return apperror.NewNotAllowDeleteClaim()
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

	if claim.Status != entity.ClaimStatusReviewing {
		return apperror.NewNotAllowUpdateClaim()
	}

	err = s.itemRepo.UpdateStatus(tx, itemID, entity.ClaimItemStatusApproved)
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

	if claim.Status != entity.ClaimStatusReviewing {
		return apperror.NewNotAllowUpdateClaim()
	}

	err = s.itemRepo.UpdateStatus(tx, itemID, entity.ClaimItemStatusRejected)
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
