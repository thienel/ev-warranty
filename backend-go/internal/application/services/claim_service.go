package services

import (
	"context"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type ClaimService interface {
	Create(ctx context.Context, cmd *CreateClaimCommand) (*entities.Claim, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error)
	GetAll(ctx context.Context, filters ClaimFilters, pagination Pagination) (*ClaimListResult, error)
	Update(ctx context.Context, id uuid.UUID, cmd *UpdateClaimCommand) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, cmd *UpdateClaimStatusCommand) error
}
