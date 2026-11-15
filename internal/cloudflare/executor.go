package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 120 * time.Second}

const cfBaseURL = "https://api.cloudflare.com/client/v4/accounts"

func ExecuteChat(cred *Credential, model string, body json.RawMessage) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/ai/v1/chat/completions", cfBaseURL, cred.AccountID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cred.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudflare request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("cloudflare returned %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func ExecuteStream(cred *Credential, model string, body json.RawMessage) (io.ReadCloser, string, error) {
	url := fmt.Sprintf("%s/%s/ai/v1/chat/completions", cfBaseURL, cred.AccountID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+cred.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("cloudflare request failed: %w", err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, "", fmt.Errorf("cloudflare returned %d: %s", resp.StatusCode, string(body))
	}

	return resp.Body, resp.Header.Get("Content-Type"), nil
}
