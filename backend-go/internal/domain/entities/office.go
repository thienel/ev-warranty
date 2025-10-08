package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OfficeTypeEVM = "EVM"
	OfficeTypeSC  = "SC"
)

type Office struct {
	ID         uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OfficeName string          `gorm:"not null;" json:"office_name"`
	OfficeType string          `gorm:"not null;" json:"office_type"`
	Address    string          `gorm:"not null" json:"address"`
	IsActive   bool            `gorm:"not null;" json:"is_active"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  *gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewOffice(officeName string, officeType string, address string, isActive bool) *Office {
	return &Office{
		ID:         uuid.New(),
		OfficeName: officeName,
		OfficeType: officeType,
		Address:    address,
		IsActive:   isActive,
	}
}

func (o *Office) Active() {
	o.IsActive = true
}

func (o *Office) Inactive() {
	o.IsActive = false
}

func IsValidOfficeType(officeType string) bool {
	switch officeType {
	case OfficeTypeEVM, OfficeTypeSC:
		return true
	default:
		return false
	}
}
