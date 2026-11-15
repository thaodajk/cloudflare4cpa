package cloudflare

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudflare4cpa/internal/wire"
	"github.com/router-for-me/CLIProxyAPI/v8/sdk/pluginabi"
	"github.com/router-for-me/CLIProxyAPI/v8/sdk/pluginapi"
)

const Version = "1.0.0"

type registration struct {
	SchemaVersion uint32                 `json:"schema_version"`
	Metadata      pluginapi.Metadata     `json:"metadata"`
	Capabilities  registrationCapability `json:"capabilities"`
}

type registrationCapability struct {
	ModelRegistrar        bool                         `json:"model_registrar"`
	ModelProvider         bool                         `json:"model_provider"`
	AuthProvider          bool                         `json:"auth_provider"`
	Executor              bool                         `json:"executor"`
	ExecutorModelScope    pluginapi.ExecutorModelScope `json:"executor_model_scope"`
	ExecutorInputFormats  []string                     `json:"executor_input_formats,omitempty"`
	ExecutorOutputFormats []string                     `json:"executor_output_formats,omitempty"`
	QuotaProvider         bool                         `json:"quota_provider"`
	ManagementAPI         bool                         `json:"management_api"`
}

type identifierResponse struct {
	Identifier string `json:"identifier"`
}

func getRegistration() registration {
	return registration{
		SchemaVersion: pluginabi.SchemaVersion,
		Metadata: pluginapi.Metadata{
			Name:             "Cloudflare Workers AI",
			Version:          Version,
			Author:           "thaodajk",
			GitHubRepository: "https://github.com/thaodajk/cloudflare4cpa",
			ConfigFields: []pluginapi.ConfigField{
				{Name: "enabled", Type: pluginapi.ConfigFieldTypeBoolean, Description: "Enable Cloudflare Workers AI plugin"},
			},
		},
		Capabilities: registrationCapability{
			ModelRegistrar:        true,
			ModelProvider:         true,
			AuthProvider:          true,
			Executor:              true,
			ExecutorModelScope:    pluginapi.ExecutorModelScopeOAuth,
			ExecutorInputFormats:  []string{"openai"},
			ExecutorOutputFormats: []string{"openai"},
			QuotaProvider:         false,
			ManagementAPI:         false,
		},
	}
}

func HandleMethod(method string, request []byte) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("cloudflare-ai: %s panic: %v", method, r)
			out, err = nil, fmt.Errorf("cloudflare-ai internal error in %s: %v", method, r)
		}
	}()

	switch method {
	case pluginabi.MethodPluginRegister, pluginabi.MethodPluginReconfigure:
		return wire.OK(getRegistration()), nil

	case pluginabi.MethodPluginQuiesce, pluginabi.MethodPluginShutdown:
		return wire.OK(map[string]bool{"acknowledged": true}), nil

	case pluginabi.MethodModelRegister:
		return wire.OK(map[string]any{
			"provider": ProviderName,
			"models":   GetModelInfoMaps(),
		}), nil

	case pluginabi.MethodModelStatic:
		return wire.OK(map[string]any{
			"provider": ProviderName,
			"models":   GetModelInfoMaps(),
		}), nil

	case pluginabi.MethodModelForAuth:
		return handleModelForAuth(request)

	case pluginabi.MethodAuthIdentifier:
		return wire.OK(identifierResponse{Identifier: ProviderName}), nil

	case pluginabi.MethodAuthParse:
		return handleAuthParse(request)

	case pluginabi.MethodAuthRefresh:
		return handleAuthRefresh(request)

	case pluginabi.MethodExecutorIdentifier:
		return wire.OK(identifierResponse{Identifier: ProviderName}), nil

	case pluginabi.MethodExecutorExecute:
		return handleExecute(request)

	case pluginabi.MethodExecutorExecuteStream:
		return handleExecuteStream(request)

	default:
		return wire.Error("unsupported_method", "method not supported: "+method), nil
	}
}

func handleModelForAuth(request []byte) ([]byte, error) {
	return wire.OK(map[string]any{
		"provider": ProviderName,
		"models":   GetModelInfoMaps(),
	}), nil
}

func handleAuthParse(request []byte) ([]byte, error) {
	var req pluginapi.AuthParseRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, err
	}
	if req.Provider != ProviderName {
		return wire.OK(pluginapi.AuthParseResponse{Handled: false}), nil
	}

	cred, err := ParseCredential(req.RawJSON)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", req.FileName, err)
	}

	return wire.OK(pluginapi.AuthParseResponse{
		Handled: true,
		Auth: pluginapi.AuthData{
			Provider:         ProviderName,
			FileName:         req.FileName,
			Label:            AuthLabel(cred, req.FileName),
			StorageJSON:      req.RawJSON,
			Metadata:         AuthMetadata(cred),
			NextRefreshAfter: NextRefresh(),
		},
	}), nil
}

func handleAuthRefresh(request []byte) ([]byte, error) {
	var req pluginapi.AuthRefreshRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, err
	}

	return wire.OK(pluginapi.AuthRefreshResponse{
		Auth: pluginapi.AuthData{
			Provider:         ProviderName,
			ID:               req.AuthID,
			StorageJSON:      req.StorageJSON,
			Metadata:         map[string]any{"refresh_interval_seconds": 86400},
			NextRefreshAfter: NextRefresh(),
		},
		NextRefreshAfter: NextRefresh(),
	}), nil
}

type executeRequest struct {
	AuthStorageJSON json.RawMessage `json:"auth_storage_json"`
	Model           string          `json:"model"`
	Body            json.RawMessage `json:"body"`
}

func handleExecute(request []byte) ([]byte, error) {
	var req executeRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, err
	}

	cred, err := ParseCredential(req.AuthStorageJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid credential: %w", err)
	}

	respBody, err := ExecuteChat(cred, req.Model, req.Body)
	if err != nil {
		return wire.ErrorStatus("upstream_error", err.Error(), 502, true), nil
	}

	return wire.OK(map[string]any{
		"status_code":  200,
		"body":         json.RawMessage(respBody),
		"content_type": "application/json",
	}), nil
}

func handleExecuteStream(request []byte) ([]byte, error) {
	var req executeRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, err
	}

	cred, err := ParseCredential(req.AuthStorageJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid credential: %w", err)
	}

	respBody, err := ExecuteChat(cred, req.Model, req.Body)
	if err != nil {
		return wire.ErrorStatus("upstream_error", err.Error(), 502, true), nil
	}

	return wire.OK(map[string]any{
		"status_code":  200,
		"body":         json.RawMessage(respBody),
		"content_type": "application/json",
	}), nil
}
