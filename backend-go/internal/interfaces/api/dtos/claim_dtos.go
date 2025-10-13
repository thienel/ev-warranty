package dtos

import (
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type CreateClaimRequest struct {
	VehicleID   uuid.UUID `json:"vehicle_id" binding:"required"`
	CustomerID  uuid.UUID `json:"customer_id" binding:"required"`
	Description string    `json:"description" binding:"required,min=10,max=1000"`
}

type UpdateClaimRequest struct {
	Description string `json:"description" binding:"required,min=10,max=1000"`
}

type ClaimResponse struct {
	ID          uuid.UUID `json:"id"`
	VehicleID   uuid.UUID `json:"vehicle_id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TotalCost   float64   `json:"total_cost"`
	ApprovedBy  uuid.UUID `json:"approved_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClaimListResponse struct {
	Claims     []ClaimResponse `json:"claims"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

type ClaimHistoryResponse struct {
	ID        uuid.UUID `json:"id"`
	ClaimID   uuid.UUID `json:"claim_id"`
	Status    string    `json:"status"`
	ChangedBy uuid.UUID `json:"changed_by"`
	ChangedAt time.Time `json:"changed_at"`
}

type ClaimHistoryListResponse struct {
	Histories []ClaimHistoryResponse `json:"histories"`
}

func ToClaimResponse(claim *entities.Claim) *ClaimResponse {
	return &ClaimResponse{
		ID:          claim.ID,
		VehicleID:   claim.VehicleID,
		CustomerID:  claim.CustomerID,
		Description: claim.Description,
		Status:      claim.Status,
		TotalCost:   claim.TotalCost,
		ApprovedBy:  claim.ApprovedBy,
		CreatedAt:   claim.CreatedAt,
		UpdatedAt:   claim.UpdatedAt,
	}
}

func ToClaimListResponse(result *services.ClaimListResult) *ClaimListResponse {
	claims := make([]ClaimResponse, len(result.Claims))
	for i, claim := range result.Claims {
		claims[i] = *ToClaimResponse(claim)
	}

	return &ClaimListResponse{
		Claims:     claims,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

func ToClaimHistoryResponse(history *entities.ClaimHistory) ClaimHistoryResponse {
	return ClaimHistoryResponse{
		ID:        history.ID,
		ClaimID:   history.ClaimID,
		Status:    history.Status,
		ChangedBy: history.ChangedBy,
		ChangedAt: history.ChangedAt,
	}
}

func ToClaimHistoryListResponse(histories []*entities.ClaimHistory) *ClaimHistoryListResponse {
	historyResponses := make([]ClaimHistoryResponse, len(histories))
	for i, history := range histories {
		historyResponses[i] = ToClaimHistoryResponse(history)
	}

	return &ClaimHistoryListResponse{
		Histories: historyResponses,
	}
}
