package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// TypeTestSuite is the test suite for the TypeScript type check.
type TypeTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *TypeTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-node-type-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *TypeTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *TypeTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	dir := filepath.Dir(filePath)
	if dir != suite.tempDir {
		err := os.MkdirAll(dir, 0755)
		suite.NoError(err)
	}
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestTypeCheck_ID tests that TypeCheck returns correct ID.
func (suite *TypeTestSuite) TestTypeCheck_ID() {
	check := &TypeCheck{}
	suite.Equal("node:type", check.ID())
}

// TestTypeCheck_Name tests that TypeCheck returns correct name.
func (suite *TypeTestSuite) TestTypeCheck_Name() {
	check := &TypeCheck{}
	suite.Equal("TypeScript Type Check", check.Name())
}

// TestTypeCheck_Run_NoPackageJSON tests that TypeCheck returns Fail when package.json doesn't exist.
func (suite *TypeTestSuite) TestTypeCheck_Run_NoPackageJSON() {
	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "package.json not found")
	suite.Equal("node:type", result.ID)
	suite.Equal("TypeScript Type Check", result.Name)
	suite.Equal(checker.LangNode, result.Language)
}

// TestTypeCheck_Run_NotTypeScriptProject tests that TypeCheck passes when not a TypeScript project.
func (suite *TypeTestSuite) TestTypeCheck_Run_NotTypeScriptProject() {
	// Create package.json without TypeScript
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "Not a TypeScript project")
}

// TestTypeCheck_Run_TypeScriptInDevDependencies tests TypeScript detection via devDependencies.
func (suite *TypeTestSuite) TestTypeCheck_Run_TypeScriptInDevDependencies() {
	// Create package.json with TypeScript in devDependencies
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0",
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}`
	suite.createTempFile("package.json", packageJSON)

	check := &TypeCheck{}

	// This will detect TypeScript but may fail if tsc is not available
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// The result depends on whether tsc is available
	// If npx is not available, it should pass with "npx not available" message
	// If TypeScript project but tsc fails, it should warn
	suite.NotEmpty(result.Message)
}

// TestTypeCheck_Run_TsconfigExists tests TypeScript detection via tsconfig.json.
func (suite *TypeTestSuite) TestTypeCheck_Run_TsconfigExists() {
	// Create package.json and tsconfig.json
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	tsConfig := `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "strict": true
  }
}`
	suite.createTempFile("package.json", packageJSON)
	suite.createTempFile("tsconfig.json", tsConfig)

	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// The result depends on whether npx/tsc is available
	suite.NotEmpty(result.Message)
	suite.Equal(checker.LangNode, result.Language)
}

// TestTypeCheck_isTypeScriptProject tests the TypeScript project detection logic.
func (suite *TypeTestSuite) TestTypeCheck_isTypeScriptProject() {
	check := &TypeCheck{}

	// No tsconfig.json, no TypeScript dependency
	packageJSON := `{
  "name": "test-package",
  "version": "1.0.0"
}`
	suite.createTempFile("package.json", packageJSON)
	suite.False(check.isTypeScriptProject(suite.tempDir))

	// Add tsconfig.json
	suite.createTempFile("tsconfig.json", "{}")
	suite.True(check.isTypeScriptProject(suite.tempDir))
}

// TestTypeCheck_countTypeErrors tests the error counting logic.
func (suite *TypeTestSuite) TestTypeCheck_countTypeErrors() {
	check := &TypeCheck{}

	// Test "Found X errors." format
	output1 := `src/index.ts(10,5): error TS2322: Type 'string' is not assignable to type 'number'.
src/utils.ts(20,3): error TS2304: Cannot find name 'foo'.
Found 2 errors.`
	suite.Equal(2, check.countTypeErrors(output1))

	// Test single error format
	output2 := `src/index.ts(10,5): error TS2322: Type 'string' is not assignable to type 'number'.
Found 1 error.`
	suite.Equal(1, check.countTypeErrors(output2))

	// Test no errors
	output3 := ``
	suite.Equal(0, check.countTypeErrors(output3))

	// Test fallback counting
	output4 := `src/index.ts(10,5): error TS2322: message
src/utils.ts(20,3): error TS2304: message`
	suite.Equal(2, check.countTypeErrors(output4))
}

// TestTypeTestSuite runs all the tests in the suite.
func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}
