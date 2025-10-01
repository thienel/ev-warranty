package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/interfaces/api/dtos"
	"auth-service/pkg/logger"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type authHandler struct {
	log          logger.Logger
	authService  services.AuthService
	tokenService services.TokenService
	userService  services.UserService
}

func NewAuthHandler(log logger.Logger, authService services.AuthService, tokenService services.TokenService,
	userService services.UserService) AuthHandler {

	return &authHandler{
		log:          log,
		authService:  authService,
		tokenService: tokenService,
		userService:  userService,
	}
}

func (h *authHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("login request", "remote_addr", c.ClientIP())

	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid login request", "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON LoginRequest")
		return
	}

	h.log.Info("attempting login", "email", req.Email)

	accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.log.Error("login failed", "email", req.Email, "error", err.Error())
		handleError(h.log, c, err, "login failed")
		return
	}

	userID, err := h.extractUserIDFromToken(ctx, accessToken)
	if err != nil {
		h.log.Error("failed to extract user ID", "error", err.Error())
		handleError(h.log, c, err, "failed to extract user info")
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		h.log.Error("failed to get user", "user_id", userID, "error", err.Error())
		handleError(h.log, c, err, "failed to get user info")
		return
	}

	c.SetCookie("refreshToken", refreshToken, 60*60*24*7, "/", "localhost", false, true)

	response := dtos.LoginResponse{
		Token: accessToken,
		User:  *dtos.GenerateUserDTO(user),
	}

	h.log.Info("login successful", "user_id", userID, "email", user.Email)
	writeSuccessResponse(c, http.StatusOK, "login successful", response)
}

func (h *authHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("logout request", "remote_addr", c.ClientIP())

	token, err := c.Cookie("refreshToken")
	if err != nil {
		h.log.Error("refresh token not found", "error", err.Error())
		err = apperrors.NewUnauthorized("refresh token not found")
		handleError(h.log, c, err, "refresh token not found in cookie")
		return
	}

	if err = h.authService.Logout(ctx, token); err != nil {
		h.log.Error("logout failed", "error", err.Error())
		handleError(h.log, c, err, "logout failed")
		return
	}

	h.log.Info("logout successful")
	writeSuccessResponse(c, http.StatusOK, "logout successful", nil)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("token refresh request", "remote_addr", c.ClientIP())

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		h.log.Error("refresh token not found", "error", err.Error())
		err = apperrors.NewUnauthorized("refresh token not found")
		handleError(h.log, c, err, "refresh token not found in cookie")
		return
	}

	newAccessToken, err := h.tokenService.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		h.log.Error("token refresh failed", "error", err.Error())
		handleError(h.log, c, err, "token refresh failed")
		return
	}

	response := map[string]string{"access_token": newAccessToken}

	h.log.Info("token refresh successful")
	writeSuccessResponse(c, http.StatusOK, "token refreshed successfully", response)
}

func (h *authHandler) ValidateToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("token validation request", "remote_addr", c.ClientIP())

	token := h.extractBearerToken(c)
	if token == "" {
		h.log.Error("missing authorization header")
		handleError(h.log, c, apperrors.ErrInvalidAuthenticationHeader, "missing or invalid authorization header")
		return
	}

	claims, err := h.tokenService.ValidateAccessToken(ctx, token)
	if err != nil {
		h.log.Error("token validation failed", "error", err.Error())
		handleError(h.log, c, err, "token validation failed")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		h.log.Error("invalid user ID in token", "user_id", claims.UserID, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid user ID"), "invalid user ID in token")
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		h.log.Error("failed to get user", "user_id", userID, "error", err.Error())
		handleError(h.log, c, err, "failed to get user info")
		return
	}

	response := map[string]any{
		"valid": true,
		"user":  dtos.GenerateUserDTO(user),
	}

	c.Header("X-User-ID", userID.String())
	c.Header("X-User-Role", user.Role)

	h.log.Info("token validation successful", "user_id", userID, "role", user.Role)
	writeSuccessResponse(c, http.StatusOK, "token is valid", response)
}

func (h *authHandler) extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}
	return strings.TrimPrefix(authHeader, bearerPrefix)
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
