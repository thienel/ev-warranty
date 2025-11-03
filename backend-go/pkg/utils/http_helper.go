package utils

import (
	"bytes"
	"encoding/json"
	"ev-warranty-go/pkg/apperror"
	"fmt"
	"io"
	"net/http"
)

func PostJSON(client *http.Client, url string, payload any, respData any) error {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return apperror.NewInternalServerError(err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return apperror.NewInternalServerError(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return apperror.NewInternalServerError(fmt.Errorf("request failed with status %d", resp.StatusCode))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return apperror.NewInternalServerError(err)
	}

	err = json.Unmarshal(data, respData)
	if err != nil {
		return apperror.NewInternalServerError(err)
	}

	return nil
}
