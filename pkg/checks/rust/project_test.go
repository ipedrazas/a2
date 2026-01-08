package rustcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ProjectTestSuite struct {
	suite.Suite
	tempDir string
	check   *ProjectCheck
}

func (s *ProjectTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "rust-project-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &ProjectCheck{}
}

func (s *ProjectTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *ProjectTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *ProjectTestSuite) TestID() {
	s.Equal("rust:project", s.check.ID())
}

func (s *ProjectTestSuite) TestName() {
	s.Equal("Rust Project", s.check.Name())
}

func (s *ProjectTestSuite) TestRun_NoCargoToml() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Cargo.toml found")
}

func (s *ProjectTestSuite) TestRun_BasicCargoToml() {
	s.writeFile("Cargo.toml", `
[package]
name = "myapp"
version = "1.0.0"
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "myapp")
	s.Contains(result.Message, "v1.0.0")
}

func (s *ProjectTestSuite) TestRun_Workspace() {
	s.writeFile("Cargo.toml", `
[workspace]
members = ["crate1", "crate2"]
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "workspace")
}

func (s *ProjectTestSuite) TestRun_MinimalCargoToml() {
	s.writeFile("Cargo.toml", `
[package]
name = "test"
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *ProjectTestSuite) TestRun_ResultLanguage() {
	s.writeFile("Cargo.toml", `[package]
name = "test"
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangRust, result.Language)
}

func (s *ProjectTestSuite) TestExtractTomlValue() {
	tests := []struct {
		content  string
		key      string
		expected string
	}{
		{`name = "myapp"`, "name", "myapp"},
		{`version = "1.0.0"`, "version", "1.0.0"},
		{`name="nospace"`, "name", "nospace"},
		{`name = 'single'`, "name", "single"},
		{`other = "value"`, "name", ""},
	}

	for _, tc := range tests {
		result := extractTomlValue(tc.content, tc.key)
		s.Equal(tc.expected, result, "for content: %s, key: %s", tc.content, tc.key)
	}
}

func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
