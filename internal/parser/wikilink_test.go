package parser

import "testing"

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantCnt  int
		want     []WikiLink
	}{
		{
			name:    "simple link",
			line:    "Vedi [[Nota]] per dettagli",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[Nota]]", Target: "Nota", Line: 1},
			},
		},
		{
			name:    "link with alias",
			line:    "Vedi [[Nota|testo mostrato]] per dettagli",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[Nota|testo mostrato]]", Target: "Nota", Alias: "testo mostrato", Line: 1},
			},
		},
		{
			name:    "link with heading",
			line:    "[[Nota#Heading]]",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[Nota#Heading]]", Target: "Nota", Heading: "Heading", Line: 1},
			},
		},
		{
			name:    "link with heading and alias",
			line:    "[[Nota#Heading|alias]]",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[Nota#Heading|alias]]", Target: "Nota", Heading: "Heading", Alias: "alias", Line: 1},
			},
		},
		{
			name:    "link with block ref",
			line:    "[[Nota^block123]]",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[Nota^block123]]", Target: "Nota", BlockRef: "block123", Line: 1},
			},
		},
		{
			name:    "embed link",
			line:    "![[Immagine.png]]",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "![[Immagine.png]]", Target: "Immagine.png", IsEmbed: true, Line: 1},
			},
		},
		{
			name:    "relative path link",
			line:    "[[sottodir/Nota]]",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[sottodir/Nota]]", Target: "sottodir/Nota", Line: 1},
			},
		},
		{
			name:    "multiple links in one line",
			line:    "[[A]] e [[B]]",
			wantCnt: 2,
			want: []WikiLink{
				{Raw: "[[A]]", Target: "A", Line: 1},
				{Raw: "[[B]]", Target: "B", Line: 1},
			},
		},
		{
			name:    "no link",
			line:    "Una linea senza link",
			wantCnt: 0,
			want:    nil,
		},
		{
			name:    "link with spaces in heading",
			line:    "[[Nota#Heading con spazi]]",
			wantCnt: 1,
			want: []WikiLink{
				{Raw: "[[Nota#Heading con spazi]]", Target: "Nota", Heading: "Heading con spazi", Line: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links := ParseLine(tt.line, 1)
			if len(links) != tt.wantCnt {
				t.Errorf("expected %d links, got %d", tt.wantCnt, len(links))
				return
			}
			for i := range links {
				if links[i].Raw != tt.want[i].Raw {
					t.Errorf("Raw: want %q, got %q", tt.want[i].Raw, links[i].Raw)
				}
				if links[i].Target != tt.want[i].Target {
					t.Errorf("Target: want %q, got %q", tt.want[i].Target, links[i].Target)
				}
				if links[i].Heading != tt.want[i].Heading {
					t.Errorf("Heading: want %q, got %q", tt.want[i].Heading, links[i].Heading)
				}
				if links[i].BlockRef != tt.want[i].BlockRef {
					t.Errorf("BlockRef: want %q, got %q", tt.want[i].BlockRef, links[i].BlockRef)
				}
				if links[i].Alias != tt.want[i].Alias {
					t.Errorf("Alias: want %q, got %q", tt.want[i].Alias, links[i].Alias)
				}
				if links[i].IsEmbed != tt.want[i].IsEmbed {
					t.Errorf("IsEmbed: want %v, got %v", tt.want[i].IsEmbed, links[i].IsEmbed)
				}
			}
		})
	}
}

func TestParseContent(t *testing.T) {
	content := "# Titolo\n\nLink a [[Nota1]].\nAltro link a [[Nota2|alias]].\n"
	links := ParseContent(content)
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
	if links[0].Line != 3 {
		t.Errorf("expected first link on line 3, got %d", links[0].Line)
	}
	if links[1].Line != 4 {
		t.Errorf("expected second link on line 4, got %d", links[1].Line)
	}
}
