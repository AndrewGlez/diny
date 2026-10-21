package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dinoDanic/diny/config"
)

var toneRules = map[config.Tone]string{
	config.Professional: "formal and precise",
	config.Casual:       "relaxed and conversational",
	config.Friendly:     "warm and approachable",
}

var lengthRules = map[config.Length]string{
	config.Short:  "subject only, 1-2 sentences max, at most 60 characters, imperative verb first, no bullets",
	config.Normal: "subject at most 70 characters (imperative); if needed, add 1-3 terse bullets for WHY/impact",
	config.Long:   "subject at most 80 characters (imperative); then 2-6 terse bullets for context/impact",
}

func commitStyle(cfg *config.Config) config.CommitConfig {
	if cfg == nil {
		return config.CommitConfig{}
	}
	return cfg.Commit
}

// styleRules renders the commit message settings from the diny config file.
func styleRules(c config.CommitConfig) string {
	var b strings.Builder
	if c.Conventional {
		b.WriteString("Use Conventional Commits: type(optional scope): subject. Common types: feat, fix, refactor, docs, test, build, ci, chore, perf, style.\n")
	} else {
		b.WriteString("Write a plain imperative subject line with no type prefix.\n")
	}
	if c.Emoji {
		b.WriteString("Start the subject with one fitting emoji.\n")
	}
	if t, ok := toneRules[c.Tone]; ok {
		fmt.Fprintf(&b, "Tone: %s.\n", t)
	}
	if l, ok := lengthRules[c.Length]; ok {
		fmt.Fprintf(&b, "Length: %s.\n", l)
	}
	if s := strings.TrimSpace(c.CustomInstructions); s != "" {
		fmt.Fprintf(&b, "Additional instructions: %s\n", s)
	}
	return b.String()
}

func commitSystemPrompt(cfg *config.Config) string {
	return "You write git commit messages from staged changes.\n" +
		"Reply with the commit message only: no preamble, no code fences, no quotes.\n" +
		styleRules(commitStyle(cfg))
}

func splitSystemPrompt(cfg *config.Config) string {
	return "You split staged changes into an ordered series of commits.\n" +
		`Reply with a JSON object only: {"groups":[{"order":1,"type":"feat","message":"...","files":["path"]}]}` + "\n" +
		"Rules:\n" +
		"- Every staged file appears in exactly one group.\n" +
		"- Use file paths exactly as they appear in the diff headers. Do not add a/ or b/ prefixes.\n" +
		"- Keep related changes together. Order groups so earlier commits do not depend on later ones.\n" +
		"- Each \"message\" is a complete commit message that follows these rules:\n" +
		styleRules(commitStyle(cfg))
}

func splitUserPrompt(diff string, extras *RequestExtras) string {
	var b strings.Builder
	b.WriteString(diff)
	if extras == nil {
		return b.String()
	}
	if len(extras.PreviousPlans) > 0 {
		plans, _ := json.Marshal(extras.PreviousPlans)
		fmt.Fprintf(&b, "\n\nPrevious plans were rejected. Fix them: %s", plans)
	}
	if fb := strings.TrimSpace(extras.Feedback); fb != "" {
		fmt.Fprintf(&b, "\n\nUser feedback: %s", fb)
	}
	return b.String()
}
