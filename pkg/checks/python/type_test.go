package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// TypeTestSuite is the test suite for the Python type check.
type TypeTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *TypeTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-python-type-test-*")
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
	suite.Equal("python:type", check.ID())
}

// TestTypeCheck_Name tests that TypeCheck returns correct name.
func (suite *TypeTestSuite) TestTypeCheck_Name() {
	check := &TypeCheck{}
	suite.Equal("Python Type Check", check.Name())
}

// TestTypeCheck_Run_NotTypedProject tests that TypeCheck passes when not a typed project.
func (suite *TypeTestSuite) TestTypeCheck_Run_NotTypedProject() {
	// Create a simple Python file without any type config
	suite.createTempFile("main.py", "print('hello')")

	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Reason, "Not a typed Python project")
	suite.Equal(checker.LangPython, result.Language)
}

// TestTypeCheck_Run_MypyIniExists tests detection via mypy.ini.
func (suite *TypeTestSuite) TestTypeCheck_Run_MypyIniExists() {
	suite.createTempFile("main.py", "print('hello')")
	suite.createTempFile("mypy.ini", "[mypy]\nstrict = True")

	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Result depends on whether mypy is installed
	suite.NotEmpty(result.Reason)
	suite.Equal(checker.LangPython, result.Language)
}

// TestTypeCheck_Run_PyTypedExists tests detection via py.typed marker.
func (suite *TypeTestSuite) TestTypeCheck_Run_PyTypedExists() {
	suite.createTempFile("main.py", "print('hello')")
	suite.createTempFile("py.typed", "")

	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Result depends on whether mypy is installed
	suite.NotEmpty(result.Reason)
	suite.Equal(checker.LangPython, result.Language)
}

// TestTypeCheck_Run_PyprojectTomlWithMypy tests detection via pyproject.toml.
func (suite *TypeTestSuite) TestTypeCheck_Run_PyprojectTomlWithMypy() {
	suite.createTempFile("main.py", "print('hello')")
	suite.createTempFile("pyproject.toml", `[tool.mypy]
strict = true
`)

	check := &TypeCheck{}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Result depends on whether mypy is installed
	suite.NotEmpty(result.Reason)
	suite.Equal(checker.LangPython, result.Language)
}

// TestTypeCheck_isTypedProject tests the typed project detection logic.
func (suite *TypeTestSuite) TestTypeCheck_isTypedProject() {
	check := &TypeCheck{}

	// No type config
	suite.createTempFile("main.py", "print('hello')")
	suite.False(check.isTypedProject(suite.tempDir))

	// Add mypy.ini
	suite.createTempFile("mypy.ini", "[mypy]")
	suite.True(check.isTypedProject(suite.tempDir))
}

// TestTypeCheck_isTypedProject_PyTyped tests detection via py.typed marker.
func (suite *TypeTestSuite) TestTypeCheck_isTypedProject_PyTyped() {
	check := &TypeCheck{}

	suite.createTempFile("main.py", "print('hello')")
	suite.False(check.isTypedProject(suite.tempDir))

	suite.createTempFile("py.typed", "")
	suite.True(check.isTypedProject(suite.tempDir))
}

// TestTypeCheck_countTypeErrors tests the error counting logic.
func (suite *TypeTestSuite) TestTypeCheck_countTypeErrors() {
	check := &TypeCheck{}

	// Test "Found X errors in Y files" format
	output1 := `main.py:10: error: Incompatible types in assignment
utils.py:20: error: Name 'foo' is not defined
Found 2 errors in 2 files (checked 5 source files)`
	suite.Equal(2, check.countTypeErrors(output1))

	// Test single error format
	output2 := `main.py:10: error: Incompatible types
Found 1 error in 1 file (checked 3 source files)`
	suite.Equal(1, check.countTypeErrors(output2))

	// Test no errors
	output3 := `Success: no issues found in 5 source files`
	suite.Equal(0, check.countTypeErrors(output3))

	// Test fallback counting
	output4 := `main.py:10: error: message
utils.py:20: error: another message`
	suite.Equal(2, check.countTypeErrors(output4))
}

// TestTypeTestSuite runs all the tests in the suite.
func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}
