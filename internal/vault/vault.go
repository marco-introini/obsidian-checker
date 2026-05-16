package vault

import (
	"path/filepath"
	"strings"
)

type Note struct {
	Path     string
	RelPath  string
	Name     string
	BaseName string
}

type Asset struct {
	Path    string
	RelPath string
	Name    string
}

type Vault struct {
	Root       string
	Notes      map[string]*Note
	Assets     map[string]*Asset
	noteIndex  map[string]*Note
	assetIndex map[string]*Asset
}

func New(root string) *Vault {
	return &Vault{
		Root:       root,
		Notes:      make(map[string]*Note),
		Assets:     make(map[string]*Asset),
		noteIndex:  make(map[string]*Note),
		assetIndex: make(map[string]*Asset),
	}
}

func (v *Vault) AddNote(relPath string) {
	name := filepath.Base(relPath)
	baseName := strings.TrimSuffix(name, filepath.Ext(name))

	note := &Note{
		Path:     filepath.Join(v.Root, relPath),
		RelPath:  relPath,
		Name:     name,
		BaseName: baseName,
	}

	v.Notes[relPath] = note
	v.noteIndex[strings.ToLower(relPath)] = note
	v.noteIndex[strings.ToLower(baseName)] = note

	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".md" {
		v.noteIndex[strings.ToLower(baseName+".md")] = note
	}
}

func (v *Vault) AddAsset(relPath string) {
	name := filepath.Base(relPath)
	asset := &Asset{
		Path:    filepath.Join(v.Root, relPath),
		RelPath: relPath,
		Name:    name,
	}
	v.Assets[relPath] = asset
	v.assetIndex[strings.ToLower(relPath)] = asset
	v.assetIndex[strings.ToLower(name)] = asset
}

func (v *Vault) ResolveNote(target string, sourceDir string, caseInsensitive bool) *Note {
	target = filepath.Clean(target)

	if caseInsensitive {
		return v.resolveNoteCI(target, sourceDir)
	}
	return v.resolveNoteCS(target, sourceDir)
}

func (v *Vault) resolveNoteCS(target string, sourceDir string) *Note {
	if strings.Contains(target, "/") || strings.Contains(target, string(filepath.Separator)) {
		relPath := filepath.Join(sourceDir, target)
		if note, ok := v.Notes[relPath]; ok {
			return note
		}
		if note, ok := v.Notes[relPath+".md"]; ok {
			return note
		}
	}

	if note, ok := v.Notes[target]; ok {
		return note
	}
	if note, ok := v.Notes[target+".md"]; ok {
		return note
	}

	if strings.Contains(target, "/") || strings.Contains(target, string(filepath.Separator)) {
		relPath := filepath.Join(sourceDir, target)
		for path, note := range v.Notes {
			if strings.EqualFold(path, relPath) || strings.EqualFold(path, relPath+".md") {
				return note
			}
		}
	}

	for path, note := range v.Notes {
		if strings.EqualFold(filepath.Base(path), target) || strings.EqualFold(filepath.Base(path), target+".md") {
			return note
		}
	}

	return nil
}

func (v *Vault) resolveNoteCI(target string, sourceDir string) *Note {
	if strings.Contains(target, "/") || strings.Contains(target, string(filepath.Separator)) {
		relPath := filepath.Join(sourceDir, target)
		if note, ok := v.noteIndex[strings.ToLower(relPath)]; ok {
			return note
		}
		if note, ok := v.noteIndex[strings.ToLower(relPath+".md")]; ok {
			return note
		}
	}

	if note, ok := v.noteIndex[strings.ToLower(target)]; ok {
		return note
	}
	if note, ok := v.noteIndex[strings.ToLower(target+".md")]; ok {
		return note
	}

	return nil
}

func (v *Vault) ResolveAsset(target string, sourceDir string, caseInsensitive bool) *Asset {
	target = filepath.Clean(target)

	if caseInsensitive {
		if strings.Contains(target, "/") || strings.Contains(target, string(filepath.Separator)) {
			relPath := filepath.Join(sourceDir, target)
			if asset, ok := v.assetIndex[strings.ToLower(relPath)]; ok {
				return asset
			}
		}
		if asset, ok := v.assetIndex[strings.ToLower(target)]; ok {
			return asset
		}
		return nil
	}

	if strings.Contains(target, "/") || strings.Contains(target, string(filepath.Separator)) {
		relPath := filepath.Join(sourceDir, target)
		if asset, ok := v.Assets[relPath]; ok {
			return asset
		}
	}
	if asset, ok := v.Assets[target]; ok {
		return asset
	}
	return nil
}

func (v *Vault) TotalNotes() int {
	return len(v.Notes)
}
