package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClaimHistory struct {
	ID        uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ClaimID   uuid.UUID       `gorm:"not null;type:uuid" json:"claim_id"`
	Claim     Claim           `gorm:"foreignKey:ClaimID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Status    string          `gorm:"not null" json:"status"`
	ChangedBy uuid.UUID       `gorm:"not null;type:uuid" json:"changed_by"`
	ChangedAt time.Time       `gorm:"autoCreateTime"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`
}
