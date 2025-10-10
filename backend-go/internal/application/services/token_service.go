package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (*CustomClaims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error)
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
		return "", apperrors.NewFailedSignAccessToken(err)
	}

	return signedToken, nil
}

func (t *tokenService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperrors.NewFailedGenerateRefreshToken(err)
	}

	rawToken := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	hashedToken, err := hashToken(rawToken)
	if err != nil {
		return "", apperrors.NewFailedGenerateRefreshToken(err)
	}

	rfToken := &entities.RefreshToken{
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: time.Now().UTC().Add(t.refreshTTL),
	}

	if err = t.repoRefreshToken.Create(ctx, rfToken); err != nil {
		return "", apperrors.NewFailedGenerateRefreshToken(err)
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
			return nil, apperrors.NewUnexpectedSigningMethod(token.Header["alg"])
		}
		return t.publicKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, apperrors.NewExpiredAccessToken()
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, apperrors.NewInvalidAccessToken()
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, apperrors.NewInvalidAccessToken()
		default:
			return nil, apperrors.NewInvalidAccessToken()
		}
	}

	if token == nil || !token.Valid {
		return nil, apperrors.NewInvalidAccessToken()
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, apperrors.NewInvalidAccessToken()
	}

	if claims.UserID == "" {
		return nil, apperrors.NewInvalidAccessToken()
	}

	return claims, nil
}

func (t *tokenService) ValidateRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	hashedToken, err := hashToken(token)
	if err != nil {
		return nil, apperrors.NewFailedHashToken()
	}

	rfToken, err := t.repoRefreshToken.Find(ctx, hashedToken)
	if err != nil {
		return nil, apperrors.NewInvalidRefreshToken()
	}

	if rfToken.IsExpired() {
		return nil, apperrors.NewExpiredRefreshToken()
	}

	if rfToken.IsRevoked {
		return nil, apperrors.NewRevokedRefreshToken()
	}

	return rfToken, nil
}

func (t *tokenService) RevokeRefreshToken(ctx context.Context, token string) error {
	hashedToken, err := hashToken(token)
	if err != nil {
		return apperrors.NewFailedHashToken()
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
