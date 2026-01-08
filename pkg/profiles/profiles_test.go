package profiles

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProfilesTestSuite struct {
	suite.Suite
}

func (s *ProfilesTestSuite) TestGet_PoCProfile() {
	profile, ok := Get("poc")
	s.True(ok)
	s.Equal("poc", profile.Name)
	s.Contains(profile.Description, "Proof of Concept")
	s.NotEmpty(profile.Disabled)

	// Verify key checks are disabled
	s.Contains(profile.Disabled, "common:license")
	s.Contains(profile.Disabled, "common:k8s")
	s.Contains(profile.Disabled, "common:health")
	s.Contains(profile.Disabled, "go:coverage")
}

func (s *ProfilesTestSuite) TestGet_LibraryProfile() {
	profile, ok := Get("library")
	s.True(ok)
	s.Equal("library", profile.Name)
	s.Contains(profile.Description, "Library")
	s.NotEmpty(profile.Disabled)

	// Verify deployment checks are disabled
	s.Contains(profile.Disabled, "common:dockerfile")
	s.Contains(profile.Disabled, "common:k8s")
	s.Contains(profile.Disabled, "common:health")

	// Verify security checks are NOT disabled
	s.NotContains(profile.Disabled, "common:sast")
	s.NotContains(profile.Disabled, "common:secrets")
}

func (s *ProfilesTestSuite) TestGet_ProductionProfile() {
	profile, ok := Get("production")
	s.True(ok)
	s.Equal("production", profile.Name)
	s.Contains(profile.Description, "Production")
	s.Empty(profile.Disabled) // All checks enabled
}

func (s *ProfilesTestSuite) TestGet_UnknownProfile() {
	_, ok := Get("unknown")
	s.False(ok)
}

func (s *ProfilesTestSuite) TestList_ReturnsAllProfiles() {
	profiles := List()
	s.Len(profiles, 3)

	// Verify order
	s.Equal("poc", profiles[0].Name)
	s.Equal("library", profiles[1].Name)
	s.Equal("production", profiles[2].Name)
}

func (s *ProfilesTestSuite) TestNames_ReturnsAllNames() {
	names := Names()
	s.Len(names, 3)
	s.Contains(names, "poc")
	s.Contains(names, "library")
	s.Contains(names, "production")
}

func (s *ProfilesTestSuite) TestPoCProfile_DisablesNonCriticalChecks() {
	profile, _ := Get("poc")

	// Common non-critical checks
	expectedDisabled := []string{
		"common:license", "common:sast", "common:k8s", "common:shutdown",
		"common:health", "common:api_docs", "common:changelog",
		"common:integration", "common:metrics", "common:errors",
		"common:precommit", "common:env",
	}

	for _, check := range expectedDisabled {
		s.Contains(profile.Disabled, check, "PoC profile should disable %s", check)
	}

	// Critical checks should NOT be disabled
	s.NotContains(profile.Disabled, "go:module")
	s.NotContains(profile.Disabled, "go:build")
	s.NotContains(profile.Disabled, "go:tests")
	s.NotContains(profile.Disabled, "python:project")
	s.NotContains(profile.Disabled, "python:build")
	s.NotContains(profile.Disabled, "python:tests")
}

func (s *ProfilesTestSuite) TestLibraryProfile_KeepsCodeQualityChecks() {
	profile, _ := Get("library")

	// Code quality checks should NOT be disabled
	s.NotContains(profile.Disabled, "go:format")
	s.NotContains(profile.Disabled, "go:vet")
	s.NotContains(profile.Disabled, "go:coverage")
	s.NotContains(profile.Disabled, "python:format")
	s.NotContains(profile.Disabled, "python:lint")
	s.NotContains(profile.Disabled, "common:sast")
	s.NotContains(profile.Disabled, "common:license")
}

func TestProfilesTestSuite(t *testing.T) {
	suite.Run(t, new(ProfilesTestSuite))
}
