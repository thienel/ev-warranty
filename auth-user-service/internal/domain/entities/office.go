package entities

import (
	"auth-service/internal/errors/apperrors"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type OfficeType string

const (
	OfficeTypeEVM OfficeType = "EVM"
	OfficeTypeSC  OfficeType = "SC"
)

type Office struct {
	id         uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	officeName string     `gorm:"type:varchar(255);not null"`
	officeType OfficeType `gorm:"type:varchar(255);not null"`
	address    string     `gorm:"type:varchar(255);not null"`
	isActive   bool       `gorm:"type:boolean;default:true"`
	createdAt  time.Time  `gorm:"autoCreateTime"`
	updatedAt  time.Time  `gorm:"autoUpdateTime"`
}

func NewOffice(officeName string, officeType OfficeType, address string, isActive bool) *Office {
	return &Office{
		id:         uuid.New(),
		officeName: officeName,
		officeType: officeType,
		address:    address,
		isActive:   isActive,
	}
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
