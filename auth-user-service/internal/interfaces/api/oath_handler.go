package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/infrastructure/oauth"
	"auth-service/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthHandler interface {
	InitiateOAuth(c *gin.Context)
	HandleCallback(c *gin.Context)
}

type oauthHandler struct {
	log          logger.Logger
	oauthService oauth.OAuthService
	authService  services.AuthService
}

func NewOAuthHandler(log logger.Logger, oauthService oauth.OAuthService, authService services.AuthService) OAuthHandler {
	return &oauthHandler{
		log:          log,
		oauthService: oauthService,
		authService:  authService,
	}
}

func (h *oauthHandler) InitiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	authURL, err := h.oauthService.GenerateAuthURL(provider)
	if err != nil {
		handleError(h.log, c, err, "Failed to generate auth URL")
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

func (h *oauthHandler) HandleCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		err := apperrors.ErrInvalidCredentials("Missing code or state in callback")
		handleError(h.log, c, err, err.Error())
		return
	}

	userInfo, err := h.oauthService.HandleCallback(c.Request.Context(), provider, code, state)
	if err != nil {
		handleError(h.log, c, err, "Failed to handle callback")
		return
	}
	accessToken, refreshToken, err := h.authService.HandleOAuthUser(c.Request.Context(), userInfo)
	if err != nil {
		handleError(h.log, c, err, "Failed to handle OAuth user")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "OAuth login successful", dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *dtos.GenerateUserDTO(*user),
	})
	c.Redirect(http.StatusFound, fmt.Sprintf("http://localhost:3000/auth/callback?token=%s", accessToken))
}
