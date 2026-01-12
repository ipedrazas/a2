package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type TestsTestSuite struct {
	suite.Suite
	tempDir string
	check   *TestsCheck
}

func (s *TestsTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "python-tests-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &TestsCheck{}
}

func (s *TestsTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *TestsTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *TestsTestSuite) TestID() {
	s.Equal("python:tests", s.check.ID())
}

func (s *TestsTestSuite) TestName() {
	s.Equal("Python Tests", s.check.Name())
}

func (s *TestsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func (s *TestsTestSuite) TestDetectTestRunner_Pytest_Ini() {
	s.writeFile("pytest.ini", "")
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("pytest", runner)
}

func (s *TestsTestSuite) TestDetectTestRunner_Conftest() {
	s.writeFile("conftest.py", "")
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("pytest", runner)
}

func (s *TestsTestSuite) TestDetectTestRunner_Default() {
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("pytest", runner)
}

func (s *TestsTestSuite) TestTruncateMessage_Short() {
	msg := "Short message"
	result := truncateMessage(msg, 100)
	s.Equal(msg, result)
}

func (s *TestsTestSuite) TestTruncateMessage_Long() {
	msg := "This is a very long message that exceeds the maximum length"
	result := truncateMessage(msg, 20)
	s.Equal("This is a very long ...", result)
}

func TestTestsTestSuite(t *testing.T) {
	suite.Run(t, new(TestsTestSuite))
}
