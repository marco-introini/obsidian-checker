package output

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/marco-introini/obsidian-checker/internal/checker"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(issues []checker.Issue, summary checker.Summary) (string, error) {
	var sb strings.Builder

	w := csv.NewWriter(&sb)

	w.Write([]string{"file", "line", "link", "target", "issue", "message"})

	for _, issue := range issues {
		w.Write([]string{
			issue.File,
			fmt.Sprintf("%d", issue.Line),
			issue.RawLink,
			issue.Target,
			string(issue.Code),
			issue.Message,
		})
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	sb.WriteString(fmt.Sprintf("\n# Riepilogo: %d file, %d link, %d rotti\n", summary.TotalFiles, summary.TotalLinks, summary.IssueCount))

	return sb.String(), nil
}
