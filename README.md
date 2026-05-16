# obsidian-checker

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-00ADD8?logo=go)](https://go.dev)

A fast, extensible CLI tool for analyzing [Obsidian](https://obsidian.md) vaults
and detecting issues like broken links, missing headings, and more.

## Features

- **Broken link detection** — find `[[wikilinks]]` pointing to non-existent notes
- **Heading validation** — optionally verify that referenced headings exist
- **Case-insensitive resolution** — matches Obsidian's behavior on macOS/Windows (configurable)
- **Multiple output formats** — `table`, `json`, `csv`
- **Configurable exclusions** — ignore specific directories or patterns
- **Extensible checker interface** — add custom checks with minimal effort

## Installation

```bash
go install github.com/marco-introini/obsidian-checker/cmd/obsidian-checker@latest
```

Or clone and build:

```bash
git clone https://github.com/marco-introini/obsidian-checker.git
cd obsidian-checker
go build -o obsidian-checker ./cmd/obsidian-checker
```

## Usage

```bash
# Check for broken links (default format: table)
obsidian-checker check links /path/to/vault

# JSON output
obsidian-checker check links /path/to/vault -f json

# Include heading validation
obsidian-checker check links /path/to/vault --check-headings

# Force case-sensitive resolution
obsidian-checker check links /path/to/vault --case-sensitive

# Exclude additional directories
obsidian-checker check links /path/to/vault -e _templates -e _archive

# Quiet mode (no progress output)
obsidian-checker check links /path/to/vault -q
```

## Configuration

Create a `.obsidian-checker.yaml` file in your vault root:

```yaml
exclude_dirs:
  - .obsidian
  - .trash
  - _templates
  - _archive

exclude_patterns:
  - "**/*.excalidraw.md"

check_headings: false
case_sensitive: false
```

CLI flags override config file values.

## Output Formats

### Table (default)

```
 N.   File                     Link                   Message
 ---  -----------------------  ---------------------  ---------------------------
 1    Home.md                  [[Broken Note]]        Note 'Broken Note' not found
```

### JSON

```json
{
  "vault": "/path/to/vault",
  "check": "broken-links",
  "issues": [...],
  "summary": { "total_files": 150, "total_links": 1200, "broken_links": 3 }
}
```

### CSV

```
file,line,link,target,issue,message
Home.md,42,[[Broken Note]],Broken Note,note_not_found,"Note 'Broken Note' not found"
```

## Exit Codes

| Code | Meaning                  |
|------|--------------------------|
| 0    | No issues found          |
| 1    | Issues detected          |
| 2    | Runtime error            |

## Roadmap

- [x] Broken link detection (v1)
- [x] Heading validation
- [ ] Block reference validation (`[[Note^block-id]]`)
- [ ] Orphan tag detection
- [ ] Orphan note detection (no incoming links)
- [ ] Ambiguous link detection (duplicate note names)
- [ ] GitHub Action integration

## License

MIT © [Marco Introini](https://github.com/marco-introini)
