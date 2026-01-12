package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type FormatTestSuite struct {
	suite.Suite
	tempDir string
	check   *FormatCheck
}

func (s *FormatTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "python-format-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &FormatCheck{}
}

func (s *FormatTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *FormatTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *FormatTestSuite) TestID() {
	s.Equal("python:format", s.check.ID())
}

func (s *FormatTestSuite) TestName() {
	s.Equal("Python Format", s.check.Name())
}

func (s *FormatTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func (s *FormatTestSuite) TestDetectFormatter_Ruff() {
	s.writeFile("ruff.toml", "")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("ruff", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_RuffDot() {
	s.writeFile(".ruff.toml", "")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("ruff", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_Black() {
	s.writeFile("pyproject.toml", "[tool.black]")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("black", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_RuffInPyproject() {
	s.writeFile("pyproject.toml", "[tool.ruff]")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("ruff", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_Auto() {
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("auto", formatter)
}

func (s *FormatTestSuite) TestPluralize() {
	s.Equal("1 file", pluralize(1, "file", "files"))
	s.Contains(pluralize(2, "file", "files"), "files")
}

func TestFormatTestSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}
