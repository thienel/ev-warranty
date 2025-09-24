package services

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"context"

	"github.com/google/uuid"
)

type OfficeService interface {
	Create(ctx context.Context, officeName string, officeType string, address string, isActive bool) (*entities.Office, error)
	GetByID(ctx context.Context, officeID uuid.UUID) (*entities.Office, error)
	GetAll(ctx context.Context) ([]*entities.Office, error)
	ActiveByID(ctx context.Context, officeID uuid.UUID) error
	DeActiveByID(ctx context.Context, officeID uuid.UUID) error
	UpdateByID(ctx context.Context, officeID uuid.UUID, officeName string, officeType string, address string) (*entities.Office, error)
	DeleteByID(ctx context.Context, officeID uuid.UUID) error
}

type officeService struct {
	repo repositories.OfficeRepository
}

func NewOfficeService(repo repositories.OfficeRepository) OfficeService {
	return &officeService{repo}
}

func (s *officeService) Create(ctx context.Context, officeName string, officeType string, address string, isActive bool) (
	*entities.Office, error) {

	if !entities.IsValidOfficeType(officeType) {
		return nil, apperrors.ErrInvalidCredentials("invalid office type")
	}

	office := entities.NewOffice(officeName, officeType, address, isActive)
	if err := s.repo.Create(ctx, office); err != nil {
		return nil, err
	}

	return office, nil
}

func (s *officeService) GetByID(ctx context.Context, officeID uuid.UUID) (*entities.Office, error) {
	return s.repo.FindByID(ctx, officeID)
}

func (s *officeService) GetAll(ctx context.Context) ([]*entities.Office, error) {
	return s.repo.FindAll(ctx)
}

func (s *officeService) ActiveByID(ctx context.Context, officeID uuid.UUID) error {
	office, err := s.repo.FindByID(ctx, officeID)
	if err != nil {
		return err
	}

	office.IsActive = true
	return s.repo.Update(ctx, office)
}

func (s *officeService) DeActiveByID(ctx context.Context, officeID uuid.UUID) error {
	office, err := s.repo.FindByID(ctx, officeID)
	if err != nil {
		return err
	}

	office.IsActive = false
	return s.repo.Update(ctx, office)
}

func (s *officeService) UpdateByID(ctx context.Context, officeID uuid.UUID, officeName string, officeType string, address string) (*entities.Office, error) {
	office, err := s.repo.FindByID(ctx, officeID)
	if err != nil {
		return nil, err
	}

	if !entities.IsValidOfficeType(officeType) {
		return nil, apperrors.ErrInvalidCredentials("invalid office type")
	}

	office.OfficeName = officeName
	office.OfficeType = officeType
	office.Address = address

	if err := s.repo.Update(ctx, office); err != nil {
		return nil, err
	}

	return office, nil
}

func (s *officeService) DeleteByID(ctx context.Context, officeID uuid.UUID) error {
	return s.repo.SoftDelete(ctx, officeID)
}
