package typescriptcheck

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
	dir, err := os.MkdirTemp("", "ts-tests-test-*")
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
	s.Equal("typescript:tests", s.check.ID())
}

func (s *TestsTestSuite) TestName() {
	s.Equal("TypeScript Tests", s.check.Name())
}

func (s *TestsTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No tsconfig.json found")
}

func (s *TestsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *TestsTestSuite) TestDetectTestRunner_Jest() {
	s.writeFile("jest.config.js", "module.exports = {}")
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("jest", runner)
}

func (s *TestsTestSuite) TestDetectTestRunner_Vitest() {
	s.writeFile("vitest.config.ts", "export default {}")
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("vitest", runner)
}

func (s *TestsTestSuite) TestDetectTestRunner_Mocha() {
	s.writeFile(".mocharc.json", "{}")
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("mocha", runner)
}

func (s *TestsTestSuite) TestDetectTestRunner_NoRunner() {
	runner := s.check.detectTestRunner(s.tempDir)
	s.Equal("", runner)
}

func (s *TestsTestSuite) TestDetectPackageManager_Default() {
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("npm", pm)
}

func TestTestsTestSuite(t *testing.T) {
	suite.Run(t, new(TestsTestSuite))
}
