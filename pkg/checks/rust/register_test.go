package rustcheck

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

	// Should have 8 Rust checks
	s.Len(checks, 8)
}

func (s *RegisterTestSuite) TestRegister_CheckIDs() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedIDs := []string{
		"rust:project",
		"rust:build",
		"rust:tests",
		"rust:format",
		"rust:lint",
		"rust:coverage",
		"rust:deps",
		"rust:logging",
	}

	for i, check := range checks {
		s.Equal(expectedIDs[i], check.Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_CheckOrder() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedOrders := []int{100, 110, 120, 200, 210, 220, 230, 250}

	for i, check := range checks {
		s.Equal(expectedOrders[i], check.Meta.Order)
	}
}

func (s *RegisterTestSuite) TestRegister_CriticalChecks() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	// First 3 checks (project, build, tests) should be critical
	s.True(checks[0].Meta.Critical, "project should be critical")
	s.True(checks[1].Meta.Critical, "build should be critical")
	s.True(checks[2].Meta.Critical, "tests should be critical")

	// Rest should not be critical
	for i := 3; i < len(checks); i++ {
		s.False(checks[i].Meta.Critical, "check %s should not be critical", checks[i].Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_AllRustLanguage() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		s.Contains(check.Meta.Languages, checker.LangRust)
	}
}

func (s *RegisterTestSuite) TestRegister_CoverageThreshold() {
	cfg := config.DefaultConfig()
	cfg.Language.Rust.CoverageThreshold = 90.0

	checks := Register(cfg)

	// Find coverage check
	for _, check := range checks {
		if check.Meta.ID == "rust:coverage" {
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
		if check.Meta.ID == "rust:coverage" {
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
