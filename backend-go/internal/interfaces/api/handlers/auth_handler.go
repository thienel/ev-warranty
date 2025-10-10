package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
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
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	h.log.Info("attempting login", "email", req.Email)

	accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	userID, err := h.extractUserIDFromToken(ctx, accessToken)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.SetCookie("refreshToken", refreshToken, 60*60*24*7, "/", "localhost", false, true)

	response := dtos.LoginResponse{
		Token: accessToken,
		User:  *dtos.GenerateUserDTO(user),
	}

	h.log.Info("login successful", "user_id", userID, "email", user.Email)
	writeSuccessResponse(c, http.StatusOK, response)
}

func (h *authHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("logout request", "remote_addr", c.ClientIP())

	token, err := c.Cookie("refreshToken")
	if err != nil {
		err = apperrors.NewRefreshTokenNotFound()
		handleError(h.log, c, err)
		return
	}

	if err = h.authService.Logout(ctx, token); err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("logout successful")
	writeSuccessResponse(c, http.StatusOK, nil)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("token refresh request", "remote_addr", c.ClientIP())

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		err = apperrors.NewRefreshTokenNotFound()
		handleError(h.log, c, err)
		return
	}

	newAccessToken, err := h.tokenService.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	response := map[string]string{"access_token": newAccessToken}

	h.log.Info("token refresh successful")
	writeSuccessResponse(c, http.StatusOK, response)
}

func (h *authHandler) ValidateToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("token validation request", "remote_addr", c.ClientIP())

	token := h.extractBearerToken(c)
	if token == "" {
		handleError(h.log, c, apperrors.NewInvalidAuthHeader())
		return
	}

	claims, err := h.tokenService.ValidateAccessToken(ctx, token)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidCredentials())
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	response := map[string]any{
		"valid": true,
		"user":  dtos.GenerateUserDTO(user),
	}

	c.Header("X-User-ID", userID.String())
	c.Header("X-User-Role", user.Role)

	h.log.Info("token validation successful", "user_id", userID, "role", user.Role)
	writeSuccessResponse(c, http.StatusOK, response)
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
		return uuid.Nil, apperrors.NewInvalidCredentials()
	}

	return userID, nil
}
