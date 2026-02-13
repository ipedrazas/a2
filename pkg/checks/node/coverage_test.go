package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

// CoverageTestSuite is the test suite for the coverage check.
type CoverageTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *CoverageTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-coverage-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *CoverageTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *CoverageTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestCoverageCheck_Run_NoPackageJSON tests that CoverageCheck returns Fail when package.json doesn't exist.
func (suite *CoverageTestSuite) TestCoverageCheck_Run_NoPackageJSON() {
	check := &CoverageCheck{Threshold: 80.0}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
	suite.Equal("node:coverage", result.ID)
	suite.Equal("Node Coverage", result.Name)
}

// TestCoverageCheck_Run_NoTestScript tests that CoverageCheck returns Warn when no test script exists.
func (suite *CoverageTestSuite) TestCoverageCheck_Run_NoTestScript() {
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", packageJSON)

	check := &CoverageCheck{Threshold: 80.0}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Reason, "No test script configured")
	suite.Contains(result.Reason, "0%")
}

// TestCoverageCheck_Run_DefaultNoTestScript tests that CoverageCheck handles default "no test specified" script.
func (suite *CoverageTestSuite) TestCoverageCheck_Run_DefaultNoTestScript() {
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &CoverageCheck{Threshold: 80.0}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Reason, "No test script configured")
}

// TestCoverageCheck_ID tests that CoverageCheck returns correct ID.
func (suite *CoverageTestSuite) TestCoverageCheck_ID() {
	check := &CoverageCheck{}
	suite.Equal("node:coverage", check.ID())
}

// TestCoverageCheck_Name tests that CoverageCheck returns correct name.
func (suite *CoverageTestSuite) TestCoverageCheck_Name() {
	check := &CoverageCheck{}
	suite.Equal("Node Coverage", check.Name())
}

// TestCoverageCheck_DetectTestRunner_Jest tests Jest detection from config file.
func (suite *CoverageTestSuite) TestCoverageCheck_DetectTestRunner_Jest() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("jest.config.js", "module.exports = {};")

	check := &CoverageCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("jest", runner)
}

// TestCoverageCheck_DetectTestRunner_Vitest tests Vitest detection from config file.
func (suite *CoverageTestSuite) TestCoverageCheck_DetectTestRunner_Vitest() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("vitest.config.ts", "export default {};")

	check := &CoverageCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("vitest", runner)
}

// TestCoverageCheck_DetectTestRunner_C8 tests c8 detection from devDependencies.
func (suite *CoverageTestSuite) TestCoverageCheck_DetectTestRunner_C8() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "c8": "^7.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &CoverageCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("c8", runner)
}

// TestCoverageCheck_DetectTestRunner_NYC tests nyc detection from devDependencies.
func (suite *CoverageTestSuite) TestCoverageCheck_DetectTestRunner_NYC() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "nyc": "^15.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &CoverageCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("nyc", runner)
}

// TestCoverageCheck_DetectTestRunner_ConfigOverride tests config override.
func (suite *CoverageTestSuite) TestCoverageCheck_DetectTestRunner_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("jest.config.js", "module.exports = {};") // Would normally detect jest

	check := &CoverageCheck{
		Config: &config.NodeLanguageConfig{
			TestRunner: "vitest",
		},
	}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("vitest", runner) // Config override takes precedence
}

// TestCoverageCheck_DetectTestRunner_Default tests default runner.
func (suite *CoverageTestSuite) TestCoverageCheck_DetectTestRunner_Default() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &CoverageCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("jest", runner) // Default
}

// TestParseNodeCoverage tests parsing coverage from various outputs.
func (suite *CoverageTestSuite) TestParseNodeCoverage() {
	// Jest text-summary format
	jestOutput := `
---------|---------|----------|---------|---------|-------------------
File      | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s
---------|---------|----------|---------|---------|-------------------
All files |   85.71 |    66.67 |     100 |   85.71 |
---------|---------|----------|---------|---------|-------------------
`
	coverage := parseNodeCoverage(jestOutput)
	suite.InDelta(85.71, coverage, 0.01)

	// c8/nyc format
	c8Output := `
----------------------|---------|----------|---------|---------|-------------------
File                  | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s
----------------------|---------|----------|---------|---------|-------------------
All files             |   75.00 |    50.00 |   66.67 |   75.00 |
Statements   : 75.00% ( 15/20 )
Branches     : 50.00% ( 5/10 )
Functions    : 66.67% ( 4/6 )
Lines        : 75.00% ( 15/20 )
`
	coverage = parseNodeCoverage(c8Output)
	suite.InDelta(75.0, coverage, 0.01)

	// Empty output
	coverage = parseNodeCoverage("")
	suite.Equal(-1.0, coverage)

	// Output without coverage
	coverage = parseNodeCoverage("All tests passed!")
	suite.Equal(-1.0, coverage)
}

// TestCoverageCheck_Run_NonExistentPath tests that CoverageCheck handles non-existent paths.
func (suite *CoverageTestSuite) TestCoverageCheck_Run_NonExistentPath() {
	check := &CoverageCheck{Threshold: 80.0}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
}

// TestCoverageTestSuite runs all the tests in the suite.
func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
