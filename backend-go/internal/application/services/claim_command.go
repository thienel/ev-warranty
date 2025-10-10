package services

import (
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type CreateClaimCommand struct {
	VehicleID   uuid.UUID
	CustomerID  uuid.UUID
	Description string
}

type UpdateClaimCommand struct {
	Description string
}

type UpdateClaimStatusCommand struct {
	Status    string
	ChangedBy uuid.UUID
}

type AddClaimItemCommand struct {
	PartCategoryID    int
	FaultyPartID      uuid.UUID
	ReplacementPartID *uuid.UUID
	IssueDescription  string
	Type              string
	Cost              float64
}

type UpdateClaimItemCommand struct {
	IssueDescription  string
	ReplacementPartID *uuid.UUID
	Cost              float64
}

type UpdateClaimItemStatusCommand struct {
	Status    string
	ChangedBy uuid.UUID
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
