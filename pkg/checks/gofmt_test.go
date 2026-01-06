package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// GofmtTestSuite is the test suite for the gofmt check package.
type GofmtTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *GofmtTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-gofmt-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *GofmtTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *GofmtTestSuite) createTempFile(name, content string) string {
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

// TestGofmtCheck_ID tests that GofmtCheck returns correct ID.
func (suite *GofmtTestSuite) TestGofmtCheck_ID() {
	check := &GofmtCheck{}
	suite.Equal("gofmt", check.ID())
}

// TestGofmtCheck_Name tests that GofmtCheck returns correct name.
func (suite *GofmtTestSuite) TestGofmtCheck_Name() {
	check := &GofmtCheck{}
	suite.Equal("Code Formatting", check.Name())
}

// TestGofmtCheck_Run_AllFormatted tests that GofmtCheck returns Pass when all files are formatted.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_AllFormatted() {
	// Create a properly formatted Go file
	formattedCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World")
}
`
	suite.createTempFile("main.go", formattedCode)

	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "All Go files are properly formatted")
}

// TestGofmtCheck_Run_UnformattedFiles tests that GofmtCheck returns Warn when files need formatting.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_UnformattedFiles() {
	// Create an unformatted Go file
	unformattedCode := `package main
import "fmt"
func main(){
fmt.Println("Hello, World")
}
`
	suite.createTempFile("main.go", unformattedCode)

	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "Unformatted files")
	suite.Contains(result.Message, "main.go")
	suite.Contains(result.Message, "Run 'gofmt -w .' to fix")
}

// TestGofmtCheck_Run_MultipleUnformattedFiles tests that GofmtCheck lists all unformatted files.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_MultipleUnformattedFiles() {
	// Create multiple unformatted files
	unformattedCode1 := `package main
import "fmt"
func main(){
fmt.Println("Hello")
}
`
	unformattedCode2 := `package utils
import "fmt"
func helper(){
fmt.Println("Helper")
}
`
	suite.createTempFile("main.go", unformattedCode1)
	suite.createTempFile("utils.go", unformattedCode2)

	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "Unformatted files")
	// Should mention both files
	suite.Contains(result.Message, "main.go")
	suite.Contains(result.Message, "utils.go")
}

// TestGofmtCheck_Run_NoGoFiles tests that GofmtCheck handles directories with no Go files.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_NoGoFiles() {
	// Create a non-Go file
	suite.createTempFile("README.md", "# Test")

	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// gofmt with no Go files should return empty output (pass)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestGofmtCheck_Run_EmptyDirectory tests that GofmtCheck handles empty directories.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_EmptyDirectory() {
	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	// Empty directory should pass
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestGofmtCheck_Run_MixedFiles tests that GofmtCheck handles mix of formatted and unformatted files.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_MixedFiles() {
	// Create one formatted file
	formattedCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`
	// Create one unformatted file
	unformattedCode := `package utils
import "fmt"
func helper(){
fmt.Println("Helper")
}
`
	suite.createTempFile("main.go", formattedCode)
	suite.createTempFile("utils.go", unformattedCode)

	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	// Should only mention the unformatted file
	suite.Contains(result.Message, "utils.go")
	// Should not mention the formatted file
	suite.NotContains(result.Message, "main.go")
}

// TestGofmtCheck_Run_NestedDirectories tests that GofmtCheck checks nested directories.
func (suite *GofmtTestSuite) TestGofmtCheck_Run_NestedDirectories() {
	// Create unformatted file in subdirectory
	unformattedCode := `package sub
import "fmt"
func test(){
fmt.Println("Test")
}
`
	suite.createTempFile("sub/helper.go", unformattedCode)

	check := &GofmtCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "sub/helper.go")
}

// TestGofmtTestSuite runs all the tests in the suite.
func TestGofmtTestSuite(t *testing.T) {
	suite.Run(t, new(GofmtTestSuite))
}
