package checks

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/stretchr/testify/suite"
)

// RegistryTestSuite is the test suite for the registry package.
type RegistryTestSuite struct {
	suite.Suite
}

// SetupTest is called before each test method.
func (suite *RegistryTestSuite) SetupTest() {
	// Setup code if needed
}

// TestGetChecks_AllEnabled tests that GetChecks returns all checks when none are disabled.
func (suite *RegistryTestSuite) TestGetChecks_AllEnabled() {
	cfg := config.DefaultConfig()
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}
	checks := GetChecks(cfg, detected)

	suite.NotEmpty(checks)
	// Should have at least the default Go checks + common checks
	suite.GreaterOrEqual(len(checks), 7) // Module, Build, Tests, Format, Vet, Coverage, Deps + Files

	// Verify critical checks are present
	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	suite.True(checkIDs["go:module"], "go:module check should be present")
	suite.True(checkIDs["go:build"], "go:build check should be present")
	suite.True(checkIDs["go:tests"], "go:tests check should be present")
}

// TestGetChecks_FiltersDisabled tests that GetChecks filters out disabled checks.
func (suite *RegistryTestSuite) TestGetChecks_FiltersDisabled() {
	cfg := &config.Config{
		Checks: config.ChecksConfig{
			Disabled: []string{"go:format", "go:vet", "go:coverage"},
		},
		Files: config.FilesConfig{
			Required: []string{"README.md"},
		},
		Coverage: config.CoverageConfig{
			Threshold: 80.0,
		},
		Language: config.LanguageConfig{
			AutoDetect: true,
		},
	}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}

	checks := GetChecks(cfg, detected)

	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	suite.False(checkIDs["go:format"], "go:format should be disabled")
	suite.False(checkIDs["go:vet"], "go:vet should be disabled")
	suite.False(checkIDs["go:coverage"], "go:coverage should be disabled")
	suite.True(checkIDs["go:module"], "go:module should still be enabled")
	suite.True(checkIDs["go:build"], "go:build should still be enabled")
}

// TestGetChecks_BackwardCompatibility tests that old check IDs still work for disabling.
func (suite *RegistryTestSuite) TestGetChecks_BackwardCompatibility() {
	cfg := &config.Config{
		Checks: config.ChecksConfig{
			// Use old check IDs
			Disabled: []string{"gofmt", "govet", "coverage"},
		},
		Files: config.FilesConfig{
			Required: []string{"README.md"},
		},
		Coverage: config.CoverageConfig{
			Threshold: 80.0,
		},
		Language: config.LanguageConfig{
			AutoDetect: true,
		},
	}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}

	checks := GetChecks(cfg, detected)

	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	// Old IDs should disable new checks
	suite.False(checkIDs["go:format"], "go:format should be disabled via alias 'gofmt'")
	suite.False(checkIDs["go:vet"], "go:vet should be disabled via alias 'govet'")
	suite.False(checkIDs["go:coverage"], "go:coverage should be disabled via alias 'coverage'")
}

// TestGetChecks_IncludesExternal tests that GetChecks includes external checks from config.
func (suite *RegistryTestSuite) TestGetChecks_IncludesExternal() {
	cfg := &config.Config{
		External: []config.ExternalCheck{
			{
				ID:       "custom-1",
				Name:     "Custom Check 1",
				Command:  "echo",
				Args:     []string{"test"},
				Severity: "warn",
			},
			{
				ID:       "custom-2",
				Name:     "Custom Check 2",
				Command:  "true",
				Args:     []string{},
				Severity: "fail",
			},
		},
		Files: config.FilesConfig{
			Required: []string{"README.md"},
		},
		Coverage: config.CoverageConfig{
			Threshold: 80.0,
		},
		Language: config.LanguageConfig{
			AutoDetect: true,
		},
	}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}

	checks := GetChecks(cfg, detected)

	checkIDs := make(map[string]bool)
	checkNames := make(map[string]string)
	for _, check := range checks {
		checkIDs[check.ID()] = true
		checkNames[check.ID()] = check.Name()
	}

	suite.True(checkIDs["custom-1"], "custom-1 should be included")
	suite.True(checkIDs["custom-2"], "custom-2 should be included")
	suite.Equal("Custom Check 1", checkNames["custom-1"])
	suite.Equal("Custom Check 2", checkNames["custom-2"])
}

