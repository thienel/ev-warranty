package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ClaimStatusDraft             = "DRAFT"
	ClaimStatusSubmitted         = "SUBMITTED"
	ClaimStatusReviewing         = "REVIEWING"
	ClaimStatusRequestInfo       = "REQUEST_INFO"
	ClaimStatusApproved          = "APPROVED"
	ClaimStatusPartiallyApproved = "PARTIALLY_APPROVED"
	ClaimStatusRejected          = "REJECTED"
	ClaimStatusCancelled         = "CANCELLED"
)

type Claim struct {
	ID          uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	VehicleID   uuid.UUID       `gorm:"not null;type:uuid" json:"vehicle_id"`
	CustomerID  uuid.UUID       `gorm:"not null;type:uuid" json:"customer_id"`
	Description string          `gorm:"not null;" json:"description"`
	Status      string          `gorm:"not null;default:draft" json:"status"`
	TotalCost   float64         `json:"total_cost"`
	ApprovedBy  uuid.UUID       `gorm:"type:uuid" json:"approved_by"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewClaim(vehicleID, customerID uuid.UUID, description, status string, approvedBy uuid.UUID) *Claim {
	return &Claim{
		ID:          uuid.New(),
		VehicleID:   vehicleID,
		CustomerID:  customerID,
		Description: description,
		Status:      status,
		ApprovedBy:  approvedBy,
	}
}

func IsValidClaimStatus(status string) bool {
	switch status {
	case ClaimStatusDraft, ClaimStatusSubmitted, ClaimStatusApproved, ClaimStatusPartiallyApproved,
		ClaimStatusCancelled, ClaimStatusReviewing, ClaimStatusRequestInfo, ClaimStatusRejected:
		return true
	default:
		return false
	}
}

func IsValidClaimStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		ClaimStatusDraft: {
			ClaimStatusSubmitted,
			ClaimStatusCancelled,
		},
		ClaimStatusSubmitted: {
			ClaimStatusReviewing,
			ClaimStatusCancelled,
		},
		ClaimStatusReviewing: {
			ClaimStatusRequestInfo,
			ClaimStatusApproved,
			ClaimStatusPartiallyApproved,
			ClaimStatusRejected,
		},
		ClaimStatusRequestInfo: {
			ClaimStatusSubmitted,
			ClaimStatusCancelled,
		},
		ClaimStatusApproved:          {},
		ClaimStatusPartiallyApproved: {},
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
