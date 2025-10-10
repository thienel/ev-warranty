package repositories

import (
	"context"
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type ClaimRepository interface {
	Create(ctx context.Context, claim *entities.Claim) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error)
	FindAll(ctx context.Context, filters ClaimFilters, pagination Pagination) ([]*entities.Claim, int64, error)
	Update(ctx context.Context, claim *entities.Claim) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Claim, error)
	FindByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]*entities.Claim, error)
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
