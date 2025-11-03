package dto

import (
	"ev-warranty-go/internal/domain/entities"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required,email,max=100"`
	Password string    `json:"password" binding:"required"`
	Role     string    `json:"role" binding:"required,max=20"`
	IsActive bool      `json:"is_active"`
	OfficeID uuid.UUID `json:"office_id" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
	OfficeID uuid.UUID `json:"office_id"`
}

type UserDTO struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required,email,max=100"`
	Role     string    `json:"role" binding:"required,max=20"`
	IsActive bool      `json:"is_active"`
	OfficeID uuid.UUID `json:"office_id" binding:"required"`
}

func GenerateUserDTO(user *entities.User) *UserDTO {
	return &UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		OfficeID: user.OfficeID,
	}
}

func GenerateUserDTOList(user []*entities.User) []*UserDTO {
	users := make([]*UserDTO, len(user))
	for i, u := range user {
		users[i] = GenerateUserDTO(u)
	}
	return users
}
