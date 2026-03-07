package config

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/safepath"
	"gopkg.in/yaml.v3"
)

// SourceDirEntry represents a single source directory with an optional profile.
type SourceDirEntry struct {
	Path    string `yaml:"path"`
	Profile string `yaml:"profile,omitempty"`
	// Disabled is populated at runtime by resolving the Profile name.
	// It is not read from YAML.
	Disabled []string `yaml:"-"`
}

// SourceDirConfig is a custom YAML type that accepts three formats:
//
//	source_dir: api                              → [{Path: "api"}]
//	source_dir: [api, cli]                       → [{Path: "api"}, {Path: "cli"}]
//	source_dir:
//	  - path: api
//	    profile: api                             → [{Path: "api", Profile: "api"}]
//	  - path: cli
//	    profile: cli                             → [{Path: "cli", Profile: "cli"}]
type SourceDirConfig []SourceDirEntry

// Paths returns just the directory paths (for backward-compatible use).
func (s SourceDirConfig) Paths() []string {
	if len(s) == 0 {
		return nil
	}
	paths := make([]string, len(s))
	for i, e := range s {
		paths[i] = e.Path
	}
	return paths
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s *SourceDirConfig) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// source_dir: api
		var single string
		if err := value.Decode(&single); err != nil {
			return err
		}
		if single != "" {
			*s = SourceDirConfig{{Path: single}}
		}
		return nil
	case yaml.SequenceNode:
		// Could be a list of strings or a list of objects
		if len(value.Content) == 0 {
			return nil
		}
		// Check first element to determine format
		if value.Content[0].Kind == yaml.ScalarNode {
			// source_dir: [api, cli]
			var list []string
			if err := value.Decode(&list); err != nil {
				return err
			}
			entries := make(SourceDirConfig, len(list))
			for i, p := range list {
				entries[i] = SourceDirEntry{Path: p}
			}
			*s = entries
			return nil
		}
		// source_dir: [{path: api, profile: api}, ...]
		var entries []SourceDirEntry
		if err := value.Decode(&entries); err != nil {
			return err
		}
		*s = SourceDirConfig(entries)
		return nil
	default:
		return fmt.Errorf("source_dir must be a string, list of strings, or list of {path, profile} objects")
	}
}

// StringOrSlice is an alias for backward compatibility.
type StringOrSlice = SourceDirConfig

// Config represents the .a2.yaml configuration file.
type Config struct {
	Coverage  CoverageConfig        `yaml:"coverage"`
	Files     FilesConfig           `yaml:"files"`
	Checks    ChecksConfig          `yaml:"checks"`
	Security  SecurityConfig        `yaml:"security,omitempty"`
	External  []ExternalCheck       `yaml:"external"`
	Execution ExecutionConfig       `yaml:"execution"`
	Language  LanguageConfig        `yaml:"language"`
	Tools     map[string]ToolConfig `yaml:"tools,omitempty"`
}

// ToolConfig allows per-tool configuration overrides.
type ToolConfig struct {
	// RunByDefault overrides the tool's default run behavior.
	// nil = use tool's default, true = run, false = don't run
	RunByDefault *bool `yaml:"run_by_default,omitempty"`
}

// LanguageConfig handles language detection and settings.
type LanguageConfig struct {
	Explicit   []string                 `yaml:"explicit,omitempty"` // Override auto-detect
	AutoDetect bool                     `yaml:"auto_detect"`        // Default: true
	Go         GoLanguageConfig         `yaml:"go,omitempty"`
	Python     PythonLanguageConfig     `yaml:"python,omitempty"`
	Node       NodeLanguageConfig       `yaml:"node,omitempty"`
	Java       JavaLanguageConfig       `yaml:"java,omitempty"`
	Rust       RustLanguageConfig       `yaml:"rust,omitempty"`
	TypeScript TypeScriptLanguageConfig `yaml:"typescript,omitempty"`
	Swift      SwiftLanguageConfig      `yaml:"swift,omitempty"`
}

// GoLanguageConfig contains Go-specific settings.
type GoLanguageConfig struct {
	SourceDir           StringOrSlice `yaml:"source_dir,omitempty"` // Subdirectory(ies) containing Go code
	CoverageThreshold   float64       `yaml:"coverage_threshold,omitempty"`
	CyclomaticThreshold int           `yaml:"cyclomatic_threshold,omitempty"`
}

// PythonLanguageConfig contains Python-specific settings.
type PythonLanguageConfig struct {
	SourceDir           StringOrSlice `yaml:"source_dir,omitempty"`      // Subdirectory(ies) containing Python code
	PackageManager      string        `yaml:"package_manager,omitempty"` // auto, pip, poetry, pipenv
	TestRunner          string        `yaml:"test_runner,omitempty"`     // auto, pytest, unittest
	Formatter           string        `yaml:"formatter,omitempty"`       // auto, black, ruff
	Linter              string        `yaml:"linter,omitempty"`          // auto, pylint, ruff, flake8
	CoverageThreshold   float64       `yaml:"coverage_threshold,omitempty"`
	CyclomaticThreshold int           `yaml:"cyclomatic_threshold,omitempty"`
}

