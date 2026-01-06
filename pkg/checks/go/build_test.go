package gocheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// BuildTestSuite is the test suite for the build check package.
type BuildTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *BuildTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-build-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *BuildTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *BuildTestSuite) createTempFile(name, content string) string {
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
func (suite *BuildTestSuite) createGoMod() {
	goModContent := `module test

go 1.21
`
	suite.createTempFile("go.mod", goModContent)
}

// TestBuildCheck_ID tests that BuildCheck returns correct ID.
func (suite *BuildTestSuite) TestBuildCheck_ID() {
	check := &BuildCheck{}
	suite.Equal("go:build", check.ID())
}

// TestBuildCheck_Name tests that BuildCheck returns correct name.
func (suite *BuildTestSuite) TestBuildCheck_Name() {
	check := &BuildCheck{}
	suite.Equal("Go Build", check.Name())
}

// TestBuildCheck_Run_SuccessfulBuild tests that BuildCheck returns Pass when build succeeds.
func (suite *BuildTestSuite) TestBuildCheck_Run_SuccessfulBuild() {
	suite.createGoMod()

	// Create a valid Go file that compiles
	validCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World")
}
`
	suite.createTempFile("main.go", validCode)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "Build successful")
	suite.Equal("go:build", result.ID)
	suite.Equal("Go Build", result.Name)
}

// TestBuildCheck_Run_FailedBuild tests that BuildCheck returns Fail when build fails.
func (suite *BuildTestSuite) TestBuildCheck_Run_FailedBuild() {
	suite.createGoMod()

	// Create an invalid Go file that won't compile
	invalidCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World"
	// Missing closing parenthesis
}
`
	suite.createTempFile("main.go", invalidCode)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status) // Critical - should be Fail
	suite.Contains(result.Message, "Build failed")
	suite.NotEmpty(result.Message)
}

// TestBuildCheck_Run_EmptyDirectory tests that BuildCheck handles empty directories.
func (suite *BuildTestSuite) TestBuildCheck_Run_EmptyDirectory() {
	suite.createGoMod()

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Empty directory with go.mod should build successfully (no code to build)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestBuildCheck_Run_MultiplePackages tests that BuildCheck checks multiple packages.
func (suite *BuildTestSuite) TestBuildCheck_Run_MultiplePackages() {
	suite.createGoMod()

	// Create main package
	mainCode := `package main

import "fmt"
import "test/sub"

func main() {
	fmt.Println("Hello")
	sub.Helper()
}
`
	// Create sub package
	subCode := `package sub

import "fmt"

func Helper() {
	fmt.Println("Helper")
}
`
	suite.createTempFile("main.go", mainCode)
	suite.createTempFile("sub/helper.go", subCode)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestBuildCheck_Run_NoGoMod tests that BuildCheck handles directories without go.mod.
func (suite *BuildTestSuite) TestBuildCheck_Run_NoGoMod() {
	// Create a Go file without go.mod
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	suite.createTempFile("main.go", code)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	// go build requires go.mod, so this should fail
	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "Build failed")
}

// TestBuildCheck_Run_StderrOutput tests that BuildCheck captures stderr output.
func (suite *BuildTestSuite) TestBuildCheck_Run_StderrOutput() {
	suite.createGoMod()

	// Create code with build error (should output to stderr)
	invalidCode := `package main

import "nonexistent/package"

func main() {
}
`
	suite.createTempFile("main.go", invalidCode)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	// Should capture error message
	suite.NotEmpty(result.Message)
}

// TestBuildCheck_Run_StdoutFallback tests that BuildCheck falls back to stdout when stderr is empty.
func (suite *BuildTestSuite) TestBuildCheck_Run_StdoutFallback() {
	suite.createGoMod()

	// Create code that might output to stdout
	invalidCode := `package main

func main() {
	var x int = "string" // Type error
}
`
	suite.createTempFile("main.go", invalidCode)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	// Should have error message from either stderr or stdout
	suite.NotEmpty(result.Message)
}

// TestBuildCheck_Run_ComplexProject tests that BuildCheck handles complex projects.
func (suite *BuildTestSuite) TestBuildCheck_Run_ComplexProject() {
	suite.createGoMod()

	// Create a more complex project structure
	mainCode := `package main

import (
	"fmt"
	"test/utils"
)

func main() {
	fmt.Println("Main")
	utils.DoSomething()
}
`
	utilsCode := `package utils

import "fmt"

func DoSomething() {
	fmt.Println("Utils")
}
`
	suite.createTempFile("main.go", mainCode)
	suite.createTempFile("utils/helper.go", utilsCode)

	check := &BuildCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestBuildTestSuite runs all the tests in the suite.
func TestBuildTestSuite(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
