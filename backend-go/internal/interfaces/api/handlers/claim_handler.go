package handlers

import "C"
import (
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/pkg/logger"

	"github.com/gin-gonic/gin"
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
	log     logger.Logger
	service services.ClaimService
}

func NewClaimHandler(log logger.Logger, claimService services.ClaimService) ClaimHandler {
	return &claimHandler{
		log:     log,
		service: claimService,
	}
}

func (h *claimHandler) GetByID(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) GetAll(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Create(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Update(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Submit(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Review(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) RequestInfo(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Cancel(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) Complete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *claimHandler) History(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
