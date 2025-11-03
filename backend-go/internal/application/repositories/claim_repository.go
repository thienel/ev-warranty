package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type ClaimRepository interface {
	Create(tx application.Tx, claim *entity.Claim) error
	Update(tx application.Tx, claim *entity.Claim) error
	UpdateStatus(tx application.Tx, id uuid.UUID, status string) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDelete(tx application.Tx, id uuid.UUID) error

	FindByID(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
	FindAll(ctx context.Context) ([]*entity.Claim, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entity.Claim, error)
	FindByVehicleID(ctx context.Context, vehicleID uuid.UUID) ([]*entity.Claim, error)
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
