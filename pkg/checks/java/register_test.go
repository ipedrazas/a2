package javacheck

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

	// Should have 8 Java checks
	s.Len(checks, 8)
}

func (s *RegisterTestSuite) TestRegister_CheckIDs() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	expectedIDs := []string{
		"java:project",
		"java:build",
		"java:tests",
		"java:format",
		"java:lint",
		"java:coverage",
		"java:deps",
		"java:logging",
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

	// First 3 checks should be critical (project, build, tests)
	s.True(checks[0].Meta.Critical, "java:project should be critical")
	s.True(checks[1].Meta.Critical, "java:build should be critical")
	s.True(checks[2].Meta.Critical, "java:tests should be critical")

	// Rest should not be critical
	for i := 3; i < len(checks); i++ {
		s.False(checks[i].Meta.Critical, "check %s should not be critical", checks[i].Meta.ID)
	}
}

func (s *RegisterTestSuite) TestRegister_AllJavaLanguage() {
	cfg := config.DefaultConfig()

	checks := Register(cfg)

	for _, check := range checks {
		s.Contains(check.Meta.Languages, checker.LangJava)
	}
}

func (s *RegisterTestSuite) TestRegister_CoverageThreshold() {
	cfg := config.DefaultConfig()
	cfg.Language.Java.CoverageThreshold = 90.0

	checks := Register(cfg)

	// Find coverage check
	var coverageCheck *CoverageCheck
	for _, check := range checks {
		if check.Meta.ID == "java:coverage" {
			coverageCheck = check.Checker.(*CoverageCheck)
			break
		}
	}

	s.Require().NotNil(coverageCheck)
	s.Equal(90.0, coverageCheck.Threshold)
}

func (s *RegisterTestSuite) TestRegister_DefaultCoverageThreshold() {
	cfg := config.DefaultConfig()
	cfg.Language.Java.CoverageThreshold = 0 // Use global

	checks := Register(cfg)

	var coverageCheck *CoverageCheck
	for _, check := range checks {
		if check.Meta.ID == "java:coverage" {
			coverageCheck = check.Checker.(*CoverageCheck)
			break
		}
	}

	s.Require().NotNil(coverageCheck)
	s.Equal(cfg.Coverage.Threshold, coverageCheck.Threshold)
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
