package dotnet

import (
	"time"

	"github.com/google/uuid"
)

type BaseResponse struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	ErrorCode string `json:"error,omitempty"`
}

type BaseDataResponse[T any] struct {
	BaseResponse
	Data *T `json:"data,omitempty"`
}

type ReservePartRequest struct {
	OfficeLocationID uuid.UUID `json:"office_location_id"`
	CategoryID       uuid.UUID `json:"category_id"`
}

type PartResponse struct {
	ID                   uuid.UUID  `json:"id"`
	SerialNumber         string     `json:"serial_number"`
	PartName             string     `json:"part_name"`
	UnitPrice            float64    `json:"unit_price"`
	CategoryID           uuid.UUID  `json:"category_id"`
	CategoryName         *string    `json:"category_name,omitempty"`
	OfficeLocationID     *uuid.UUID `json:"office_location_id,omitempty"`
	Status               string     `json:"status"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
	CanBeUsedInWorkOrder bool       `json:"can_be_used_in_work_order"`
	IsInStock            bool       `json:"is_in_stock"`
}
