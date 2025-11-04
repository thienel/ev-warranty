package persistence

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type officeRepository struct {
	db *gorm.DB
}

func NewOfficeRepository(db *gorm.DB) repository.OfficeRepository {
	return &officeRepository{db}
}

func (o *officeRepository) Create(ctx context.Context, user *entity.Office) error {
	if err := o.db.WithContext(ctx).Create(user).Error; err != nil {
		if dup := getDuplicateKeyConstraint(err); dup != "" {
			return apperror.ErrDuplicateKey.WithMessage(dup + " already existed").WithError(err)
		}
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (o *officeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Office, error) {
	var office entity.Office
	if err := o.db.WithContext(ctx).Where("id = ?", id).First(&office).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrNotFoundError.WithMessage("Office not found").WithError(err)
		}
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return &office, nil
}

func (o *officeRepository) FindAll(ctx context.Context) ([]*entity.Office, error) {
	var offices []*entity.Office
	if err := o.db.WithContext(ctx).Find(&offices).Error; err != nil {
		return nil, apperror.ErrDBOperation.WithError(err)
	}
	return offices, nil
}

func (o *officeRepository) Update(ctx context.Context, office *entity.Office) error {
	if err := o.db.WithContext(ctx).Model(office).
		Select("office_name", "office_type", "address", "is_active").
		Updates(office).Error; err != nil {
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}

func (o *officeRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := o.db.WithContext(ctx).Delete(&entity.Office{}, "id = ?", id).Error; err != nil {
		return apperror.ErrDBOperation.WithError(err)
	}
	return nil
}
