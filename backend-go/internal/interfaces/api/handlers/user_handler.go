package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandler struct {
	log         logger.Logger
	userService services.UserService
}

func NewUserHandler(log logger.Logger, userService services.UserService) UserHandler {
	return &userHandler{
		log:         log,
		userService: userService,
	}
}

func (h userHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("create user request", "remote_addr", c.ClientIP())

	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	h.log.Info("creating user", "email", req.Email, "role", req.Role)

	params := &services.UserCreateCommand{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.TrimSpace(req.Email),
		Role:     strings.TrimSpace(req.Role),
		Password: req.Password,
		IsActive: req.IsActive,
		OfficeID: req.OfficeID,
	}

	user, err := h.userService.Create(ctx, params)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	response := *dtos.GenerateUserDTO(user)
	h.log.Info("user created", "user_id", user.ID, "email", user.Email)
	writeSuccessResponse(c, http.StatusCreated, response)
}

func (h userHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("update user request", "remote_addr", c.ClientIP())

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidCredentials())
		return
	}

	var req dtos.UpdateUserRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.NewInvalidJsonRequest())
		return
	}

	h.log.Info("updating user", "user_id", id, "role", req.Role)

	cmd := &services.UserUpdateCommand{
		Name:     strings.TrimSpace(req.Name),
		Role:     strings.TrimSpace(req.Role),
		IsActive: req.IsActive,
		OfficeID: req.OfficeID,
	}

	err = h.userService.Update(ctx, id, cmd)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("user updated", "user_id", id)
	writeSuccessResponse(c, http.StatusOK, nil)
}

func (h userHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get user request", "remote_addr", c.ClientIP())

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidCredentials())
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("user retrieved", "user_id", userID, "email", user.Email)
	writeSuccessResponse(c, http.StatusOK, user)
}

func (h userHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get all users request", "remote_addr", c.ClientIP())

	users, err := h.userService.GetAll(ctx)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	usersDto := dtos.GenerateUserDTOList(users)
	h.log.Info("users retrieved", "count", len(usersDto))
	writeSuccessResponse(c, http.StatusOK, usersDto)
}

func (h userHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("delete user request", "remote_addr", c.ClientIP())

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidCredentials())
		return
	}

	err = h.userService.Delete(ctx, userID)
	if err != nil {
		handleError(h.log, c, err)
		return
	}

	h.log.Info("user deleted", "user_id", userID)
	writeSuccessResponse(c, http.StatusNoContent, nil)
}
