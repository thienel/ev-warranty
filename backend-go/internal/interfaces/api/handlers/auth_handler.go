package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/interfaces/api/dto"
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

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.SuccessResponse{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("login request", "remote_addr", c.ClientIP())

	var req dto.LoginRequest
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

	response := dto.LoginResponse{
		Token: accessToken,
		User:  *dto.GenerateUserDTO(user),
	}

	h.log.Info("login successful", "user_id", userID, "email", user.Email)
	writeSuccessResponse(c, http.StatusOK, response)
}

// Logout godoc
// @Summary User logout
// @Description Logout user by invalidating refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 204 "Logout successful"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/logout [post]
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

	c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)

	h.log.Info("logout successful")
	c.Status(http.StatusNoContent)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token from cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse{data=dto.RefreshTokenResponse} "New access token"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/token [post]
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

	h.log.Info("token refresh successful")
	writeSuccessResponse(c, http.StatusOK, dto.RefreshTokenResponse{
		Token: newAccessToken,
	})
}

// ValidateToken godoc
// @Summary Validate access token
// @Description Validate the provided JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.SuccessResponse{data=dto.ValidateTokenResponse} "Token is valid"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/token [get]
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
		handleError(h.log, c, apperrors.NewInvalidUserID())
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Header("X-User-ID", userID.String())
	c.Header("X-User-Role", user.Role)

	h.log.Info("token validation successful", "user_id", userID, "role", user.Role)
	writeSuccessResponse(c, http.StatusOK,
		&dto.ValidateTokenResponse{
			Valid: true,
			User:  *dto.GenerateUserDTO(user),
		})
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
		return uuid.Nil, apperrors.NewInvalidUserID()
	}

	return userID, nil
}
