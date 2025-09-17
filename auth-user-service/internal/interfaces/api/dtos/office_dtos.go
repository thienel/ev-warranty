package dtos

type CreateOfficeRequest struct {
	OfficeName string `json:"office_name"`
	OfficeType string `json:"office_type"`
	IsActive   bool   `json:"is_active"`
	Address    string `json:"address"`
}

type UpdateOfficeRequest struct {
	OfficeName string `json:"office_name"`
	OfficeType string `json:"office_type"`
	Address    string `json:"address"`
}
