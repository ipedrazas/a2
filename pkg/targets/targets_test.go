package targets

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TargetsTestSuite struct {
	suite.Suite
}

func (s *TargetsTestSuite) TestGet_PoCTarget() {
	target, ok := Get("poc")
	s.True(ok)
	s.Equal("poc", target.Name)
	s.Contains(target.Description, "Proof of Concept")
	s.NotEmpty(target.Disabled)

	// Verify key checks are disabled
	s.Contains(target.Disabled, "common:license")
	s.Contains(target.Disabled, "common:sast")
	s.Contains(target.Disabled, "go:coverage")
}

func (s *TargetsTestSuite) TestGet_ProductionTarget() {
	target, ok := Get("production")
	s.True(ok)
	s.Equal("production", target.Name)
	s.Contains(target.Description, "Production")
	s.Empty(target.Disabled) // All checks enabled
}

func (s *TargetsTestSuite) TestGet_UnknownTarget() {
	_, ok := Get("unknown")
	s.False(ok)
}

func (s *TargetsTestSuite) TestList_ReturnsAllTargets() {
	targets := List()
	s.Len(targets, 2)

	// Verify order
	s.Equal("poc", targets[0].Name)
	s.Equal("production", targets[1].Name)
}

func (s *TargetsTestSuite) TestNames_ReturnsAllNames() {
	names := Names()
	s.Len(names, 2)
	s.Contains(names, "poc")
	s.Contains(names, "production")
}

func (s *TargetsTestSuite) TestPoCTarget_DisablesNonCriticalChecks() {
	target, _ := Get("poc")

	// Common non-critical checks
	expectedDisabled := []string{
		"common:license", "common:sast", "common:changelog",
		"common:precommit", "common:env", "common:contributing",
	}

	for _, check := range expectedDisabled {
		s.Contains(target.Disabled, check, "PoC target should disable %s", check)
	}

	// Critical checks should NOT be disabled
	s.NotContains(target.Disabled, "go:module")
	s.NotContains(target.Disabled, "go:build")
	s.NotContains(target.Disabled, "go:tests")
	s.NotContains(target.Disabled, "python:project")
	s.NotContains(target.Disabled, "python:build")
	s.NotContains(target.Disabled, "python:tests")
}

func TestTargetsTestSuite(t *testing.T) {
	suite.Run(t, new(TargetsTestSuite))
}
