# Obsidian Checker — Specifiche dell'Applicazione

## Panoramica

`obsidian-checker` è un'applicazione CLI scritta in Go che analizza un vault
Obsidian per rilevare problemi e incoerenze. Il vault Obsidian è una directory
contenente note in formato Markdown collegate tra loro tramite wiki-link nella
forma `[[nome-nota]]`.

L'obiettivo è fornire uno strumento veloce, estendibile e adatto all'uso in CI/CD
o in locale.

---

## Funzionalità

### V1 — Controllo Link Rotti (broken links)

Rileva tutti i wiki-link che puntano a note inesistenti all'interno del vault.

#### Regole di Risoluzione dei Link

Obsidian risolve i link con le seguenti regole:

1. **Nome base**: `[[Nome Nota]]` cerca un file `Nome Nota.md` (case-insensitive su
   macOS/Windows, case-sensitive su Linux).

2. **Percorso relativo**: `[[sottodir/Nome Nota]]` cerca in un percorso relativo
   alla nota corrente. Se non trovato, cerca anche il solo nome file.

3. **Heading**: `[[Nota#Heading]]` verifica che esista la nota e che l'heading
   esista nel file (il check degli heading è opzionale, controllabile via flag).

4. **Block reference**: `[[Nota^block-id]]` verifica che esista la nota.
   Il check del block-id è fuori scope per la V1.

5. **Alias (display text)**: `[[Nota|testo visualizzato]]` — l'alias è ignorato
   ai fini del controllo, vale solo il nome nota.

6. **Embed**: `![[Nota]]` — trattato come un link normale ai fini del controllo.

7. **Link a immagini/allegati**: `[[immagine.png]]`, `[[documento.pdf]]` — sono
   considerati link validi se il file esiste nel vault.

#### Note Inesistenti

Una nota è considerata inesistente quando:
- Non esiste alcun file `.md` con quel nome (case-insensitive) nel vault.
- Il percorso specificato non corrisponde a nessun file.
- Il file esiste ma ha un'estensione diversa da `.md` e non è un allegato noto.

#### Esclusioni

Directory da escludere dalla scansione:
- `.obsidian/` (configurazione di Obsidian)
- `.trash/` (cestino)
- Qualsiasi directory che inizia con `.` (directory nascoste)
- Directory configurabili tramite file di configurazione

File da escludere:
- Template (configurabile tramite cartella template)

#### Strategia Case-Insensitive

Poiché Obsidian su macOS/Windows tratta i nomi file in modo case-insensitive,
l'applicazione deve:
- Di default, usare la strategia del sistema operativo su cui gira (case-insensitive
  su macOS/Windows, case-sensitive su Linux).
- Offrire un flag `--case-sensitive` per forzare il comportamento case-sensitive.
- Offrire un flag `--case-insensitive` per forzare il comportamento case-insensitive.

---

## Interfaccia CLI

### Comandi

```
obsidian-checker [command] [flags]
```

#### Comando: `check links` (o `check broken-links`)

```
obsidian-checker check links [flags] <vault-path>
```

**Flags:**

| Flag                 | Short | Default                          | Descrizione                                          |
|----------------------|-------|----------------------------------|------------------------------------------------------|
| `--vault`            | `-v`  | `.`                              | Percorso del vault Obsidian                          |
| `--config`           | `-c`  | `.obsidian-checker.yaml`         | Percorso del file di configurazione                  |
| `--exclude`          | `-e`  | `[]`                             | Pattern glob aggiuntivi da escludere                 |
| `--check-headings`   |       | `false`                          | Verifica anche che gli heading referenziati esistano |
| `--case-sensitive`   |       | auto (OS-dependent)              | Forza risoluzione case-sensitive                     |
| `--case-insensitive` |       | auto (OS-dependent)              | Forza risoluzione case-insensitive                   |
| `--format`           | `-f`  | `table`                          | Formato output: `table`, `json`, `csv`               |
| `--quiet`            | `-q`  | `false`                          | Mostra solo errori, niente output di progresso       |

#### Comando: `check all`

Esegue tutti i controlli disponibili.

### Formati di Output

#### Table (default)
```
  N.  File                        Link                    Issue
 ---  --------------------------  ----------------------  ---------------------------
  1    Note Importante.md         [[Nota Inesistente]]     Nota non trovata
  2    Progetti/Attiva.md         [[Archiviata]]           Nota non trovata
  3    giornale/2024-01-15.md     [[Allegato perso.pdf]]   File non trovato
```

