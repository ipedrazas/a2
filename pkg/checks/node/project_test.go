package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// ProjectTestSuite is the test suite for the project check.
type ProjectTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *ProjectTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-project-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *ProjectTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *ProjectTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestProjectCheck_Run_NoPackageJSON tests that ProjectCheck returns Fail when package.json doesn't exist.
func (suite *ProjectTestSuite) TestProjectCheck_Run_NoPackageJSON() {
	check := &ProjectCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
	suite.Equal("node:project", result.ID)
	suite.Equal("Node Project", result.Name)
}

// TestProjectCheck_Run_ValidPackageJSON tests that ProjectCheck returns Pass when package.json is valid.
func (suite *ProjectTestSuite) TestProjectCheck_Run_ValidPackageJSON() {
	validPackageJSON := `{
  "name": "test-package",
  "version": "1.0.0",
  "description": "A test package"
}`
	suite.createTempFile("package.json", validPackageJSON)

	check := &ProjectCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "test-package")
	suite.Contains(result.Message, "1.0.0")
}

// TestProjectCheck_Run_InvalidJSON tests that ProjectCheck returns Fail when package.json is invalid JSON.
func (suite *ProjectTestSuite) TestProjectCheck_Run_InvalidJSON() {
	invalidJSON := `{ "name": "test", invalid json }`
	suite.createTempFile("package.json", invalidJSON)

	check := &ProjectCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "invalid JSON")
}

// TestProjectCheck_Run_MissingName tests that ProjectCheck returns Fail when name field is missing.
func (suite *ProjectTestSuite) TestProjectCheck_Run_MissingName() {
	missingName := `{
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", missingName)

	check := &ProjectCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "missing required 'name' field")
}

// TestProjectCheck_Run_MissingVersion tests that ProjectCheck returns Warn when version field is missing.
func (suite *ProjectTestSuite) TestProjectCheck_Run_MissingVersion() {
	missingVersion := `{
  "name": "test-package"
}`
	suite.createTempFile("package.json", missingVersion)

	check := &ProjectCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "missing 'version' field")
}

// TestProjectCheck_ID tests that ProjectCheck returns correct ID.
func (suite *ProjectTestSuite) TestProjectCheck_ID() {
	check := &ProjectCheck{}
	suite.Equal("node:project", check.ID())
}

// TestProjectCheck_Name tests that ProjectCheck returns correct name.
func (suite *ProjectTestSuite) TestProjectCheck_Name() {
	check := &ProjectCheck{}
	suite.Equal("Node Project", check.Name())
}

// TestProjectCheck_Run_ComplexPackageJSON tests that ProjectCheck handles complex package.json files.
func (suite *ProjectTestSuite) TestProjectCheck_Run_ComplexPackageJSON() {
	complexPackageJSON := `{
  "name": "@scope/test-package",
  "version": "2.0.0-beta.1",
  "description": "A complex test package",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "test": "jest"
  },
  "dependencies": {
    "lodash": "^4.17.21"
  },
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}`
	suite.createTempFile("package.json", complexPackageJSON)

	check := &ProjectCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "@scope/test-package")
	suite.Contains(result.Message, "2.0.0-beta.1")
}

// TestProjectCheck_Run_NonExistentPath tests that ProjectCheck handles non-existent paths.
func (suite *ProjectTestSuite) TestProjectCheck_Run_NonExistentPath() {
	check := &ProjectCheck{}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
}

// TestProjectTestSuite runs all the tests in the suite.
func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
