package main

import (
	"fmt"
	"os"

	"github.com/marco-introini/obsidian-checker/internal/checker"
	"github.com/marco-introini/obsidian-checker/internal/config"
	"github.com/marco-introini/obsidian-checker/internal/output"
	"github.com/marco-introini/obsidian-checker/internal/vault"
	"github.com/spf13/cobra"
)

var (
	cfgFile         string
	outputFormat    string
	vaultPath       string
	quiet           bool
	checkHeadings   bool
	caseSensitive   bool
	caseInsensitive bool
	exclude         []string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

var rootCmd = &cobra.Command{
	Use:   "obsidian-checker",
	Short: "Strumento CLI per analizzare vault Obsidian",
	Long: `obsidian-checker analizza un vault Obsidian per rilevare problemi
e incoerenze come link rotti verso note inesistenti.`,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Esegue controlli sul vault",
}

var linksCmd = &cobra.Command{
	Use:   "links [vault-path]",
	Short: "Controlla i link rotti nel vault",
	Long: `Analizza tutti i file markdown del vault e verifica che
ogni wiki-link ([[...]]) punti a una nota o file esistente.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCheckLinks,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "File di configurazione (default: .obsidian-checker.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Formato output: table, json, csv")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Output minimale")
	rootCmd.PersistentFlags().BoolVar(&checkHeadings, "check-headings", false, "Verifica heading referenziati")
	rootCmd.PersistentFlags().BoolVar(&caseSensitive, "case-sensitive", false, "Forza risoluzione case-sensitive")
	rootCmd.PersistentFlags().BoolVar(&caseInsensitive, "case-insensitive", false, "Forza risoluzione case-insensitive")
	rootCmd.PersistentFlags().StringSliceVarP(&exclude, "exclude", "e", nil, "Directory aggiuntive da escludere")

	checkCmd.AddCommand(linksCmd)
	rootCmd.AddCommand(checkCmd)
}

func runCheckLinks(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		vaultPath = args[0]
	}
	if vaultPath == "" {
		vaultPath = "."
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("errore caricamento configurazione: %w", err)
	}

	cfg.VaultPath = vaultPath
	cfg.Format = outputFormat
	cfg.Quiet = quiet
	cfg.CheckHeadings = checkHeadings

	if caseSensitive {
		t := true
		cfg.CaseSensitive = &t
	} else if caseInsensitive {
		f := false
		cfg.CaseSensitive = &f
	}

	if len(exclude) > 0 {
		cfg.ExcludeExtra = exclude
	}

	if !cfg.Quiet {
		fmt.Fprintf(os.Stderr, "Scansione del vault: %s\n", cfg.VaultPath)
	}

	v, err := vault.Scan(cfg.VaultPath, cfg.AllExcludeDirs(), cfg.ExcludePatterns)
	if err != nil {
		return fmt.Errorf("errore scansione vault: %w", err)
	}

	if !cfg.Quiet {
		fmt.Fprintf(os.Stderr, "Trovate %d note e %d allegati\n", len(v.Notes), len(v.Assets))
		fmt.Fprintf(os.Stderr, "Analisi link in corso...\n")
	}

	blChecker := checker.NewBrokenLinksChecker(cfg.IsCaseSensitive(), cfg.CheckHeadings)
	issues, summary, err := blChecker.Check(v)
	if err != nil {
		return fmt.Errorf("errore controllo link: %w", err)
	}

	var formatter output.Formatter
	switch cfg.Format {
	case "json":
		formatter = &output.JSONFormatter{VaultPath: cfg.VaultPath}
	case "csv":
		formatter = &output.CSVFormatter{}
	default:
		formatter = &output.TableFormatter{}
	}

	out, err := formatter.Format(issues, summary)
	if err != nil {
		return fmt.Errorf("errore formattazione output: %w", err)
	}

	fmt.Print(out)

	if len(issues) > 0 {
		os.Exit(1)
	}

	return nil
}
