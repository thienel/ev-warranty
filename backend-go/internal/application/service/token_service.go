package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (*CustomClaims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type tokenService struct {
	repoRefreshToken repositories.RefreshTokenRepository
	accessTTL        time.Duration
	refreshTTL       time.Duration
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
}

func NewTokenService(repoRefreshToken repositories.RefreshTokenRepository, accessTokenTTL, refreshTokenTTL time.Duration, pri *rsa.PrivateKey, pub *rsa.PublicKey) TokenService {
	return &tokenService{
		repoRefreshToken: repoRefreshToken,
		accessTTL:        accessTokenTTL,
		refreshTTL:       refreshTokenTTL,
		privateKey:       pri,
		publicKey:        pub,
	}
}

func (t *tokenService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(t.accessTTL)

	claims := CustomClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "auth-service",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", apperror.NewFailedSignAccessToken(err)
	}

	return signedToken, nil
}

func (t *tokenService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperror.NewFailedGenerateRefreshToken(err)
	}

	rawToken := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	hashedToken, err := hashToken(rawToken)
	if err != nil {
		return "", apperror.NewFailedGenerateRefreshToken(err)
	}

	rfToken := &entity.RefreshToken{
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: time.Now().UTC().Add(t.refreshTTL),
	}

	if err = t.repoRefreshToken.Create(ctx, rfToken); err != nil {
		return "", apperror.NewFailedGenerateRefreshToken(err)
	}

	return rawToken, nil
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (t *tokenService) ValidateAccessToken(ctx context.Context, tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, apperror.NewUnexpectedSigningMethod(token.Header["alg"])
		}
		return t.publicKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, apperror.NewExpiredAccessToken()
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, apperror.NewInvalidAccessToken()
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, apperror.NewInvalidAccessToken()
		default:
			return nil, apperror.NewInvalidAccessToken()
		}
	}

	if token == nil || !token.Valid {
		return nil, apperror.NewInvalidAccessToken()
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, apperror.NewInvalidAccessToken()
	}

	if claims.UserID == "" {
		return nil, apperror.NewInvalidAccessToken()
	}

	return claims, nil
}

func (t *tokenService) ValidateRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	hashedToken, err := hashToken(token)
	if err != nil {
		return nil, apperror.NewFailedHashToken()
	}

	rfToken, err := t.repoRefreshToken.Find(ctx, hashedToken)
	if err != nil {
		return nil, apperror.NewInvalidRefreshToken()
	}

	if rfToken.IsExpired() {
		return nil, apperror.NewExpiredRefreshToken()
	}

	if rfToken.IsRevoked {
		return nil, apperror.NewRevokedRefreshToken()
	}

	return rfToken, nil
}

func (t *tokenService) RevokeRefreshToken(ctx context.Context, token string) error {
	hashedToken, err := hashToken(token)
	if err != nil {
		return apperror.NewFailedHashToken()
	}

	return t.repoRefreshToken.Revoke(ctx, hashedToken)
}

func (t *tokenService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	rfToken, err := t.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := t.GenerateAccessToken(rfToken.UserID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func hashToken(token string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
