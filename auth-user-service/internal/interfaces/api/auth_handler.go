package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/interfaces/api/dtos"
	"auth-service/pkg/logger"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestTimeout = 30 * time.Second
	bearerPrefix   = "Bearer "
)

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type authHandler struct {
	authService  services.AuthService
	tokenService services.TokenService
	logger       logger.Logger
}

func NewAuthHandler(authService services.AuthService, tokenService services.TokenService, logger logger.Logger) AuthHandler {
	return &authHandler{
		authService:  authService,
		tokenService: tokenService,
		logger:       logger,
	}
}

func (h *authHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidJSONRequest, "invalid JSON LoginRequest")
		return
	}

	h.logger.Info("attempting login", "username_or_email", req.UsernameOrEmail)

	accessToken, refreshToken, err := h.authService.Login(ctx, req.UsernameOrEmail, req.Password)
	if err != nil {
		handleError(h.logger, c, err, "login failed")
		return
	}

	userID, err := h.extractUserIDFromToken(ctx, accessToken)
	if err != nil {
		handleError(h.logger, c, err, "failed to extract user info")
		return
	}

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		handleError(h.logger, c, err, "failed to get user info")
		return
	}

	response := dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *dtos.GenerateUserDTO(*user),
	}

	h.logger.Info("login successful", "user_id", userID)
	writeSuccessResponse(c, http.StatusOK, "login successful", response)
}

func (h *authHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidJSONRequest, "invalid JSON RegisterRequest")
		return
	}

	h.logger.Info("attempting registration", "username", req.Username, "email", req.Email)

	user, err := h.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		handleError(h.logger, c, err, "registration failed")
		return
	}

	response := dtos.APIResponse{
		Success: true,
		Message: "registration successful",
		Data:    *dtos.GenerateUserDTO(*user),
	}

	h.logger.Info("registration successful", "user_id", user.ID)
	writeSuccessResponse(c, http.StatusCreated, "registration successful", response)
}

func (h *authHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	refreshToken := h.extractRefreshTokenFromBody(c)
	if refreshToken == "" {
		return
	}

	h.logger.Info("attempting logout")

	if err := h.authService.Logout(ctx, refreshToken); err != nil {
		handleError(h.logger, c, err, "logout failed")
		return
	}

	h.logger.Info("logout successful")
	writeSuccessResponse(c, http.StatusOK, "logout successful", nil)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	refreshToken := h.extractRefreshTokenFromBody(c)
	if refreshToken == "" {
		return
	}

	h.logger.Info("attempting token refresh")

	newAccessToken, err := h.tokenService.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		handleError(h.logger, c, err, "token refresh failed")
		return
	}

	response := map[string]string{
		"access_token": newAccessToken,
	}

	h.logger.Info("token refresh successful")
	writeSuccessResponse(c, http.StatusOK, "token refreshed successfully", response)
}

func (h *authHandler) ValidateToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	token := h.extractBearerToken(c)
	if token == "" {
		handleError(h.logger, c, apperrors.ErrInvalidAuthenticationHeader, "missing or invalid authorization header")
		return
	}

	h.logger.Info("attempting token validation")

	claims, err := h.tokenService.ValidateAccessToken(ctx, token)
	if err != nil {
		handleError(h.logger, c, err, "token validation failed")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid user ID"), "invalid user ID in token")
		return
	}

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		handleError(h.logger, c, err, "failed to get user info")
		return
	}

	response := map[string]any{
		"valid": true,
		"user":  dtos.GenerateUserDTO(*user),
	}

	h.logger.Info("token validation successful", "user_id", userID)
	writeSuccessResponse(c, http.StatusOK, "token is valid", response)
}

func (h *authHandler) extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}
	return strings.TrimPrefix(authHeader, bearerPrefix)
}

func (h *authHandler) extractRefreshTokenFromBody(c *gin.Context) string {
	var tokenReq struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	err := c.ShouldBindJSON(&tokenReq)
	if err != nil {
		return ""
	}

	return tokenReq.RefreshToken
}

func (h *authHandler) extractUserIDFromToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	claims, err := h.tokenService.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, apperrors.NewBadRequest("invalid user ID format")
	}

	return userID, nil
}
