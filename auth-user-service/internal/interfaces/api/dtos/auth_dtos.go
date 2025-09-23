package dtos

import (
	"auth-service/internal/domain/entities"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type UserDTO struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Email    string    `json:"email" binding:"required,email,max=100"`
	IsActive bool      `json:"is_active" binding:"required"`
}

func GenerateUserDTO(user entities.User) *UserDTO {
	return &UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}
