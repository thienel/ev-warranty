package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
)

type OAuthService interface {
	GenerateAuthURL(providerName string) (string, error)
	HandleCallback(ctx context.Context, providerName, code, state string) (*providers.UserInfo, error)
	HandleCallbackError(ctx context.Context, state string)
}

type oauthService struct {
	provider providers.Provider
	userRepo repositories.UserRepository
	states   map[string]bool
}

func NewOAuthService(provider providers.Provider, userRepo repositories.UserRepository) OAuthService {
	return &oauthService{
		provider: provider,
		userRepo: userRepo,
		states:   make(map[string]bool),
	}
}

func (s *oauthService) GenerateAuthURL(providerName string) (string, error) {
	state := s.generateState()
	s.states[state] = true

	authURL := s.provider.GetAuthURL(state)
	return authURL, nil
}

func (s *oauthService) HandleCallback(ctx context.Context, providerName, code, state string) (*providers.UserInfo, error) {
	if !s.states[state] {
		return nil, apperrors.NewInvalidCredentials()
	}
	delete(s.states, state)

	token, err := s.provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return s.provider.GetUserInfo(ctx, token)
}

func (s *oauthService) HandleCallbackError(ctx context.Context, state string) {
	delete(s.states, state)
}

func (s *oauthService) generateState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
