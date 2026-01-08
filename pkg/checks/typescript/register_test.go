package typescriptcheck

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

	// Should have 9 TypeScript checks
	s.Len(checks, 9)
}

func (s *RegisterTestSuite) TestRegister_CheckIDs() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedIDs := []string{
		"typescript:project",
		"typescript:build",
		"typescript:tests",
		"typescript:format",
		"typescript:lint",
		"typescript:type",
		"typescript:coverage",
		"typescript:deps",
		"typescript:logging",
	}

	for i, check := range checks {
		s.Equal(expectedIDs[i], check.Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_CheckOrder() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedOrders := []int{100, 110, 120, 200, 210, 215, 220, 230, 250}

	for i, check := range checks {
		s.Equal(expectedOrders[i], check.Meta.Order)
	}
}

func (s *RegisterTestSuite) TestRegister_CriticalChecks() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// project, build, tests, and type should be critical
	s.True(checks[0].Meta.Critical, "project should be critical")
	s.True(checks[1].Meta.Critical, "build should be critical")
	s.True(checks[2].Meta.Critical, "tests should be critical")
	s.False(checks[3].Meta.Critical, "format should not be critical")
	s.False(checks[4].Meta.Critical, "lint should not be critical")
	s.True(checks[5].Meta.Critical, "type should be critical")
	s.False(checks[6].Meta.Critical, "coverage should not be critical")
	s.False(checks[7].Meta.Critical, "deps should not be critical")
	s.False(checks[8].Meta.Critical, "logging should not be critical")
}

func (s *RegisterTestSuite) TestRegister_AllTypeScriptLanguage() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		s.Contains(check.Meta.Languages, checker.LangTypeScript)
	}
}

func (s *RegisterTestSuite) TestRegister_CoverageThreshold() {
	cfg := config.DefaultConfig()
	cfg.Language.TypeScript.CoverageThreshold = 90.0

	checks := Register(cfg)

	// Find coverage check
	for _, check := range checks {
		if check.Meta.ID == "typescript:coverage" {
			coverageCheck, ok := check.Checker.(*CoverageCheck)
			s.True(ok)
			s.Equal(90.0, coverageCheck.Threshold)
			return
		}
	}
	s.Fail("coverage check not found")
}

func (s *RegisterTestSuite) TestRegister_DefaultCoverageThreshold() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// Find coverage check
	for _, check := range checks {
		if check.Meta.ID == "typescript:coverage" {
			coverageCheck, ok := check.Checker.(*CoverageCheck)
			s.True(ok)
			s.Equal(80.0, coverageCheck.Threshold) // Default threshold
			return
		}
	}
	s.Fail("coverage check not found")
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
