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

func (h *claimItemHandler) Create(c *gin.Context) {
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
		handleError(h.log, c, apperrors.NewInvalidCredentials())
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

func (h *claimItemHandler) Delete(c *gin.Context) {
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

	writeSuccessResponse(c, http.StatusOK, gin.H{"message": "Claim item deleted successfully"})
}

func (h *claimItemHandler) Approve(c *gin.Context) {
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

	writeSuccessResponse(c, http.StatusOK, gin.H{"message": "Claim item approved successfully"})
}

func (h *claimItemHandler) Reject(c *gin.Context) {
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

	writeSuccessResponse(c, http.StatusOK, gin.H{"message": "Claim item rejected successfully"})
}

func parseItemIDParam(c *gin.Context) (uuid.UUID, error) {
	itemIDStr := c.Param("itemID")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidCredentials()
	}
	return itemID, nil
}

func parseClaimIDParam(c *gin.Context) (uuid.UUID, error) {
	claimIDStr := c.Param("id")
	claimID, err := uuid.Parse(claimIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidCredentials()
	}
	return claimID, nil
}
