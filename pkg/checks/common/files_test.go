package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// FilesTestSuite is the test suite for the files check package.
type FilesTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *FilesTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-files-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *FilesTestSuite) TearDownTest() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with the given content.
func (suite *FilesTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestFileExistsCheck_Run_AllFilesExist tests that FileExistsCheck returns Pass when all files exist.
func (suite *FilesTestSuite) TestFileExistsCheck_Run_AllFilesExist() {
	suite.createTempFile("README.md", "# Test README")
	suite.createTempFile("LICENSE", "MIT License")

	check := &FileExistsCheck{
		Files: []string{"README.md", "LICENSE"},
	}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Equal("All required files present", result.Message)
	suite.Equal("file_exists", result.ID)
	suite.Equal("Required Files", result.Name)
}

// TestFileExistsCheck_Run_MissingFiles tests that FileExistsCheck returns Warn when files are missing.
func (suite *FilesTestSuite) TestFileExistsCheck_Run_MissingFiles() {
	suite.createTempFile("README.md", "# Test README")
	// LICENSE is missing

	check := &FileExistsCheck{
		Files: []string{"README.md", "LICENSE"},
	}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "Missing files")
	suite.Contains(result.Message, "LICENSE")
}

// TestFileExistsCheck_Run_MultipleMissingFiles tests that FileExistsCheck handles multiple missing files.
func (suite *FilesTestSuite) TestFileExistsCheck_Run_MultipleMissingFiles() {
	check := &FileExistsCheck{
		Files: []string{"README.md", "LICENSE", "CONTRIBUTING.md"},
	}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	suite.Contains(result.Message, "Missing files")
	suite.Contains(result.Message, "README.md")
	suite.Contains(result.Message, "LICENSE")
	suite.Contains(result.Message, "CONTRIBUTING.md")
}

// TestFileExistsCheck_Run_EmptyFileList tests that FileExistsCheck handles empty file list.
func (suite *FilesTestSuite) TestFileExistsCheck_Run_EmptyFileList() {
	check := &FileExistsCheck{
		Files: []string{},
	}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Equal("All required files present", result.Message)
}

// TestFileExistsCheck_Run_NestedFiles tests that FileExistsCheck handles nested file paths.
func (suite *FilesTestSuite) TestFileExistsCheck_Run_NestedFiles() {
	// Create nested directory
	nestedDir := filepath.Join(suite.tempDir, "docs")
	err := os.MkdirAll(nestedDir, 0755)
	suite.NoError(err)

	suite.createTempFile("docs/README.md", "# Nested README")

	check := &FileExistsCheck{
		Files: []string{"docs/README.md"},
	}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
}

// TestFileExistsCheck_ID tests that FileExistsCheck returns correct ID.
func (suite *FilesTestSuite) TestFileExistsCheck_ID() {
	check := &FileExistsCheck{}
	suite.Equal("file_exists", check.ID())
}

// TestFileExistsCheck_Name tests that FileExistsCheck returns correct name.
func (suite *FilesTestSuite) TestFileExistsCheck_Name() {
	check := &FileExistsCheck{}
	suite.Equal("Required Files", check.Name())
}

// TestDefaultFileExistsCheck tests that DefaultFileExistsCheck returns check with default files.
func (suite *FilesTestSuite) TestDefaultFileExistsCheck() {
	check := DefaultFileExistsCheck()

	suite.NotNil(check)
	suite.Equal([]string{"README.md", "LICENSE"}, check.Files)
	suite.Equal("file_exists", check.ID())
	suite.Equal("Required Files", check.Name())
}

// TestFileExistsCheck_Run_PartialMatch tests that FileExistsCheck correctly identifies which files are missing.
func (suite *FilesTestSuite) TestFileExistsCheck_Run_PartialMatch() {
	suite.createTempFile("README.md", "# Test")
	// Missing LICENSE and CONTRIBUTING.md

	check := &FileExistsCheck{
		Files: []string{"README.md", "LICENSE", "CONTRIBUTING.md"},
	}

	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status)
	// Should mention missing files but not README.md
	suite.Contains(result.Message, "LICENSE")
	suite.Contains(result.Message, "CONTRIBUTING.md")
	suite.NotContains(result.Message, "README.md")
}

// TestFilesTestSuite runs all the tests in the suite.
func TestFilesTestSuite(t *testing.T) {
	suite.Run(t, new(FilesTestSuite))
}
