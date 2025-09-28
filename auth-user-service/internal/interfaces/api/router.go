package api

import (
	"auth-service/internal/infrastructure/database"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	DB *database.Database
}

func NewRouter(db *database.Database, authHandler AuthHandler, oauthHandler OAuthHandler,
	officeHandler OfficeHandler, userHandler UserHandler) *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	router.Use(gin.Recovery(), gin.Logger())
	healthHandler := HealthHandler{DB: db}

	router.GET("/health", healthHandler.Check)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/token", authHandler.RefreshToken)
		auth.GET("/token", authHandler.ValidateToken)

		auth.GET("/:provider", oauthHandler.InitiateOAuth)
		auth.GET("/:provider/callback", oauthHandler.HandleCallback)
	}

	users := router.Group("/users")
	{
		users.POST("/", userHandler.Create)
		users.GET("/", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
		users.PATCH("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	office := router.Group("/offices")
	{
		office.POST("/", officeHandler.Create)
		office.GET("/", officeHandler.GetAll)
		office.GET("/:id", officeHandler.GetById)
		office.PATCH("/:id", officeHandler.Update)
		office.DELETE("/:id", officeHandler.Delete)
	}

	return router
}

func (h *HealthHandler) Check(c *gin.Context) {
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
