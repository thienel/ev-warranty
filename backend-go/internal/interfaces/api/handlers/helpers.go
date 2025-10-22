package handlers

import (
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
	"strconv"
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
	c.JSON(appErr.HttpCode, dtos.APIResponse{
		Error: appErr.ErrorCode,
	})
}

func writeSuccessResponse(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, dtos.APIResponse{
		Data: data,
	})
}

func parseClaimFilters(c *gin.Context) (services.ClaimFilters, error) {
	filters := services.ClaimFilters{}

	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		customerID, err := uuid.Parse(customerIDStr)
		if err != nil {
			return filters, apperrors.NewInvalidQueryParameter("customer_id")
		}
		filters.CustomerID = &customerID
	}

	if vehicleIDStr := c.Query("vehicle_id"); vehicleIDStr != "" {
		vehicleID, err := uuid.Parse(vehicleIDStr)
		if err != nil {
			return filters, apperrors.NewInvalidQueryParameter("vehicle_id")
		}
		filters.VehicleID = &vehicleID
	}

	if status := c.Query("status"); status != "" {
		if !entities.IsValidClaimStatus(status) {
			return filters, apperrors.NewInvalidQueryParameter("status")
		}
		filters.Status = &status
	}

	if fromDateStr := c.Query("from_date"); fromDateStr != "" {
		fromDate, err := time.Parse(time.RFC3339, fromDateStr)
		if err != nil {
			return filters, apperrors.NewInvalidQueryParameter("from_date")
		}
		filters.FromDate = &fromDate
	}

	if toDateStr := c.Query("to_date"); toDateStr != "" {
		toDate, err := time.Parse(time.RFC3339, toDateStr)
		if err != nil {
			return filters, apperrors.NewInvalidQueryParameter("to_date")
		}
		filters.ToDate = &toDate
	}

	return filters, nil
}

func parsePagination(c *gin.Context) services.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return services.Pagination{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}
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
	userRole := c.GetHeader(headerUserRole)
	if userRole == "" {
		return "", apperrors.NewMissingUserRole()
	}
	return userRole, nil
}

func allowedRoles(c *gin.Context, allowedRoles ...string) error {
	userRole := c.GetHeader(headerUserRole)
	for _, role := range allowedRoles {
		if userRole == role {
			return nil
		}
	}
	return apperrors.NewUnauthorizedRole()
}
