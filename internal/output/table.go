package output

import (
	"fmt"
	"strings"

	"github.com/marco-introini/obsidian-checker/internal/checker"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(issues []checker.Issue, summary checker.Summary) (string, error) {
	var sb strings.Builder

	if len(issues) == 0 {
		sb.WriteString("\nNessun link rotto trovato.\n")
		sb.WriteString(fmt.Sprintf("\nFile analizzati: %d, Link controllati: %d\n", summary.TotalFiles, summary.TotalLinks))
		return sb.String(), nil
	}

	sb.WriteString(fmt.Sprintf("\n%-4s  %-26s  %-24s  %s\n", "N.", "File", "Link", "Messaggio"))
	sb.WriteString(fmt.Sprintf("%s  %s  %s  %s\n",
		strings.Repeat("-", 4),
		strings.Repeat("-", 26),
		strings.Repeat("-", 24),
		strings.Repeat("-", 30),
	))

	for i, issue := range issues {
		file := truncate(issue.File, 26)
		link := truncate(issue.RawLink, 24)
		sb.WriteString(fmt.Sprintf("%-4d  %-26s  %-24s  %s\n", i+1, file, link, issue.Message))
	}

	sb.WriteString(fmt.Sprintf("\nRiepilogo: %d file, %d link, %d rotti\n",
		summary.TotalFiles, summary.TotalLinks, summary.IssueCount))

	return sb.String(), nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
