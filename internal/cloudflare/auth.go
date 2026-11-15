package cloudflare

import (
	"encoding/json"
	"fmt"
	"time"
)

type Credential struct {
	Type      string `json:"type"`
	AccountID string `json:"account_id"`
	APIToken  string `json:"api_token"`
	Email     string `json:"email,omitempty"`
	Label     string `json:"label,omitempty"`
}

func ParseCredential(raw json.RawMessage) (*Credential, error) {
	var cred Credential
	if err := json.Unmarshal(raw, &cred); err != nil {
		return nil, err
	}
	if cred.Type != ProviderName {
		return nil, fmt.Errorf("unexpected type %q", cred.Type)
	}
	if cred.AccountID == "" {
		return nil, fmt.Errorf("account_id is required")
	}
	if cred.APIToken == "" {
		return nil, fmt.Errorf("api_token is required")
	}
	return &cred, nil
}

func AuthLabel(cred *Credential, fileName string) string {
	if cred.Label != "" {
		return cred.Label
	}
	if cred.Email != "" {
		return cred.Email
	}
	return fileName
}

func AuthMetadata(cred *Credential) map[string]any {
	m := map[string]any{
		"refresh_interval_seconds": 86400,
	}
	if cred.Email != "" {
		m["email"] = cred.Email
	}
	return m
}

func NextRefresh() time.Time {
	return time.Now().Add(24 * time.Hour)
}
