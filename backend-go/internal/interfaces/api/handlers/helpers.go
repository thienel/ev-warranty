package handlers

import (
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/interfaces/api/dtos"
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

func handleError(log logger.Logger, c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		appErr = apperrors.NewInternalServerError(err)
	}

	log.Error(appErr.ErrorCode, "error", appErr.Error())
	c.JSON(appErr.HttpCode, dtos.ErrorResponse{
		Error: appErr.ErrorCode,
	})
}

func writeSuccessResponse(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, dtos.SuccessResponse{
		Data: data,
	})
}

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, error) {
	userIDStr := c.GetHeader(headerUserIDKey)
	if userIDStr == "" {
		return uuid.Nil, apperrors.NewMissingUserID()
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidUserID()
	}

	return userID, nil
}

func getUserRoleFromHeader(c *gin.Context) (string, error) {
	role := c.GetHeader(headerUserRole)
	if role == "" {
		return "", apperrors.NewMissingUserRole()
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

	return apperrors.NewUnauthorizedRole()
}
