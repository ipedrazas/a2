package safepath

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// SafePathTestSuite is the test suite for the safepath package.
type SafePathTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest creates a temporary directory for each test.
func (suite *SafePathTestSuite) SetupTest() {
	var err error
	suite.tempDir, err = os.MkdirTemp("", "safepath_test")
	suite.Require().NoError(err)

	// Create some test files
	err = os.WriteFile(filepath.Join(suite.tempDir, "test.txt"), []byte("test content"), 0644)
	suite.Require().NoError(err)

	// Create a subdirectory with a file
	subDir := filepath.Join(suite.tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	suite.Require().NoError(err)

	err = os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested content"), 0644)
	suite.Require().NoError(err)
}

// TearDownTest removes the temporary directory after each test.
func (suite *SafePathTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// TestSafeJoin_ValidPath tests SafeJoin with a valid relative path.
func (suite *SafePathTestSuite) TestSafeJoin_ValidPath() {
	result, err := SafeJoin(suite.tempDir, "test.txt")
	suite.NoError(err)
	suite.Equal(filepath.Join(suite.tempDir, "test.txt"), result)
}

// TestSafeJoin_NestedPath tests SafeJoin with a valid nested path.
func (suite *SafePathTestSuite) TestSafeJoin_NestedPath() {
	result, err := SafeJoin(suite.tempDir, "subdir/nested.txt")
	suite.NoError(err)
	expected := filepath.Join(suite.tempDir, "subdir", "nested.txt")
	suite.Equal(expected, result)
}

// TestSafeJoin_RootPath tests SafeJoin with the root itself.
func (suite *SafePathTestSuite) TestSafeJoin_RootPath() {
	result, err := SafeJoin(suite.tempDir, ".")
	suite.NoError(err)
	absRoot, _ := filepath.Abs(suite.tempDir)
	suite.Equal(absRoot, result)
}

// TestSafeJoin_TraversalAttack tests SafeJoin rejects directory traversal.
func (suite *SafePathTestSuite) TestSafeJoin_TraversalAttack() {
	_, err := SafeJoin(suite.tempDir, "../../../etc/passwd")
	suite.Error(err)
	suite.Contains(err.Error(), "path escapes root directory")
}

// TestSafeJoin_TraversalInMiddle tests SafeJoin rejects traversal in middle of path.
func (suite *SafePathTestSuite) TestSafeJoin_TraversalInMiddle() {
	_, err := SafeJoin(suite.tempDir, "subdir/../../etc/passwd")
	suite.Error(err)
	suite.Contains(err.Error(), "path escapes root directory")
}

// TestSafeJoin_AbsolutePath tests SafeJoin rejects absolute paths.
func (suite *SafePathTestSuite) TestSafeJoin_AbsolutePath() {
	_, err := SafeJoin(suite.tempDir, "/etc/passwd")
	suite.Error(err)
	suite.Contains(err.Error(), "absolute paths not allowed")
}

// TestSafeJoin_DotDotOnly tests SafeJoin rejects just ".."
func (suite *SafePathTestSuite) TestSafeJoin_DotDotOnly() {
	_, err := SafeJoin(suite.tempDir, "..")
	suite.Error(err)
	suite.Contains(err.Error(), "path escapes root directory")
}

// TestSafeJoin_EmptyPath tests SafeJoin with empty path returns root.
func (suite *SafePathTestSuite) TestSafeJoin_EmptyPath() {
	result, err := SafeJoin(suite.tempDir, "")
	suite.NoError(err)
	absRoot, _ := filepath.Abs(suite.tempDir)
	suite.Equal(absRoot, result)
}

// TestReadFile_ValidFile tests ReadFile with a valid file.
func (suite *SafePathTestSuite) TestReadFile_ValidFile() {
	data, err := ReadFile(suite.tempDir, "test.txt")
	suite.NoError(err)
	suite.Equal("test content", string(data))
}

// TestReadFile_NestedFile tests ReadFile with a nested file.
func (suite *SafePathTestSuite) TestReadFile_NestedFile() {
	data, err := ReadFile(suite.tempDir, "subdir/nested.txt")
	suite.NoError(err)
	suite.Equal("nested content", string(data))
}

// TestReadFile_NonExistent tests ReadFile with a non-existent file.
func (suite *SafePathTestSuite) TestReadFile_NonExistent() {
	_, err := ReadFile(suite.tempDir, "nonexistent.txt")
	suite.Error(err)
}

// TestReadFile_TraversalAttack tests ReadFile rejects directory traversal.
func (suite *SafePathTestSuite) TestReadFile_TraversalAttack() {
	_, err := ReadFile(suite.tempDir, "../../../etc/passwd")
	suite.Error(err)
	suite.Contains(err.Error(), "path escapes root directory")
}

// TestExists_ExistingFile tests Exists with an existing file.
func (suite *SafePathTestSuite) TestExists_ExistingFile() {
	exists := Exists(suite.tempDir, "test.txt")
	suite.True(exists)
}

// TestExists_NonExistentFile tests Exists with a non-existent file.
func (suite *SafePathTestSuite) TestExists_NonExistentFile() {
	exists := Exists(suite.tempDir, "nonexistent.txt")
	suite.False(exists)
}

// TestExists_NestedFile tests Exists with a nested file.
func (suite *SafePathTestSuite) TestExists_NestedFile() {
	exists := Exists(suite.tempDir, "subdir/nested.txt")
	suite.True(exists)
}

// TestExists_TraversalAttack tests Exists rejects directory traversal.
func (suite *SafePathTestSuite) TestExists_TraversalAttack() {
	exists := Exists(suite.tempDir, "../../../etc/passwd")
	suite.False(exists) // Returns false on error
}

// TestExists_Directory tests Exists with a directory.
func (suite *SafePathTestSuite) TestExists_Directory() {
	exists := Exists(suite.tempDir, "subdir")
	suite.True(exists)
}

// TestStat_ExistingFile tests Stat with an existing file.
func (suite *SafePathTestSuite) TestStat_ExistingFile() {
	info, err := Stat(suite.tempDir, "test.txt")
	suite.NoError(err)
	suite.NotNil(info)
	suite.Equal("test.txt", info.Name())
	suite.False(info.IsDir())
}

// TestStat_Directory tests Stat with a directory.
func (suite *SafePathTestSuite) TestStat_Directory() {
	info, err := Stat(suite.tempDir, "subdir")
	suite.NoError(err)
	suite.NotNil(info)
	suite.True(info.IsDir())
}

// TestStat_NonExistent tests Stat with a non-existent file.
func (suite *SafePathTestSuite) TestStat_NonExistent() {
	_, err := Stat(suite.tempDir, "nonexistent.txt")
	suite.Error(err)
}

// TestStat_TraversalAttack tests Stat rejects directory traversal.
func (suite *SafePathTestSuite) TestStat_TraversalAttack() {
	_, err := Stat(suite.tempDir, "../../../etc/passwd")
	suite.Error(err)
	suite.Contains(err.Error(), "path escapes root directory")
}

// TestHasPrefix_Match tests hasPrefix with matching prefix.
func (suite *SafePathTestSuite) TestHasPrefix_Match() {
	result := hasPrefix("/foo/bar/baz", "/foo/bar")
	suite.True(result)
}

// TestHasPrefix_NoMatch tests hasPrefix with non-matching prefix.
func (suite *SafePathTestSuite) TestHasPrefix_NoMatch() {
	result := hasPrefix("/foo/bar", "/baz")
	suite.False(result)
}

// TestHasPrefix_PartialMatch tests hasPrefix behavior with similar paths.
func (suite *SafePathTestSuite) TestHasPrefix_PartialMatch() {
	// Note: hasPrefix compares cleaned paths. /foo/ becomes /foo after cleaning.
	// The SafeJoin function handles the "foobar vs foo" case by adding a separator.
	// This test verifies the raw hasPrefix behavior.
	result := hasPrefix("/foo/bar", "/foo/")
	suite.True(result) // /foo/bar does start with /foo (after cleaning /foo/ to /foo)
}

// TestHasPrefix_ExactMatch tests hasPrefix with exact match.
func (suite *SafePathTestSuite) TestHasPrefix_ExactMatch() {
	result := hasPrefix("/foo/bar", "/foo/bar")
	suite.True(result)
}

// TestHasPrefix_ShorterPath tests hasPrefix with shorter path than prefix.
func (suite *SafePathTestSuite) TestHasPrefix_ShorterPath() {
	result := hasPrefix("/foo", "/foo/bar")
	suite.False(result)
}

// TestSafeJoin_SymlinkInPath tests SafeJoin with symlinks.
func (suite *SafePathTestSuite) TestSafeJoin_SymlinkInPath() {
	// Create a symlink inside tempDir pointing outside
	symlinkPath := filepath.Join(suite.tempDir, "link")
	err := os.Symlink("/tmp", symlinkPath)
	if err != nil {
		suite.T().Skip("Cannot create symlinks (may require elevated privileges)")
	}

	// SafeJoin should still work as it validates the logical path
	result, err := SafeJoin(suite.tempDir, "link")
	suite.NoError(err)
	suite.Contains(result, "link")
}

// TestSafePathTestSuite runs all the tests in the suite.
func TestSafePathTestSuite(t *testing.T) {
	suite.Run(t, new(SafePathTestSuite))
}
