package typescriptcheck

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
	dir, err := os.MkdirTemp("", "ts-build-test-*")
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
	s.Equal("typescript:build", s.check.ID())
}

func (s *BuildTestSuite) TestName() {
	s.Equal("TypeScript Build", s.check.Name())
}

func (s *BuildTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No tsconfig.json found")
}

func (s *BuildTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *BuildTestSuite) TestDetectPackageManager_Npm() {
	s.writeFile("package.json", `{}`)
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("npm", pm)
}

func (s *BuildTestSuite) TestDetectPackageManager_Yarn() {
	s.writeFile("yarn.lock", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("yarn", pm)
}

func (s *BuildTestSuite) TestDetectPackageManager_Pnpm() {
	s.writeFile("pnpm-lock.yaml", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("pnpm", pm)
}

func (s *BuildTestSuite) TestDetectPackageManager_Bun() {
	s.writeFile("bun.lockb", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("bun", pm)
}

func TestBuildTestSuite(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
