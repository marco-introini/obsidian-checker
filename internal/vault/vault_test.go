package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.MkdirAll(filepath.Join(dir, ".obsidian"), 0755)
	os.MkdirAll(filepath.Join(dir, ".trash"), 0755)

	os.WriteFile(filepath.Join(dir, "Home.md"), []byte("# Home"), 0644)
	os.WriteFile(filepath.Join(dir, "subdir", "Nota.md"), []byte("# Nota"), 0644)
	os.WriteFile(filepath.Join(dir, "foto.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, ".obsidian", "config.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(dir, ".trash", "Old.md"), []byte("# Old"), 0644)

	v, err := Scan(dir, []string{".obsidian", ".trash"}, nil)
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}

	if len(v.Notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(v.Notes))
	}

	if len(v.Assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(v.Assets))
	}

	if _, ok := v.Notes["Home.md"]; !ok {
		t.Error("expected Home.md in vault")
	}

	if _, ok := v.Notes["subdir/Nota.md"]; !ok {
		t.Error("expected subdir/Nota.md in vault")
	}

	if _, ok := v.Assets["foto.png"]; !ok {
		t.Error("expected foto.png in vault")
	}
}

func TestResolveNote(t *testing.T) {
	v := New("/vault")
	v.AddNote("Home.md")
	v.AddNote("subdir/Nota.md")
	v.AddNote("subdir/Archivio.md")

	tests := []struct {
		name            string
		target          string
		sourceDir       string
		caseInsensitive bool
		wantNil         bool
		wantRelPath     string
	}{
		{
			name:            "direct match",
			target:          "Home",
			sourceDir:       "",
			caseInsensitive: false,
			wantNil:         false,
			wantRelPath:     "Home.md",
		},
		{
			name:            "direct match with extension",
			target:          "Home.md",
			sourceDir:       "",
			caseInsensitive: false,
			wantNil:         false,
			wantRelPath:     "Home.md",
		},
		{
			name:            "relative path match",
			target:          "subdir/Nota",
			sourceDir:       "",
			caseInsensitive: false,
			wantNil:         false,
			wantRelPath:     "subdir/Nota.md",
		},
		{
			name:            "case insensitive match",
			target:          "home",
			sourceDir:       "",
			caseInsensitive: true,
			wantNil:         false,
			wantRelPath:     "Home.md",
		},
		{
			name:            "no match",
			target:          "inesistente",
			sourceDir:       "",
			caseInsensitive: false,
			wantNil:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := v.ResolveNote(tt.target, tt.sourceDir, tt.caseInsensitive)
			if tt.wantNil {
				if note != nil {
					t.Errorf("expected nil, got %v", note)
				}
			} else {
				if note == nil {
					t.Error("expected non-nil note")
					return
				}
				if note.RelPath != tt.wantRelPath {
					t.Errorf("expected RelPath %q, got %q", tt.wantRelPath, note.RelPath)
				}
			}
		})
	}
}

func TestTotalNotes(t *testing.T) {
	v := New("/vault")
	v.AddNote("a.md")
	v.AddNote("b.md")
	v.AddNote("c.md")

	if v.TotalNotes() != 3 {
		t.Errorf("expected 3, got %d", v.TotalNotes())
	}
}
