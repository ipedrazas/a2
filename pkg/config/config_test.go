package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite is the test suite for the config package.
type ConfigTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *ConfigTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-config-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *ConfigTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *ConfigTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestDefaultConfig tests that DefaultConfig returns correct defaults.
func (suite *ConfigTestSuite) TestDefaultConfig() {
	cfg := DefaultConfig()

	suite.NotNil(cfg)
	suite.Equal(80.0, cfg.Coverage.Threshold)
	suite.Equal([]string{"README.md", "LICENSE"}, cfg.Files.Required)
	suite.Empty(cfg.Checks.Disabled)
	suite.Empty(cfg.External)
}

// TestLoad_NoConfigFile tests that Load returns default config when .a2.yaml doesn't exist.
func (suite *ConfigTestSuite) TestLoad_NoConfigFile() {
	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	suite.Equal(80.0, cfg.Coverage.Threshold)
	suite.Equal([]string{"README.md", "LICENSE"}, cfg.Files.Required)
}

// TestLoad_ValidConfig tests that Load parses a valid YAML config file.
func (suite *ConfigTestSuite) TestLoad_ValidConfig() {
	configContent := `
coverage:
  threshold: 90.0
files:
  required:
    - README.md
    - LICENSE
    - CONTRIBUTING.md
checks:
  disabled:
    - gofmt
    - govet
external:
  - id: custom-check
    name: Custom Check
    command: echo
    args: ["test"]
    severity: warn
  - id: security
    name: Security Scan
    command: gosec
    args: ["./..."]
    severity: fail
    source_dir: backend
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	suite.Equal(90.0, cfg.Coverage.Threshold)
	suite.Equal([]string{"README.md", "LICENSE", "CONTRIBUTING.md"}, cfg.Files.Required)
	suite.Equal([]string{"gofmt", "govet"}, cfg.Checks.Disabled)
	suite.Len(cfg.External, 2)
	suite.Equal("custom-check", cfg.External[0].ID)
	suite.Equal("Custom Check", cfg.External[0].Name)
	suite.Equal("echo", cfg.External[0].Command)
	suite.Equal([]string{"test"}, cfg.External[0].Args)
	suite.Equal("warn", cfg.External[0].Severity)
	suite.Equal("", cfg.External[0].SourceDir)
	suite.Equal("security", cfg.External[1].ID)
	suite.Equal("backend", cfg.External[1].SourceDir)
}

// TestLoad_InvalidYAML tests that Load handles invalid YAML gracefully.
func (suite *ConfigTestSuite) TestLoad_InvalidYAML() {
	configContent := `
coverage:
  threshold: invalid
  - broken yaml
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.Error(err)
	suite.Nil(cfg)
}

