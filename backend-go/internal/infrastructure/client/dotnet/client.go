package dotnet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	defaultTimeout = 30 * time.Second
)

type Client interface {
	ReservePart(ctx context.Context, officeLocationID, categoryID uuid.UUID) (*PartResponse, error)
	UnreservePart(ctx context.Context, partID uuid.UUID) error
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) Client {
	return &client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *client) ReservePart(ctx context.Context, officeLocationID, categoryID uuid.UUID) (*PartResponse, error) {
	url := fmt.Sprintf("%s/api/v1/parts/reserve", c.baseURL)

	reqBody := ReservePartRequest{
		OfficeLocationID: officeLocationID,
		CategoryID:       categoryID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response BaseDataResponse[PartResponse]
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.IsSuccess || resp.StatusCode != http.StatusOK {
		if response.Message != "" {
			return nil, fmt.Errorf("failed to reserve part: %s (code: %s)", response.Message, response.ErrorCode)
		}
		return nil, fmt.Errorf("failed to reserve part: unexpected status code %d", resp.StatusCode)
	}

	if response.Data == nil {
		return nil, fmt.Errorf("no part data in response")
	}

	return response.Data, nil
}

func (c *client) UnreservePart(ctx context.Context, partID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/parts/%s/unreserve", c.baseURL, partID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var response BaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.IsSuccess {
		if response.Message != "" {
			return fmt.Errorf("failed to unreserve part: %s (code: %s)", response.Message, response.ErrorCode)
		}
		return fmt.Errorf("failed to unreserve part: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
