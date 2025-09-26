package services

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"context"
	"strings"

	"github.com/google/uuid"
)

type CreateOfficeCommand struct {
	OfficeName string
	OfficeType string
	Address    string
	IsActive   bool
}

type UpdateOfficeCommand struct {
	OfficeName *string
	OfficeType *string
	Address    *string
	IsActive   *bool
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
		return nil, apperrors.ErrInvalidCredentials("invalid office type")
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

	if cmd.OfficeName != nil {
		office.OfficeName = strings.TrimSpace(*cmd.OfficeName)
	}

	if cmd.OfficeType != nil {
		*cmd.OfficeType = strings.TrimSpace(*cmd.OfficeType)
		if !entities.IsValidOfficeType(*cmd.OfficeType) {
			return apperrors.ErrInvalidCredentials("invalid office type")
		}
		office.OfficeType = *cmd.OfficeType
	}

	if cmd.Address != nil {
		office.Address = strings.TrimSpace(*cmd.Address)
	}

	if cmd.IsActive != nil {
		office.IsActive = *cmd.IsActive
	}

	if err = s.repo.Update(ctx, office); err != nil {
		return err
	}

	return nil
}

func (s *officeService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}
