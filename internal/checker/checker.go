package checker

import "github.com/marco-introini/obsidian-checker/internal/vault"

type IssueCode string

const (
	CodeNoteNotFound IssueCode = "note_not_found"
	CodeFileNotFound IssueCode = "file_not_found"
)

type Issue struct {
	File    string    `json:"file"`
	Line    int       `json:"line"`
	RawLink string    `json:"link"`
	Target  string    `json:"target"`
	Code    IssueCode `json:"issue"`
	Message string    `json:"message"`
}

type Summary struct {
	TotalFiles int `json:"total_files"`
	TotalLinks int `json:"total_links"`
	IssueCount int `json:"broken_links"`
}

type Checker interface {
	Name() string
	Check(v *vault.Vault) ([]Issue, Summary, error)
}