// NodeLanguageConfig contains Node.js-specific settings.
type NodeLanguageConfig struct {
	SourceDir         StringOrSlice `yaml:"source_dir,omitempty"`      // Subdirectory(ies) containing Node.js code
	PackageManager    string        `yaml:"package_manager,omitempty"` // auto, npm, yarn, pnpm, bun
	TestRunner        string        `yaml:"test_runner,omitempty"`     // auto, jest, vitest, mocha, npm-test
	Formatter         string        `yaml:"formatter,omitempty"`       // auto, prettier, biome
	Linter            string        `yaml:"linter,omitempty"`          // auto, eslint, biome, oxlint
	CoverageThreshold float64       `yaml:"coverage_threshold,omitempty"`
}

// JavaLanguageConfig contains Java-specific settings.
type JavaLanguageConfig struct {
	SourceDir         StringOrSlice `yaml:"source_dir,omitempty"`         // Subdirectory(ies) containing Java code
	BuildTool         string        `yaml:"build_tool,omitempty"`         // auto, maven, gradle
	TestRunner        string        `yaml:"test_runner,omitempty"`        // auto, junit, testng
	CoverageThreshold float64       `yaml:"coverage_threshold,omitempty"` // default 80
}

// RustLanguageConfig contains Rust-specific settings.
type RustLanguageConfig struct {
	SourceDir         StringOrSlice `yaml:"source_dir,omitempty"`         // Subdirectory(ies) containing Rust code
	CoverageThreshold float64       `yaml:"coverage_threshold,omitempty"` // default 80
}

// TypeScriptLanguageConfig contains TypeScript-specific settings.
type TypeScriptLanguageConfig struct {
	SourceDir         StringOrSlice `yaml:"source_dir,omitempty"`      // Subdirectory(ies) containing TypeScript code
	PackageManager    string        `yaml:"package_manager,omitempty"` // auto, npm, yarn, pnpm, bun
	TestRunner        string        `yaml:"test_runner,omitempty"`     // auto, jest, vitest, mocha
	Formatter         string        `yaml:"formatter,omitempty"`       // auto, prettier, biome, dprint
	Linter            string        `yaml:"linter,omitempty"`          // auto, eslint, biome, oxlint
	CoverageThreshold float64       `yaml:"coverage_threshold,omitempty"`
}

// SwiftLanguageConfig contains Swift-specific settings.
type SwiftLanguageConfig struct {
	SourceDir         StringOrSlice `yaml:"source_dir,omitempty"`         // Subdirectory(ies) containing Swift code
	Formatter         string        `yaml:"formatter,omitempty"`          // auto, swift-format, swiftformat
	Linter            string        `yaml:"linter,omitempty"`             // auto, swiftlint
	CoverageThreshold float64       `yaml:"coverage_threshold,omitempty"` // default 80
}

// ExecutionConfig configures how checks are executed.
type ExecutionConfig struct {
	Parallel bool `yaml:"parallel"`
}

// ExternalCheck defines a custom external check.
type ExternalCheck struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Command   string   `yaml:"command"`
	Args      []string `yaml:"args"`
	Severity  string   `yaml:"severity"`             // "warn" or "fail"
	SourceDir string   `yaml:"source_dir,omitempty"` // Subdirectory to run the command in (like language source_dir)
}

// CoverageConfig configures the coverage check.
type CoverageConfig struct {
	Threshold float64 `yaml:"threshold"`
}

// FilesConfig configures the file existence check.
type FilesConfig struct {
	Required []string `yaml:"required"`
}

// SecurityConfig configures security-related checks.
type SecurityConfig struct {
	Filesystem FileSystemConfig `yaml:"filesystem,omitempty"`
}

// FileSystemConfig configures the filesystem security check.
type FileSystemConfig struct {
	// Allow is a list of allowlist rules for filesystem findings.
	// Examples:
	// - "pkg/checks/common/k8s.go:94"
	// - "pkg/checks/common/k8s.go:os.ReadDir(chartsDir)"
	// - "pkg/checks/common/**"
	Allow []string `yaml:"allow,omitempty"`
}

