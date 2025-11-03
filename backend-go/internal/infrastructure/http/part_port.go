package http

import (
	"ev-warranty-go/internal/application/port"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/utils"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type partHTTPAdapter struct {
	baseURL string
	client  *http.Client
}

func NewPartHTTPAdapter(baseURL string, client *http.Client) port.PartPort {
	return &partHTTPAdapter{
		baseURL: baseURL,
		client:  client,
	}
}

func (a *partHTTPAdapter) ReserveByOfficeIDAndCategoryID(officeID, categoryID uuid.UUID) (*entity.Part, error) {
	request := struct {
		OfficeLocationID uuid.UUID `json:"office_location_id"`
		CategoryID       uuid.UUID `json:"category_id"`
	}{
		OfficeLocationID: officeID,
		CategoryID:       categoryID,
	}

	var respDto struct {
		IsSuccess bool        `json:"is_success"`
		Message   string      `json:"message"`
		Error     string      `json:"error"`
		Data      entity.Part `json:"data"`
	}

	err := utils.PostJSON(a.client, fmt.Sprintf("%s/parts/reserve", a.baseURL), request, &respDto)
	if err != nil {
		return nil, err
	}

	if !respDto.IsSuccess {
		return nil, apperror.NewInternalServerError(fmt.Errorf("reserve failed: %s", respDto.Message))
	}

	return &respDto.Data, nil
}

func (a *partHTTPAdapter) UnReserveByID(id uuid.UUID) error {
	request := struct{}{}
	var respDto struct {
		IsSuccess bool   `json:"is_success"`
		Message   string `json:"message"`
		Error     string `json:"error"`
	}

	err := utils.PostJSON(a.client, fmt.Sprintf("%s/parts/%s/unreserve",
		a.baseURL, id.String()), request, &respDto)
	if err != nil {
		return err
	}

	if !respDto.IsSuccess {
		return apperror.NewInternalServerError(fmt.Errorf("reserve failed: %s", respDto.Message))
	}

	return nil
}
