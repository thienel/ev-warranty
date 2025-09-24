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

type UserService interface {
	Create(ctx context.Context, name, email, role, password string, isActive bool, officeID uuid.UUID) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}

type userService struct {
	userRepo      repositories.UserRepository
	officeService OfficeService
}

func NewUserService(userRepo repositories.UserRepository, officeService OfficeService) UserService {
	return &userService{userRepo, officeService}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) Create(ctx context.Context, name, email, role, password string, isActive bool, officeID uuid.UUID) (
	*entities.User, error) {

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
