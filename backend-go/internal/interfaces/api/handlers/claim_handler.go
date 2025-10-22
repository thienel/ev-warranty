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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

func (h *claimHandler) Complete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

func (h *claimHandler) History(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
		return
	}

	histories, err := h.service.GetHistory(ctx, id)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, histories)
}
