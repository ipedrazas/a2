package nodecheck

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

// RegisterTestSuite is the test suite for the register functionality.
type RegisterTestSuite struct {
	suite.Suite
}

// TestRegister_ReturnsAllChecks tests that Register returns all 7 checks.
func (suite *RegisterTestSuite) TestRegister_ReturnsAllChecks() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	suite.Len(checks, 7)
}

// TestRegister_CheckIDs tests that all check IDs are correct.
func (suite *RegisterTestSuite) TestRegister_CheckIDs() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedIDs := []string{
		"node:project",
		"node:build",
		"node:tests",
		"node:format",
		"node:lint",
		"node:coverage",
		"node:deps",
	}

	for i, check := range checks {
		suite.Equal(expectedIDs[i], check.Meta.ID)
	}
}

// TestRegister_CriticalChecks tests that critical checks are properly marked.
func (suite *RegisterTestSuite) TestRegister_CriticalChecks() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// Project, Build, and Tests should be critical
	suite.True(checks[0].Meta.Critical, "node:project should be critical")
	suite.True(checks[1].Meta.Critical, "node:build should be critical")
	suite.True(checks[2].Meta.Critical, "node:tests should be critical")

	// Format, Lint, Coverage, and Deps should not be critical
	suite.False(checks[3].Meta.Critical, "node:format should not be critical")
	suite.False(checks[4].Meta.Critical, "node:lint should not be critical")
	suite.False(checks[5].Meta.Critical, "node:coverage should not be critical")
	suite.False(checks[6].Meta.Critical, "node:deps should not be critical")
}

// TestRegister_CheckOrder tests that checks are ordered correctly.
func (suite *RegisterTestSuite) TestRegister_CheckOrder() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// Critical checks should have order 100-120
	suite.Equal(100, checks[0].Meta.Order)
	suite.Equal(110, checks[1].Meta.Order)
	suite.Equal(120, checks[2].Meta.Order)

	// Non-critical checks should have order 200-230
	suite.Equal(200, checks[3].Meta.Order)
	suite.Equal(210, checks[4].Meta.Order)
	suite.Equal(220, checks[5].Meta.Order)
	suite.Equal(230, checks[6].Meta.Order)
}

// TestRegister_LanguageIsNode tests that all checks are for Node language.
func (suite *RegisterTestSuite) TestRegister_LanguageIsNode() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		suite.Contains(check.Meta.Languages, checker.LangNode)
	}
}

// TestRegister_CoverageThresholdFromConfig tests that coverage threshold is read from config.
func (suite *RegisterTestSuite) TestRegister_CoverageThresholdFromConfig() {
	cfg := config.DefaultConfig()
	cfg.Language.Node.CoverageThreshold = 75.0

	checks := Register(cfg)

	// Find coverage check
	var coverageCheck *CoverageCheck
	for _, check := range checks {
		if check.Meta.ID == "node:coverage" {
			coverageCheck = check.Checker.(*CoverageCheck)
			break
		}
	}

	suite.NotNil(coverageCheck)
	suite.Equal(75.0, coverageCheck.Threshold)
}

// TestRegister_CoverageThresholdFromGlobal tests that global threshold is used when language-specific is not set.
func (suite *RegisterTestSuite) TestRegister_CoverageThresholdFromGlobal() {
	cfg := config.DefaultConfig()
	cfg.Coverage.Threshold = 60.0
	cfg.Language.Node.CoverageThreshold = 0 // Not set

	checks := Register(cfg)

	// Find coverage check
	var coverageCheck *CoverageCheck
	for _, check := range checks {
		if check.Meta.ID == "node:coverage" {
			coverageCheck = check.Checker.(*CoverageCheck)
			break
		}
	}

	suite.NotNil(coverageCheck)
	suite.Equal(60.0, coverageCheck.Threshold)
}

// TestRegister_CheckerImplementsInterface tests that all checkers implement the Checker interface.
func (suite *RegisterTestSuite) TestRegister_CheckerImplementsInterface() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		// Verify each checker implements the interface
		var _ checker.Checker = check.Checker
		suite.NotEmpty(check.Checker.ID())
		suite.NotEmpty(check.Checker.Name())
	}
}

// TestRegisterTestSuite runs all the tests in the suite.
func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
