package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OfficeTypeEVM = "evm"
	OfficeTypeSC  = "sc"
)

type Office struct {
	ID         uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OfficeName string          `gorm:"not null;"`
	OfficeType string          `gorm:"not null;size:10"`
	Address    string          `gorm:"not null"`
	IsActive   bool            `gorm:"not null;"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
	DeletedAt  *gorm.DeletedAt `gorm:"index"`
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
