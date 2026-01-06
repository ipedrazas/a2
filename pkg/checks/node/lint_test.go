package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

// LintTestSuite is the test suite for the lint check.
type LintTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *LintTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-lint-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *LintTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *LintTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestLintCheck_Run_NoPackageJSON tests that LintCheck returns Fail when package.json doesn't exist.
func (suite *LintTestSuite) TestLintCheck_Run_NoPackageJSON() {
	check := &LintCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
	suite.Equal("node:lint", result.ID)
	suite.Equal("Node Lint", result.Name)
}

// TestLintCheck_Run_NoLinterConfigured tests that LintCheck returns Pass when no linter is configured.
func (suite *LintTestSuite) TestLintCheck_Run_NoLinterConfigured() {
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", packageJSON)

	check := &LintCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	// Should indicate no linter configured or npx not available
}

// TestLintCheck_ID tests that LintCheck returns correct ID.
func (suite *LintTestSuite) TestLintCheck_ID() {
	check := &LintCheck{}
	suite.Equal("node:lint", check.ID())
}

// TestLintCheck_Name tests that LintCheck returns correct name.
func (suite *LintTestSuite) TestLintCheck_Name() {
	check := &LintCheck{}
	suite.Equal("Node Lint", check.Name())
}

// TestLintCheck_DetectLinter_ESLint tests ESLint detection from config file.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_ESLint() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile(".eslintrc.json", "{}")

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("eslint", linter)
}

// TestLintCheck_DetectLinter_ESLintFlatConfig tests ESLint flat config detection.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_ESLintFlatConfig() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("eslint.config.js", "export default [];")

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("eslint", linter)
}

// TestLintCheck_DetectLinter_Biome tests Biome detection from config file.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_Biome() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("biome.json", "{}")

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("biome", linter)
}

// TestLintCheck_DetectLinter_Oxlint tests Oxlint detection from config file.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_Oxlint() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("oxlint.json", "{}")

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("oxlint", linter)
}

// TestLintCheck_DetectLinter_DevDependencies tests detection from devDependencies.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_DevDependencies() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "eslint": "^8.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("eslint", linter)
}

// TestLintCheck_DetectLinter_BiomeDevDependencies tests Biome detection from devDependencies.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_BiomeDevDependencies() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "@biomejs/biome": "^1.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("biome", linter)
}

// TestLintCheck_DetectLinter_ConfigOverride tests config override.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile(".eslintrc.json", "{}") // Would normally detect eslint

	check := &LintCheck{
		Config: &config.NodeLanguageConfig{
			Linter: "biome",
		},
	}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("biome", linter) // Config override takes precedence
}

// TestLintCheck_DetectLinter_Default tests default (auto) linter.
func (suite *LintTestSuite) TestLintCheck_DetectLinter_Default() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &LintCheck{}
	linter := check.detectLinter(suite.tempDir)

	suite.Equal("auto", linter)
}

// TestLintCheck_HasESLintConfig tests ESLint config detection.
func (suite *LintTestSuite) TestLintCheck_HasESLintConfig() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	check := &LintCheck{}

	// No config - should return false
	suite.False(check.hasESLintConfig(suite.tempDir))

	// Add .eslintrc.json - should return true
	suite.createTempFile(".eslintrc.json", "{}")
	suite.True(check.hasESLintConfig(suite.tempDir))
}

// TestLintCheck_HasBiomeConfig tests Biome config detection.
func (suite *LintTestSuite) TestLintCheck_HasBiomeConfig() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	check := &LintCheck{}

	// No config - should return false
	suite.False(check.hasBiomeConfig(suite.tempDir))

	// Add biome.json - should return true
	suite.createTempFile("biome.json", "{}")
	suite.True(check.hasBiomeConfig(suite.tempDir))
}

// TestLintCheck_HasOxlintConfig tests Oxlint config detection.
func (suite *LintTestSuite) TestLintCheck_HasOxlintConfig() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	check := &LintCheck{}

	// No config - should return false
	suite.False(check.hasOxlintConfig(suite.tempDir))

	// Add oxlint.json - should return true
	suite.createTempFile("oxlint.json", "{}")
	suite.True(check.hasOxlintConfig(suite.tempDir))
}

// TestCountLintIssues tests the countLintIssues helper function.
func (suite *LintTestSuite) TestCountLintIssues() {
	output := `
src/index.js:1:1: error: Missing semicolon
src/index.js:2:1: warning: Unexpected console statement
src/index.js:3:1: error: 'foo' is not defined
`
	count := countLintIssues(output)
	suite.Equal(3, count)

	// Empty output
	suite.Equal(0, countLintIssues(""))

	// Output with no issues
	suite.Equal(0, countLintIssues("All files passed linting"))
}

// TestLintCheck_Run_NonExistentPath tests that LintCheck handles non-existent paths.
func (suite *LintTestSuite) TestLintCheck_Run_NonExistentPath() {
	check := &LintCheck{}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
}

// TestLintTestSuite runs all the tests in the suite.
func TestLintTestSuite(t *testing.T) {
	suite.Run(t, new(LintTestSuite))
}
