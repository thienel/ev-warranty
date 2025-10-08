package api

import (
	"errors"
	"ev-warranty-go/internal/errors/apperrors"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	requestTimeout = 30 * time.Second
	bearerPrefix   = "Bearer "
)

func handleError(log logger.Logger, c *gin.Context, err error, message string) {
	log.Error(message, "error", err)

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		writeErrorResponse(c, appErr.Code, appErr.Message)
		return
	}

	writeErrorResponse(c, http.StatusInternalServerError, "internal server error")
}

func writeSuccessResponse(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, dtos.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func writeErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, dtos.APIResponse{
		Success: false,
		Message: message,
	})
}
