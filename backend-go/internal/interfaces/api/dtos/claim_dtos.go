package dtos

import (
	"ev-warranty-go/internal/domain/entities"

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

type CreateClaimItemRequest struct {
	PartCategoryID    int        `json:"part_category_id" binding:"required"`
	FaultyPartID      uuid.UUID  `json:"faulty_part_id" binding:"required"`
	ReplacementPartID *uuid.UUID `json:"replacement_part_id"`
	IssueDescription  string     `json:"issue_description" binding:"required,min=10,max=1000"`
	Type              string     `json:"type" binding:"required"`
	Cost              float64    `json:"cost" binding:"required,min=0"`
}

type ClaimItemListResponse struct {
	Items []entities.ClaimItem `json:"items"`
	Total int                  `json:"total"`
}

type ClaimAttachmentListResponse struct {
	Attachments []entities.ClaimAttachment `json:"attachments"`
	Total       int                        `json:"total"`
}
