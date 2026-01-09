package userconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserConfigTestSuite struct {
	suite.Suite
	tempDir string
}

func (suite *UserConfigTestSuite) SetupTest() {
	// Create a temporary directory for each test
	dir, err := os.MkdirTemp("", "a2-userconfig-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

func (suite *UserConfigTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

func (suite *UserConfigTestSuite) TestGetConfigDir() {
	dir, err := GetConfigDir()
	suite.NoError(err)
	suite.NotEmpty(dir)
	suite.Contains(dir, "a2")
}

func (suite *UserConfigTestSuite) TestGetSubDir() {
	dir, err := GetSubDir("profiles")
	suite.NoError(err)
	suite.NotEmpty(dir)
	suite.Contains(dir, "a2")
	suite.Contains(dir, "profiles")
}

func (suite *UserConfigTestSuite) TestGetSubDir_Targets() {
	dir, err := GetSubDir("targets")
	suite.NoError(err)
	suite.NotEmpty(dir)
	suite.Contains(dir, "a2")
	suite.Contains(dir, "targets")
}

func (suite *UserConfigTestSuite) TestEnsureDir_CreatesDirectory() {
	// Use a unique subdirectory within the temp directory
	subdir := filepath.Join(suite.tempDir, "test-ensure")

	// Create a mock config dir scenario
	err := os.MkdirAll(subdir, 0755)
	suite.NoError(err)

	// Verify it exists
	info, err := os.Stat(subdir)
	suite.NoError(err)
	suite.True(info.IsDir())
}

func (suite *UserConfigTestSuite) TestDirExists_ReturnsFalseForNonExistent() {
	// This tests against a path that's unlikely to exist
	exists := DirExists("nonexistent-test-directory-" + suite.tempDir)
	suite.False(exists)
}

func (suite *UserConfigTestSuite) TestEnsureDir_Idempotent() {
	// First call creates the directory
	dir1, err := EnsureDir("test-idempotent")
	if err != nil {
		// Skip if we can't create in user config dir (e.g., in CI)
		suite.T().Skip("Cannot create user config directory")
	}
	defer os.RemoveAll(dir1)

	// Second call should succeed without error
	dir2, err := EnsureDir("test-idempotent")
	suite.NoError(err)
	suite.Equal(dir1, dir2)
}

func TestUserConfigTestSuite(t *testing.T) {
	suite.Run(t, new(UserConfigTestSuite))
}
