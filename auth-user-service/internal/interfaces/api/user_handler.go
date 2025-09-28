package api

import (
	"auth-service/internal/application/services"
	"auth-service/internal/errors/apperrors"
	"auth-service/internal/interfaces/api/dtos"
	"auth-service/pkg/logger"
	"context"
	"net/http"

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

	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON RegisterRequest")
		return
	}

	h.log.Info("creating user", "email", req.Email)

	params := &services.UserCreateCommand{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
		IsActive: req.IsActive,
		OfficeID: req.OfficeID,
	}
	user, err := h.userService.Create(ctx, params)
	if err != nil {
		handleError(h.log, c, err, "user creation failed")
		return
	}

	h.log.Info("registration successful", "user_id", user.ID)
	writeSuccessResponse(c, http.StatusCreated, "registration successful", *dtos.GenerateUserDTO(*user))
}

func (h userHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(h.log, c, apperrors.ErrInvalidJSONRequest, "invalid JSON RegisterRequest")
		return
	}

	h.log.Info("updating user", "email", req.Email)

	cmd := &services.UserUpdateCommand{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
		OfficeID: req.OfficeID,
	}
	err := h.userService.Update(ctx, req.ID, cmd)
	if err != nil {
		handleError(h.log, c, err, "user update failed")
		return
	}

	h.log.Info("registration successful", "user_id", req.ID)
	writeSuccessResponse(c, http.StatusOK, "user update successful", nil)
}

func (h userHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		err = apperrors.ErrInvalidCredentials("invalid id")
		handleError(h.log, c, err, "invalid id in get user by id")
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		handleError(h.log, c, err, "user not found in get user by id")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "user retrieve successfully", user)
}

func (h userHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	users, err := h.userService.GetAll(ctx)
	if err != nil {
		handleError(h.log, c, err, "error getting all users")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "users retrieve successfully", users)
}

func (h userHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		err = apperrors.ErrInvalidCredentials("invalid user id")
		handleError(h.log, c, err, "invalid user ID deleting user")
		return
	}

	err = h.userService.Delete(ctx, userID)
	if err != nil {
		handleError(h.log, c, err, "error deleting user")
		return
	}

	writeSuccessResponse(c, http.StatusNoContent, "delete user successfully", nil)
}
