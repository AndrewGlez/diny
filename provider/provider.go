// Package provider calls an OpenAI-compatible chat completions endpoint.
package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dinoDanic/diny/config"
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []message       `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Chat sends one system+user exchange and returns the reply text. jsonMode asks
// for a JSON object when the provider has json_mode enabled in config.
func Chat(cfg config.ProviderConfig, system, user string, jsonMode bool) (string, error) {
	var key string
	if cfg.APIKeyEnv != "" {
		key = strings.TrimSpace(os.Getenv(cfg.APIKeyEnv))
		if key == "" {
			return "", fmt.Errorf("API key env var %s is not set (provider.api_key_env)", cfg.APIKeyEnv)
		}
	}

	payload := chatRequest{
		Model: cfg.Model,
		Messages: []message{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	}
	if jsonMode && cfg.JSONMode {
		payload.ResponseFormat = &responseFormat{Type: "json_object"}
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	client := &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call %s: %w", cfg.BaseURL, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var out chatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response (HTTP %d): %w", res.StatusCode, err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("provider error: %s", out.Error.Message)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("provider returned HTTP %d", res.StatusCode)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("provider returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}
