package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClaimAttachment struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ClaimID        uuid.UUID       `gorm:"not null;type:uuid" json:"claimID"`
	Claim          Claim           `gorm:"foreignKey:ClaimID;references:ID;constrains:OnDelete:CASCADE"`
	AttachmentType string          `gorm:"not null" json:"attachment_type"`
	URL            string          `gorm:"not null;type:text" json:"url"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"-"`
}
