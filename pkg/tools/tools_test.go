package tools

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ToolsTestSuite struct {
	suite.Suite
}

func (s *ToolsTestSuite) TestByName_Found() {
	tool := ByName("gitleaks")
	s.NotNil(tool)
	s.Equal("gitleaks", tool.Name)
	s.Equal(checker.LangCommon, tool.Language)
}

func (s *ToolsTestSuite) TestByName_NotFound() {
	tool := ByName("nonexistent-tool")
	s.Nil(tool)
}

func (s *ToolsTestSuite) TestByLanguage() {
	goTools := ByLanguage(checker.LangGo)
	s.NotEmpty(goTools)

	for _, t := range goTools {
		s.Equal(checker.LangGo, t.Language)
	}
}

func (s *ToolsTestSuite) TestForLanguages() {
	langs := []checker.Language{checker.LangGo, checker.LangPython}
	tools := ForLanguages(langs)
	s.NotEmpty(tools)

	// Should include Go, Python, and Common tools
	hasGo := false
	hasPython := false
	hasCommon := false
	for _, t := range tools {
		switch t.Language {
		case checker.LangGo:
			hasGo = true
		case checker.LangPython:
			hasPython = true
		case checker.LangCommon:
			hasCommon = true
		}
	}
	s.True(hasGo, "should include Go tools")
	s.True(hasPython, "should include Python tools")
	s.True(hasCommon, "should include Common tools")
}

func (s *ToolsTestSuite) TestShouldRunByDefault_WithConfigOverrideTrue() {
	override := true
	result := ShouldRunByDefault("gitleaks", &override)
	s.True(result)
}

func (s *ToolsTestSuite) TestShouldRunByDefault_WithConfigOverrideFalse() {
	override := false
	result := ShouldRunByDefault("gitleaks", &override)
	s.False(result)
}

func (s *ToolsTestSuite) TestShouldRunByDefault_NoOverride_SecurityTool() {
	// gitleaks has RunByDefault: true in registry
	result := ShouldRunByDefault("gitleaks", nil)
	s.True(result, "security tools should run by default")
}

func (s *ToolsTestSuite) TestShouldRunByDefault_NoOverride_SlowTool() {
	// cargo-tarpaulin has RunByDefault: false in registry
	result := ShouldRunByDefault("cargo-tarpaulin", nil)
	s.False(result, "slow tools should not run by default")
}

func (s *ToolsTestSuite) TestShouldRunByDefault_UnknownTool() {
	result := ShouldRunByDefault("unknown-tool", nil)
	s.False(result, "unknown tools should not run by default")
}

func (s *ToolsTestSuite) TestRunByDefaultValues() {
	// Verify expected RunByDefault values for key tools
	testCases := []struct {
		name     string
		expected bool
		reason   string
	}{
		// Security tools - should run by default
		{"gitleaks", true, "security tool"},
		{"semgrep", true, "security tool"},
		{"trivy", true, "security tool"},
		{"govulncheck", true, "security tool"},
		{"pip-audit", true, "security tool"},
		{"cargo-audit", true, "security tool"},

		// Fast linters - should run by default
		{"ruff", true, "fast linter"},
		{"eslint", true, "linter"},
		{"biome", true, "fast linter"},
		{"swiftlint", true, "linter"},

		// Formatters - should run by default
		{"black", true, "formatter"},
		{"prettier", true, "formatter"},
		{"swift-format", true, "formatter"},

		// Slow/heavy tools - should NOT run by default
		{"cargo-tarpaulin", false, "slow coverage tool"},

		// Tools needing config - should NOT run by default
		{"mypy", false, "needs project config"},
		{"pytest", false, "tests should be explicit"},
	}

	for _, tc := range testCases {
		tool := ByName(tc.name)
		s.NotNil(tool, "tool %s should exist in registry", tc.name)
		s.Equal(tc.expected, tool.RunByDefault, "%s (%s) RunByDefault mismatch", tc.name, tc.reason)
	}
}

func TestToolsTestSuite(t *testing.T) {
	suite.Run(t, new(ToolsTestSuite))
}
