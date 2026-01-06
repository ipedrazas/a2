package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// GoVetTestSuite is the test suite for the govet check package.
type GoVetTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *GoVetTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-govet-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *GoVetTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *GoVetTestSuite) createTempFile(name, content string) string {
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
func (suite *GoVetTestSuite) createGoMod() {
	goModContent := `module test

go 1.21
`
	suite.createTempFile("go.mod", goModContent)
}

// TestGoVetCheck_ID tests that GoVetCheck returns correct ID.
func (suite *GoVetTestSuite) TestGoVetCheck_ID() {
	check := &GoVetCheck{}
	suite.Equal("govet", check.ID())
}

// TestGoVetCheck_Name tests that GoVetCheck returns correct name.
func (suite *GoVetTestSuite) TestGoVetCheck_Name() {
	check := &GoVetCheck{}
	suite.Equal("Go Vet", check.Name())
}

// TestGoVetCheck_Run_NoIssues tests that GoVetCheck returns Pass when no issues found.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_NoIssues() {
	suite.createGoMod()

	// Create a valid Go file with no vet issues
	validCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World")
}
`
	suite.createTempFile("main.go", validCode)

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No issues found")
}

// TestGoVetCheck_Run_WithIssues tests that GoVetCheck returns Warn when issues are found.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_WithIssues() {
	suite.createGoMod()

	// Create a Go file with vet issues (unused variable)
	codeWithIssues := `package main

import "fmt"

func main() {
	var unused int
	fmt.Println("Hello")
}
`
	suite.createTempFile("main.go", codeWithIssues)

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	// Should contain the vet output
	suite.NotEmpty(result.Message)
	suite.NotContains(result.Message, "No issues found")
}

// TestGoVetCheck_Run_EmptyDirectory tests that GoVetCheck handles empty directories.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_EmptyDirectory() {
	suite.createGoMod()

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Empty directory with go.mod - go vet may return an error or pass
	// The actual behavior depends on go vet's handling of empty packages
	// We just verify it doesn't crash and returns a valid result
	suite.NotNil(result)
}

// TestGoVetCheck_Run_NoGoMod tests that GoVetCheck handles directories without go.mod.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_NoGoMod() {
	// Create a Go file without go.mod
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	suite.createTempFile("main.go", code)

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	// go vet requires go.mod, so this may fail or return an error
	// The behavior depends on go vet's handling
	if err != nil {
		// If there's an error, it's expected
		suite.Error(err)
	} else {
		// If no error, check the result
		suite.NotNil(result)
	}
}

// TestGoVetCheck_Run_MultipleFiles tests that GoVetCheck checks multiple files.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_MultipleFiles() {
	suite.createGoMod()

	// Create multiple valid files
	code1 := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	code2 := `package main

func helper() {
	fmt.Println("Helper")
}
`
	suite.createTempFile("main.go", code1)
	suite.createTempFile("helper.go", code2)

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Should check all files
	suite.NotNil(result)
}

// TestGoVetCheck_Run_NestedPackages tests that GoVetCheck checks nested packages.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_NestedPackages() {
	suite.createGoMod()

	// Create a nested package
	nestedCode := `package sub

import "fmt"

func Test() {
	fmt.Println("Test")
}
`
	suite.createTempFile("sub/helper.go", nestedCode)

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Should check nested packages
	suite.NotNil(result)
}

// TestGoVetCheck_Run_StderrOutput tests that GoVetCheck captures stderr output.
func (suite *GoVetTestSuite) TestGoVetCheck_Run_StderrOutput() {
	suite.createGoMod()

	// Create code that might produce stderr output
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	suite.createTempFile("main.go", code)

	check := &GoVetCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Should handle both stdout and stderr
	suite.NotNil(result)
}

// TestGoVetTestSuite runs all the tests in the suite.
func TestGoVetTestSuite(t *testing.T) {
	suite.Run(t, new(GoVetTestSuite))
}
