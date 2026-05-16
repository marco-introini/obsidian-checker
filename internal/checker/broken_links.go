package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marco-introini/obsidian-checker/internal/parser"
	"github.com/marco-introini/obsidian-checker/internal/resolver"
	"github.com/marco-introini/obsidian-checker/internal/vault"
)

type BrokenLinksChecker struct {
	CaseInsensitive bool
	CheckHeadings   bool
}

func NewBrokenLinksChecker(caseSensitive, checkHeadings bool) *BrokenLinksChecker {
	return &BrokenLinksChecker{
		CaseInsensitive: !caseSensitive,
		CheckHeadings:   checkHeadings,
	}
}

func (c *BrokenLinksChecker) Name() string {
	return "broken-links"
}

func (c *BrokenLinksChecker) Check(v *vault.Vault) ([]Issue, Summary, error) {
	var issues []Issue
	res := resolver.New(v, c.CaseInsensitive)
	totalLinks := 0
	fileCount := 0

	for _, note := range v.Notes {
		fileCount++
		content, err := os.ReadFile(note.Path)
		if err != nil {
			return nil, Summary{}, fmt.Errorf("error reading %s: %w", note.Path, err)
		}

		links := parser.ParseContent(string(content))
		totalLinks += len(links)

		for _, link := range links {
			resolved := res.Resolve(link, note.RelPath)

			if !resolved.Resolved {
				isImage := isImageLink(link.Target)

				issue := Issue{
					File:    note.RelPath,
					Line:    link.Line,
					RawLink: link.Raw,
					Target:  link.Target,
				}

				if isImage {
					issue.Code = CodeFileNotFound
					issue.Message = fmt.Sprintf("File '%s' non trovato", link.Target)
				} else {
					issue.Code = CodeNoteNotFound
					issue.Message = fmt.Sprintf("Nota '%s' non trovata nel vault", link.Target)
				}

				issues = append(issues, issue)
				continue
			}

			if c.CheckHeadings && resolved.TargetNote != nil && link.Heading != "" {
				if !c.headingExists(resolved.TargetNote, link.Heading) {
					issues = append(issues, Issue{
						File:    note.RelPath,
						Line:    link.Line,
						RawLink: link.Raw,
						Target:  link.Target + "#" + link.Heading,
						Code:    "heading_not_found",
						Message: fmt.Sprintf("Heading '#%s' non trovato in '%s'", link.Heading, resolved.TargetNote.RelPath),
					})
				}
			}
		}
	}

	return issues, Summary{
		TotalFiles: fileCount,
		TotalLinks: totalLinks,
		IssueCount: len(issues),
	}, nil
}

func (c *BrokenLinksChecker) headingExists(note *vault.Note, heading string) bool {
	content, err := os.ReadFile(note.Path)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	normalized := strings.ToLower(strings.TrimSpace(heading))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			hText := strings.TrimLeft(trimmed, "#")
			hText = strings.TrimSpace(hText)
			if strings.ToLower(hText) == normalized {
				return true
			}
		}
	}
	return false
}

func isImageLink(target string) bool {
	ext := strings.ToLower(filepath.Ext(target))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".bmp":
		return true
	}
	return false
}
