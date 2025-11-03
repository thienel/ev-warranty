package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	AttachmentTypeVideo = "video"
	AttachmentTypeImage = "image"

	AttachmentRequirePerClaim = 3
)

type ClaimAttachment struct {
	ID        uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ClaimID   uuid.UUID       `gorm:"not null;type:uuid" json:"claim_id"`
	Claim     Claim           `gorm:"foreignKey:ClaimID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Type      string          `gorm:"not null" json:"type"`
	URL       string          `gorm:"not null;type:text" json:"url"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewClaimAttachment(claimID uuid.UUID, attachmentType, url string) *ClaimAttachment {
	return &ClaimAttachment{
		ID:      uuid.New(),
		ClaimID: claimID,
		Type:    attachmentType,
		URL:     url,
	}
}

func IsValidAttachmentType(attachmentType string) bool {
	switch attachmentType {
	case AttachmentTypeVideo, AttachmentTypeImage:
		return true
	default:
		return false
	}
}
