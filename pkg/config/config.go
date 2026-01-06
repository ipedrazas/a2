package config

import (
	"github.com/ipedrazas/a2/pkg/safepath"
	"gopkg.in/yaml.v3"
)

// Config represents the .a2.yaml configuration file.
type Config struct {
	Coverage  CoverageConfig  `yaml:"coverage"`
	Files     FilesConfig     `yaml:"files"`
	Checks    ChecksConfig    `yaml:"checks"`
	External  []ExternalCheck `yaml:"external"`
	Execution ExecutionConfig `yaml:"execution"`
}

// ExecutionConfig configures how checks are executed.
type ExecutionConfig struct {
	Parallel bool `yaml:"parallel"`
}

// ExternalCheck defines a custom external check.
type ExternalCheck struct {
	ID       string   `yaml:"id"`
	Name     string   `yaml:"name"`
	Command  string   `yaml:"command"`
	Args     []string `yaml:"args"`
	Severity string   `yaml:"severity"` // "warn" or "fail"
}

// CoverageConfig configures the coverage check.
type CoverageConfig struct {
	Threshold float64 `yaml:"threshold"`
}

// FilesConfig configures the file existence check.
type FilesConfig struct {
	Required []string `yaml:"required"`
}

// ChecksConfig configures which checks to run.
type ChecksConfig struct {
	Disabled []string `yaml:"disabled"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Coverage: CoverageConfig{
			Threshold: 80.0,
		},
		Files: FilesConfig{
			Required: []string{"README.md", "LICENSE"},
		},
		Checks: ChecksConfig{
			Disabled: []string{},
		},
		Execution: ExecutionConfig{
			Parallel: true, // Run checks in parallel by default
		},
	}
}

// Load reads configuration from .a2.yaml in the given path.
// If the file doesn't exist, returns default configuration.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	// Use safepath to prevent directory traversal attacks
	data, err := safepath.ReadFile(path, ".a2.yaml")
	if err != nil {
		// Check if file doesn't exist (safepath wraps the error)
		if !safepath.Exists(path, ".a2.yaml") {
			// No config file, use defaults
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsCheckDisabled returns true if the given check ID is disabled.
func (c *Config) IsCheckDisabled(checkID string) bool {
	for _, disabled := range c.Checks.Disabled {
		if disabled == checkID {
			return true
		}
	}
	return false
}
