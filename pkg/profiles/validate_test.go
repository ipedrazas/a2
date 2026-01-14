package profiles

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
	validCheckIDs map[string]bool
	validIDList   []string
}

func (suite *ValidateTestSuite) SetupTest() {
	// Set up a sample list of valid check IDs
	suite.validCheckIDs = map[string]bool{
		"go:build":        true,
		"go:tests":        true,
		"go:race":         true,
		"go:format":       true,
		"go:vet":          true,
		"go:coverage":     true,
		"common:health":   true,
		"common:k8s":      true,
		"common:secrets":  true,
		"common:ci":       true,
		"common:license":  true,
		"common:sast":     true,
		"python:build":    true,
		"python:tests":    true,
		"python:coverage": true,
	}
	suite.validIDList = make([]string, 0, len(suite.validCheckIDs))
	for id := range suite.validCheckIDs {
		suite.validIDList = append(suite.validIDList, id)
	}
}

func (suite *ValidateTestSuite) TestValidateProfile_ValidProfile() {
	profile := Profile{
		Name:        "test",
		Description: "Test profile",
		Disabled:    []string{"common:health", "common:k8s"},
		Source:      SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	suite.Empty(result.Errors)
	suite.Equal("test", result.Name)
}

func (suite *ValidateTestSuite) TestValidateProfile_EmptyDisabled() {
	profile := Profile{
		Name:        "empty",
		Description: "Profile with no disabled checks",
		Disabled:    []string{},
		Source:      SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	suite.Empty(result.Errors)
	suite.Empty(result.Warnings)
}

func (suite *ValidateTestSuite) TestValidateProfile_UnknownCheckID() {
	profile := Profile{
		Name:        "bad",
		Description: "Profile with unknown check",
		Disabled:    []string{"common:health", "unknown:check"},
		Source:      SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.False(result.Valid)
	suite.Len(result.Errors, 1)
	suite.Contains(result.Errors[0], "unknown check ID: unknown:check")
}

func (suite *ValidateTestSuite) TestValidateProfile_TypoWithSuggestion() {
	profile := Profile{
		Name:        "typo",
		Description: "Profile with typo in check ID",
		Disabled:    []string{"common:helth"}, // typo: should be common:health
		Source:      SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.False(result.Valid)
	suite.Len(result.Errors, 1)
	suite.Contains(result.Errors[0], "unknown check ID: common:helth")
	// Should have a suggestion warning
	suite.NotEmpty(result.Warnings)
	suite.Contains(result.Warnings[0], "did you mean")
}

func (suite *ValidateTestSuite) TestValidateProfile_DuplicateDisabled() {
	profile := Profile{
		Name:        "dupe",
		Description: "Profile with duplicate disabled checks",
		Disabled:    []string{"common:health", "common:k8s", "common:health"}, // duplicate
		Source:      SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid) // Duplicates are warnings, not errors
	suite.Empty(result.Errors)
	suite.NotEmpty(result.Warnings)
	suite.Contains(result.Warnings[0], "duplicate disabled check: common:health")
}

func (suite *ValidateTestSuite) TestValidateProfile_OverridesBuiltIn() {
	// Test that overriding a built-in profile produces a warning
	profile := Profile{
		Name:        "cli", // This is a built-in profile name
		Description: "Custom CLI profile",
		Disabled:    []string{"common:health"},
		Source:      SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	suite.Empty(result.Errors)
	suite.NotEmpty(result.Warnings)
	suite.Contains(result.Warnings[0], "overrides built-in profile: cli")
}

func (suite *ValidateTestSuite) TestValidateProfile_BuiltInNotWarning() {
	// Built-in profile should not warn about overriding itself
	profile := Profile{
		Name:        "cli",
		Description: "CLI profile",
		Disabled:    []string{"common:health"},
		Source:      SourceBuiltIn,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	// Should not have override warning since Source is BuiltIn
	for _, w := range result.Warnings {
		suite.NotContains(w, "overrides built-in")
	}
}

func (suite *ValidateTestSuite) TestValidateProfile_MultipleErrors() {
	profile := Profile{
		Name:        "multi",
		Description: "Profile with multiple issues",
		Disabled: []string{
			"unknown:one",
			"common:health",
			"unknown:two",
			"common:health", // duplicate
		},
		Source: SourceUser,
	}

	result := ValidateProfile(profile, suite.validCheckIDs, suite.validIDList)

	suite.False(result.Valid)
	suite.Len(result.Errors, 2) // two unknown checks
}

func (suite *ValidateTestSuite) TestGetValidCheckIDs_ReturnsNonEmpty() {
	ids, list := GetValidCheckIDs()

	suite.NotEmpty(ids)
	suite.NotEmpty(list)
	// Should contain at least some known check IDs
	suite.True(ids["go:build"])
	suite.True(ids["common:health"])
}

func TestValidateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateTestSuite))
}
