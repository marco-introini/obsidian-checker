package output

import "github.com/marco-introini/obsidian-checker/internal/checker"

type Formatter interface {
	Format(issues []checker.Issue, summary checker.Summary) (string, error)
}
