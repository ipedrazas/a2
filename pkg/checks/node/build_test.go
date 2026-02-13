package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

// BuildTestSuite is the test suite for the build check.
type BuildTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *BuildTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-build-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *BuildTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *BuildTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestBuildCheck_Run_NoPackageJSON tests that BuildCheck returns Fail when package.json doesn't exist.
func (suite *BuildTestSuite) TestBuildCheck_Run_NoPackageJSON() {
	check := &BuildCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
	suite.Equal("node:build", result.ID)
	suite.Equal("Node Build", result.Name)
}

// TestBuildCheck_ID tests that BuildCheck returns correct ID.
func (suite *BuildTestSuite) TestBuildCheck_ID() {
	check := &BuildCheck{}
	suite.Equal("node:build", check.ID())
}

// TestBuildCheck_Name tests that BuildCheck returns correct name.
func (suite *BuildTestSuite) TestBuildCheck_Name() {
	check := &BuildCheck{}
	suite.Equal("Node Build", check.Name())
}

// TestBuildCheck_DetectPackageManager_NPM tests npm detection (default).
func (suite *BuildTestSuite) TestBuildCheck_DetectPackageManager_NPM() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &BuildCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("npm", pm)
}

// TestBuildCheck_DetectPackageManager_Yarn tests yarn detection from yarn.lock.
func (suite *BuildTestSuite) TestBuildCheck_DetectPackageManager_Yarn() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("yarn.lock", "# yarn lockfile v1")

	check := &BuildCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("yarn", pm)
}

// TestBuildCheck_DetectPackageManager_PNPM tests pnpm detection from pnpm-lock.yaml.
func (suite *BuildTestSuite) TestBuildCheck_DetectPackageManager_PNPM() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("pnpm-lock.yaml", "lockfileVersion: 5.4")

	check := &BuildCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("pnpm", pm)
}

// TestBuildCheck_DetectPackageManager_Bun tests bun detection from bun.lockb.
func (suite *BuildTestSuite) TestBuildCheck_DetectPackageManager_Bun() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("bun.lockb", "binary content")

	check := &BuildCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("bun", pm)
}

// TestBuildCheck_DetectPackageManager_ConfigOverride tests config override.
func (suite *BuildTestSuite) TestBuildCheck_DetectPackageManager_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("yarn.lock", "# yarn lockfile v1") // Would normally detect yarn

	check := &BuildCheck{
		Config: &config.NodeLanguageConfig{
			PackageManager: "pnpm",
		},
	}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("pnpm", pm) // Config override takes precedence
}

// TestBuildCheck_Run_NonExistentPath tests that BuildCheck handles non-existent paths.
func (suite *BuildTestSuite) TestBuildCheck_Run_NonExistentPath() {
	check := &BuildCheck{}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
}

// TestTruncateMessage tests the TruncateMessage helper function.
func (suite *BuildTestSuite) TestTruncateMessage() {
	// Short message should not be truncated
	short := "short message"
	suite.Equal(short, checkutil.TruncateMessage(short, 100))

	// Long message should be truncated
	long := "this is a very long message that should be truncated at some point"
	truncated := checkutil.TruncateMessage(long, 20)
	suite.Equal("this is a very long ...", truncated)
	suite.Len(truncated, 23) // 20 + "..."

	// Message with leading/trailing whitespace
	whitespace := "  trimmed  "
	suite.Equal("trimmed", checkutil.TruncateMessage(whitespace, 100))
}

// TestBuildTestSuite runs all the tests in the suite.
func TestBuildTestSuite(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
