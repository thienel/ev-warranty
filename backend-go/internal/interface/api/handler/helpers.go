package handler

import (
	"errors"
	"ev-warranty-go/internal/interface/api/dto"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestTimeout  = 30 * time.Second
	bearerPrefix    = "Bearer "
	headerUserIDKey = "X-User-ID"
	headerUserRole  = "X-User-Role"
)

func writeErrorResponse(log logger.Logger, c *gin.Context, err error) {
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		appErr = apperror.ErrInternalServerError.WithError(err)
	}

	log.Error(appErr.ErrorCode, "error", appErr.Error())

	log.Error(err.Error(), "method", c.Request.Method, "path", c.FullPath(), "error", err)

	c.JSON(appErr.HttpCode, dto.APIResponse{
		IsSuccess: false,
		Error:     appErr.ErrorCode,
		Message:   appErr.Message,
	})
}

func writeSuccessResponse(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, dto.APIResponse{
		IsSuccess: true,
		Data:      data,
	})
}

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, error) {
	userIDStr := c.GetHeader(headerUserIDKey)
	if userIDStr == "" {
		return uuid.Nil, apperror.ErrMissingUserID
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, apperror.ErrInvalidUserID
	}

	return userID, nil
}

func getUserRoleFromHeader(c *gin.Context) (string, error) {
	role := c.GetHeader(headerUserRole)
	if role == "" {
		return "", apperror.ErrMissingUserRole
	}

	return role, nil
}

func allowedRoles(c *gin.Context, allowedRoles ...string) error {
	userRole, err := getUserRoleFromHeader(c)
	if err != nil {
		return err
	}

	for _, role := range allowedRoles {
		if userRole == role {
			return nil
		}
	}

	return apperror.ErrUnauthorizedRole
}
