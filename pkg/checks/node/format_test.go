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

// FormatTestSuite is the test suite for the format check.
type FormatTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *FormatTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-format-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *FormatTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *FormatTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestFormatCheck_Run_NoPackageJSON tests that FormatCheck returns Fail when package.json doesn't exist.
func (suite *FormatTestSuite) TestFormatCheck_Run_NoPackageJSON() {
	check := &FormatCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
	suite.Equal("node:format", result.ID)
	suite.Equal("Node Format", result.Name)
}

// TestFormatCheck_Run_NoFormatterConfigured tests that FormatCheck returns Pass when no formatter is configured.
func (suite *FormatTestSuite) TestFormatCheck_Run_NoFormatterConfigured() {
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", packageJSON)

	check := &FormatCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	// Should indicate no formatter configured or npx not available
}

// TestFormatCheck_ID tests that FormatCheck returns correct ID.
func (suite *FormatTestSuite) TestFormatCheck_ID() {
	check := &FormatCheck{}
	suite.Equal("node:format", check.ID())
}

// TestFormatCheck_Name tests that FormatCheck returns correct name.
func (suite *FormatTestSuite) TestFormatCheck_Name() {
	check := &FormatCheck{}
	suite.Equal("Node Format", check.Name())
}

// TestFormatCheck_DetectFormatter_Prettier tests Prettier detection from config file.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_Prettier() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile(".prettierrc", "{}")

	check := &FormatCheck{}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("prettier", formatter)
}

// TestFormatCheck_DetectFormatter_PrettierJS tests Prettier detection from JS config file.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_PrettierJS() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("prettier.config.js", "module.exports = {};")

	check := &FormatCheck{}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("prettier", formatter)
}

// TestFormatCheck_DetectFormatter_Biome tests Biome detection from config file.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_Biome() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile("biome.json", "{}")

	check := &FormatCheck{}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("biome", formatter)
}

// TestFormatCheck_DetectFormatter_DevDependencies tests detection from devDependencies.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_DevDependencies() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "prettier": "^3.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &FormatCheck{}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("prettier", formatter)
}

// TestFormatCheck_DetectFormatter_BiomeDevDependencies tests Biome detection from devDependencies.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_BiomeDevDependencies() {
	packageJSON := `{
  "name": "test",
  "version": "1.0.0",
  "devDependencies": {
    "@biomejs/biome": "^1.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &FormatCheck{}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("biome", formatter)
}

// TestFormatCheck_DetectFormatter_ConfigOverride tests config override.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_ConfigOverride() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	suite.createTempFile(".prettierrc", "{}") // Would normally detect prettier

	check := &FormatCheck{
		Config: &config.NodeLanguageConfig{
			Formatter: "biome",
		},
	}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("biome", formatter) // Config override takes precedence
}

// TestFormatCheck_DetectFormatter_Default tests default (auto) formatter.
func (suite *FormatTestSuite) TestFormatCheck_DetectFormatter_Default() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)

	check := &FormatCheck{}
	formatter := check.detectFormatter(suite.tempDir)

	suite.Equal("auto", formatter)
}

// TestFormatCheck_HasPrettierConfig tests prettier config detection.
func (suite *FormatTestSuite) TestFormatCheck_HasPrettierConfig() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	check := &FormatCheck{}

	// No config - should return false
	suite.False(check.hasPrettierConfig(suite.tempDir))

	// Add .prettierrc - should return true
	suite.createTempFile(".prettierrc", "{}")
	suite.True(check.hasPrettierConfig(suite.tempDir))
}

// TestFormatCheck_HasBiomeConfig tests biome config detection.
func (suite *FormatTestSuite) TestFormatCheck_HasBiomeConfig() {
	suite.createTempFile("package.json", `{"name": "test", "version": "1.0.0"}`)
	check := &FormatCheck{}

	// No config - should return false
	suite.False(check.hasBiomeConfig(suite.tempDir))

	// Add biome.json - should return true
	suite.createTempFile("biome.json", "{}")
	suite.True(check.hasBiomeConfig(suite.tempDir))
}

// TestPluralize tests the Pluralize helper function.
func (suite *FormatTestSuite) TestPluralize() {
	suite.Equal("file", checkutil.Pluralize(1, "file", "files"))
	suite.Equal("files", checkutil.Pluralize(0, "file", "files"))
	suite.Equal("files", checkutil.Pluralize(2, "file", "files"))
	suite.Equal("files", checkutil.Pluralize(10, "file", "files"))
}

// TestFormatCheck_Run_NonExistentPath tests that FormatCheck handles non-existent paths.
func (suite *FormatTestSuite) TestFormatCheck_Run_NonExistentPath() {
	check := &FormatCheck{}

	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Reason, "package.json not found")
}

// TestFormatTestSuite runs all the tests in the suite.
func TestFormatTestSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}
