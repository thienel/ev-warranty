package handlers

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
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

// Create godoc
// @Summary Create a new user
// @Description Create a new user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param createUserRequest body dtos.CreateUserRequest true "User creation data"
// @Success 201 {object} dtos.APIResponse{data=dtos.UserDTO} "User created successfully"
// @Failure 400 {object} dtos.APIResponse "Bad request"
// @Failure 401 {object} dtos.APIResponse "Unauthorized"
// @Failure 403 {object} dtos.APIResponse "Forbidden"
// @Failure 500 {object} dtos.APIResponse "Internal server error"
// @Router /users [post]
func (h userHandler) Create(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

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

// Update godoc
// @Summary Update a user
// @Description Update an existing user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param updateUserRequest body dtos.UpdateUserRequest true "User update data"
// @Success 200 {object} dtos.APIResponse{data=dtos.UserDTO} "User updated successfully"
// @Failure 400 {object} dtos.APIResponse "Bad request"
// @Failure 401 {object} dtos.APIResponse "Unauthorized"
// @Failure 403 {object} dtos.APIResponse "Forbidden"
// @Failure 404 {object} dtos.APIResponse "User not found"
// @Failure 500 {object} dtos.APIResponse "Internal server error"
// @Router /users/{id} [put]
func (h userHandler) Update(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("update user request", "remote_addr", c.ClientIP())

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
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
	writeSuccessResponse(c, http.StatusNoContent, nil)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} dtos.APIResponse{data=dtos.UserDTO} "User retrieved successfully"
// @Failure 400 {object} dtos.APIResponse "Bad request"
// @Failure 401 {object} dtos.APIResponse "Unauthorized"
// @Failure 404 {object} dtos.APIResponse "User not found"
// @Failure 500 {object} dtos.APIResponse "Internal server error"
// @Router /users/{id} [get]
func (h userHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("get user request", "remote_addr", c.ClientIP())

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
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

// GetAll godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dtos.APIResponse{data=[]dtos.UserDTO} "Users retrieved successfully"
// @Failure 401 {object} dtos.APIResponse "Unauthorized"
// @Failure 500 {object} dtos.APIResponse "Internal server error"
// @Router /users [get]
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

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 204 {object} dtos.APIResponse "User deleted successfully"
// @Failure 400 {object} dtos.APIResponse "Bad request"
// @Failure 401 {object} dtos.APIResponse "Unauthorized"
// @Failure 403 {object} dtos.APIResponse "Forbidden"
// @Failure 404 {object} dtos.APIResponse "User not found"
// @Failure 500 {object} dtos.APIResponse "Internal server error"
// @Router /users/{id} [delete]
func (h userHandler) Delete(c *gin.Context) {
	if err := allowedRoles(c, entities.UserRoleAdmin); err != nil {
		handleError(h.log, c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	h.log.Info("delete user request", "remote_addr", c.ClientIP())

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		handleError(h.log, c, apperrors.NewInvalidUUID())
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
