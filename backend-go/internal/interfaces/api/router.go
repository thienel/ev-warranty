package api

import (
	"ev-warranty-go/internal/infrastructure/database"
	"ev-warranty-go/internal/interfaces/api/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	DB *database.Database
}

func NewRouter(db *database.Database, authHandler handlers.AuthHandler,
	oauthHandler handlers.OAuthHandler, officeHandler handlers.OfficeHandler,
	userHandler handlers.UserHandler, claimHandler handlers.ClaimHandler,
	itemHandler handlers.ClaimItemHandler, attachmentHandler handlers.ClaimAttachmentHandler) *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery(), gin.Logger())
	healthHandler := HealthHandler{DB: db}

	router.GET("/health", healthHandler.Check)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/token", authHandler.RefreshToken)
		auth.GET("/token", authHandler.ValidateToken)

		auth.GET("/google", oauthHandler.InitiateOAuth)
		auth.GET("/google/callback", oauthHandler.HandleCallback)
	}

	users := router.Group("/users")
	{
		users.POST("", userHandler.Create)
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	office := router.Group("/offices")
	{
		office.POST("", officeHandler.Create)
		office.GET("", officeHandler.GetAll)
		office.GET("/:id", officeHandler.GetById)
		office.PUT("/:id", officeHandler.Update)
		office.DELETE("/:id", officeHandler.Delete)
	}

	claim := router.Group("/claims")
	{
		claim.GET("", claimHandler.GetAll)
		claim.POST("", claimHandler.Create)
		claim.GET("/:id", claimHandler.GetByID)
		claim.PUT("/:id", claimHandler.Update)
		claim.DELETE("/:id", claimHandler.Delete)

		claim.POST("/:id/submit", claimHandler.Submit)
		claim.POST("/:id/review", claimHandler.Review)
		claim.POST("/:id/requestinfo", claimHandler.RequestInfo)
		claim.POST("/:id/cancel", claimHandler.Cancel)
		claim.POST("/:id/complete", claimHandler.Complete)
		claim.POST("/:id/history", claimHandler.History)
	}

	claimItem := router.Group("/claims/:id/items")
	{
		claimItem.GET("", itemHandler.GetByClaimID)
		claimItem.POST("", itemHandler.Create)
		claimItem.GET("/:itemID", itemHandler.GetByID)
		claimItem.DELETE("/:itemID", itemHandler.Delete)
		claimItem.POST("/:itemID/approve", itemHandler.Approve)
		claimItem.POST("/:itemID/reject", itemHandler.Reject)
	}

	claimAttachment := router.Group("/claims/:id/attachments")
	{
		claimAttachment.GET("", attachmentHandler.GetByClaimID)
		claimAttachment.POST("", attachmentHandler.Create)
		claimAttachment.GET("/:attachmentID", attachmentHandler.GetByID)
		claimAttachment.DELETE("/:attachmentID", attachmentHandler.Delete)
	}

	return router
}

func (h *HealthHandler) Check(c *gin.Context) {
	err := h.DB.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "unhealthy",
			"service": "backend-go",
			"db":      "down",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "backend-go",
		"db":      "up",
	})
}
