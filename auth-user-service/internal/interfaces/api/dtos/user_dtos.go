package dtos

import (
	"auth-service/internal/domain/entities"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required,email,max=100"`
	Password string    `json:"password" binding:"required"`
	Role     string    `json:"role" binding:"required,max=20"`
	IsActive bool      `json:"is_active" binding:"required"`
	OfficeID uuid.UUID `json:"office_id" binding:"required"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID  `json:"id" binding:"required"`
	Name     *string    `json:"name"`
	Email    *string    `json:"email" binding:"email,max=100"`
	Role     *string    `json:"role" binding:"max=20"`
	IsActive *bool      `json:"is_active"`
	OfficeID *uuid.UUID `json:"office_id"`
}

type UserDTO struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required,email,max=100"`
	Role     string    `json:"role" binding:"required,max=20"`
	IsActive bool      `json:"is_active" binding:"required"`
	OfficeID uuid.UUID `json:"office_id" binding:"required"`
}

func GenerateUserDTO(user entities.User) *UserDTO {
	return &UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		OfficeID: user.OfficeID,
	}
}
