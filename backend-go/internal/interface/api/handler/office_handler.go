package handler

import (
	"context"
	"ev-warranty-go/internal/application/service"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/interface/api/dto"
	"ev-warranty-go/pkg/apperror"
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
	service service.OfficeService
}

func NewOfficeHandler(log logger.Logger, service service.OfficeService) OfficeHandler {
	return &officeHandler{
		log:     log,
		service: service,
	}
}

// Create godoc
// @Summary Create a new office
// @Description Create a new office (Admin only)
// @Tags offices
// @Accept json
// @Produce json
// @Security Bearer
// @Param createOfficeRequest body dto.CreateOfficeRequest true "Office creation data"
// @Success 201 {object} dto.SuccessResponse{data=entity.Office} "Office created successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /offices [post]
func (h *officeHandler) Create(c *gin.Context) {
	if err := allowedRoles(c, entity.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("create office request", "remote_addr", c.ClientIP())

	var req dto.CreateOfficeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperror.ErrInvalidJsonRequest)
		return
	}

	if !entity.IsValidOfficeType(req.OfficeType) {
		handleError(h.log, c, apperror.ErrInvalidInput)
		return
	}

	h.log.Info("creating office", "name", req.OfficeName, "type", req.OfficeType)

	cmd := &service.CreateOfficeCommand{
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

// GetByID godoc
// @Summary Get office by ID
// @Description Retrieve a specific office by its ID
// @Tags offices
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Office ID"
// @Success 200 {object} dto.SuccessResponse{data=entity.Office} "Office retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Office not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /offices/{id} [get]
func (h *officeHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get office request", "remote_addr", c.ClientIP())

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.log, c, apperror.ErrInvalidParams)
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

// GetAll godoc
// @Summary Get all offices
// @Description Retrieve a list of all offices
// @Tags offices
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.SuccessResponse{data=[]entity.Office} "Offices retrieved successfully"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /offices [get]
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

// Update godoc
// @Summary Update an office
// @Description Update an existing office by ID (Admin only)
// @Tags offices
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Office ID"
// @Param updateOfficeRequest body dto.UpdateOfficeRequest true "Office update data"
// @Success 204 "Office updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Office not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /offices/{id} [put]
func (h *officeHandler) Update(c *gin.Context) {
	if err := allowedRoles(c, entity.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("update office request", "remote_addr", c.ClientIP())

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(h.log, c, apperror.ErrInvalidParams)
		return
	}

	var req dto.UpdateOfficeRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperror.ErrInvalidJsonRequest)
		return
	}

	h.log.Info("updating office", "office_id", id, "name", req.OfficeName)

	cmd := &service.UpdateOfficeCommand{
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
	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete an office
// @Description Delete an office by ID (Admin only)
// @Tags offices
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Office ID"
// @Success 204 "Office deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Office not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /offices/{id} [delete]
func (h *officeHandler) Delete(c *gin.Context) {
	if err := allowedRoles(c, entity.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("delete office request", "remote_addr", c.ClientIP())

	officeIDStr := c.Param("id")
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		handleError(h.log, c, apperror.ErrInvalidParams)
		return
	}

	err = h.service.DeleteByID(ctx, officeID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("office deleted", "office_id", officeID)
	c.Status(http.StatusNoContent)
}
