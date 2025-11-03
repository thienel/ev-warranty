package port

import (
	"ev-warranty-go/internal/domain/entity"

	"github.com/google/uuid"
)

type PartPort interface {
	ReserveByOfficeIDAndCategoryID(office, categoryID uuid.UUID) (*entity.Part, error)
	UnReserveByID(id uuid.UUID) error
}
