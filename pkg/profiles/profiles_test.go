package profiles

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProfilesTestSuite struct {
	suite.Suite
}

func (s *ProfilesTestSuite) TestGet_CLIProfile() {
	profile, ok := Get("cli")
	s.True(ok)
	s.Equal("cli", profile.Name)
	s.Contains(profile.Description, "Command-line")
	s.NotEmpty(profile.Disabled)

	// Verify server-related checks are disabled
	s.Contains(profile.Disabled, "common:health")
	s.Contains(profile.Disabled, "common:k8s")
	s.Contains(profile.Disabled, "common:metrics")
	s.Contains(profile.Disabled, "common:api_docs")
	s.Contains(profile.Disabled, "common:tracing")
}

func (s *ProfilesTestSuite) TestGet_APIProfile() {
	profile, ok := Get("api")
	s.True(ok)
	s.Equal("api", profile.Name)
	s.Contains(profile.Description, "API")

	// API should have minimal disabled checks
	s.Len(profile.Disabled, 1)
	s.Contains(profile.Disabled, "common:e2e")

	// Operational checks should NOT be disabled
	s.NotContains(profile.Disabled, "common:health")
	s.NotContains(profile.Disabled, "common:metrics")
	s.NotContains(profile.Disabled, "common:api_docs")
	s.NotContains(profile.Disabled, "common:tracing")
}

func (s *ProfilesTestSuite) TestGet_LibraryProfile() {
	profile, ok := Get("library")
	s.True(ok)
	s.Equal("library", profile.Name)
	s.Contains(profile.Description, "library")
	s.NotEmpty(profile.Disabled)

	// Verify deployment checks are disabled
	s.Contains(profile.Disabled, "common:dockerfile")
	s.Contains(profile.Disabled, "common:k8s")
	s.Contains(profile.Disabled, "common:health")
	s.Contains(profile.Disabled, "common:metrics")

	// Code quality checks should NOT be disabled
	s.NotContains(profile.Disabled, "common:sast")
	s.NotContains(profile.Disabled, "common:secrets")
	s.NotContains(profile.Disabled, "common:license")
}

func (s *ProfilesTestSuite) TestGet_DesktopProfile() {
	profile, ok := Get("desktop")
	s.True(ok)
	s.Equal("desktop", profile.Name)
	s.Contains(profile.Description, "Desktop")
	s.NotEmpty(profile.Disabled)

	// Verify server-related checks are disabled
	s.Contains(profile.Disabled, "common:health")
	s.Contains(profile.Disabled, "common:k8s")
	s.Contains(profile.Disabled, "common:api_docs")
	s.Contains(profile.Disabled, "common:tracing")

	// E2E tests should NOT be disabled (desktop apps need E2E)
	s.NotContains(profile.Disabled, "common:e2e")
}

func (s *ProfilesTestSuite) TestGet_UnknownProfile() {
	_, ok := Get("unknown")
	s.False(ok)
}

func (s *ProfilesTestSuite) TestList_ReturnsAllProfiles() {
	profiles := List()
	s.Len(profiles, 4)

	// Verify all profiles are present (now sorted alphabetically)
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	s.Contains(names, "cli")
	s.Contains(names, "api")
	s.Contains(names, "library")
	s.Contains(names, "desktop")
}

func (s *ProfilesTestSuite) TestNames_ReturnsAllNames() {
	names := Names()
	s.Len(names, 4)
	s.Contains(names, "cli")
	s.Contains(names, "api")
	s.Contains(names, "library")
	s.Contains(names, "desktop")
}

func (s *ProfilesTestSuite) TestCLIProfile_KeepsCodeQualityChecks() {
	profile, _ := Get("cli")

	// Code quality checks should NOT be disabled
	s.NotContains(profile.Disabled, "common:sast")
	s.NotContains(profile.Disabled, "common:secrets")
	s.NotContains(profile.Disabled, "common:license")
	s.NotContains(profile.Disabled, "common:ci")
	s.NotContains(profile.Disabled, "common:precommit")
}

func (s *ProfilesTestSuite) TestLibraryProfile_KeepsDocumentation() {
	profile, _ := Get("library")

	// Documentation checks should NOT be disabled
	s.NotContains(profile.Disabled, "common:changelog")
	s.NotContains(profile.Disabled, "common:contributing")
}

func TestProfilesTestSuite(t *testing.T) {
	suite.Run(t, new(ProfilesTestSuite))
}
