package persistence

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type officeRepository struct {
	db *gorm.DB
}

const entityOfficeName = "office"

func NewOfficeRepository(db *gorm.DB) repositories.OfficeRepository {
	return &officeRepository{db}
}

func (o *officeRepository) Create(ctx context.Context, user *entities.Office) error {
	if err := o.db.WithContext(ctx).Create(user).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperrors.ErrDuplicateKey(dup)
		}
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (o *officeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Office, error) {
	var office entities.Office
	if err := o.db.WithContext(ctx).Where("id = ?", id).First(&office).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound(entityOfficeName)
		}
		return nil, apperrors.ErrDBOperation(err)
	}
	return &office, nil
}

func (o *officeRepository) Update(ctx context.Context, office *entities.Office) error {
	if err := o.db.WithContext(ctx).Model(office).
		Select("office_name", "office_type", "address", "is_active").
		Updates(office).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}

func (o *officeRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := o.db.WithContext(ctx).Delete(&entities.Office{}, "id = ?", id).Error; err != nil {
		return apperrors.ErrDBOperation(err)
	}
	return nil
}
