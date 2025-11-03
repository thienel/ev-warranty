package port

import (
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type PartPort interface {
	FindByOfficeIDAndCategoryID(office, categoryID uuid.UUID) (*entities.Part, error)
	UpdateStatus(id uuid.UUID, status string) error
}
