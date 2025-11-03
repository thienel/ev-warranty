package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid bool    `json:"valid"`
	User  UserDTO `json:"user,omitempty"`
}
