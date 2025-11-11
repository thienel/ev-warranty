package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ClaimStatusDraft             = "DRAFT"
	ClaimStatusSubmitted         = "SUBMITTED"
	ClaimStatusReviewing         = "REVIEWING"
	ClaimStatusApproved          = "APPROVED"
	ClaimStatusPartiallyApproved = "PARTIALLY_APPROVED"
	ClaimStatusRejected          = "REJECTED"
	ClaimStatusCancelled         = "CANCELLED"
	ClaimStatusCompleted         = "COMPLETED"
)

const (
	MinItemPerClaim        = 1
	MinAttachmentPerClaim  = 2
	MaxClaimsPerTechnician = 3
)

type Claim struct {
	ID           uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerID   uuid.UUID       `gorm:"not null;type:uuid" json:"customer_id"`
	VehicleID    uuid.UUID       `gorm:"not null;type:uuid" json:"vehicle_id"`
	Kilometers   int             `gorm:"not null;" json:"kilometers"`
	Description  string          `gorm:"not null;" json:"description"`
	Status       string          `gorm:"not null;default:DRAFT" json:"status"`
	TotalCost    float64         `json:"total_cost"`
	StaffID      uuid.UUID       `gorm:"type:uuid" json:"staff_id"`
	TechnicianID uuid.UUID       `gorm:"type:uuid" json:"technician_id"`
	ApprovedBy   *uuid.UUID      `gorm:"type:uuid" json:"approved_by,omitempty"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewClaim(vehicleID, customerID uuid.UUID, kilometers int, description string, staffID, technicianID uuid.UUID) *Claim {

	return &Claim{
		ID:           uuid.New(),
		CustomerID:   customerID,
		VehicleID:    vehicleID,
		Kilometers:   kilometers,
		Description:  description,
		StaffID:      staffID,
		TechnicianID: technicianID,
	}
}

func IsValidClaimStatus(status string) bool {
	switch status {
	case ClaimStatusDraft, ClaimStatusSubmitted, ClaimStatusApproved, ClaimStatusPartiallyApproved,
		ClaimStatusCancelled, ClaimStatusReviewing, ClaimStatusRejected, ClaimStatusCompleted:
		return true
	default:
		return false
	}
}

func IsValidClaimStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		ClaimStatusDraft: {
			ClaimStatusSubmitted,
		},
		ClaimStatusSubmitted: {
			ClaimStatusReviewing,
			ClaimStatusCancelled,
		},
		ClaimStatusReviewing: {
			ClaimStatusApproved,
			ClaimStatusPartiallyApproved,
			ClaimStatusRejected,
		},
		ClaimStatusApproved:          {ClaimStatusCompleted},
		ClaimStatusPartiallyApproved: {ClaimStatusCompleted},
		ClaimStatusRejected:          {},
		ClaimStatusCancelled:         {},
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
