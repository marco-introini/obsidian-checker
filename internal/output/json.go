package output

import (
	"encoding/json"

	"github.com/marco-introini/obsidian-checker/internal/checker"
)

type JSONFormatter struct {
	VaultPath string
}

type jsonOutput struct {
	Vault   string          `json:"vault"`
	Check   string          `json:"check"`
	Issues  []checker.Issue `json:"issues"`
	Summary checker.Summary `json:"summary"`
}

func (f *JSONFormatter) Format(issues []checker.Issue, summary checker.Summary) (string, error) {
	out := jsonOutput{
		Vault:   f.VaultPath,
		Check:   "broken-links",
		Issues:  issues,
		Summary: summary,
	}

	if out.Issues == nil {
		out.Issues = []checker.Issue{}
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
