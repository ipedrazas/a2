package pythoncheck

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
	dir, err := os.MkdirTemp("", "python-project-test-*")
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
	s.Equal("python:project", s.check.ID())
}

func (s *ProjectTestSuite) TestName() {
	s.Equal("Python Project", s.check.Name())
}

func (s *ProjectTestSuite) TestRun_NoProject() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Python project configuration found")
}

func (s *ProjectTestSuite) TestRun_PyprojectToml() {
	s.writeFile("pyproject.toml", `[project]
name = "test"`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "pyproject.toml")
}

func (s *ProjectTestSuite) TestRun_SetupPy() {
	s.writeFile("setup.py", `from setuptools import setup
setup(name='test')`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "setup.py")
	s.Contains(result.Message, "consider migrating")
}

func (s *ProjectTestSuite) TestRun_RequirementsTxt() {
	s.writeFile("requirements.txt", `flask==2.0.0`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "requirements.txt")
}

func (s *ProjectTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
