package port

import (
	"ev-warranty-go/internal/domain/entity"

	"github.com/google/uuid"
)

type PartPort interface {
	FindByOfficeIDAndCategoryID(office, categoryID uuid.UUID) (*entity.Part, error)
	UpdateStatus(id uuid.UUID, status string) error
}
