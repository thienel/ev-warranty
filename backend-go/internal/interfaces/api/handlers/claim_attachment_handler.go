package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
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

// GetByID godoc
// @Summary Get claim attachment by ID
// @Description Retrieve a specific claim attachment by its ID
// @Tags claim-attachments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param attachmentID path string true "Attachment ID"
// @Success 200 {object} dto.SuccessResponse{data=entities.ClaimAttachment} "Claim attachment retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Claim attachment not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /claims/{id}/attachments/{attachmentID} [get]
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

// GetByClaimID godoc
// @Summary Get claim attachments by claim ID
// @Description Retrieve all attachments for a specific claim
// @Tags claim-attachments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.ClaimAttachmentListResponse} "Claim attachments retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Claim not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /claims/{id}/attachments [get]
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

// Create godoc
// @Summary Upload claim attachments
// @Description Upload files as attachments to a claim (SC Technician only)
// @Tags claim-attachments
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param files formData file true "Files to upload"
// @Success 201 {object} dto.SuccessResponse{data=[]entities.ClaimAttachment} "Claim attachments uploaded successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Claim not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /claims/{id}/attachments [post]
func (h *claimAttachmentHandler) Create(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScTechnician); err != nil {
		handleError(h.log, c, err)
		return
	}

	claimID, err := parseClaimIDParam(c)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil || form == nil {
		handleError(h.log, c, apperrors.NewInvalidMultipartFormRequest())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		handleError(h.log, c, apperrors.NewInvalidMultipartFormRequest())
		return
	}

	var attachments []*entities.ClaimAttachment
	err = h.txManager.Do(c.Request.Context(), func(tx application.Tx) error {
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				return apperrors.NewInvalidMultipartFormRequest()
			}
			attachment, err := h.service.Create(tx, claimID, file)
			if err != nil {
				return err
			}
			if err = file.Close(); err != nil {
				h.log.Error("Failed to close file", "error", err)
			}
			attachments = append(attachments, attachment)
		}
		return nil
	})

	if err != nil {
		handleError(h.log, c, err)
		return
	}

	writeSuccessResponse(c, http.StatusCreated, attachments)
}

// Delete godoc
// @Summary Delete a claim attachment
// @Description Remove an attachment from a claim (SC Technician only)
// @Tags claim-attachments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Claim ID"
// @Param attachmentID path string true "Attachment ID"
// @Success 204 {object} dto.ErrorResponse "Claim attachment deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Claim attachment not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /claims/{id}/attachments/{attachmentID} [delete]
func (h *claimAttachmentHandler) Delete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleScTechnician); err != nil {
		handleError(h.log, c, err)
		return
	}

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

	writeSuccessResponse(c, http.StatusNoContent, nil)
}

func parseAttachmentIDParam(c *gin.Context) (uuid.UUID, error) {
	attachmentIDStr := c.Param("attachmentID")
	attachmentID, err := uuid.Parse(attachmentIDStr)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidUUID()
	}
	return attachmentID, nil
}
