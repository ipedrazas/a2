package gocheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// TestsTestSuite is the test suite for the tests check package.
type TestsTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *TestsTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-tests-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *TestsTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *TestsTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	// Create directory if needed
	dir := filepath.Dir(filePath)
	if dir != suite.tempDir {
		err := os.MkdirAll(dir, 0755)
		suite.NoError(err)
	}
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// createGoMod creates a basic go.mod file for testing.
func (suite *TestsTestSuite) createGoMod() {
	goModContent := `module test

go 1.21
`
	suite.createTempFile("go.mod", goModContent)
}

// TestTestRunnerCheck_ID tests that TestRunnerCheck returns correct ID.
func (suite *TestsTestSuite) TestTestsCheck_ID() {
	check := &TestsCheck{}
	suite.Equal("go:tests", check.ID())
}

// TestTestRunnerCheck_Name tests that TestRunnerCheck returns correct name.
func (suite *TestsTestSuite) TestTestsCheck_Name() {
	check := &TestsCheck{}
	suite.Equal("Go Tests", check.Name())
}

// TestTestRunnerCheck_Run_AllTestsPass tests that TestRunnerCheck returns Pass when all tests pass.
func (suite *TestsTestSuite) TestTestsCheck_Run_AllTestsPass() {
	suite.createGoMod()

	// Create a Go file
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	// Create test file
	testCode := `package main

import "testing"

func TestPass(t *testing.T) {
	if 1+1 != 2 {
		t.Error("Math is broken")
	}
}
`
	suite.createTempFile("main.go", code)
	suite.createTempFile("main_test.go", testCode)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "All tests passed")
	suite.Equal("go:tests", result.ID)
	suite.Equal("Go Tests", result.Name)
}

// TestTestRunnerCheck_Run_TestsFail tests that TestRunnerCheck returns Fail when tests fail.
func (suite *TestsTestSuite) TestTestsCheck_Run_TestsFail() {
	suite.createGoMod()

	// Create a Go file with failing tests
	code := `package main

import "testing"

func TestFail(t *testing.T) {
	t.Error("This test fails")
}
`
	suite.createTempFile("main.go", `package main`)
	suite.createTempFile("main_test.go", code)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status) // Critical - should be Fail
	suite.Contains(result.Message, "Tests failed")
	suite.NotEmpty(result.Message)
}

// TestTestRunnerCheck_Run_NoTestFiles tests that TestRunnerCheck returns Pass when no test files exist.
func (suite *TestsTestSuite) TestTestsCheck_Run_NoTestFiles() {
	suite.createGoMod()

	// Create a Go file without tests
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	suite.createTempFile("main.go", code)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// No test files should pass (not a failure)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No test files found")
}

// TestTestCheck_Run_EmptyDirectory tests that TestsCheck handles empty directories.
func (suite *TestsTestSuite) TestTestsCheck_Run_EmptyDirectory() {
	suite.createGoMod()

	// Create at least one Go file (empty directory with just go.mod causes different error)
	code := `package main
`
	suite.createTempFile("main.go", code)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Directory with code but no test files should pass
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No test files found")
}

// TestTestRunnerCheck_Run_MultipleTestFiles tests that TestsCheck handles multiple test files.
func (suite *TestsTestSuite) TestTestsCheck_Run_MultipleTestFiles() {
	suite.createGoMod()

	// Create multiple test files
	code1 := `package main

import "testing"

func TestOne(t *testing.T) {
	if 1 != 1 {
		t.Error("Fail")
	}
}
`
	code2 := `package main

import "testing"

func TestTwo(t *testing.T) {
	if 2 != 2 {
		t.Error("Fail")
	}
}
`
	suite.createTempFile("main.go", `package main`)
	suite.createTempFile("main_test.go", code1)
	suite.createTempFile("helper_test.go", code2)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "All tests passed")
}

// TestTestRunnerCheck_Run_NestedPackages tests that TestsCheck checks nested packages.
func (suite *TestsTestSuite) TestTestsCheck_Run_NestedPackages() {
	suite.createGoMod()

	// Create nested package with tests
	nestedCode := `package sub

import "testing"

func TestSub(t *testing.T) {
	if 1 != 1 {
		t.Error("Fail")
	}
}
`
	suite.createTempFile("main.go", `package main`)
	suite.createTempFile("sub/helper.go", `package sub`)
	suite.createTempFile("sub/helper_test.go", nestedCode)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestTestRunnerCheck_Run_StderrOutput tests that TestsCheck captures stderr output.
func (suite *TestsTestSuite) TestTestsCheck_Run_StderrOutput() {
	suite.createGoMod()

	// Create a test that outputs to stderr
	testCode := `package main

import (
	"os"
	"testing"
)

func TestStderr(t *testing.T) {
	os.Stderr.WriteString("Test output")
}
`
	suite.createTempFile("main.go", `package main`)
	suite.createTempFile("main_test.go", testCode)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestTestRunnerCheck_Run_StdoutFallback tests that TestRunnerCheck falls back to stdout.
func (suite *TestsTestSuite) TestTestsCheck_Run_StdoutFallback() {
	suite.createGoMod()

	// Create a failing test
	testCode := `package main

import "testing"

func TestFail(t *testing.T) {
	t.Error("Test failed")
}
`
	suite.createTempFile("main.go", `package main`)
	suite.createTempFile("main_test.go", testCode)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	// Should have error message
	suite.NotEmpty(result.Message)
}

// TestTestRunnerCheck_Run_PartialFailure tests that TestRunnerCheck handles partial test failures.
func (suite *TestsTestSuite) TestTestsCheck_Run_PartialFailure() {
	suite.createGoMod()

	// Create one passing and one failing test
	testCode := `package main

import "testing"

func TestPass(t *testing.T) {
	// Pass
}

func TestFail(t *testing.T) {
	t.Error("This fails")
}
`
	suite.createTempFile("main.go", `package main`)
	suite.createTempFile("main_test.go", testCode)

	check := &TestsCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "Tests failed")
}

// TestTestsTestSuite runs all the tests in the suite.
func TestTestsTestSuite(t *testing.T) {
	suite.Run(t, new(TestsTestSuite))
}
