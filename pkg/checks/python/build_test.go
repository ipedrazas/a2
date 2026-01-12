package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type BuildTestSuite struct {
	suite.Suite
	tempDir string
	check   *BuildCheck
}

func (s *BuildTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "python-build-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &BuildCheck{}
}

func (s *BuildTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *BuildTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *BuildTestSuite) TestID() {
	s.Equal("python:build", s.check.ID())
}

func (s *BuildTestSuite) TestName() {
	s.Equal("Python Build", s.check.Name())
}

func (s *BuildTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func (s *BuildTestSuite) TestDetectPackageManager_Pip() {
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("pip", pm)
}

func (s *BuildTestSuite) TestDetectPackageManager_Poetry() {
	s.writeFile("poetry.lock", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("poetry", pm)
}

func (s *BuildTestSuite) TestDetectPackageManager_Pipenv() {
	s.writeFile("Pipfile", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("pipenv", pm)
}

func TestBuildTestSuite(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
