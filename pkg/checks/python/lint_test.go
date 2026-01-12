package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type LintTestSuite struct {
	suite.Suite
	tempDir string
	check   *LintCheck
}

func (s *LintTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "python-lint-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &LintCheck{}
}

func (s *LintTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *LintTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *LintTestSuite) TestID() {
	s.Equal("python:lint", s.check.ID())
}

func (s *LintTestSuite) TestName() {
	s.Equal("Python Lint", s.check.Name())
}

func (s *LintTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func (s *LintTestSuite) TestDetectLinter_Ruff() {
	s.writeFile("ruff.toml", "")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("ruff", linter)
}

func (s *LintTestSuite) TestDetectLinter_RuffDot() {
	s.writeFile(".ruff.toml", "")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("ruff", linter)
}

func (s *LintTestSuite) TestDetectLinter_Flake8() {
	s.writeFile(".flake8", "")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("flake8", linter)
}

func (s *LintTestSuite) TestDetectLinter_Pylint() {
	s.writeFile(".pylintrc", "")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("pylint", linter)
}

func (s *LintTestSuite) TestDetectLinter_RuffInPyproject() {
	s.writeFile("pyproject.toml", "[tool.ruff]")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("ruff", linter)
}

func (s *LintTestSuite) TestDetectLinter_Flake8InPyproject() {
	s.writeFile("pyproject.toml", "[tool.flake8]")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("flake8", linter)
}

func (s *LintTestSuite) TestDetectLinter_Auto() {
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("auto", linter)
}

func TestLintTestSuite(t *testing.T) {
	suite.Run(t, new(LintTestSuite))
}
