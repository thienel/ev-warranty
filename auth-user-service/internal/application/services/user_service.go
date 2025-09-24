package services

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/security"
	"context"
	"strings"

	"github.com/google/uuid"
)

type UserCreateParams struct {
	Name     string
	Email    string
	Role     string
	Password string
	IsActive bool
	OfficeID uuid.UUID
}

type UserUpdateParams struct {
	Id       uuid.UUID
	Name     string
	Email    string
	Role     string
	IsActive bool
	OfficeID uuid.UUID
}

type UserService interface {
	Create(ctx context.Context, params *UserCreateParams) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetAll(ctx context.Context) ([]*entities.User, error)
	Update(ctx context.Context, params *UserUpdateParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo      repositories.UserRepository
	officeService OfficeService
}

func NewUserService(userRepo repositories.UserRepository, officeService OfficeService) UserService {
	return &userService{userRepo, officeService}
}

func (s *userService) Create(ctx context.Context, params *UserCreateParams) (
	*entities.User, error) {

	name := params.Name
	email := params.Email
	role := params.Role
	password := params.Password
	isActive := params.IsActive
	officeID := params.OfficeID

	email = strings.TrimSpace(email)
	if !entities.IsValidEmail(email) {
		return nil, apperrors.ErrInvalidCredentials("invalid email")
	}
	if !entities.IsValidPassword(password) {
		return nil, apperrors.ErrInvalidCredentials("invalid password")
	}
	if !entities.IsValidUserRole(role) {
		return nil, apperrors.ErrInvalidCredentials("invalid role")
	}
	if _, err := s.officeService.GetByID(ctx, officeID); err != nil {
		return nil, err
	}
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return nil, apperrors.ErrHashPassword(err)
	}

	name = strings.TrimSpace(name)
	user := entities.NewUser(name, email, role, passwordHash, isActive, officeID)

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

func (s *userService) Update(ctx context.Context, params *UserUpdateParams) error {
	user, err := s.userRepo.FindByID(ctx, params.Id)
	if err != nil {
		return err
	}

	name := params.Name
	email := params.Email
	role := params.Role
	isActive := params.IsActive
	officeID := params.OfficeID
	if user, _ := s.userRepo.FindByEmail(ctx, email); user != nil {
		return apperrors.ErrInvalidCredentials("invalid email")
	}
	if !entities.IsValidEmail(email) {
		return apperrors.ErrInvalidCredentials("invalid email")
	}
	if !entities.IsValidUserRole(role) {
		return apperrors.ErrInvalidCredentials("invalid role")
	}
	if office, err := s.officeService.GetByID(ctx, officeID); err != nil || office == nil {
		return err
	}

	user.Name = name
	user.Email = email
	user.Role = role
	user.IsActive = isActive
	user.OfficeID = officeID

	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.SoftDelete(ctx, id)
}
