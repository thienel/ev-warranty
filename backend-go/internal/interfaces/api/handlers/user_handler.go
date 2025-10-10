package handlers

import (
	"context"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/errors/apperrors"
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
		h.log.Error("invalid create user request", "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON RegisterRequest")
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
		h.log.Error("user creation failed", "email", req.Email, "error", err.Error())
		handleError(h.log, c, err, "user creation failed")
		return
	}

	response := *dtos.GenerateUserDTO(user)
	h.log.Info("user created", "user_id", user.ID, "email", user.Email)
	writeSuccessResponse(c, http.StatusCreated, "registration successful", response)
}

func (h userHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("update user request", "remote_addr", c.ClientIP())

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error("invalid user ID", "id", idStr, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidCredentials("invalid user id"), "invalid user id")
		return
	}

	var req dtos.UpdateUserRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid update user request", "user_id", id, "error", err.Error())
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON RegisterRequest")
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
		h.log.Error("user update failed", "user_id", id, "error", err.Error())
		handleError(h.log, c, err, "user update failed")
		return
	}

	h.log.Info("user updated", "user_id", id)
	writeSuccessResponse(c, http.StatusOK, "user update successful", nil)
}

func (h userHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get user request", "remote_addr", c.ClientIP())

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.log.Error("invalid user ID", "id", userIDStr, "error", err.Error())
		err = apperrors.ErrInvalidCredentials("invalid id")
		handleError(h.log, c, err, "invalid id in get user by id")
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		h.log.Error("failed to get user", "user_id", userID, "error", err.Error())
		handleError(h.log, c, err, "user not found in get user by id")
		return
	}

	h.log.Info("user retrieved", "user_id", userID, "email", user.Email)
	writeSuccessResponse(c, http.StatusOK, "user retrieve successfully", user)
}

func (h userHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get all users request", "remote_addr", c.ClientIP())

	users, err := h.userService.GetAll(ctx)
	if err != nil {
		h.log.Error("failed to get users", "error", err.Error())
		handleError(h.log, c, err, "error getting all users")
		return
	}

	usersDto := dtos.GenerateUserDTOList(users)
	h.log.Info("users retrieved", "count", len(usersDto))
	writeSuccessResponse(c, http.StatusOK, "users retrieve successfully", usersDto)
}

func (h userHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("delete user request", "remote_addr", c.ClientIP())

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.log.Error("invalid user ID", "id", userIDStr, "error", err.Error())
		err = apperrors.ErrInvalidCredentials("invalid user id")
		handleError(h.log, c, err, "invalid user ID deleting user")
		return
	}

	err = h.userService.Delete(ctx, userID)
	if err != nil {
		h.log.Error("user deletion failed", "user_id", userID, "error", err.Error())
		handleError(h.log, c, err, "error deleting user")
		return
	}

	h.log.Info("user deleted", "user_id", userID)
	writeSuccessResponse(c, http.StatusNoContent, "delete user successfully", nil)
}
