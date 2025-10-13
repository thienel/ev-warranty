package repositories

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type ClaimRepository interface {
	Create(tx application.Tx, claim *entities.Claim) error
	Update(tx application.Tx, claim *entities.Claim) error
	UpdateStatus(tx application.Tx, id uuid.UUID, status string) error
	HardDelete(tx application.Tx, id uuid.UUID) error
	SoftDelete(tx application.Tx, id uuid.UUID) error

	FindByID(ctx context.Context, id uuid.UUID) (*entities.Claim, error)
	FindAll(ctx context.Context, filters ClaimFilters, pagination Pagination) ([]*entities.Claim, int64, error)
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
