package external

import (
	"bytes"
	"context"
	"encoding/json"
	"ev-warranty-go/pkg/apperror"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreateWorkOrderRequest struct {
	ClaimID              uuid.UUID `json:"claim_id"`
	AssignedTechnicianID uuid.UUID `json:"assigned_technician_id"`
}

type WorkOrderResponse struct {
	ID                   uuid.UUID `json:"id"`
	ClaimID              uuid.UUID `json:"claim_id"`
	AssignedTechnicianID uuid.UUID `json:"assigned_technician_id"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
}

type BaseResponseDto struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type WorkOrderResponseDto struct {
	BaseResponseDto
	Data *WorkOrderResponse `json:"data,omitempty"`
}

type DotnetClient interface {
	CreateWorkOrder(ctx context.Context, req *CreateWorkOrderRequest, authToken string) (*WorkOrderResponse, error)
}

type dotnetClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDotnetClient(baseURL string) DotnetClient {
	return &dotnetClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *dotnetClient) CreateWorkOrder(ctx context.Context, req *CreateWorkOrderRequest, authToken string) (*WorkOrderResponse, error) {
	url := fmt.Sprintf("%s/work-orders", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, apperror.ErrInternal.WithMessage("Failed to marshal request: " + err.Error())
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, apperror.ErrInternal.WithMessage("Failed to create request: " + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Add Authorization header if provided
	if authToken != "" {
		httpReq.Header.Set("Authorization", authToken)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, apperror.ErrExternalServiceUnavailable.WithMessage("Failed to call .NET backend: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperror.ErrInternal.WithMessage("Failed to read response body: " + err.Error())
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse error response
		var errorResp BaseResponseDto
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, apperror.ErrExternalServiceError.WithMessage(
				fmt.Sprintf("Backend .NET returned error (status: %d): %s", resp.StatusCode, string(body)))
		}

		return nil, apperror.ErrExternalServiceError.WithMessage(
			fmt.Sprintf("Backend .NET error (status: %d): %s", resp.StatusCode, errorResp.Message))
	}

	var successResp WorkOrderResponseDto
	if err := json.Unmarshal(body, &successResp); err != nil {
		return nil, apperror.ErrInternal.WithMessage("Failed to parse success response: " + err.Error())
	}

	if !successResp.Success || successResp.Data == nil {
		return nil, apperror.ErrExternalServiceError.WithMessage("Backend .NET returned invalid response")
	}

	return successResp.Data, nil
}