// TestGetChecks_EmptyConfig tests that GetChecks handles empty config.
func (suite *RegistryTestSuite) TestGetChecks_EmptyConfig() {
	cfg := &config.Config{}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}

	checks := GetChecks(cfg, detected)

	suite.NotEmpty(checks)
	// Should still have default checks
	suite.GreaterOrEqual(len(checks), 3) // At least Module, Build, Tests
}

// TestGetChecks_Ordering tests that critical checks come first.
func (suite *RegistryTestSuite) TestGetChecks_Ordering() {
	cfg := config.DefaultConfig()
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}
	checks := GetChecks(cfg, detected)

	// Critical checks should be first
	suite.GreaterOrEqual(len(checks), 3)

	// First three should be critical (go:module, go:build, go:tests)
	criticalIDs := []string{"go:module", "go:build", "go:tests"}
	for i, expectedID := range criticalIDs {
		if i < len(checks) {
			suite.Equal(expectedID, checks[i].ID(), "Critical check %s should be at position %d", expectedID, i)
		}
	}
}

// TestGetChecks_AllDisabled tests that GetChecks handles all checks disabled.
func (suite *RegistryTestSuite) TestGetChecks_AllDisabled() {
	cfg := &config.Config{
		Checks: config.ChecksConfig{
			Disabled: []string{
				"go:module", "go:build", "go:tests", "file_exists",
				"go:format", "go:vet", "go:coverage", "go:deps",
			},
		},
		Files: config.FilesConfig{
			Required: []string{"README.md"},
		},
		Coverage: config.CoverageConfig{
			Threshold: 80.0,
		},
		Language: config.LanguageConfig{
			AutoDetect: true,
		},
	}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}

	checks := GetChecks(cfg, detected)

	// Should return empty or only external checks
	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	suite.False(checkIDs["go:module"])
	suite.False(checkIDs["go:build"])
	suite.False(checkIDs["go:tests"])
}

// TestDefaultChecks tests that DefaultChecks returns checks based on detected language.
// When run from a directory without language indicator files (like go.mod),
// only common checks are returned.
func (suite *RegistryTestSuite) TestDefaultChecks() {
	checks := DefaultChecks()

	// Should at least have common checks
	suite.NotEmpty(checks)

	// Verify common checks are present
	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	// Common checks should always be present
	suite.True(checkIDs["file_exists"], "file_exists should be present as common check")
}

// TestGetChecks_ExternalWithDisabled tests that external checks respect disabled list.
func (suite *RegistryTestSuite) TestGetChecks_ExternalWithDisabled() {
	cfg := &config.Config{
		External: []config.ExternalCheck{
			{
				ID:       "external-1",
				Name:     "External 1",
				Command:  "echo",
				Args:     []string{"test"},
				Severity: "warn",
			},
		},
		Checks: config.ChecksConfig{
			Disabled: []string{"external-1"},
		},
		Files: config.FilesConfig{
			Required: []string{"README.md"},
		},
		Coverage: config.CoverageConfig{
			Threshold: 80.0,
		},
		Language: config.LanguageConfig{
			AutoDetect: true,
		},
	}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
		Primary:   checker.LangGo,
	}

	checks := GetChecks(cfg, detected)

	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	suite.False(checkIDs["external-1"], "external-1 should be disabled")
}

// TestGetChecks_NoLanguageDetected tests that only common checks are returned when no language is detected.
func (suite *RegistryTestSuite) TestGetChecks_NoLanguageDetected() {
	cfg := config.DefaultConfig()
	detected := language.DetectionResult{
		Languages: []checker.Language{}, // No languages detected
	}

	// GetChecks returns only common checks if no language specified
	checks := GetChecks(cfg, detected)

	// Should only have common checks (file_exists)
	checkIDs := make(map[string]bool)
	for _, check := range checks {
		checkIDs[check.ID()] = true
	}

	suite.True(checkIDs["file_exists"], "file_exists should be present as common check")
	suite.False(checkIDs["go:module"], "go:module should not be present without language detection")
}

// TestRegistryTestSuite runs all the tests in the suite.
func TestRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(RegistryTestSuite))
}
