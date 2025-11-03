package http

import (
	"encoding/json"
	"ev-warranty-go/internal/application/port"
	"ev-warranty-go/internal/domain/entity"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type partHTTPAdapter struct {
	baseURL string
	client  *http.Client
}

func NewPartHTTPAdapter() port.PartPort {
	return &partHTTPAdapter{}
}

func (a *partHTTPAdapter) FindByOfficeIDAndCategoryID(office, category uuid.UUID) (*entity.Part, error) {
	resp, _ := a.client.Get(fmt.Sprintf("%s/parts/%s%s", a.baseURL, office, category))
	var dto struct {
		IsSuccess bool        `json:"is_success"`
		Message   string      `json:"message"`
		Error     string      `json:"error"`
		Data      entity.Part `json:"data"`
	}

	err := json.NewDecoder(resp.Body).Decode(&dto)
	if err != nil {
		return nil, err
	}

	return &dto.Data, nil
}

func (a *partHTTPAdapter) UpdateStatus(id uuid.UUID, status string) error {
	resp, _ := a.client.Get(fmt.Sprintf("%s/users/%s", a.baseURL, id))
	var dto struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&dto)
	return nil
}