#### JSON
```json
{
  "vault": "/path/to/vault",
  "check": "broken-links",
  "issues": [
    {
      "file": "Note Importante.md",
      "line": 42,
      "link": "[[Nota Inesistente]]",
      "target": "Nota Inesistente",
      "issue": "note_not_found",
      "message": "Nota 'Nota Inesistente' non trovata nel vault"
    }
  ],
  "summary": {
    "total_files": 150,
    "total_links": 1200,
    "broken_links": 3
  }
}
```

#### CSV
```
file,line,link,target,issue,message
Note Importante.md,42,[[Nota Inesistente]],Nota Inesistente,note_not_found,"Nota 'Nota Inesistente' non trovata nel vault"
```

### Codici di Uscita

| Codice | Significato                                 |
|--------|---------------------------------------------|
| 0      | Nessun problema trovato                     |
| 1      | Trovati uno o più problemi                  |
| 2      | Errore di esecuzione (vault non valido etc) |

---

## Configurazione

File `.obsidian-checker.yaml` nella root del vault:

```yaml
# V1
exclude_dirs:
  - .obsidian
  - .trash
  - _templates
  - _archivio

exclude_patterns:
  - "**/*.excalidraw.md"   # File generati da plugin Excalidraw

check_headings: false
case_sensitive: false   # true = on, false = off, omesso = auto
```

I flag CLI sovrascrivono i valori di configurazione.

---

## Architettura

```
obsidian-checker/
├── cmd/
│   └── obsidian-checker/
│       └── main.go
├── internal/
│   ├── vault/           # Scansione e indicizzazione del vault
│   │   ├── vault.go     # Struttura dati del vault
│   │   └── scanner.go   # Scanner per file .md e allegati
│   ├── parser/          # Parsing dei wiki-link
│   │   └── wikilink.go  # Estrazione e parsing dei [[link]]
│   ├── checker/         # Interfaccia e implementazioni dei check
│   │   ├── checker.go   # Interfaccia Checker
│   │   └── broken_links.go  # Implementazione controllo link rotti
│   ├── resolver/        # Risoluzione link → file
│   │   └── resolver.go
│   ├── output/          # Formattazione output
│   │   ├── table.go
│   │   ├── json.go
│   │   └── csv.go
│   └── config/          # Gestione configurazione
│       └── config.go
├── go.mod
├── go.sum
└── docs/
    └── SPECS.md
```

---

## Dipendenze Esterne (Go)

| Libreria                          | Scopo                                  |
|-----------------------------------|----------------------------------------|
| `github.com/spf13/cobra`          | Gestione comandi e flag CLI            |
| `github.com/spf13/viper`          | Gestione file di configurazione        |
| `gopkg.in/yaml.v3`                | Parsing YAML (usato da Viper)          |
| `github.com/stretchr/testify`     | Testing                                |

---

## Criteri di Qualità

- **Performance**: Un vault di 10.000 note deve essere analizzato in < 5 secondi.
- **Testabilità**: Il parser e il resolver devono essere puri e testabili senza filesystem.
- **Estendibilità**: L'interfaccia `Checker` permette di aggiungere nuovi tipi di
  controllo senza modificare il core.

### Interfaccia Checker

```go
type Issue struct {
    File    string
    Line    int
    Link    string
    Target  string
    Code    string
    Message string
}

type Summary struct {
    TotalFiles int
    TotalLinks int
    IssueCount int
}

type Checker interface {
    Name() string
    Check(v *vault.Vault) ([]Issue, Summary, error)
}
```

---

## Sviluppo Futuro (Post-V1)

- **V2**: Controllo heading inesistenti (`[[Nota#HeadingCheNonEsiste]]`)
- **V3**: Controllo block reference inesistenti (`[[Nota^block-id]]`)
- **V4**: Tag orfani (tag non usati)
- **V5**: Note orfane (note senza link entranti)
- **V6**: Rilevazione link ambigui (più note con lo stesso nome)
- **V7**: Integrazione CI/CD come GitHub Action

---

## Riferimenti

- [Obsidian Help — Linking notes](https://help.obsidian.md/Linking+notes+and+files/Internal+links)
- [Obsidian Help — Embed files](https://help.obsidian.md/Linking+notes+and+files/Embed+files)
