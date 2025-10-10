package providers

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type GoogleProvider struct {
	Config *oauth2.Config
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) Provider {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
	return &GoogleProvider{Config: cfg}
}

func (g *GoogleProvider) Name() string {
	return "google"
}

func (g *GoogleProvider) GetAuthURL(state string) string {
	url := g.Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return url
}

func (g *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.Config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New("exchange code failed")
	}
	return token, nil
}

func (g *GoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	verifier := oidc.NewVerifier(
		"https://accounts.google.com",
		oidc.NewRemoteKeySet(ctx, "https://www.googleapis.com/oauth2/v3/certs"),
		&oidc.Config{ClientID: g.Config.ClientID},
	)

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.New("id token verification failed")
	}

	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}

	if err = idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &UserInfo{
		Provider:   g.Name(),
		ProviderID: claims.Sub,
		Email:      claims.Email,
	}, nil
}
