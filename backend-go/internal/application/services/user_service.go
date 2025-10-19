package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/security"

	"github.com/google/uuid"
)

type UserCreateCommand struct {
	Name     string
	Email    string
	Role     string
	Password string
	IsActive bool
	OfficeID uuid.UUID
}

type UserUpdateCommand struct {
	Name     string
	Role     string
	IsActive bool
	OfficeID uuid.UUID
}

type UserService interface {
	Create(ctx context.Context, cmd *UserCreateCommand) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetAll(ctx context.Context) ([]*entities.User, error)
	Update(ctx context.Context, id uuid.UUID, cmd *UserUpdateCommand) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo      repositories.UserRepository
	officeService OfficeService
}

func NewUserService(userRepo repositories.UserRepository, officeService OfficeService) UserService {
	return &userService{userRepo, officeService}
}

func (s *userService) Create(ctx context.Context, cmd *UserCreateCommand) (*entities.User, error) {
	if !entities.IsValidName(cmd.Name) || !entities.IsValidEmail(cmd.Email) ||
		!entities.IsValidPassword(cmd.Password) || !entities.IsValidUserRole(cmd.Role) {
		return nil, apperrors.NewInvalidUserInput()
	}

	if _, err := s.officeService.GetByID(ctx, cmd.OfficeID); err != nil {
		return nil, err
	}

	passwordHash, err := security.HashPassword(cmd.Password)
	if err != nil {
		return nil, apperrors.NewHashPasswordError(err)
	}

	user := entities.NewUser(cmd.Name, cmd.Email, cmd.Role, passwordHash, cmd.IsActive, cmd.OfficeID)
	office, err := s.officeService.GetByID(ctx, cmd.OfficeID)
	if err != nil || office == nil {
		return nil, err
	}
	if !user.IsValidOfficeByRole(office.OfficeType) {
		return nil, apperrors.NewInvalidOfficeType()
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) GetAll(ctx context.Context) ([]*entities.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, cmd *UserUpdateCommand) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !entities.IsValidName(cmd.Name) || !entities.IsValidUserRole(cmd.Role) {
		return apperrors.NewInvalidUserInput()
	}
	user.Role = cmd.Role
	user.Name = cmd.Name
	user.IsActive = cmd.IsActive

	office, err := s.officeService.GetByID(ctx, cmd.OfficeID)
	if err != nil || office == nil {
		return err
	}
	if !user.IsValidOfficeByRole(office.OfficeType) {
		return apperrors.NewInvalidOfficeType()
	}
	user.OfficeID = cmd.OfficeID

	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.SoftDelete(ctx, id)
}
