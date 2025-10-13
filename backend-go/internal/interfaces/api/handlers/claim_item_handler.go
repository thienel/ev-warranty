package handlers

import (
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/pkg/logger"

	"github.com/gin-gonic/gin"
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
	log     logger.Logger
	service services.ClaimItemService
}

func NewClaimItemHandler(log logger.Logger, service services.ClaimItemService) ClaimItemHandler {
	return &claimItemHandler{
		log:     log,
		service: service,
	}
}

func (h *claimItemHandler) GetByID(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimItemHandler) GetByClaimID(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimItemHandler) Create(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimItemHandler) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimItemHandler) Approve(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimItemHandler) Reject(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
