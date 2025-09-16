package services

import (
	"auth-service/internal/domain/entities"
	"auth-service/internal/domain/repositories"
	"auth-service/internal/errors/apperrors"
	"auth-service/pkg/logger"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (*CustomClaims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type tokenService struct {
	log              logger.Logger
	repoRefreshToken repositories.RefreshTokenRepository
	accessTTL        time.Duration
	refreshTTL       time.Duration
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
}

func NewTokenService(log logger.Logger, repoRefreshToken repositories.RefreshTokenRepository, accessTokenTTL, refreshTokenTTL time.Duration, pri *rsa.PrivateKey, pub *rsa.PublicKey) TokenService {
	return &tokenService{
		log:              log,
		repoRefreshToken: repoRefreshToken,
		accessTTL:        accessTokenTTL,
		refreshTTL:       refreshTokenTTL,
		privateKey:       pri,
		publicKey:        pub,
	}
}

func (t *tokenService) GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
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

	select {
	case <-ctx.Done():
		return "", apperrors.ErrRequestTimeout(ctx.Err())
	default:
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", apperrors.ErrFailedSignAccessToken(err)
	}

	return signedToken, nil
}

func (t *tokenService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	select {
	case <-ctx.Done():
		return "", apperrors.ErrRequestTimeout(ctx.Err())
	default:
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperrors.ErrFailedGenerateRefreshToken(err)
	}

	rawToken := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	hashedToken, err := hashToken(rawToken)
	if err != nil {
		return "", apperrors.ErrFailedGenerateRefreshToken(err)
	}

	if err := t.repoRefreshToken.RevokeAllByUserID(ctx, userID); err != nil {
		t.log.Warn("failed to revoke existing tokens", "user_id", userID, "error", err)
	}

	rfToken := &entities.RefreshToken{
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: time.Now().UTC().Add(t.refreshTTL),
	}

	if err := t.repoRefreshToken.Create(ctx, rfToken); err != nil {
		return "", apperrors.ErrFailedGenerateRefreshToken(err)
	}

	return rawToken, nil
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (t *tokenService) ValidateAccessToken(ctx context.Context, tokenStr string) (*CustomClaims, error) {
	select {
	case <-ctx.Done():
		return nil, apperrors.ErrRequestTimeout(ctx.Err())
	default:
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.publicKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, apperrors.ErrExpiredAccessToken
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, apperrors.ErrInvalidAccessToken
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, apperrors.ErrInvalidAccessToken
		default:
			return nil, apperrors.ErrInvalidAccessToken
		}
	}

	if token == nil || !token.Valid {
		return nil, apperrors.ErrInvalidAccessToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, apperrors.ErrInvalidAccessToken
	}

	if claims.UserID == "" {
		return nil, apperrors.ErrInvalidAccessToken
	}

	return claims, nil
}

func (t *tokenService) ValidateRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	hashedToken, err := hashToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to hash token: %w", err)
	}

	rfToken, err := t.repoRefreshToken.FindByToken(ctx, hashedToken)
	if err != nil {
		return nil, apperrors.ErrInvalidRefreshToken
	}

	if rfToken.IsExpired() {
		return nil, apperrors.ErrExpiredRefreshToken
	}

	if rfToken.IsRevoked {
		return nil, apperrors.ErrRevokedRefreshToken
	}

	return rfToken, nil
}

func (t *tokenService) RevokeRefreshToken(ctx context.Context, token string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	hashedToken, err := hashToken(token)
	if err != nil {
		return fmt.Errorf("failed to hash token: %w", err)
	}

	return t.repoRefreshToken.RevokeByToken(ctx, hashedToken)
}

func (t *tokenService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	rfToken, err := t.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := t.GenerateAccessToken(ctx, rfToken.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func hashToken(token string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
