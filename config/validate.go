package config

import (
	"fmt"
	"slices"
)

func (c *Config) Validate() error {
	if c.Theme == "" {
		return fmt.Errorf("theme is required")
	}

	validTones := []Tone{Professional, Casual, Friendly}
	if c.Commit.Tone == "" {
		return fmt.Errorf("tone is required")
	}
	if !slices.Contains(validTones, c.Commit.Tone) {
		return fmt.Errorf("invalid tone '%s', must be one of: professional, casual, friendly", c.Commit.Tone)
	}

	validLengths := []Length{Short, Normal, Long}
	if c.Commit.Length == "" {
		return fmt.Errorf("length is required")
	}
	if !slices.Contains(validLengths, c.Commit.Length) {
		return fmt.Errorf("invalid length '%s', must be one of: short, normal, long", c.Commit.Length)
	}
	if c.Provider.BaseURL == "" {
		return fmt.Errorf("provider.base_url is required")
	}
	if c.Provider.Model == "" {
		return fmt.Errorf("provider.model is required")
	}
	if c.Provider.TimeoutSeconds <= 0 {
		return fmt.Errorf("provider.timeout_seconds must be positive")
	}

	return nil
}
