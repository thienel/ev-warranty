package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type officeRepository struct {
	db *gorm.DB
}

func NewOfficeRepository(db *gorm.DB) repositories.OfficeRepository {
	return &officeRepository{db}
}

func (o *officeRepository) Create(ctx context.Context, user *entities.Office) error {
	if err := o.db.WithContext(ctx).Create(user).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperror.NewDBDuplicateKeyError(dup)
		}
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (o *officeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Office, error) {
	var office entities.Office
	if err := o.db.WithContext(ctx).Where("id = ?", id).First(&office).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewOfficeNotFound()
		}
		return nil, apperror.NewDBOperationError(err)
	}
	return &office, nil
}

func (o *officeRepository) FindAll(ctx context.Context) ([]*entities.Office, error) {
	var offices []*entities.Office
	if err := o.db.WithContext(ctx).Find(&offices).Error; err != nil {
		return nil, apperror.NewDBOperationError(err)
	}
	return offices, nil
}

func (o *officeRepository) Update(ctx context.Context, office *entities.Office) error {
	if err := o.db.WithContext(ctx).Model(office).
		Select("office_name", "office_type", "address", "is_active").
		Updates(office).Error; err != nil {
		return apperror.NewDBOperationError(err)
	}
	return nil
}

func (o *officeRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := o.db.WithContext(ctx).Delete(&entities.Office{}, "id = ?", id).Error; err != nil {
		return apperror.NewDBOperationError(err)
	}
	return nil
}
