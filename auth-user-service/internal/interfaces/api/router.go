package api

import (
	"auth-service/internal/infrastructure/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
	DB *database.Database
}

func NewRouter(db *database.Database, authHandler AuthHandler, oauthHandler OAuthHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	healthHandler := healthHandler{DB: db}

	router.GET("/health", healthHandler.Check)
	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)
	router.POST("/logout", authHandler.Logout)
	router.POST("/token/refresh", authHandler.RefreshToken)
	router.GET("/token/validate", authHandler.ValidateToken)
	router.GET("/:provider", oauthHandler.InitiateOAuth)
	router.GET("/:provider/callback", oauthHandler.HandleCallback)
	return router
}

func (h *healthHandler) Check(c *gin.Context) {
	err := h.DB.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "unhealthy",
			"service": "auth-service",
			"db":      "down",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
		"db":      "up",
	})
}