// TestLoad_PartialConfig tests that Load merges partial config with defaults.
func (suite *ConfigTestSuite) TestLoad_PartialConfig() {
	configContent := `
coverage:
  threshold: 75.0
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	suite.Equal(75.0, cfg.Coverage.Threshold)
	// Should still have default files
	suite.Equal([]string{"README.md", "LICENSE"}, cfg.Files.Required)
}

// TestLoad_EmptyConfig tests that Load handles empty config file.
func (suite *ConfigTestSuite) TestLoad_EmptyConfig() {
	suite.createTempFile(".a2.yaml", "")

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	// Should use defaults
	suite.Equal(80.0, cfg.Coverage.Threshold)
}

// TestIsCheckDisabled_DisabledCheck tests that IsCheckDisabled returns true for disabled checks.
func (suite *ConfigTestSuite) TestIsCheckDisabled_DisabledCheck() {
	cfg := &Config{
		Checks: ChecksConfig{
			Disabled: []string{"gofmt", "govet", "coverage"},
		},
	}

	suite.True(cfg.IsCheckDisabled("gofmt"))
	suite.True(cfg.IsCheckDisabled("govet"))
	suite.True(cfg.IsCheckDisabled("coverage"))
}

// TestIsCheckDisabled_EnabledCheck tests that IsCheckDisabled returns false for enabled checks.
func (suite *ConfigTestSuite) TestIsCheckDisabled_EnabledCheck() {
	cfg := &Config{
		Checks: ChecksConfig{
			Disabled: []string{"gofmt"},
		},
	}

	suite.False(cfg.IsCheckDisabled("govet"))
	suite.False(cfg.IsCheckDisabled("coverage"))
	suite.False(cfg.IsCheckDisabled("build"))
}

// TestIsCheckDisabled_EmptyList tests that IsCheckDisabled handles empty disabled list.
func (suite *ConfigTestSuite) TestIsCheckDisabled_EmptyList() {
	cfg := &Config{
		Checks: ChecksConfig{
			Disabled: []string{},
		},
	}

	suite.False(cfg.IsCheckDisabled("gofmt"))
	suite.False(cfg.IsCheckDisabled("any-check"))
}

// TestIsCheckDisabled_CaseSensitive tests that IsCheckDisabled is case-sensitive.
func (suite *ConfigTestSuite) TestIsCheckDisabled_CaseSensitive() {
	cfg := &Config{
		Checks: ChecksConfig{
			Disabled: []string{"gofmt"},
		},
	}

	suite.True(cfg.IsCheckDisabled("gofmt"))
	suite.False(cfg.IsCheckDisabled("Gofmt"))
	suite.False(cfg.IsCheckDisabled("GOFMT"))
}

// TestIsCheckDisabled_CommonAliases tests that short names work for common checks.
func (suite *ConfigTestSuite) TestIsCheckDisabled_CommonAliases() {
	cfg := &Config{
		Checks: ChecksConfig{
			Disabled: []string{"license", "k8s", "health"},
		},
	}

	// Short names should disable the full check IDs
	suite.True(cfg.IsCheckDisabled("common:license"))
	suite.True(cfg.IsCheckDisabled("devops:k8s"))
	suite.True(cfg.IsCheckDisabled("common:health"))

	// Other checks should still be enabled
	suite.False(cfg.IsCheckDisabled("common:secrets"))
	suite.False(cfg.IsCheckDisabled("common:sast"))
}

// TestIsCheckDisabled_GoAliases tests that legacy Go aliases work.
func (suite *ConfigTestSuite) TestIsCheckDisabled_GoAliases() {
	cfg := &Config{
		Checks: ChecksConfig{
			Disabled: []string{"gofmt", "govet"},
		},
	}

	// Legacy names should disable the new check IDs
	suite.True(cfg.IsCheckDisabled("go:format"))
	suite.True(cfg.IsCheckDisabled("go:vet"))

	// Other checks should still be enabled
	suite.False(cfg.IsCheckDisabled("go:build"))
	suite.False(cfg.IsCheckDisabled("go:tests"))
}

// TestGetSourceDirsForLang tests that GetSourceDirsForLang returns the correct source directories per language.
func (suite *ConfigTestSuite) TestGetSourceDirsForLang() {
	cfg := &Config{
		Language: LanguageConfig{
			Go:         GoLanguageConfig{SourceDir: StringOrSlice{"backend/go"}},
			Rust:       RustLanguageConfig{SourceDir: StringOrSlice{"src-tauri"}},
			TypeScript: TypeScriptLanguageConfig{SourceDir: StringOrSlice{"frontend"}},
		},
	}

	suite.Equal([]string{"backend/go"}, cfg.GetSourceDirsForLang("go"))
	suite.Equal([]string{"src-tauri"}, cfg.GetSourceDirsForLang("rust"))
	suite.Equal([]string{"frontend"}, cfg.GetSourceDirsForLang("typescript"))
	suite.Nil(cfg.GetSourceDirsForLang("python"))  // Not configured
	suite.Nil(cfg.GetSourceDirsForLang("node"))    // Not configured
	suite.Nil(cfg.GetSourceDirsForLang("java"))    // Not configured
	suite.Nil(cfg.GetSourceDirsForLang("unknown")) // Unknown language
}

// TestGetSourceDirsForLang_Multiple tests that GetSourceDirsForLang returns multiple directories.
func (suite *ConfigTestSuite) TestGetSourceDirsForLang_Multiple() {
	cfg := &Config{
		Language: LanguageConfig{
			Go: GoLanguageConfig{SourceDir: StringOrSlice{"api", "agent"}},
		},
	}

	suite.Equal([]string{"api", "agent"}, cfg.GetSourceDirsForLang("go"))
}

// TestGetSourceDirs tests that GetSourceDirs returns a map of all configured source directories.
func (suite *ConfigTestSuite) TestGetSourceDirs() {
	cfg := &Config{
		Language: LanguageConfig{
			Go:   GoLanguageConfig{SourceDir: StringOrSlice{"backend"}},
			Rust: RustLanguageConfig{SourceDir: StringOrSlice{"src-tauri"}},
		},
	}

	dirs := cfg.GetSourceDirs()

	suite.Len(dirs, 2)
	suite.Equal([]string{"backend"}, dirs["go"])
	suite.Equal([]string{"src-tauri"}, dirs["rust"])

	// Languages without source_dir should not be in the map
	_, hasNode := dirs["node"]
	suite.False(hasNode)
}

// TestGetSourceDirs_Empty tests that GetSourceDirs returns empty map when no source directories are configured.
func (suite *ConfigTestSuite) TestGetSourceDirs_Empty() {
	cfg := DefaultConfig()

	dirs := cfg.GetSourceDirs()

	suite.Empty(dirs)
}

// TestLoad_WithSourceDir_String tests that Load parses source_dir as a single string.
func (suite *ConfigTestSuite) TestLoad_WithSourceDir_String() {
	configContent := `
language:
  rust:
    source_dir: src-tauri
  node:
    source_dir: frontend
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	suite.Equal(StringOrSlice{"src-tauri"}, cfg.Language.Rust.SourceDir)
	suite.Equal(StringOrSlice{"frontend"}, cfg.Language.Node.SourceDir)
	suite.Nil(cfg.Language.Go.SourceDir) // Not configured
}

