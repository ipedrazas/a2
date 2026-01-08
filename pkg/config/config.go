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
	Language  LanguageConfig  `yaml:"language"`
}

// LanguageConfig handles language detection and settings.
type LanguageConfig struct {
	Explicit   []string             `yaml:"explicit,omitempty"` // Override auto-detect
	AutoDetect bool                 `yaml:"auto_detect"`        // Default: true
	Go         GoLanguageConfig     `yaml:"go,omitempty"`
	Python     PythonLanguageConfig `yaml:"python,omitempty"`
	Node       NodeLanguageConfig   `yaml:"node,omitempty"`
	Java       JavaLanguageConfig   `yaml:"java,omitempty"`
	Rust       RustLanguageConfig   `yaml:"rust,omitempty"`
}

// GoLanguageConfig contains Go-specific settings.
type GoLanguageConfig struct {
	CoverageThreshold   float64 `yaml:"coverage_threshold,omitempty"`
	CyclomaticThreshold int     `yaml:"cyclomatic_threshold,omitempty"`
}

// PythonLanguageConfig contains Python-specific settings.
type PythonLanguageConfig struct {
	PackageManager      string  `yaml:"package_manager,omitempty"` // auto, pip, poetry, pipenv
	TestRunner          string  `yaml:"test_runner,omitempty"`     // auto, pytest, unittest
	Formatter           string  `yaml:"formatter,omitempty"`       // auto, black, ruff
	Linter              string  `yaml:"linter,omitempty"`          // auto, pylint, ruff, flake8
	CoverageThreshold   float64 `yaml:"coverage_threshold,omitempty"`
	CyclomaticThreshold int     `yaml:"cyclomatic_threshold,omitempty"`
}

// NodeLanguageConfig contains Node.js-specific settings.
type NodeLanguageConfig struct {
	PackageManager    string  `yaml:"package_manager,omitempty"` // auto, npm, yarn, pnpm, bun
	TestRunner        string  `yaml:"test_runner,omitempty"`     // auto, jest, vitest, mocha, npm-test
	Formatter         string  `yaml:"formatter,omitempty"`       // auto, prettier, biome
	Linter            string  `yaml:"linter,omitempty"`          // auto, eslint, biome, oxlint
	CoverageThreshold float64 `yaml:"coverage_threshold,omitempty"`
}

// JavaLanguageConfig contains Java-specific settings.
type JavaLanguageConfig struct {
	BuildTool         string  `yaml:"build_tool,omitempty"`         // auto, maven, gradle
	TestRunner        string  `yaml:"test_runner,omitempty"`        // auto, junit, testng
	CoverageThreshold float64 `yaml:"coverage_threshold,omitempty"` // default 80
}

// RustLanguageConfig contains Rust-specific settings.
type RustLanguageConfig struct {
	CoverageThreshold float64 `yaml:"coverage_threshold,omitempty"` // default 80
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
		Language: LanguageConfig{
			AutoDetect: true,
			Go: GoLanguageConfig{
				CoverageThreshold:   80.0,
				CyclomaticThreshold: 15,
			},
			Python: PythonLanguageConfig{
				PackageManager:      "auto",
				TestRunner:          "auto",
				Formatter:           "auto",
				Linter:              "auto",
				CoverageThreshold:   80.0,
				CyclomaticThreshold: 15,
			},
			Node: NodeLanguageConfig{
				PackageManager:    "auto",
				TestRunner:        "auto",
				Formatter:         "auto",
				Linter:            "auto",
				CoverageThreshold: 80.0,
			},
			Java: JavaLanguageConfig{
				BuildTool:         "auto",
				TestRunner:        "auto",
				CoverageThreshold: 80.0,
			},
			Rust: RustLanguageConfig{
				CoverageThreshold: 80.0,
			},
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

// checkAliases maps old check IDs to new language-prefixed IDs for backward compatibility.
// Also maps short names to full IDs for convenience with --skip flag.
var checkAliases = map[string]string{
	// Go check aliases (legacy)
	"go_mod":   "go:module",
	"build":    "go:build",
	"tests":    "go:tests",
	"gofmt":    "go:format",
	"govet":    "go:vet",
	"coverage": "go:coverage",
	"deps":     "go:deps",
	// Common check short names
	"dockerfile":  "common:dockerfile",
	"ci":          "common:ci",
	"health":      "common:health",
	"secrets":     "common:secrets",
	"env":         "common:env",
	"license":     "common:license",
	"sast":        "common:sast",
	"api_docs":    "common:api_docs",
	"changelog":   "common:changelog",
	"integration": "common:integration",
	"metrics":     "common:metrics",
	"errors":      "common:errors",
	"precommit":   "common:precommit",
	"k8s":         "common:k8s",
	"shutdown":    "common:shutdown",
}

// IsCheckDisabled returns true if the given check ID is disabled.
func (c *Config) IsCheckDisabled(checkID string) bool {
	for _, disabled := range c.Checks.Disabled {
		// Direct match
		if disabled == checkID {
			return true
		}
		// Check if disabled ID is an alias for the check ID
		if alias, ok := checkAliases[disabled]; ok && alias == checkID {
			return true
		}
		// Check if check ID is an alias for the disabled ID
		if alias, ok := checkAliases[checkID]; ok && alias == disabled {
			return true
		}
	}
	return false
}
