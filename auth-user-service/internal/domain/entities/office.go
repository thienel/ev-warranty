package entities

import (
	"auth-service/internal/errors/apperrors"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfficeType string

const (
	OfficeTypeEVM OfficeType = "EVM"
	OfficeTypeSC  OfficeType = "SC"
)

type Office struct {
	Id         uuid.UUID       `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	OfficeName string          `gorm:"type:varchar(255);not null"`
	OfficeType OfficeType      `gorm:"not null"`
	Address    string          `gorm:"type:varchar(255);not null"`
	IsActive   bool            `gorm:"type:boolean;default:true"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
	DeletedAt  *gorm.DeletedAt `gorm:"index"`
}

func NewOffice(officeName string, officeType OfficeType, address string, isActive bool) *Office {
	return &Office{
		Id:         uuid.New(),
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

func (ot OfficeType) Value() (driver.Value, error) {
	return string(ot), nil
}

func (ot *OfficeType) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch s := value.(type) {
	case string:
		*ot = OfficeType(s)
		return nil
	case []byte:
		*ot = OfficeType(s)
		return nil
	default:
		return apperrors.ErrScanValue
	}
}

func (ot *OfficeType) IsValid() bool {
	switch *ot {
	case OfficeTypeEVM, OfficeTypeSC:
		return true
	default:
		return false
	}
}
