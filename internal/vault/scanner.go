package vault

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	extMarkdown = ".md"
	extCanvas   = ".canvas"
)

var noteExtensions = map[string]bool{
	".md":     true,
	".canvas": true,
}

var assetExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".webp": true,
	".bmp":  true,
	".pdf":  true,
	".mp3":  true,
	".mp4":  true,
	".webm": true,
	".ogg":  true,
	".wav":  true,
	".flac": true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".txt":  true,
	".csv":  true,
	".json": true,
	".xml":  true,
	".zip":  true,
	".mov":  true,
	".avi":  true,
}

func Scan(root string, excludeDirs []string, excludePatterns []string) (*Vault, error) {
	v := New(root)

	excludeSet := make(map[string]bool)
	for _, d := range excludeDirs {
		excludeSet[filepath.Clean(d)] = true
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		parts := strings.Split(filepath.ToSlash(relPath), "/")
		for _, part := range parts {
			if strings.HasPrefix(part, ".") {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		for _, part := range parts {
			if excludeSet[part] {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(relPath))
		if noteExtensions[ext] {
			v.AddNote(filepath.ToSlash(relPath))
		} else if assetExtensions[ext] {
			v.AddAsset(filepath.ToSlash(relPath))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return v, nil
}

func ReadNote(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
