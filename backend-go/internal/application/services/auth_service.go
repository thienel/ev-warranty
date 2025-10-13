package services

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/internal/security"
	"strings"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	Logout(ctx context.Context, token string) error
	HandleOAuthUser(ctx context.Context, userInfo *providers.UserInfo) (string, string, error)
}

type authService struct {
	userRepo     repositories.UserRepository
	tokenService TokenService
}

func NewAuthService(userRepo repositories.UserRepository, tokenService TokenService) AuthService {
	return &authService{userRepo, tokenService}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	email = strings.TrimSpace(email)
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", apperrors.NewInvalidCredentials()
	}
	if !user.IsActive {
		return "", "", apperrors.NewUserInactive()
	}
	if !security.VerifyPassword(password, user.PasswordHash) {
		return "", "", apperrors.NewUserPasswordInvalid()
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	err := s.tokenService.RevokeRefreshToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) HandleOAuthUser(ctx context.Context, userInfo *providers.UserInfo) (string, string, error) {
	user, err := s.userRepo.FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID)
	if err != nil {
		if err.Error() == "user not found" {
			user, err = s.userRepo.FindByEmail(ctx, userInfo.Email)
			if err != nil {
				return "", "", err
			}
		}
	}

	if !user.IsOAuthUser() {
		user.LinkToOAuth(userInfo.Provider, userInfo.ProviderID)
		if err = s.userRepo.Update(ctx, user); err != nil {
			return "", "", err
		}
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
