package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/pkg/apperror"
)

type OAuthService interface {
	GenerateAuthURL() (string, error)
	HandleCallback(ctx context.Context, code, state string) (*providers.UserInfo, error)
}

type oauthService struct {
	provider providers.Provider
	userRepo repository.UserRepository
	states   map[string]bool
}

func NewOAuthService(provider providers.Provider, userRepo repository.UserRepository) OAuthService {
	return &oauthService{
		provider: provider,
		userRepo: userRepo,
		states:   make(map[string]bool),
	}
}

func (s *oauthService) GenerateAuthURL() (string, error) {
	state := s.generateState()
	s.states[state] = true

	authURL := s.provider.GetAuthURL(state)
	return authURL, nil
}

func (s *oauthService) HandleCallback(ctx context.Context, code, state string) (*providers.UserInfo, error) {
	if !s.states[state] {
		return nil, apperror.ErrInvalidAuthHeader
	}
	delete(s.states, state)

	token, err := s.provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return s.provider.GetUserInfo(ctx, token)
}

func (s *oauthService) generateState() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
