package service

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/security"
	"strings"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	Logout(ctx context.Context, token string) error
	HandleOAuthUser(ctx context.Context, userInfo *providers.UserInfo) (string, string, error)
}

type authService struct {
	userRepo     repository.UserRepository
	tokenService TokenService
}

func NewAuthService(userRepo repository.UserRepository, tokenService TokenService) AuthService {
	return &authService{userRepo, tokenService}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	email = strings.TrimSpace(email)
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if !user.IsActive {
		return "", "", apperror.ErrUserInactive
	}
	if !security.VerifyPassword(password, user.PasswordHash) {
		return "", "", apperror.ErrInvalidCredentials
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
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.ErrorCode == apperror.ErrNotFoundError.ErrorCode {
			user, err = s.userRepo.FindByEmail(ctx, userInfo.Email)
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
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
