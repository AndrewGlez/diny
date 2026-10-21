package commit

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/llm"
)

func CreateCommitMessage(gitDiff string, cfg *config.Config) (string, error) {
	commitMessage, err := llm.CreateCommitMessageWithLLM(gitDiff, cfg)

	if err != nil {
		return "", err
	}

	return commitMessage, nil
}
