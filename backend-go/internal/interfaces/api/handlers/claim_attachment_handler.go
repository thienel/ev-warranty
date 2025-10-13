package handlers

import (
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/pkg/logger"

	"github.com/gin-gonic/gin"
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

func (c2 claimAttachmentHandler) GetByID(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c2 claimAttachmentHandler) GetByClaimID(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c2 claimAttachmentHandler) Create(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c2 claimAttachmentHandler) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
