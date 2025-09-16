package services

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/infrastructure/oauth/providers"
	"auth-service/internal/security"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	loginTypeEmail    = "email"
	loginTypeUsername = "username"
)

type AuthService interface {
	Login(ctx context.Context, usernameOrEmail, password string) (string, string, error)
	Register(ctx context.Context, username, email, password string) (*entities.User, error)
	Logout(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	HandleOAuthUser(ctx context.Context, userInfo *providers.UserInfo) (*entities.User, string, string, error)
}

type authService struct {
	userRepo     repositories.UserRepository
	tokenService TokenService
}

func NewAuthService(userRepo repositories.UserRepository, tokenService TokenService) AuthService {
	return &authService{userRepo, tokenService}
}

func (s *authService) Login(ctx context.Context, usernameOrEmail, password string) (string, string, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)
	loginType := s.identifyLoginType(usernameOrEmail)
	user, err := s.findUser(ctx, loginType, usernameOrEmail)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", apperrors.ErrInvalidCredentials("invalid username or email")
	}
	if !user.IsActive() {
		return "", "", apperrors.ErrUserInactive
	}
	if !security.VerifyPassword(password, *user.PasswordHash) {
		return "", "", apperrors.ErrInvalidCredentials("invalid password")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) identifyLoginType(usernameOrEmail string) string {
	if entities.IsValidEmail(usernameOrEmail) {
		return loginTypeEmail
	}
	if entities.IsValidUserName(usernameOrEmail) {
		return loginTypeUsername
	}
	return ""
}

func (s *authService) findUser(ctx context.Context, loginType, usernameOrEmail string) (*entities.User, error) {
	switch loginType {
	case loginTypeEmail:
		return s.userRepo.FindByEmail(ctx, usernameOrEmail)
	case loginTypeUsername:
		return s.userRepo.FindByUsername(ctx, usernameOrEmail)
	default:
		return nil, apperrors.ErrInvalidCredentials("invalid username or email")
	}
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *authService) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	if !entities.IsValidUserName(username) {
		return nil, apperrors.ErrInvalidCredentials("invalid username")
	}
	return s.userRepo.FindByUsername(ctx, username)
}

func (s *authService) Register(ctx context.Context, username, email, password string) (*entities.User, error) {
	email = strings.TrimSpace(email)
	if !entities.IsValidEmail(email) {
		return nil, apperrors.ErrInvalidCredentials("invalid email")
	}
	username = strings.TrimSpace(username)
	if !entities.IsValidUserName(username) {
		return nil, apperrors.ErrInvalidCredentials("invalid username")
	}
	if !entities.IsValidPassword(password) {
		return nil, apperrors.ErrInvalidCredentials("invalid password")
	}
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return nil, apperrors.ErrHashPassword(err)
	}

	user := entities.NewUser(username, email, passwordHash)

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	err := s.tokenService.RevokeRefreshToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) HandleOAuthUser(ctx context.Context, userInfo *providers.UserInfo) (*entities.User, string, string, error) {
	user, err := s.userRepo.FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID)
	if err != nil {
		if err.Error() == "user not found" {
			user, err = s.userRepo.FindByEmail(ctx, userInfo.Email)
			if err != nil && err.Error() != "user not found" {
				return nil, "", "", err
			}
		}
	}

	if user == nil {
		user = entities.NewOAuthUser(s.generateUsernameFromEmail(ctx, userInfo.Email), userInfo.Email, userInfo.Provider, userInfo.ProviderID)
		if err = s.userRepo.Create(ctx, user); err != nil {
			return nil, "", "", err
		}
	}

	if !user.IsOAuthUser() {
		user.LinkToOAuth(userInfo.Provider, userInfo.ProviderID)
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, "", "", err
		}
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) generateUsernameFromEmail(ctx context.Context, email string) string {
	parts := strings.Split(email, "@")
	baseUsername := parts[0]

	for i := 0; i < 10; i++ {
		username := baseUsername
		if i > 0 {
			username = fmt.Sprintf("%s%d", baseUsername, i)
		}

		if _, err := s.userRepo.FindByUsername(ctx, username); err != nil {
			return username
		}
	}

	return fmt.Sprintf("%s_%s", baseUsername, uuid.New().String()[:8])
}
