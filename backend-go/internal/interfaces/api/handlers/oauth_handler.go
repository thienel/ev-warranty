package handlers

import (
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/infrastructure/oauth"
	"ev-warranty-go/pkg/logger"
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
	h.log.Info("OAuth initiate", "remote_addr", c.ClientIP())

	authURL, err := h.oauthService.GenerateAuthURL()
	if err != nil {
		h.log.Error("OAuth auth URL generation failed", "error", err.Error())
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	h.log.Info("OAuth redirect")
	c.Redirect(http.StatusFound, authURL)
}

func (h *oauthHandler) HandleCallback(c *gin.Context) {
	h.log.Info("OAuth callback", "remote_addr", c.ClientIP())

	state := c.Query("state")
	if state == "" {
		h.log.Error("OAuth missing state")
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	errParam := c.Query("error")
	if errParam != "" {
		h.log.Error("OAuth provider error", "oauth_error", errParam)
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	code := c.Query("code")
	if code == "" {
		h.log.Error("OAuth missing code")
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	userInfo, err := h.oauthService.HandleCallback(c.Request.Context(), code, state)
	if err != nil {
		h.log.Error("OAuth callback failed", "error", err.Error())
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	h.log.Info("OAuth user info retrieved", "email", userInfo.Email)

	accessToken, refreshToken, err := h.authService.HandleOAuthUser(c.Request.Context(), userInfo)
	if err != nil {
		h.log.Error("OAuth user handling failed", "email", userInfo.Email, "error", err.Error())
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=%s", h.frontendBaseURL, errMsg))
		return
	}

	c.SetCookie("refreshToken", refreshToken, 60*60*24*7, "/", "localhost", false, true)

	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.frontendBaseURL, accessToken)
	h.log.Info("OAuth login successful", "email", userInfo.Email)
	c.Redirect(http.StatusFound, redirectURL)
}