// ChecksConfig configures which checks to run.
type ChecksConfig struct {
	// Disabled is a list of check IDs or wildcard patterns to skip.
	// Wildcard patterns (*:logging, node:*, *:*) must be quoted in YAML
	// (e.g. "*:logging") because unquoted * is YAML's alias character.
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
		Security: SecurityConfig{
			Filesystem: FileSystemConfig{
				Allow: []string{},
			},
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
			TypeScript: TypeScriptLanguageConfig{
				PackageManager:    "auto",
				TestRunner:        "auto",
				Formatter:         "auto",
				Linter:            "auto",
				CoverageThreshold: 80.0,
			},
			Swift: SwiftLanguageConfig{
				Formatter:         "auto",
				Linter:            "auto",
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
		errStr := err.Error()
		if strings.Contains(errStr, "alphabetic or numeric") || strings.Contains(errStr, "alias") {
			return nil, fmt.Errorf("%w (wildcard patterns in checks.disabled must be quoted, e.g. \"*:logging\")", err)
		}
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
	"k8s":         "devops:k8s",
	"shutdown":    "common:shutdown",
}

// GetSourceDirsForLang returns the configured source directory paths for a language.
// Returns nil if not configured (meaning use root path).
func (c *Config) GetSourceDirsForLang(lang string) []string {
	entries := c.GetSourceDirEntriesForLang(lang)
	if len(entries) == 0 {
		return nil
	}
	return SourceDirConfig(entries).Paths()
}

// GetSourceDirEntriesForLang returns the full source directory entries for a language,
// including optional profile associations.
func (c *Config) GetSourceDirEntriesForLang(lang string) []SourceDirEntry {
	switch lang {
	case "go":
		return c.Language.Go.SourceDir
	case "python":
		return c.Language.Python.SourceDir
	case "node":
		return c.Language.Node.SourceDir
	case "java":
		return c.Language.Java.SourceDir
	case "rust":
		return c.Language.Rust.SourceDir
	case "typescript":
		return c.Language.TypeScript.SourceDir
	case "swift":
		return c.Language.Swift.SourceDir
	default:
		return nil
	}
}

// GetSourceDirs returns a map of all configured source directory paths.
// Only languages with non-empty source directories are included.
func (c *Config) GetSourceDirs() map[string][]string {
	dirs := make(map[string][]string)
	if len(c.Language.Go.SourceDir) > 0 {
		dirs["go"] = c.Language.Go.SourceDir.Paths()
	}
	if len(c.Language.Python.SourceDir) > 0 {
		dirs["python"] = c.Language.Python.SourceDir.Paths()
	}
	if len(c.Language.Node.SourceDir) > 0 {
		dirs["node"] = c.Language.Node.SourceDir.Paths()
	}
	if len(c.Language.Java.SourceDir) > 0 {
		dirs["java"] = c.Language.Java.SourceDir.Paths()
	}
	if len(c.Language.Rust.SourceDir) > 0 {
		dirs["rust"] = c.Language.Rust.SourceDir.Paths()
	}
	if len(c.Language.TypeScript.SourceDir) > 0 {
		dirs["typescript"] = c.Language.TypeScript.SourceDir.Paths()
	}
	if len(c.Language.Swift.SourceDir) > 0 {
		dirs["swift"] = c.Language.Swift.SourceDir.Paths()
	}
	return dirs
}

// IsCheckDisabled returns true if the given check ID is disabled.
func (c *Config) IsCheckDisabled(checkID string) bool {
	for _, disabled := range c.Checks.Disabled {
		if MatchDisabled(checkID, disabled) {
			return true
		}
	}
	return false
}

// MatchDisabled returns true if the disabled pattern applies to the check ID.
// It supports wildcard patterns and legacy aliases.
func MatchDisabled(checkID, disabled string) bool {
	// Wildcard pattern match
	if matchesPattern(checkID, disabled) {
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
	return false
}

// ResolveSourceDirProfiles resolves profile names in source_dir entries to their
// disabled check lists. The resolver function maps a profile name to its disabled checks.
// This must be called after loading profiles but before running checks.
func (c *Config) ResolveSourceDirProfiles(resolver func(string) []string) {
	resolveEntries := func(entries SourceDirConfig) {
		for i := range entries {
			if entries[i].Profile != "" {
				entries[i].Disabled = resolver(entries[i].Profile)
			}
		}
	}
	resolveEntries(c.Language.Go.SourceDir)
	resolveEntries(c.Language.Python.SourceDir)
	resolveEntries(c.Language.Node.SourceDir)
	resolveEntries(c.Language.Java.SourceDir)
	resolveEntries(c.Language.Rust.SourceDir)
	resolveEntries(c.Language.TypeScript.SourceDir)
	resolveEntries(c.Language.Swift.SourceDir)
}

// GetToolRunByDefault returns the run_by_default override for a tool.
// Returns nil if no override is configured (use tool's default).
// Returns pointer to bool if override is configured.
func (c *Config) GetToolRunByDefault(toolName string) *bool {
	if c.Tools == nil {
		return nil
	}
	if toolCfg, ok := c.Tools[toolName]; ok {
		return toolCfg.RunByDefault
	}
	return nil
}
