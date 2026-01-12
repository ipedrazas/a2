package typescriptcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type DepsTestSuite struct {
	suite.Suite
	tempDir string
	check   *DepsCheck
}

func (s *DepsTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "ts-deps-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &DepsCheck{}
}

func (s *DepsTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *DepsTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *DepsTestSuite) TestID() {
	s.Equal("typescript:deps", s.check.ID())
}

func (s *DepsTestSuite) TestName() {
	s.Equal("TypeScript Vulnerabilities", s.check.Name())
}

func (s *DepsTestSuite) TestRun_NoPackageJson() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No package.json found")
}

func (s *DepsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *DepsTestSuite) TestDetectPackageManager_Npm() {
	s.writeFile("package.json", `{}`)
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("npm", pm)
}

func (s *DepsTestSuite) TestDetectPackageManager_Yarn() {
	s.writeFile("yarn.lock", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("yarn", pm)
}

func (s *DepsTestSuite) TestDetectPackageManager_Pnpm() {
	s.writeFile("pnpm-lock.yaml", "")
	pm := s.check.detectPackageManager(s.tempDir)
	s.Equal("pnpm", pm)
}

func (s *DepsTestSuite) TestHasSnyk_WithFile() {
	s.writeFile(".snyk", "{}")
	s.True(s.check.hasSnyk(s.tempDir))
}

func (s *DepsTestSuite) TestHasSnyk_WithCI() {
	s.writeFile(".github/workflows/ci.yml", `
jobs:
  security:
    steps:
      - uses: snyk/actions/node@master
`)
	s.True(s.check.hasSnyk(s.tempDir))
}

func (s *DepsTestSuite) TestHasSnyk_NoSnyk() {
	s.False(s.check.hasSnyk(s.tempDir))
}

func (s *DepsTestSuite) TestRun_WithDependabot() {
	s.writeFile("package.json", `{"name": "test"}`)
	s.writeFile(".github/dependabot.yml", `
version: 2
updates:
  - package-ecosystem: npm
`)
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Dependabot")
}

func (s *DepsTestSuite) TestRun_WithRenovate() {
	s.writeFile("package.json", `{"name": "test"}`)
	s.writeFile("renovate.json", `{}`)
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Renovate")
}

func TestDepsTestSuite(t *testing.T) {
	suite.Run(t, new(DepsTestSuite))
}
