// Package llm builds the prompts for each request type and sends them through
// the configured provider.
package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/provider"
)

type SplitGroup struct {
	Order   int      `json:"order"`
	Type    string   `json:"type"`
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

type responseData struct {
	Message string       `json:"message"`
	Groups  []SplitGroup `json:"groups"`
}

type RequestExtras struct {
	PreviousPlans [][]SplitGroup
	Feedback      string
}

const genericSystemPrompt = "Follow the instructions in the user message exactly. Reply with the requested text only."

// doRequest sends one request type to the provider from cfg.Provider.
func doRequest(reqType string, userPrompt string, cfg *config.Config, extras *RequestExtras) (*responseData, error) {
	if cfg == nil {
		return nil, fmt.Errorf("missing config")
	}
	switch reqType {
	case "commit":
		text, err := provider.Chat(cfg.Provider, commitSystemPrompt(cfg), userPrompt, false)
		if err != nil {
			return nil, err
		}
		return &responseData{Message: stripFence(text)}, nil
	case "split":
		text, err := provider.Chat(cfg.Provider, splitSystemPrompt(cfg), splitUserPrompt(userPrompt, extras), true)
		if err != nil {
			return nil, err
		}
		var out responseData
		if err := json.Unmarshal([]byte(stripFence(text)), &out); err != nil {
			return nil, fmt.Errorf("decode split plan: %w", err)
		}
		return &out, nil
	default:
		text, err := provider.Chat(cfg.Provider, genericSystemPrompt, userPrompt, false)
		if err != nil {
			return nil, err
		}
		return &responseData{Message: stripFence(text)}, nil
	}
}

// stripFence removes a surrounding Markdown code fence that models sometimes add.
func stripFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[i+1:]
	}
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "```"))
}

func sendRequest(reqType string, userPrompt string, cfg *config.Config) (string, error) {
	data, err := doRequest(reqType, userPrompt, cfg, nil)
	if err != nil {
		return "", err
	}
	if data.Message == "" {
		return "", fmt.Errorf("empty message from provider")
	}
	return data.Message, nil
}

func CreateCommitMessageWithLLM(gitDiff string, cfg *config.Config) (string, error) {
	return sendRequest("commit", gitDiff, cfg)
}

func CreateTimelineWithLLM(prompt string, cfg *config.Config) (string, error) {
	return sendRequest("timeline", prompt, cfg)
}

func CreateSplitPlanWithLLM(gitDiff string, cfg *config.Config, extras *RequestExtras) ([]SplitGroup, error) {
	data, err := doRequest("split", gitDiff, cfg, extras)
	if err != nil {
		return nil, err
	}
	if len(data.Groups) == 0 {
		return nil, fmt.Errorf("empty split plan from provider")
	}
	return data.Groups, nil
}
