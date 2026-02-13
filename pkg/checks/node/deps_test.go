package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

// DepsTestSuite is the test suite for the deps check.
type DepsTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *DepsTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-deps-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *DepsTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *DepsTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestDepsCheck_Run_NoPackageJSON tests that DepsCheck returns Fail when package.json doesn't exist.
func (suite *DepsTestSuite) TestDepsCheck_Run_NoPackageJSON() {
	check := &DepsCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
	suite.Equal("node:deps", result.ID)
	suite.Equal("Node Vulnerabilities", result.Name)
}

// TestDepsCheck_Run_BunSkipped tests that DepsCheck skips for Bun projects.
func (suite *DepsTestSuite) TestDepsCheck_Run_BunSkipped() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("bun.lockb", "binary content")

	check := &DepsCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Reason, "Bun does not have built-in security audit")
}

// TestDepsCheck_ID tests that DepsCheck returns correct ID.
func (suite *DepsTestSuite) TestDepsCheck_ID() {
	check := &DepsCheck{}
	suite.Equal("node:deps", check.ID())
}

// TestDepsCheck_Name tests that DepsCheck returns correct name.
func (suite *DepsTestSuite) TestDepsCheck_Name() {
	check := &DepsCheck{}
	suite.Equal("Node Vulnerabilities", check.Name())
}

// TestDepsCheck_DetectPackageManager_NPM tests npm detection (default).
func (suite *DepsTestSuite) TestDepsCheck_DetectPackageManager_NPM() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &DepsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("npm", pm)
}

// TestDepsCheck_DetectPackageManager_Yarn tests yarn detection from yarn.lock.
func (suite *DepsTestSuite) TestDepsCheck_DetectPackageManager_Yarn() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("yarn.lock", "# yarn lockfile v1")

	check := &DepsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("yarn", pm)
}

// TestDepsCheck_DetectPackageManager_PNPM tests pnpm detection from pnpm-lock.yaml.
func (suite *DepsTestSuite) TestDepsCheck_DetectPackageManager_PNPM() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("pnpm-lock.yaml", "lockfileVersion: 5.4")

	check := &DepsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("pnpm", pm)
}

// TestDepsCheck_DetectPackageManager_ConfigOverride tests config override.
func (suite *DepsTestSuite) TestDepsCheck_DetectPackageManager_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("yarn.lock", "# yarn lockfile v1") // Would normally detect yarn

	check := &DepsCheck{
		Config: &config.NodeLanguageConfig{
			PackageManager: "pnpm",
		},
	}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("pnpm", pm) // Config override takes precedence
}

// TestParseAuditOutput_NPM6Format tests parsing npm 6 audit format.
func (suite *DepsTestSuite) TestParseAuditOutput_NPM6Format() {
	output := `{
  "metadata": {
    "vulnerabilities": {
      "total": 5,
      "low": 1,
      "moderate": 2,
      "high": 1,
      "critical": 1
    }
  }
}`
	count, err := parseAuditOutput(output, "npm")
	suite.NoError(err)
	suite.Equal(5, count)
}

// TestParseAuditOutput_NPM7Format tests parsing npm 7+ audit format.
func (suite *DepsTestSuite) TestParseAuditOutput_NPM7Format() {
	output := `{
  "vulnerabilities": {
    "lodash": {},
    "axios": {},
    "moment": {}
  }
}`
	count, err := parseAuditOutput(output, "npm")
	suite.NoError(err)
	suite.Equal(3, count)
}

// TestParseAuditOutput_EmptyOutput tests parsing empty output.
func (suite *DepsTestSuite) TestParseAuditOutput_EmptyOutput() {
	count, err := parseAuditOutput("", "npm")
	suite.NoError(err)
	suite.Equal(0, count)
}

// TestParseAuditOutput_NoVulnerabilities tests parsing output with no vulnerabilities.
func (suite *DepsTestSuite) TestParseAuditOutput_NoVulnerabilities() {
	output := `{
  "metadata": {
    "vulnerabilities": {
      "total": 0,
      "low": 0,
      "moderate": 0,
      "high": 0,
      "critical": 0
    }
  }
}`
	count, err := parseAuditOutput(output, "npm")
	suite.NoError(err)
	suite.Equal(0, count)
}

// TestParseYarnAudit tests parsing yarn audit NDJSON format.
func (suite *DepsTestSuite) TestParseYarnAudit() {
	output := `{"type":"auditAdvisory","data":{"advisory":{"id":1234}}}
{"type":"auditAdvisory","data":{"advisory":{"id":5678}}}
{"type":"auditSummary","data":{"vulnerabilities":{"total":2}}}`

	count, err := parseYarnAudit(output)
	suite.NoError(err)
	suite.Equal(2, count)
}

// TestParseYarnAudit_EmptyOutput tests parsing empty yarn audit output.
func (suite *DepsTestSuite) TestParseYarnAudit_EmptyOutput() {
	count, err := parseYarnAudit("")
	suite.NoError(err)
	suite.Equal(0, count)
}

// TestParseYarnAudit_DuplicateAdvisories tests that duplicate advisories are counted once.
func (suite *DepsTestSuite) TestParseYarnAudit_DuplicateAdvisories() {
	output := `{"type":"auditAdvisory","data":{"advisory":{"id":1234}}}
{"type":"auditAdvisory","data":{"advisory":{"id":1234}}}
{"type":"auditAdvisory","data":{"advisory":{"id":5678}}}`

	count, err := parseYarnAudit(output)
	suite.NoError(err)
	suite.Equal(2, count) // Only 2 unique advisories
}

// TestDepsCheck_Run_NonExistentPath tests that DepsCheck handles non-existent paths.
func (suite *DepsTestSuite) TestDepsCheck_Run_NonExistentPath() {
	check := &DepsCheck{}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
}

// TestDepsTestSuite runs all the tests in the suite.
func TestDepsTestSuite(t *testing.T) {
	suite.Run(t, new(DepsTestSuite))
}
