package handlers

import (
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	requestTimeout = 30 * time.Second
	bearerPrefix   = "Bearer "
)

func handleError(log logger.Logger, c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		appErr = apperrors.NewInternalServerError(err)
	}

	log.Error(appErr.ErrorCode, "error", appErr.Error())
	c.JSON(appErr.HttpCode, dtos.APIResponse{
		Error: appErr.ErrorCode,
	})
}

func writeSuccessResponse(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, dtos.APIResponse{
		Data: data,
	})
}
