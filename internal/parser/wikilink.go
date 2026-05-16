package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var wikiLinkRe = regexp.MustCompile(`!?\[\[([^\]|#^]+)(?:#([^\]|^]+))?(?:\^([^\]|]+))?(?:\|([^\]]+))?\]\]`)

type WikiLink struct {
	Raw      string
	IsEmbed  bool
	Target   string
	Heading  string
	BlockRef string
	Alias    string
	Line     int
}

func ParseContent(content string) []WikiLink {
	lines := strings.Split(content, "\n")
	var links []WikiLink
	seen := make(map[string]bool)

	for lineNum, line := range lines {
		if isInCodeBlock(line) {
			continue
		}
		lineLinks := ParseLine(line, lineNum+1)
		for _, l := range lineLinks {
			key := l.Raw + ":" + strconv.Itoa(l.Line)
			if !seen[key] {
				seen[key] = true
				links = append(links, l)
			}
		}
	}

	return links
}

func ParseLine(line string, lineNum int) []WikiLink {
	matches := wikiLinkRe.FindAllStringSubmatch(line, -1)
	var links []WikiLink

	for _, m := range matches {
		links = append(links, WikiLink{
			Raw:      m[0],
			IsEmbed:  strings.HasPrefix(m[0], "!"),
			Target:   strings.TrimSpace(m[1]),
			Heading:  strings.TrimSpace(m[2]),
			BlockRef: strings.TrimSpace(m[3]),
			Alias:    strings.TrimSpace(m[4]),
			Line:     lineNum,
		})
	}

	return links
}

func isInCodeBlock(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "```")
}
