package entities

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	ID               uuid.UUID `json:"id"`
	SerialNumber     string    `json:"serial_number"`
	PartName         string    `json:"part_name"`
	UnitPrice        float64   `json:"unit_price"`
	CategoryID       uuid.UUID `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	OfficeLocationID uuid.UUID `json:"office_location_id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
