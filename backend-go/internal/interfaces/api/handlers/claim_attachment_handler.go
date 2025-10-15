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

type ClaimAttachmentHandler interface {
	GetByID(c *gin.Context)
	GetByClaimID(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type claimAttachmentHandler struct {
	log       logger.Logger
	txManager application.TxManager
	service   services.ClaimAttachmentService
}

func NewClaimAttachmentHandler(log logger.Logger, txManager application.TxManager, service services.ClaimAttachmentService) ClaimAttachmentHandler {
	return &claimAttachmentHandler{
		log:       log,
		txManager: txManager,
		service:   service,
	}
}

func (h *claimAttachmentHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	attachmentID, err := parseAttachmentIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	attachment, err := h.service.GetByID(ctx, attachmentID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, attachment)
}

func (h *claimAttachmentHandler) GetByClaimID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	attachments, err := h.service.GetByClaimID(ctx, claimID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, attachments)
}

func (h *claimAttachmentHandler) Create(c *gin.Context) {
	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	var req dtos.CreateClaimAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	if !entities.IsValidAttachmentType(req.Type) {
		handleError(h.log, c, apperrors.NewInvalidCredentials())
		return
	}

	cmd := &services.CreateAttachmentCommand{
		Type: req.Type,
		URL:  req.URL,
	}

	var attachment *entities.ClaimAttachment
	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		var txErr error
		attachment, txErr = h.service.Create(tx, claimID, cmd)
		return txErr
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusCreated, attachment)
}

func (h *claimAttachmentHandler) Delete(c *gin.Context) {
	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	attachmentID, err := parseAttachmentIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		return h.service.HardDelete(tx, claimID, attachmentID)
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusOK, gin.H{"message": "Claim attachment deleted successfully"})
}

func parseAttachmentIDParam(c *gin.Context) (uuid.UUID, error) {
	attachmentIDStr := c.Param("attachmentID")
	attachmentID, err := uuid.Parse(attachmentIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidCredentials()
	}
	return attachmentID, nil
}
