package api

import (
	"auth-service/internal/infrastructure/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
	DB *database.Database
}

func NewRouter(db *database.Database, authHandler AuthHandler, oauthHandler OAuthHandler, officeHandler OfficeHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	healthHandler := healthHandler{DB: db}

	router.GET("/health", healthHandler.Check)

	auth := router.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/token/refresh", authHandler.RefreshToken)
	auth.GET("/token/validate", authHandler.ValidateToken)

	auth.GET("/:provider", oauthHandler.InitiateOAuth)
	auth.GET("/:provider/callback", oauthHandler.HandleCallback)

	office := router.Group("/offices")
	office.POST("/", officeHandler.Create)
	office.GET("/", officeHandler.GetAll)
	office.GET("/:id", officeHandler.GetById)
	office.PUT("/:id", officeHandler.Update)
	office.PUT("/:id/activate", officeHandler.Active)
	office.PUT("/:id/deactivate", officeHandler.Inactive)
	office.DELETE("/:id", officeHandler.Delete)

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
