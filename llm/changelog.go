package llm

import "github.com/dinoDanic/diny/config"

func CreateChangelogWithLLM(prompt string, cfg *config.Config) (string, error) {
	return CreateTimelineWithLLM(prompt, cfg)
}
