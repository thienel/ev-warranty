package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClaimHandler interface {
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)

	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)

	Submit(c *gin.Context)
	Review(c *gin.Context)
	RequestInfo(c *gin.Context)
	Cancel(c *gin.Context)
	Complete(c *gin.Context)

	History(c *gin.Context)
}

type claimHandler struct {
	log       logger.Logger
	txManager application.TxManager
	service   services.ClaimService
}

func NewClaimHandler(log logger.Logger, txManager application.TxManager, claimService services.ClaimService) ClaimHandler {
	return &claimHandler{
		log:       log,
		txManager: txManager,
		service:   claimService,
	}
}

// GetByID godoc
// @Summary Get claim by ID
// @Description Retrieve a specific claim by its ID
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 200 {object} dtos.SuccessResponse{data=entities.Claim} "Claim retrieved successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id} [get]
func (h *claimHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	claim, err := h.service.GetByID(ctx, id)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, claim)
}

// GetAll godoc
// @Summary Get all claims
// @Description Retrieve a list of all claims with optional filtering and pagination
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by claim status"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} dtos.SuccessResponse{data=[]entities.Claim} "Claims retrieved successfully"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims [get]
func (h *claimHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	filters, err := parseClaimFilters(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	pagination := parsePagination(c)

	result, err := h.service.GetAll(ctx, filters, pagination)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, result)
}

// Create godoc
// @Summary Create a new claim
// @Description Create a new warranty claim (SC Technician/Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param createClaimRequest body dtos.CreateClaimRequest true "Claim creation data"
// @Success 201 {object} dtos.SuccessResponse{data=entities.Claim} "Claim created successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims [post]
func (h *claimHandler) Create(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScTechnician, entities.UserRoleScStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	var req dtos.CreateClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	cmd := &services.CreateClaimCommand{
		VehicleID:   req.VehicleID,
		CustomerID:  req.CustomerID,
		CreatorID:   userID,
		Description: req.Description,
	}

	var claim *entities.Claim
	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		var txErr error
		claim, txErr = h.service.Create(tx, cmd)
		return txErr
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusCreated, claim)
}

// Update godoc
// @Summary Update a claim
// @Description Update an existing claim by ID (SC Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param updateClaimRequest body dtos.UpdateClaimRequest true "Claim update data"
// @Success 204 "Claim updated successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id} [put]
func (h *claimHandler) Update(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	var req dtos.UpdateClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	cmd := &services.UpdateClaimCommand{
		Description: req.Description,
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.Update(tx, id, cmd)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete a claim
// @Description Delete a claim by ID (SC Staff: hard delete, EVM Staff: soft delete)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 204 "Claim deleted successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id} [delete]
func (h *claimHandler) Delete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScStaff, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		role, _ := getUserRoleFromHeader(c)
		if role == entities.UserRoleScStaff {
			return h.service.HardDelete(tx, id)
		}
		return h.service.SoftDelete(tx, id)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Submit godoc
// @Summary Submit a claim for review
// @Description Submit a claim to EVM for review (SC Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 204 "Claim submitted successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/submit [post]
func (h *claimHandler) Submit(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.Submit(tx, id, userID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Review godoc
// @Summary Review a claim
// @Description Review and approve/reject a submitted claim (EVM Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 204 "Claim reviewed successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/review [post]
func (h *claimHandler) Review(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.UpdateStatus(tx, id, entities.ClaimStatusReviewing, userID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RequestInfo godoc
// @Summary Request additional information for a claim
// @Description Request additional information from SC for a claim (EVM Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 204 "Information request sent successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/request-info [post]
func (h *claimHandler) RequestInfo(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.UpdateStatus(tx, id, entities.ClaimStatusRequestInfo, userID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Cancel godoc
// @Summary Cancel a claim
// @Description Cancel a claim and update its status (SC Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 204 "Claim cancelled successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/cancel [post]
func (h *claimHandler) Cancel(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.UpdateStatus(tx, id, entities.ClaimStatusCancelled, userID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Complete godoc
// @Summary Complete a claim
// @Description Mark a claim as completed (EVM Staff only)
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 204 "Claim completed successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/complete [post]
func (h *claimHandler) Complete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.Complete(tx, id, userID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// History godoc
// @Summary Get claim history
// @Description Retrieve the history of status changes for a specific claim
// @Tags claims
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 200 {object} dtos.SuccessResponse{data=[]entities.ClaimHistory} "Claim history retrieved successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/history [get]
func (h *claimHandler) History(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	history, err := h.service.GetHistory(ctx, id)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, history)
}
