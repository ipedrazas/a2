package common

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	suite.Suite
}

func (s *RegisterTestSuite) TestRegister_ReturnsAllChecks() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// Should have 5 built-in checks (file_exists, dockerfile, ci, health, secrets)
	s.Len(checks, 5)
}

func (s *RegisterTestSuite) TestRegister_CheckIDs() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedIDs := []string{
		"file_exists",
		"common:dockerfile",
		"common:ci",
		"common:health",
		"common:secrets",
	}

	for i, check := range checks {
		s.Equal(expectedIDs[i], check.Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_CheckOrder() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedOrders := []int{900, 910, 920, 930, 940}

	for i, check := range checks {
		s.Equal(expectedOrders[i], check.Meta.Order)
	}
}

func (s *RegisterTestSuite) TestRegister_AllNonCritical() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		s.False(check.Meta.Critical, "check %s should not be critical", check.Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_AllCommonLanguage() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		s.Contains(check.Meta.Languages, checker.LangCommon)
	}
}

func (s *RegisterTestSuite) TestRegister_WithExternalChecks() {
	cfg := config.DefaultConfig()
	cfg.External = []config.ExternalCheck{
		{
			ID:       "custom:lint",
			Name:     "Custom Linter",
			Command:  "custom-lint",
			Args:     []string{"--fix"},
			Severity: "warn",
		},
		{
			ID:       "custom:security",
			Name:     "Security Scan",
			Command:  "sec-scan",
			Args:     []string{"."},
			Severity: "fail",
		},
	}

	checks := Register(cfg)

	// 5 built-in + 2 external
	s.Len(checks, 7)

	// Check external checks (indices 5 and 6 after 5 built-in checks)
	s.Equal("custom:lint", checks[5].Meta.ID)
	s.Equal("Custom Linter", checks[5].Meta.Name)
	s.False(checks[5].Meta.Critical) // severity: warn

	s.Equal("custom:security", checks[6].Meta.ID)
	s.Equal("Security Scan", checks[6].Meta.Name)
	s.True(checks[6].Meta.Critical) // severity: fail
}

func (s *RegisterTestSuite) TestRegister_ExternalCheckOrder() {
	cfg := config.DefaultConfig()
	cfg.External = []config.ExternalCheck{
		{
			ID:       "ext:test",
			Name:     "External Test",
			Command:  "test",
			Severity: "warn",
		},
	}

	checks := Register(cfg)

	// External checks should have order 1000 (index 5 after 5 built-in checks)
	s.Equal(1000, checks[5].Meta.Order)
}

func (s *RegisterTestSuite) TestRegister_FileExistsUsesConfig() {
	cfg := config.DefaultConfig()
	cfg.Files.Required = []string{"README.md", "LICENSE", "CHANGELOG.md"}

	checks := Register(cfg)

	// First check should be file_exists
	fileCheck, ok := checks[0].Checker.(*FileExistsCheck)
	s.True(ok)
	s.Equal([]string{"README.md", "LICENSE", "CHANGELOG.md"}, fileCheck.Files)
}

func (s *RegisterTestSuite) TestRegister_CheckerImplementsInterface() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		var _ checker.Checker = check.Checker
		s.NotEmpty(check.Checker.ID())
		s.NotEmpty(check.Checker.Name())
	}
}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
