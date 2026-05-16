package resolver

import (
	"path/filepath"

	"github.com/marco-introini/obsidian-checker/internal/parser"
	"github.com/marco-introini/obsidian-checker/internal/vault"
)

type Resolver struct {
	Vault           *vault.Vault
	CaseInsensitive bool
}

func New(v *vault.Vault, caseInsensitive bool) *Resolver {
	return &Resolver{
		Vault:           v,
		CaseInsensitive: caseInsensitive,
	}
}

type ResolvedLink struct {
	WikiLink    parser.WikiLink
	TargetNote  *vault.Note
	TargetAsset *vault.Asset
	SourceDir   string
	Resolved    bool
}

func (r *Resolver) Resolve(link parser.WikiLink, sourceFile string) ResolvedLink {
	sourceDir := filepath.Dir(sourceFile)

	rl := ResolvedLink{
		WikiLink:  link,
		SourceDir: sourceDir,
	}

	note := r.Vault.ResolveNote(link.Target, sourceDir, r.CaseInsensitive)
	if note != nil {
		rl.TargetNote = note
		rl.Resolved = true
		return rl
	}

	asset := r.Vault.ResolveAsset(link.Target, sourceDir, r.CaseInsensitive)
	if asset != nil {
		rl.TargetAsset = asset
		rl.Resolved = true
		return rl
	}

	return rl
}
