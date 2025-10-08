package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ClaimItemStatusPending    = "PENDING"
	ClaimItemStatusRejected   = "REJECTED"
	ClaimItemStatusApproved   = "APPROVED"
	ClaimItemStatusInProgress = "IN_PROGRESS"
	ClaimItemStatusCompleted  = "COMPLETED"

	ClaimItemTypeReplacement = "REPLACEMENT"
	ClaimItemTypeRepair      = "REPAIR"
)

type ClaimItem struct {
	ID                uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ClaimID           uuid.UUID       `gorm:"not null;type:uuid" json:"claim_id"`
	Claim             Claim           `gorm:"foreignKey:ClaimID;references:ID;constrains:OnDelete:CASCADE" json:"-"`
	PartCategoryID    int             `gorm:"not null" json:"part_category_id"`
	FaultyPartID      uuid.UUID       `gorm:"not null;type:uuid" json:"faulty_part_id"`
	ReplacementPartID *uuid.UUID      `gorm:"type:uuid" json:"replacement_part_id"`
	IssueDescription  string          `gorm:"not null;type:text" json:"issue_description"`
	LineStatus        string          `gorm:"not null" json:"line_status"`
	Type              string          `gorm:"not null" json:"type"`
	Cost              float64         `json:"cost"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         *gorm.DeletedAt `json:"-"`
}

func NewClaimItem(claimID uuid.UUID, partCategoryID int, faultyPartID uuid.UUID, replacementPartID *uuid.UUID,
	IssueDescription, lineStatus, claimItemType string, cost float64) *ClaimItem {

	return &ClaimItem{
		ID:                uuid.New(),
		ClaimID:           claimID,
		PartCategoryID:    partCategoryID,
		FaultyPartID:      faultyPartID,
		ReplacementPartID: replacementPartID,
		IssueDescription:  IssueDescription,
		LineStatus:        lineStatus,
		Type:              claimItemType,
		Cost:              cost,
	}
}

func IsValidClaimItemStatus(status string) bool {
	switch status {
	case ClaimItemStatusPending, ClaimItemStatusRejected, ClaimItemStatusApproved, ClaimItemStatusInProgress, ClaimItemStatusCompleted:
		return true
	default:
		return false
	}
}

func IsValidClaimItemType(claimItemType string) bool {
	switch claimItemType {
	case ClaimItemTypeReplacement, ClaimItemTypeRepair:
		return true
	default:
		return false
	}
}
