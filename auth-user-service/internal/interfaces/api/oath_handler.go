package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/infrastructure/oauth"
	"auth-service/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	errMsg = "Error while login with Google, please try again!"
)

type OAuthHandler interface {
	InitiateOAuth(c *gin.Context)
	HandleCallback(c *gin.Context)
}

type oauthHandler struct {
	log             logger.Logger
	frontendBaseURL string
	oauthService    oauth.OAuthService
	authService     services.AuthService
}

func NewOAuthHandler(log logger.Logger, frontendBaseURL string, oauthService oauth.OAuthService, authService services.AuthService) OAuthHandler {
	return &oauthHandler{
		log:             log,
		frontendBaseURL: frontendBaseURL,
		oauthService:    oauthService,
		authService:     authService,
	}
}

func (h *oauthHandler) InitiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	authURL, err := h.oauthService.GenerateAuthURL(provider)
	if err != nil {
		h.log.Error("Failed to generate auth URL")
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

func (h *oauthHandler) HandleCallback(c *gin.Context) {
	provider := c.Param("provider")
	state := c.Query("state")
	if state == "" {
		h.log.Error("Missing state in callback")
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	errParam := c.Query("error")
	if errParam != "" {
		h.log.Error("error in oauth callback: ", "error", errParam)
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	code := c.Query("code")
	if code == "" {
		h.log.Error("Missing code in callback")
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	userInfo, err := h.oauthService.HandleCallback(c.Request.Context(), provider, code, state)
	if err != nil {
		h.log.Error("error in handle oauth callback: ", "error", err)
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}
	accessToken, refreshToken, err := h.authService.HandleOAuthUser(c.Request.Context(), userInfo)
	if err != nil {
		h.log.Error("error in handle oauth user: ", "error", err)
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	c.SetCookie(
		"refreshToken",
		refreshToken,
		60*60*24*7,
		"/",
		"localhost",
		false,
		true,
	)

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/callback?token=%s", h.frontendBaseURL, accessToken))
}
