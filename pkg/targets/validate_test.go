package targets

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

func (suite *ValidateTestSuite) TestValidateTarget_ValidTarget() {
	target := Target{
		Name:        "staging",
		Description: "Staging environment target",
		Disabled:    []string{"common:sast", "go:coverage"},
		Source:      SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	suite.Empty(result.Errors)
	suite.Equal("staging", result.Name)
}

func (suite *ValidateTestSuite) TestValidateTarget_EmptyDisabled() {
	target := Target{
		Name:        "strict",
		Description: "Target with no disabled checks",
		Disabled:    []string{},
		Source:      SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	suite.Empty(result.Errors)
	suite.Empty(result.Warnings)
}

func (suite *ValidateTestSuite) TestValidateTarget_UnknownCheckID() {
	target := Target{
		Name:        "bad",
		Description: "Target with unknown check",
		Disabled:    []string{"go:coverage", "nonexistent:check"},
		Source:      SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.False(result.Valid)
	suite.Len(result.Errors, 1)
	suite.Contains(result.Errors[0], "unknown check ID: nonexistent:check")
}

func (suite *ValidateTestSuite) TestValidateTarget_TypoWithSuggestion() {
	target := Target{
		Name:        "typo",
		Description: "Target with typo in check ID",
		Disabled:    []string{"go:coverge"}, // typo: should be go:coverage
		Source:      SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.False(result.Valid)
	suite.Len(result.Errors, 1)
	suite.Contains(result.Errors[0], "unknown check ID: go:coverge")
	// Should have a suggestion warning
	suite.NotEmpty(result.Warnings)
	suite.Contains(result.Warnings[0], "did you mean")
}

func (suite *ValidateTestSuite) TestValidateTarget_DuplicateDisabled() {
	target := Target{
		Name:        "dupe",
		Description: "Target with duplicate disabled checks",
		Disabled:    []string{"go:coverage", "common:sast", "go:coverage"}, // duplicate
		Source:      SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid) // Duplicates are warnings, not errors
	suite.Empty(result.Errors)
	suite.NotEmpty(result.Warnings)
	suite.Contains(result.Warnings[0], "duplicate disabled check: go:coverage")
}

func (suite *ValidateTestSuite) TestValidateTarget_OverridesBuiltIn() {
	// Test that overriding a built-in target produces a warning
	target := Target{
		Name:        "poc", // This is a built-in target name
		Description: "Custom POC target",
		Disabled:    []string{"common:license"},
		Source:      SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	suite.Empty(result.Errors)
	suite.NotEmpty(result.Warnings)
	suite.Contains(result.Warnings[0], "overrides built-in target: poc")
}

func (suite *ValidateTestSuite) TestValidateTarget_BuiltInNotWarning() {
	// Built-in target should not warn about overriding itself
	target := Target{
		Name:        "poc",
		Description: "POC target",
		Disabled:    []string{"common:license"},
		Source:      SourceBuiltIn,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

	suite.True(result.Valid)
	// Should not have override warning since Source is BuiltIn
	for _, w := range result.Warnings {
		suite.NotContains(w, "overrides built-in")
	}
}

func (suite *ValidateTestSuite) TestValidateTarget_MultipleErrors() {
	target := Target{
		Name:        "multi",
		Description: "Target with multiple issues",
		Disabled: []string{
			"invalid:first",
			"go:coverage",
			"invalid:second",
			"go:coverage", // duplicate
		},
		Source: SourceUser,
	}

	result := ValidateTarget(target, suite.validCheckIDs, suite.validIDList)

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
