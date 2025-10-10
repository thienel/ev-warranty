package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/errors/apperrors"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"fmt"
)

type OAuthService interface {
	RegisterProvider(provider providers.Provider)
	GenerateAuthURL(providerName string) (string, error)
	HandleCallback(ctx context.Context, providerName, code, state string) (*providers.UserInfo, error)
	HandleCallbackError(ctx context.Context, state string)
}

type oauthService struct {
	providers map[string]providers.Provider
	userRepo  repositories.UserRepository
	states    map[string]bool
}

func NewOAuthService(userRepo repositories.UserRepository) OAuthService {
	return &oauthService{
		providers: make(map[string]providers.Provider),
		userRepo:  userRepo,
		states:    make(map[string]bool),
	}
}

func (s *oauthService) RegisterProvider(provider providers.Provider) {
	s.providers[provider.Name()] = provider
}

func (s *oauthService) GenerateAuthURL(providerName string) (string, error) {
	provider, exists := s.providers[providerName]
	if !exists {
		return "", apperrors.ErrNotFound(fmt.Sprintf("provider %s", providerName))
	}

	state := s.generateState()
	s.states[state] = true

	authURL := provider.GetAuthURL(state)
	return authURL, nil
}

func (s *oauthService) HandleCallback(ctx context.Context, providerName, code, state string) (*providers.UserInfo, error) {
	if !s.states[state] {
		return nil, apperrors.ErrInvalidCredentials("invalid state parameter")
	}
	delete(s.states, state)

	provider, exists := s.providers[providerName]
	if !exists {
		return nil, apperrors.ErrNotFound(fmt.Sprintf("provider %s", providerName))
	}

	token, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return provider.GetUserInfo(ctx, token)
}

func (s *oauthService) HandleCallbackError(ctx context.Context, state string) {
	delete(s.states, state)
}

func (s *oauthService) generateState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
