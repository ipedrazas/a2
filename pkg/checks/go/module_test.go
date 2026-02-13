package gocheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// ModuleTestSuite is the test suite for the module check package.
type ModuleTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *ModuleTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-module-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *ModuleTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *ModuleTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestModuleCheck_Run_NoGoMod tests that ModuleCheck returns Fail when go.mod doesn't exist.
func (suite *ModuleTestSuite) TestModuleCheck_Run_NoGoMod() {
	check := &ModuleCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "go.mod not found")
	suite.Equal("go:module", result.ID)
	suite.Equal("Go Module", result.Name)
}

// TestModuleCheck_Run_ValidGoMod tests that ModuleCheck returns Pass when go.mod is valid with version.
func (suite *ModuleTestSuite) TestModuleCheck_Run_ValidGoMod() {
	validGoMod := `module github.com/ipedrazas/a2

go 1.21

require (
	github.com/spf13/cobra v1.10.2
)
`
	suite.createTempFile("go.mod", validGoMod)

	check := &ModuleCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Reason, "github.com/ipedrazas/a2")
	suite.Contains(result.Reason, "Go 1.21")
}

// TestModuleCheck_Run_InvalidGoMod tests that ModuleCheck returns Fail when go.mod is invalid.
func (suite *ModuleTestSuite) TestModuleCheck_Run_InvalidGoMod() {
	invalidGoMod := `module github.com/ipedrazas/a2
invalid syntax here
`
	suite.createTempFile("go.mod", invalidGoMod)

	check := &ModuleCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "go.mod is invalid")
}

// TestModuleCheck_Run_NoGoVersion tests that ModuleCheck returns Warn when Go version is missing.
func (suite *ModuleTestSuite) TestModuleCheck_Run_NoGoVersion() {
	goModNoVersion := `module github.com/ipedrazas/a2

require (
	github.com/spf13/cobra v1.10.2
)
`
	suite.createTempFile("go.mod", goModNoVersion)

	check := &ModuleCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Reason, "does not specify a Go version")
}

// TestModuleCheck_Run_EmptyGoVersion tests that ModuleCheck handles empty Go version (treated as invalid).
func (suite *ModuleTestSuite) TestModuleCheck_Run_EmptyGoVersion() {
	goModEmptyVersion := `module github.com/ipedrazas/a2

go 

require (
	github.com/spf13/cobra v1.10.2
)
`
	suite.createTempFile("go.mod", goModEmptyVersion)

	check := &ModuleCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	// Empty go directive is treated as invalid syntax, not missing version
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "go.mod is invalid")
}

// TestModuleCheck_ID tests that ModuleCheck returns correct ID.
func (suite *ModuleTestSuite) TestModuleCheck_ID() {
	check := &ModuleCheck{}
	suite.Equal("go:module", check.ID())
}

// TestModuleCheck_Name tests that ModuleCheck returns correct name.
func (suite *ModuleTestSuite) TestModuleCheck_Name() {
	check := &ModuleCheck{}
	suite.Equal("Go Module", check.Name())
}

// TestModuleCheck_Run_ComplexGoMod tests that ModuleCheck handles complex go.mod files.
func (suite *ModuleTestSuite) TestModuleCheck_Run_ComplexGoMod() {
	complexGoMod := `module github.com/ipedrazas/a2

go 1.21

require (
	github.com/spf13/cobra v1.10.2
	golang.org/x/mod v0.31.0
)

require (
	golang.org/x/sys v0.30.0 // indirect
)
`
	suite.createTempFile("go.mod", complexGoMod)

	check := &ModuleCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Reason, "Go 1.21")
}

// TestModuleCheck_Run_FileReadError tests that ModuleCheck handles file read errors.
func (suite *ModuleTestSuite) TestModuleCheck_Run_FileReadError() {
	check := &ModuleCheck{}

	// Use a non-existent path - should return Fail result (go.mod not found)
	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	// The check treats missing go.mod as a Fail, not an error
	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "go.mod not found")
}

// TestModuleTestSuite runs all the tests in the suite.
func TestModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
