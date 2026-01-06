package gocheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// RaceTestSuite is the test suite for the race check package.
type RaceTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *RaceTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-race-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *RaceTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *RaceTestSuite) createTempFile(name, content string) string {
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
func (suite *RaceTestSuite) createGoMod() {
	goModContent := `module test

go 1.21
`
	suite.createTempFile("go.mod", goModContent)
}

// TestRaceCheck_ID tests that RaceCheck returns correct ID.
func (suite *RaceTestSuite) TestRaceCheck_ID() {
	check := &RaceCheck{}
	suite.Equal("go:race", check.ID())
}

// TestRaceCheck_Name tests that RaceCheck returns correct name.
func (suite *RaceTestSuite) TestRaceCheck_Name() {
	check := &RaceCheck{}
	suite.Equal("Go Race Detection", check.Name())
}

// TestRaceCheck_Run_NoRaceConditions tests that RaceCheck returns Pass when no races found.
func (suite *RaceTestSuite) TestRaceCheck_Run_NoRaceConditions() {
	suite.createGoMod()

	// Create a simple Go file and test without race conditions
	code := `package main

func main() {}
`
	testCode := `package main

import "testing"

func TestSimple(t *testing.T) {
	x := 1
	if x != 1 {
		t.Error("Failed")
	}
}
`
	suite.createTempFile("main.go", code)
	suite.createTempFile("main_test.go", testCode)

	check := &RaceCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No race conditions detected")
	suite.Equal("go:race", result.ID)
	suite.Equal("Go Race Detection", result.Name)
	suite.Equal(checker.LangGo, result.Language)
}

// TestRaceCheck_Run_NoTestFiles tests that RaceCheck handles no test files gracefully.
func (suite *RaceTestSuite) TestRaceCheck_Run_NoTestFiles() {
	suite.createGoMod()

	// Create a Go file without tests
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	suite.createTempFile("main.go", code)

	check := &RaceCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No test files")
}

// TestRaceCheck_Run_TestFailure tests that RaceCheck handles test failures gracefully.
func (suite *RaceTestSuite) TestRaceCheck_Run_TestFailure() {
	suite.createGoMod()

	// Create a failing test (not a race condition, just a test failure)
	code := `package main

func main() {}
`
	testCode := `package main

import "testing"

func TestFailing(t *testing.T) {
	t.Error("This test fails intentionally")
}
`
	suite.createTempFile("main.go", code)
	suite.createTempFile("main_test.go", testCode)

	check := &RaceCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Test failure (not race) should still warn but not be a critical failure
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "Tests failed during race detection")
}

// TestRaceTestSuite runs all the tests in the suite.
func TestRaceTestSuite(t *testing.T) {
	suite.Run(t, new(RaceTestSuite))
}
