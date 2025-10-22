package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OfficeHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
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
	if err := allowedRoles(c, entities.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("create office request", "remote_addr", c.ClientIP())

	var req dtos.CreateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	if !entities.IsValidOfficeType(req.OfficeType) {
		handleError(h.log, c, apperrors.NewInvalidOfficeType())
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
		handleError(h.log, c, err)
		return
	}

	h.log.Info("office created", "office_id", office.ID, "name", office.OfficeName)
	writeSuccessResponse(c, http.StatusCreated, office)
}

func (h *officeHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get office request", "remote_addr", c.ClientIP())

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	office, err := h.service.GetByID(ctx, officeID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("office retrieved", "office_id", officeID, "name", office.OfficeName)
	writeSuccessResponse(c, http.StatusOK, office)
}

func (h *officeHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get all offices request", "remote_addr", c.ClientIP())

	offices, err := h.service.GetAll(ctx)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("offices retrieved", "count", len(offices))
	writeSuccessResponse(c, http.StatusOK, offices)
}

func (h *officeHandler) Update(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("update office request", "remote_addr", c.ClientIP())

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	var req dtos.UpdateOfficeRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
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
		handleError(h.log, c, err)
		return
	}

	h.log.Info("office updated", "office_id", id)
	writeSuccessResponse(c, http.StatusNoContent, nil)
}

func (h *officeHandler) Delete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("delete office request", "remote_addr", c.ClientIP())

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	err = h.service.DeleteByID(ctx, officeID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("office deleted", "office_id", officeID)
	writeSuccessResponse(c, http.StatusNoContent, nil)
}
