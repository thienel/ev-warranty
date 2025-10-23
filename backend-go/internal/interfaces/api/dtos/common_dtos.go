package dtos

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data any `json:"data,omitempty"`
}
