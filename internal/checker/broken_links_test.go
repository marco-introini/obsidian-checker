package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-introini/obsidian-checker/internal/vault"
)

func TestBrokenLinksChecker(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "Esistente.md"), []byte("# Esistente\n"), 0644)
	os.WriteFile(filepath.Join(dir, "ConLink.md"), []byte("# ConLink\n\nVedi [[Esistente]] e [[Mancante]].\n"), 0644)
	os.WriteFile(filepath.Join(dir, "ConEmbed.md"), []byte("# ConEmbed\n![[foto.png]]\n"), 0644)
	os.WriteFile(filepath.Join(dir, "foto.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, "ConHeading.md"), []byte("# ConHeading\n\n## Intro\n\nTesto.\n\n[[ConHeading#Intro]]\n[[ConHeading#NonEsiste]]\n"), 0644)

	v, err := vault.Scan(dir, []string{".obsidian", ".trash"}, nil)
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}

	t.Run("without heading check", func(t *testing.T) {
		c := NewBrokenLinksChecker(true, false)
		issues, s, err := c.Check(v)
		if err != nil {
			t.Fatalf("Check error: %v", err)
		}

		if s.TotalFiles != 4 {
			t.Errorf("expected 4 files, got %d", s.TotalFiles)
		}

		if s.IssueCount != 1 {
			t.Errorf("expected 1 broken link, got %d", s.IssueCount)
			for _, iss := range issues {
				t.Logf("Issue: %+v", iss)
			}
		}

		if len(issues) > 0 && issues[0].Target != "Mancante" {
			t.Errorf("expected 'Mancante', got %q", issues[0].Target)
		}
	})

	t.Run("with heading check", func(t *testing.T) {
		c := NewBrokenLinksChecker(true, true)
		issues, _, err := c.Check(v)
		if err != nil {
			t.Fatalf("Check error: %v", err)
		}

		headingIssues := 0
		for _, iss := range issues {
			if iss.Code == "heading_not_found" {
				headingIssues++
			}
		}

		if headingIssues != 1 {
			t.Errorf("expected 1 heading issue, got %d", headingIssues)
			for _, iss := range issues {
				t.Logf("Issue: %+v", iss)
			}
		}
	})
}
