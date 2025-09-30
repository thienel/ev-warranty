package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/domain/entities"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/interfaces/api/dtos"
	"auth-service/pkg/logger"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OfficeHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type officeHandler struct {
	log     logger.Logger
	service services.OfficeService
}

func NewOfficeHandler(log logger.Logger, service services.OfficeService) OfficeHandler {
	return &officeHandler{
		log:     log,
		service: service,
	}
}

func (h *officeHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req dtos.CreateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON request")
		return
	}

	if !entities.IsValidOfficeType(req.OfficeType) {
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office type"), "invalid office type")
		return
	}

	cmd := &services.CreateOfficeCommand{
		OfficeName: strings.TrimSpace(req.OfficeName),
		OfficeType: strings.TrimSpace(req.OfficeType),
		Address:    strings.TrimSpace(req.Address),
		IsActive:   req.IsActive,
	}

	office, err := h.service.Create(ctx, cmd)
	if err != nil {
		handleError(h.log, c, err, "error creating office")
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
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	office, err := h.service.GetByID(ctx, officeID)
	if err != nil {
		handleError(h.log, c, err, "error getting office")
		return
	}
	writeSuccessResponse(c, http.StatusOK, "office retrieved", office)
}

func (h *officeHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	offices, err := h.service.GetAll(ctx)
	if err != nil {
		handleError(h.log, c, err, "error getting offices")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "offices retrieved", offices)
}

func (h *officeHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	var req dtos.UpdateOfficeRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON request")
		return
	}

	cmd := &services.UpdateOfficeCommand{
		OfficeName: strings.TrimSpace(req.OfficeName),
		OfficeType: strings.TrimSpace(req.OfficeType),
		Address:    strings.TrimSpace(req.Address),
		IsActive:   req.IsActive,
	}

	err = h.service.Update(ctx, id, cmd)
	if err != nil {
		handleError(h.log, c, err, "error updating office")
		return
	}

	writeSuccessResponse(c, http.StatusNoContent, "office updated", nil)
}

func (h *officeHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	err = h.service.DeleteByID(ctx, officeID)
	if err != nil {
		handleError(h.log, c, err, "error deleting office")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "office deleted", nil)
}