// TestLoad_WithSourceDir_List tests that Load parses source_dir as a list of strings.
func (suite *ConfigTestSuite) TestLoad_WithSourceDir_List() {
	configContent := `
language:
  go:
    source_dir: [api, agent]
  node:
    source_dir:
      - frontend
      - admin
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	suite.Equal(StringOrSlice{"api", "agent"}, cfg.Language.Go.SourceDir)
	suite.Equal(StringOrSlice{"frontend", "admin"}, cfg.Language.Node.SourceDir)
}

// TestGetToolRunByDefault_NotConfigured tests that GetToolRunByDefault returns nil when no tool override.
func (suite *ConfigTestSuite) TestGetToolRunByDefault_NotConfigured() {
	cfg := DefaultConfig()

	result := cfg.GetToolRunByDefault("gitleaks")
	suite.Nil(result)
}

// TestGetToolRunByDefault_ConfiguredTrue tests that GetToolRunByDefault returns true override.
func (suite *ConfigTestSuite) TestGetToolRunByDefault_ConfiguredTrue() {
	runByDefault := true
	cfg := &Config{
		Tools: map[string]ToolConfig{
			"gitleaks": {RunByDefault: &runByDefault},
		},
	}

	result := cfg.GetToolRunByDefault("gitleaks")
	suite.NotNil(result)
	suite.True(*result)
}

// TestGetToolRunByDefault_ConfiguredFalse tests that GetToolRunByDefault returns false override.
func (suite *ConfigTestSuite) TestGetToolRunByDefault_ConfiguredFalse() {
	runByDefault := false
	cfg := &Config{
		Tools: map[string]ToolConfig{
			"gitleaks": {RunByDefault: &runByDefault},
		},
	}

	result := cfg.GetToolRunByDefault("gitleaks")
	suite.NotNil(result)
	suite.False(*result)
}

// TestGetToolRunByDefault_DifferentTool tests that GetToolRunByDefault only affects the specified tool.
func (suite *ConfigTestSuite) TestGetToolRunByDefault_DifferentTool() {
	runByDefault := false
	cfg := &Config{
		Tools: map[string]ToolConfig{
			"gitleaks": {RunByDefault: &runByDefault},
		},
	}

	// The configured tool should have an override
	gitleaksResult := cfg.GetToolRunByDefault("gitleaks")
	suite.NotNil(gitleaksResult)

	// Other tools should not have an override
	semgrepResult := cfg.GetToolRunByDefault("semgrep")
	suite.Nil(semgrepResult)
}

// TestLoad_WithToolOverrides tests that Load parses tool configuration.
func (suite *ConfigTestSuite) TestLoad_WithToolOverrides() {
	configContent := `
tools:
  gitleaks:
    run_by_default: false
  semgrep:
    run_by_default: true
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)

	gitleaksOverride := cfg.GetToolRunByDefault("gitleaks")
	suite.NotNil(gitleaksOverride)
	suite.False(*gitleaksOverride)

	semgrepOverride := cfg.GetToolRunByDefault("semgrep")
	suite.NotNil(semgrepOverride)
	suite.True(*semgrepOverride)

	// Tools not in config should return nil
	trivyOverride := cfg.GetToolRunByDefault("trivy")
	suite.Nil(trivyOverride)
}

// TestConfigTestSuite runs all the tests in the suite.
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
