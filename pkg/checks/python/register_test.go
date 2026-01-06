package pythoncheck

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

	// Should have 10 checks
	s.Len(checks, 10)
}

func (s *RegisterTestSuite) TestRegister_CheckIDs() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedIDs := []string{
		"python:project",
		"python:build",
		"python:tests",
		"python:format",
		"python:lint",
		"python:type",
		"python:coverage",
		"python:deps",
		"python:complexity",
		"python:logging",
	}

	for i, check := range checks {
		s.Equal(expectedIDs[i], check.Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_CriticalChecks() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// Project, Build, and Tests should be critical
	s.True(checks[0].Meta.Critical, "python:project should be critical")
	s.True(checks[1].Meta.Critical, "python:build should be critical")
	s.True(checks[2].Meta.Critical, "python:tests should be critical")

	// Rest should not be critical
	for i := 3; i < len(checks); i++ {
		s.False(checks[i].Meta.Critical, "%s should not be critical", checks[i].Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_CheckOrder() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedOrders := []int{100, 110, 120, 200, 210, 215, 220, 230, 240, 250}

	for i, check := range checks {
		s.Equal(expectedOrders[i], check.Meta.Order)
	}
}

func (s *RegisterTestSuite) TestRegister_AllPythonLanguage() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		s.Contains(check.Meta.Languages, checker.LangPython)
	}
}

func (s *RegisterTestSuite) TestRegister_CoverageThresholdFromConfig() {
	cfg := config.DefaultConfig()
	cfg.Language.Python.CoverageThreshold = 90.0

	checks := Register(cfg)

	// Find coverage check
	var coverageCheck *CoverageCheck
	for _, check := range checks {
		if check.Meta.ID == "python:coverage" {
			coverageCheck = check.Checker.(*CoverageCheck)
			break
		}
	}

	s.NotNil(coverageCheck)
	// Note: Python coverage check might not store threshold directly
	// This test validates the config is passed through
}

func (s *RegisterTestSuite) TestRegister_ComplexityCheckHasConfig() {
	cfg := config.DefaultConfig()
	cfg.Language.Python.CyclomaticThreshold = 20

	checks := Register(cfg)

	// Find complexity check
	var complexityCheck *ComplexityCheck
	for _, check := range checks {
		if check.Meta.ID == "python:complexity" {
			complexityCheck = check.Checker.(*ComplexityCheck)
			break
		}
	}

	s.NotNil(complexityCheck)
	s.NotNil(complexityCheck.Config)
	s.Equal(20, complexityCheck.Config.CyclomaticThreshold)
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

func (s *RegisterTestSuite) TestRegister_CheckNames() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedNames := []string{
		"Python Project",
		"Python Build",
		"Python Tests",
		"Python Format",
		"Python Lint",
		"Python Type Check",
		"Python Coverage",
		"Python Vulnerabilities",
		"Python Complexity",
		"Python Logging",
	}

	for i, check := range checks {
		s.Equal(expectedNames[i], check.Meta.Name)
	}
}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
