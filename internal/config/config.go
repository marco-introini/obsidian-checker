package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ExcludeDirs     []string `yaml:"exclude_dirs"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
	CheckHeadings   bool     `yaml:"check_headings"`
	CaseSensitive   *bool    `yaml:"case_sensitive"`

	Format          string
	VaultPath       string
	Quiet           bool
	ExcludeExtra    []string
}

var defaultExcludeDirs = []string{".obsidian", ".trash", ".git"}
var defaultExcludePatterns = []string{}

func Default() *Config {
	return &Config{
		ExcludeDirs:     defaultExcludeDirs,
		ExcludePatterns: defaultExcludePatterns,
		CheckHeadings:   false,
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.ExcludeDirs == nil {
		cfg.ExcludeDirs = defaultExcludeDirs
	}
	if cfg.ExcludePatterns == nil {
		cfg.ExcludePatterns = defaultExcludePatterns
	}

	return cfg, nil
}

func (c *Config) IsCaseSensitive() bool {
	if c.CaseSensitive != nil {
		return *c.CaseSensitive
	}
	return runtime.GOOS == "linux"
}

func (c *Config) AllExcludeDirs() []string {
	dirs := make([]string, 0, len(c.ExcludeDirs)+len(c.ExcludeExtra))
	dirs = append(dirs, c.ExcludeDirs...)
	dirs = append(dirs, c.ExcludeExtra...)
	return dirs
}
