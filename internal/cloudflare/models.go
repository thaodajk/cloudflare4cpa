package cloudflare

const ProviderName = "cloudflare-ai"

type ModelInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Description string `json:"description,omitempty"`
	MaxTokens   int    `json:"max_tokens,omitempty"`
	Free        bool   `json:"free,omitempty"`
}

var DefaultModels = []ModelInfo{
	{ID: "@cf/meta/llama-4-scout-17b-16e-instruct", Name: "Llama 4 Scout 17B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/meta/llama-3.3-70b-instruct-fp8-fast", Name: "Llama 3.3 70B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/meta/llama-3.2-11b-vision-instruct", Name: "Llama 3.2 11B Vision", Provider: ProviderName, MaxTokens: 4096, Free: true},
	{ID: "@cf/meta/llama-3.2-3b-instruct", Name: "Llama 3.2 3B", Provider: ProviderName, MaxTokens: 4096, Free: true},
	{ID: "@cf/meta/llama-3.1-8b-instruct-fp8", Name: "Llama 3.1 8B", Provider: ProviderName, MaxTokens: 4096, Free: true},
	{ID: "@cf/deepseek-ai/deepseek-r1-distill-qwen-32b", Name: "DeepSeek R1 Qwen 32B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/deepseek-ai/deepseek-v4-flash-0731", Name: "DeepSeek V4 Flash", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/deepseek-ai/deepseek-v4-pro-0813", Name: "DeepSeek V4 Pro", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/deepseek-ai/deepseek-v4.1-flash-1017", Name: "DeepSeek V4.1 Flash", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/qwen/qwen2.5-coder-32b-instruct", Name: "Qwen 2.5 Coder 32B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/qwen/qwen3-30b-a3b-fp8", Name: "Qwen 3 30B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/qwen/qwen3.8-27b", Name: "Qwen 3.8 27B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/qwen/qwq-32b", Name: "QwQ 32B Reasoning", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/openai/gpt-oss-120b", Name: "GPT OSS 120B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/openai/gpt-oss-20b", Name: "GPT OSS 20B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/nvidia/nemotron-3-120b-a12b", Name: "Nemotron 120B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/moonshotai/kimi-k2.6", Name: "Kimi K2.6", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/moonshotai/kimi-k2.7-code", Name: "Kimi K2.7 Code", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/mistralai/mistral-small-3.1-24b-instruct", Name: "Mistral Small 3.1 24B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/google/gemma-4-26b-a4b-it", Name: "Gemma 4 26B", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/zai-org/glm-5.3", Name: "GLM 5.3", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/zai-org/glm-5.3-flash", Name: "GLM 5.3 Flash", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/ibm-granite/granite-4.0-h-micro", Name: "Granite 4.0 Micro", Provider: ProviderName, MaxTokens: 4096, Free: true},
	{ID: "@cf/cloudflare/clef", Name: "Cloudflare Clef", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/cloudflare/clef-flash", Name: "Clef Flash", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/cloudflare/clef-omni", Name: "Clef Omni", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/cloudflare/clairic", Name: "Clairic", Provider: ProviderName, MaxTokens: 8192, Free: true},
	{ID: "@cf/thinkingmachines/inkling-2-small", Name: "Inkling 2 Small", Provider: ProviderName, MaxTokens: 4096, Free: true},
}

func GetModelInfoMaps() []map[string]any {
	var out []map[string]any
	for _, m := range DefaultModels {
		out = append(out, map[string]any{
			"id":         m.ID,
			"name":       m.Name,
			"provider":   m.Provider,
			"max_tokens": m.MaxTokens,
		})
	}
	return out
}
