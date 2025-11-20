package handler

import (
	"ev-warranty-go/internal/application/service"
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
	authService     service.AuthService
}

func NewOAuthHandler(log logger.Logger, frontendBaseURL string, oauthService oauth.OAuthService, authService service.AuthService) OAuthHandler {
	return &oauthHandler{
		log:             log,
		frontendBaseURL: frontendBaseURL,
		oauthService:    oauthService,
		authService:     authService,
	}
}

// InitiateOAuth godoc
// @Summary Initiate Google OAuth login
// @Description Redirect to Google OAuth authorization URL
// @Tags auth
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect to Google OAuth"
// @Failure 500 {string} string "Redirect to frontend with error"
// @Router /auth/google [get]
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

// HandleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Process Google OAuth callback and authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param state query string true "OAuth state parameter"
// @Param code query string true "OAuth authorization code"
// @Param error query string false "OAuth error parameter"
// @Success 302 {string} string "Redirect to frontend with token"
// @Failure 302 {string} string "Redirect to frontend with error"
// @Router /auth/google/callback [get]
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

	redirectURL := fmt.Sprintf("%s/callback?token=%s", h.frontendBaseURL, accessToken)
	h.log.Info("OAuth login successful", "email", userInfo.Email)
	c.Redirect(http.StatusFound, redirectURL)
}
