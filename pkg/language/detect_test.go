package language

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// DetectTestSuite is the test suite for the language detection package.
type DetectTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *DetectTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-detect-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *DetectTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createFile creates a file in the temp directory.
func (suite *DetectTestSuite) createFile(name string) {
	filePath := filepath.Join(suite.tempDir, name)
	// Create parent directories if needed
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0700)
	suite.NoError(err)
	err = os.WriteFile(filePath, []byte(""), 0600)
	suite.NoError(err)
}

// TestDetect_Go tests detection of Go projects.
func (suite *DetectTestSuite) TestDetect_Go() {
	suite.createFile("go.mod")

	result := Detect(suite.tempDir)

	suite.True(result.HasLanguage(checker.LangGo))
	suite.Equal(checker.LangGo, result.Primary)
	suite.Contains(result.Indicators[checker.LangGo], "go.mod")
}

// TestDetect_Rust tests detection of Rust projects.
func (suite *DetectTestSuite) TestDetect_Rust() {
	suite.createFile("Cargo.toml")

	result := Detect(suite.tempDir)

	suite.True(result.HasLanguage(checker.LangRust))
	suite.Equal(checker.LangRust, result.Primary)
	suite.Contains(result.Indicators[checker.LangRust], "Cargo.toml")
}

// TestDetect_Node tests detection of Node.js projects.
func (suite *DetectTestSuite) TestDetect_Node() {
	suite.createFile("package.json")

	result := Detect(suite.tempDir)

	suite.True(result.HasLanguage(checker.LangNode))
	suite.Equal(checker.LangNode, result.Primary)
}

// TestDetect_MultiLanguage tests detection of multi-language projects.
func (suite *DetectTestSuite) TestDetect_MultiLanguage() {
	suite.createFile("package.json")
	suite.createFile("Cargo.toml")

	result := Detect(suite.tempDir)

	suite.True(result.MultiLang)
	suite.True(result.HasLanguage(checker.LangNode))
	suite.True(result.HasLanguage(checker.LangRust))
	suite.Len(result.Languages, 2)
}

// TestDetect_NoLanguage tests detection when no language indicators are found.
func (suite *DetectTestSuite) TestDetect_NoLanguage() {
	// Empty directory
	result := Detect(suite.tempDir)

	suite.Empty(result.Languages)
	suite.False(result.MultiLang)
}

// TestDetectWithOverride tests explicit language override.
func (suite *DetectTestSuite) TestDetectWithOverride() {
	// Even without indicator files, explicit override should work
	result := DetectWithOverride(suite.tempDir, []checker.Language{checker.LangRust, checker.LangNode})

	suite.True(result.HasLanguage(checker.LangRust))
	suite.True(result.HasLanguage(checker.LangNode))
	suite.Equal(checker.LangRust, result.Primary)
	suite.True(result.MultiLang)
	suite.Nil(result.Indicators) // No indicator files when explicit
}

// TestDetectWithSourceDirs tests detection with custom source directories.
func (suite *DetectTestSuite) TestDetectWithSourceDirs() {
	// Create a Tauri-like project structure
	// - Node.js in root
	// - Rust in src-tauri/
	suite.createFile("package.json")
	suite.createFile("src-tauri/Cargo.toml")

	// Without source_dir configuration, Rust won't be detected
	result := Detect(suite.tempDir)
	suite.True(result.HasLanguage(checker.LangNode))
	suite.False(result.HasLanguage(checker.LangRust))

	// With source_dir configuration, Rust should be detected
	sourceDirs := map[string]string{
		"rust": "src-tauri",
	}
	result = DetectWithSourceDirs(suite.tempDir, sourceDirs)

	suite.True(result.HasLanguage(checker.LangNode))
	suite.True(result.HasLanguage(checker.LangRust))
	suite.True(result.MultiLang)
}

// TestDetectWithSourceDirs_EmptyMap tests detection with empty source directories map.
func (suite *DetectTestSuite) TestDetectWithSourceDirs_EmptyMap() {
	suite.createFile("go.mod")

	result := DetectWithSourceDirs(suite.tempDir, map[string]string{})

	suite.True(result.HasLanguage(checker.LangGo))
	suite.Equal(checker.LangGo, result.Primary)
}

// TestDetectWithSourceDirs_MultipleConfigs tests detection with multiple configured source directories.
func (suite *DetectTestSuite) TestDetectWithSourceDirs_MultipleConfigs() {
	// Create a complex monorepo structure
	suite.createFile("backend/go.mod")
	suite.createFile("frontend/package.json")
	suite.createFile("desktop/Cargo.toml")

	sourceDirs := map[string]string{
		"go":   "backend",
		"node": "frontend",
		"rust": "desktop",
	}
	result := DetectWithSourceDirs(suite.tempDir, sourceDirs)

	suite.True(result.HasLanguage(checker.LangGo))
	suite.True(result.HasLanguage(checker.LangNode))
	suite.True(result.HasLanguage(checker.LangRust))
	suite.Len(result.Languages, 3)
	suite.True(result.MultiLang)
}

// TestHasLanguage tests the HasLanguage method.
func (suite *DetectTestSuite) TestHasLanguage() {
	result := DetectionResult{
		Languages: []checker.Language{checker.LangGo, checker.LangRust},
	}

	suite.True(result.HasLanguage(checker.LangGo))
	suite.True(result.HasLanguage(checker.LangRust))
	suite.False(result.HasLanguage(checker.LangNode))
	suite.False(result.HasLanguage(checker.LangPython))
}

// TestDetectTestSuite runs all the tests in the suite.
func TestDetectTestSuite(t *testing.T) {
	suite.Run(t, new(DetectTestSuite))
}
