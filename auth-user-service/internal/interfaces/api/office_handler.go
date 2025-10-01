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

	h.log.Info("create office request", "remote_addr", c.ClientIP())

	var req dtos.CreateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid create office request", "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON request")
		return
	}

	if !entities.IsValidOfficeType(req.OfficeType) {
		h.log.Error("invalid office type", "office_type", req.OfficeType)
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office type"), "invalid office type")
		return
	}

	h.log.Info("creating office", "name", req.OfficeName, "type", req.OfficeType)

	cmd := &services.CreateOfficeCommand{
		OfficeName: strings.TrimSpace(req.OfficeName),
		OfficeType: strings.TrimSpace(req.OfficeType),
		Address:    strings.TrimSpace(req.Address),
		IsActive:   req.IsActive,
	}

	office, err := h.service.Create(ctx, cmd)
	if err != nil {
		h.log.Error("office creation failed", "name", cmd.OfficeName, "error", err.Error())
		handleError(h.log, c, err, "error creating office")
		return
	}

	h.log.Info("office created", "office_id", office.ID, "name", office.OfficeName)
	writeSuccessResponse(c, http.StatusCreated, "office created", office)
}

func (h *officeHandler) GetById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get office request", "remote_addr", c.ClientIP())

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		h.log.Error("invalid office ID", "id", officeIDStr, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	office, err := h.service.GetByID(ctx, officeID)
	if err != nil {
		h.log.Error("failed to get office", "office_id", officeID, "error", err.Error())
		handleError(h.log, c, err, "error getting office")
		return
	}

	h.log.Info("office retrieved", "office_id", officeID, "name", office.OfficeName)
	writeSuccessResponse(c, http.StatusOK, "office retrieved", office)
}

func (h *officeHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get all offices request", "remote_addr", c.ClientIP())

	offices, err := h.service.GetAll(ctx)
	if err != nil {
		h.log.Error("failed to get offices", "error", err.Error())
		handleError(h.log, c, err, "error getting offices")
		return
	}

	h.log.Info("offices retrieved", "count", len(offices))
	writeSuccessResponse(c, http.StatusOK, "offices retrieved", offices)
}

func (h *officeHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("update office request", "remote_addr", c.ClientIP())

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error("invalid office ID", "id", idStr, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	var req dtos.UpdateOfficeRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid update office request", "office_id", id, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON request")
		return
	}

	h.log.Info("updating office", "office_id", id, "name", req.OfficeName)

	cmd := &services.UpdateOfficeCommand{
		OfficeName: strings.TrimSpace(req.OfficeName),
		OfficeType: strings.TrimSpace(req.OfficeType),
		Address:    strings.TrimSpace(req.Address),
		IsActive:   req.IsActive,
	}

	err = h.service.Update(ctx, id, cmd)
	if err != nil {
		h.log.Error("office update failed", "office_id", id, "error", err.Error())
		handleError(h.log, c, err, "error updating office")
		return
	}

	h.log.Info("office updated", "office_id", id)
	writeSuccessResponse(c, http.StatusNoContent, "office updated", nil)
}

func (h *officeHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("delete office request", "remote_addr", c.ClientIP())

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		h.log.Error("invalid office ID", "id", officeIDStr, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid office ID"), "invalid office ID format")
		return
	}

	err = h.service.DeleteByID(ctx, officeID)
	if err != nil {
		h.log.Error("office deletion failed", "office_id", officeID, "error", err.Error())
		handleError(h.log, c, err, "error deleting office")
		return
	}

	h.log.Info("office deleted", "office_id", officeID)
	writeSuccessResponse(c, http.StatusOK, "office deleted", nil)
}
