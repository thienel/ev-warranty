package service

import (
	"context"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/security"

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
	Create(ctx context.Context, cmd *UserCreateCommand) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, cmd *UserUpdateCommand) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo   repository.UserRepository
	officeRepo repository.OfficeRepository
}

func NewUserService(userRepo repository.UserRepository, officeRepo repository.OfficeRepository) UserService {
	return &userService{userRepo, officeRepo}
}

func (s *userService) Create(ctx context.Context, cmd *UserCreateCommand) (*entity.User, error) {
	if !entity.IsValidName(cmd.Name) || !entity.IsValidEmail(cmd.Email) ||
		!entity.IsValidPassword(cmd.Password) || !entity.IsValidUserRole(cmd.Role) {
		return nil, apperror.NewInvalidUserInput()
	}

	office, err := s.officeRepo.FindByID(ctx, cmd.OfficeID)
	if err != nil {
		return nil, err
	}

	passwordHash, err := security.HashPassword(cmd.Password)
	if err != nil {
		return nil, apperror.NewHashPasswordError(err)
	}

	user := entity.NewUser(cmd.Name, cmd.Email, cmd.Role, passwordHash, cmd.IsActive, cmd.OfficeID)
	if !user.IsValidOfficeByRole(office.OfficeType) {
		return nil, apperror.NewInvalidOfficeType()
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) GetAll(ctx context.Context) ([]*entity.User, error) {
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

	if !entity.IsValidName(cmd.Name) || !entity.IsValidUserRole(cmd.Role) {
		return apperror.NewInvalidUserInput()
	}
	user.Role = cmd.Role
	user.Name = cmd.Name
	user.IsActive = cmd.IsActive

	office, err := s.officeRepo.FindByID(ctx, cmd.OfficeID)
	if err != nil {
		return err
	}
	if !user.IsValidOfficeByRole(office.OfficeType) {
		return apperror.NewInvalidOfficeType()
	}
	user.OfficeID = cmd.OfficeID

	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.SoftDelete(ctx, id)
}
