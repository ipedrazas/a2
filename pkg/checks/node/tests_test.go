package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

// TestsCheckTestSuite is the test suite for the tests check.
type TestsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *TestsCheckTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-tests-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *TestsCheckTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *TestsCheckTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestTestsCheck_Run_NoPackageJSON tests that TestsCheck returns Fail when package.json doesn't exist.
func (suite *TestsCheckTestSuite) TestTestsCheck_Run_NoPackageJSON() {
	check := &TestsCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
	suite.Equal("node:tests", result.ID)
	suite.Equal("Node Tests", result.Name)
}

// TestTestsCheck_Run_NoTestScript tests that TestsCheck returns Pass when no test script is defined.
func (suite *TestsCheckTestSuite) TestTestsCheck_Run_NoTestScript() {
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TestsCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No test script defined")
}

// TestTestsCheck_Run_DefaultNoTestScript tests that TestsCheck handles default "no test specified" script.
func (suite *TestsCheckTestSuite) TestTestsCheck_Run_DefaultNoTestScript() {
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TestsCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No tests configured")
}

// TestTestsCheck_ID tests that TestsCheck returns correct ID.
func (suite *TestsCheckTestSuite) TestTestsCheck_ID() {
	check := &TestsCheck{}
	suite.Equal("node:tests", check.ID())
}

// TestTestsCheck_Name tests that TestsCheck returns correct name.
func (suite *TestsCheckTestSuite) TestTestsCheck_Name() {
	check := &TestsCheck{}
	suite.Equal("Node Tests", check.Name())
}

// TestTestsCheck_DetectTestRunner_Jest tests Jest detection from config file.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_Jest() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("jest.config.js", "module.exports = {};")

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("jest", runner)
}

// TestTestsCheck_DetectTestRunner_Vitest tests Vitest detection from config file.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_Vitest() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("vitest.config.ts", "export default {};")

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("vitest", runner)
}

// TestTestsCheck_DetectTestRunner_Mocha tests Mocha detection from config file.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_Mocha() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile(".mocharc.json", "{}")

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("mocha", runner)
}

// TestTestsCheck_DetectTestRunner_DevDependencies tests detection from devDependencies.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_DevDependencies() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "jest": "^29.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("jest", runner)
}

// TestTestsCheck_DetectTestRunner_ConfigOverride tests config override.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("jest.config.js", "module.exports = {};") // Would normally detect jest

	check := &TestsCheck{
		Config: &config.NodeLanguageConfig{
			TestRunner: "vitest",
		},
	}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("vitest", runner) // Config override takes precedence
}

// TestTestsCheck_DetectTestRunner_Default tests default runner.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_Default() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("npm-test", runner)
}

// TestTestsCheck_Run_NonExistentPath tests that TestsCheck handles non-existent paths.
func (suite *TestsCheckTestSuite) TestTestsCheck_Run_NonExistentPath() {
	check := &TestsCheck{}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
}

// TestTestsCheck_DetectPackageManager_NPM tests npm detection (default).
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectPackageManager_NPM() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &TestsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("npm", pm)
}

// TestTestsCheck_DetectPackageManager_Yarn tests yarn detection from yarn.lock.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectPackageManager_Yarn() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("yarn.lock", "# yarn lockfile v1")

	check := &TestsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("yarn", pm)
}

// TestTestsCheck_DetectPackageManager_PNPM tests pnpm detection from pnpm-lock.yaml.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectPackageManager_PNPM() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("pnpm-lock.yaml", "lockfileVersion: 5.4")

	check := &TestsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("pnpm", pm)
}

// TestTestsCheck_DetectPackageManager_Bun tests bun detection from bun.lockb.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectPackageManager_Bun() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("bun.lockb", "binary content")

	check := &TestsCheck{}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("bun", pm)
}

// TestTestsCheck_DetectPackageManager_ConfigOverride tests config override.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectPackageManager_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("yarn.lock", "# yarn lockfile v1") // Would normally detect yarn

	check := &TestsCheck{
		Config: &config.NodeLanguageConfig{
			PackageManager: "pnpm",
		},
	}
	pm := check.detectPackageManager(suite.tempDir)

	suite.Equal("pnpm", pm) // Config override takes precedence
}

// TestTestsCheck_DetectTestRunner_VitestDevDeps tests vitest detection from devDependencies.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_VitestDevDeps() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "vitest": "^1.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("vitest", runner)
}

// TestTestsCheck_DetectTestRunner_MochaDevDeps tests mocha detection from devDependencies.
func (suite *TestsCheckTestSuite) TestTestsCheck_DetectTestRunner_MochaDevDeps() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "mocha": "^10.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TestsCheck{}
	pkg, _ := ParsePackageJSON(suite.tempDir)
	runner := check.detectTestRunner(suite.tempDir, pkg)

	suite.Equal("mocha", runner)
}

// TestTestsCheckTestSuite runs all the tests in the suite.
func TestTestsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(TestsCheckTestSuite))
}
