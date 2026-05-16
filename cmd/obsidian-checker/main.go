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
	Short: "CLI tool for analyzing Obsidian vaults",
	Long: `obsidian-checker scans an Obsidian vault to detect issues
like broken links pointing to non-existent notes.`,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run checks on a vault",
}

var linksCmd = &cobra.Command{
	Use:   "links [vault-path]",
	Short: "Check for broken links in the vault",
	Long: `Scan all markdown files in the vault and verify
every wiki-link ([[...]]) points to an existing note or file.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCheckLinks,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default: .obsidian-checker.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table, json, csv")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output")
	rootCmd.PersistentFlags().BoolVar(&checkHeadings, "check-headings", false, "Validate referenced headings")
	rootCmd.PersistentFlags().BoolVar(&caseSensitive, "case-sensitive", false, "Force case-sensitive resolution")
	rootCmd.PersistentFlags().BoolVar(&caseInsensitive, "case-insensitive", false, "Force case-insensitive resolution")
	rootCmd.PersistentFlags().StringSliceVarP(&exclude, "exclude", "e", nil, "Additional directories to exclude")

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
		return fmt.Errorf("error loading configuration: %w", err)
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
		fmt.Fprintf(os.Stderr, "Scanning vault: %s\n", cfg.VaultPath)
	}

	v, err := vault.Scan(cfg.VaultPath, cfg.AllExcludeDirs(), cfg.ExcludePatterns)
	if err != nil {
		return fmt.Errorf("error scanning vault: %w", err)
	}

	if !cfg.Quiet {
		fmt.Fprintf(os.Stderr, "Found %d notes and %d attachments\n", len(v.Notes), len(v.Assets))
		fmt.Fprintf(os.Stderr, "Analyzing links...\n")
	}

	blChecker := checker.NewBrokenLinksChecker(cfg.IsCaseSensitive(), cfg.CheckHeadings)
	issues, summary, err := blChecker.Check(v)
	if err != nil {
		return fmt.Errorf("error checking links: %w", err)
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
		return fmt.Errorf("error formatting output: %w", err)
	}

	fmt.Print(out)

	if len(issues) > 0 {
		os.Exit(1)
	}

	return nil
}
