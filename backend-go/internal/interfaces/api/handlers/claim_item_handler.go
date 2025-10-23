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

type ClaimItemHandler interface {
	GetByID(c *gin.Context)
	GetByClaimID(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Approve(c *gin.Context)
	Reject(c *gin.Context)
}

type claimItemHandler struct {
	log       logger.Logger
	txManager application.TxManager
	service   services.ClaimItemService
}

func NewClaimItemHandler(log logger.Logger, txManager application.TxManager, service services.ClaimItemService) ClaimItemHandler {
	return &claimItemHandler{
		log:       log,
		txManager: txManager,
		service:   service,
	}
}

// GetByID godoc
// @Summary Get claim item by ID
// @Description Retrieve a specific claim item by its ID
// @Tags claim-items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param itemID path string true "Claim Item ID"
// @Success 200 {object} dtos.SuccessResponse{data=entities.ClaimItem} "Claim item retrieved successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 404 {object} dtos.ErrorResponse "Claim item not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/items/{itemID} [get]
func (h *claimItemHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	itemID, err := parseItemIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	item, err := h.service.GetByID(ctx, itemID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, item)
}

// GetByClaimID godoc
// @Summary Get claim items by claim ID
// @Description Retrieve all items for a specific claim
// @Tags claim-items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 200 {object} dtos.SuccessResponse{data=dtos.ClaimItemListResponse} "Claim items retrieved successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/items [get]
func (h *claimItemHandler) GetByClaimID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	items, err := h.service.GetByClaimID(ctx, claimID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, items)
}

// Create godoc
// @Summary Create a new claim item
// @Description Add a new item to a claim (SC Staff only)
// @Tags claim-items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param createClaimItemRequest body dtos.CreateClaimItemRequest true "Claim item creation data"
// @Success 201 {object} dtos.SuccessResponse{data=entities.ClaimItem} "Claim item created successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/items [post]
func (h *claimItemHandler) Create(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	var req dtos.CreateClaimItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	if !entities.IsValidClaimItemType(req.Type) {
		handleError(h.log, c, apperrors.NewInvalidClaimItemType())
		return
	}

	cmd := &services.CreateClaimItemCommand{
		PartCategoryID:    req.PartCategoryID,
		FaultyPartID:      req.FaultyPartID,
		ReplacementPartID: req.ReplacementPartID,
		IssueDescription:  req.IssueDescription,
		Status:            entities.ClaimItemStatusPending,
		Type:              req.Type,
		Cost:              req.Cost,
	}

	var item *entities.ClaimItem
	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		var txErr error
		item, txErr = h.service.Create(tx, claimID, cmd)
		return txErr
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusCreated, item)
}

// Delete godoc
// @Summary Delete a claim item
// @Description Remove an item from a claim (SC Staff only)
// @Tags claim-items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param itemID path string true "Claim Item ID"
// @Success 204 "Claim item deleted successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim item not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/items/{itemID} [delete]
func (h *claimItemHandler) Delete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	itemID, err := parseItemIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.HardDelete(tx, claimID, itemID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Approve godoc
// @Summary Approve a claim item
// @Description Approve a claim item for processing (EVM Staff only)
// @Tags claim-items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param itemID path string true "Claim Item ID"
// @Success 204 "Claim item approved successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim item not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/items/{itemID}/approve [post]
func (h *claimItemHandler) Approve(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	itemID, err := parseItemIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.Approve(tx, claimID, itemID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Reject godoc
// @Summary Reject a claim item
// @Description Reject a claim item (EVM Staff only)
// @Tags claim-items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param itemID path string true "Claim Item ID"
// @Success 204 "Claim item rejected successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure 403 {object} dtos.ErrorResponse "Forbidden"
// @Failure 404 {object} dtos.ErrorResponse "Claim item not found"
// @Failure 500 {object} dtos.ErrorResponse "Internal server error"
// @Router /claims/{id}/items/{itemID}/reject [post]
func (h *claimItemHandler) Reject(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleEvmStaff); err != nil {
		handleError(h.log, c, err)
		return
	}

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	itemID, err := parseItemIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.Reject(tx, claimID, itemID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseItemIDParam(c *gin.Context) (uuid.UUID, error) {
	itemIDStr := c.Param("itemID")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidUUID()
	}
	return itemID, nil
}

func parseClaimIDParam(c *gin.Context) (uuid.UUID, error) {
	claimIDStr := c.Param("id")
	claimID, err := uuid.Parse(claimIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidUUID()
	}
	return claimID, nil
}
