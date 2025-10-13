package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type CreateOfficeCommand struct {
	OfficeName string
	OfficeType string
	Address    string
	IsActive   bool
}

type UpdateOfficeCommand struct {
	OfficeName string
	OfficeType string
	Address    string
	IsActive   bool
}

type OfficeService interface {
	Create(ctx context.Context, cmd *CreateOfficeCommand) (*entities.Office, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Office, error)
	GetAll(ctx context.Context) ([]*entities.Office, error)
	Update(ctx context.Context, id uuid.UUID, cmd *UpdateOfficeCommand) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type officeService struct {
	repo repositories.OfficeRepository
}

func NewOfficeService(repo repositories.OfficeRepository) OfficeService {
	return &officeService{repo}
}

func (s *officeService) Create(ctx context.Context, cmd *CreateOfficeCommand) (*entities.Office, error) {
	if !entities.IsValidOfficeType(cmd.OfficeType) {
		return nil, apperrors.NewInvalidOfficeType()
	}

	office := entities.NewOffice(cmd.OfficeName, cmd.OfficeType, cmd.Address, cmd.IsActive)
	if err := s.repo.Create(ctx, office); err != nil {
		return nil, err
	}

	return office, nil
}

func (s *officeService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Office, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *officeService) GetAll(ctx context.Context) ([]*entities.Office, error) {
	return s.repo.FindAll(ctx)
}

func (s *officeService) Update(ctx context.Context, id uuid.UUID, cmd *UpdateOfficeCommand) error {
	office, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !entities.IsValidOfficeType(cmd.OfficeType) {
		return apperrors.NewInvalidOfficeType()
	}
	office.OfficeName = cmd.OfficeName
	office.OfficeType = cmd.OfficeType
	office.Address = cmd.Address
	office.IsActive = cmd.IsActive

	if err = s.repo.Update(ctx, office); err != nil {
		return err
	}

	return nil
}

func (s *officeService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}
