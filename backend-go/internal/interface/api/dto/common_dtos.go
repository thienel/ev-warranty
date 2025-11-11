package dto

type APIResponse struct {
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error,omitempty"`
	Data      any    `json:"data,omitempty"`
	Message   string `json:"message,omitempty"`
}
