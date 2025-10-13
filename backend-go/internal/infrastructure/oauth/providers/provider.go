package providers

import (
	"context"

	"golang.org/x/oauth2"
)

type Provider interface {
	Name() string
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error)
}

type UserInfo struct {
	Provider   string
	ProviderID string
	Email      string
	Name       string
}
