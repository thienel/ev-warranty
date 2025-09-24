package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/domain/entities"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/interfaces/api/dtos"
	"auth-service/pkg/logger"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OfficeHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	GetAll(c *gin.Context)
	Active(c *gin.Context)
	Inactive(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type officeHandler struct {
	logger  logger.Logger
	service services.OfficeService
}

func NewOfficeHandler(log logger.Logger, service services.OfficeService) OfficeHandler {
	return &officeHandler{
		logger:  log,
		service: service,
	}
}

func (h *officeHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req dtos.CreateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidJSONRequest, "invalid JSON request")
		return
	}

	if !entities.IsValidOfficeType(req.OfficeType) {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office type"), "invalid office type")
		return
	}

	office, err := h.service.Create(ctx, req.OfficeName, req.OfficeType, req.Address, req.IsActive)
	if err != nil {
		handleError(h.logger, c, err, "error creating office")
		return
	}

	writeSuccessResponse(c, http.StatusCreated, "office created", office)
}

func (h *officeHandler) GetById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	office, err := h.service.GetByID(ctx, officeID)
	if err != nil {
		handleError(h.logger, c, err, "error getting office")
		return
	}
	writeSuccessResponse(c, http.StatusOK, "office retrieved", office)
}

func (h *officeHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	offices, err := h.service.GetAll(ctx)
	if err != nil {
		handleError(h.logger, c, err, "error getting offices")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "offices retrieved", offices)
}

func (h *officeHandler) Active(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	err = h.service.ActiveByID(ctx, officeID)
	if err != nil {
		handleError(h.logger, c, err, "error activating office")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "office activated", nil)
}

func (h *officeHandler) Inactive(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	err = h.service.DeActiveByID(ctx, officeID)
	if err != nil {
		handleError(h.logger, c, err, "error deactivating office")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "office deactivated", nil)
}

func (h *officeHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	var req dtos.UpdateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidJSONRequest, "invalid JSON request")
		return
	}

	if !entities.IsValidOfficeType(req.OfficeType) {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office type"), "invalid office type")
		return
	}

	office, err := h.service.UpdateByID(ctx, officeID, req.OfficeName, req.OfficeType, req.Address)
	if err != nil {
		handleError(h.logger, c, err, "error updating office")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "office updated", office)
}

func (h *officeHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.logger, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	err = h.service.DeleteByID(ctx, officeID)
	if err != nil {
		handleError(h.logger, c, err, "error deleting office")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "office deleted", nil)
}
