package utils

import (
	"bytes"
	"encoding/json"
	"ev-warranty-go/pkg/apperror"
	"fmt"
	"io"
	"net/http"
)

func PostJSON(client *http.Client, url string, payload any, respData any) (int, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, apperror.ErrExternalServiceMarshalRequest.WithError(err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return 0, apperror.ErrExternalServiceUnavailable.WithError(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, apperror.NewExternalServiceRequestFailed(resp.StatusCode,
			fmt.Errorf("request to %s failed", url))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, apperror.ErrExternalServiceReadResponse.WithError(err)
	}

	err = json.Unmarshal(data, respData)
	if err != nil {
		return 0, apperror.ErrExternalServiceUnmarshalResponse.WithError(err)
	}

	return resp.StatusCode, nil
}
