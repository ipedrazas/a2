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
`
	suite.createTempFile(".a2.yaml", configContent)

	cfg, err := Load(suite.tempDir)

	suite.NoError(err)
	suite.NotNil(cfg)
	suite.Equal(90.0, cfg.Coverage.Threshold)
	suite.Equal([]string{"README.md", "LICENSE", "CONTRIBUTING.md"}, cfg.Files.Required)
	suite.Equal([]string{"gofmt", "govet"}, cfg.Checks.Disabled)
	suite.Len(cfg.External, 1)
	suite.Equal("custom-check", cfg.External[0].ID)
	suite.Equal("Custom Check", cfg.External[0].Name)
	suite.Equal("echo", cfg.External[0].Command)
	suite.Equal([]string{"test"}, cfg.External[0].Args)
	suite.Equal("warn", cfg.External[0].Severity)
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
	suite.True(cfg.IsCheckDisabled("common:k8s"))
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

// TestConfigTestSuite runs all the tests in the suite.
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
