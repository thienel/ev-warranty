package dtos

type CreateOfficeRequest struct {
	OfficeName string `json:"office_name" binding:"required"`
	OfficeType string `json:"office_type" binding:"required"`
	IsActive   bool   `json:"is_active" binding:"required"`
	Address    string `json:"address" binding:"required"`
}

type UpdateOfficeRequest struct {
	OfficeName *string `json:"office_name"`
	OfficeType *string `json:"office_type"`
	Address    *string `json:"address"`
	IsActive   *bool   `json:"is_active"`
}
