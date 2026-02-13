package typescriptcheck

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
	dir, err := os.MkdirTemp("", "typescript-project-test-*")
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
	s.Equal("typescript:project", s.check.ID())
}

func (s *ProjectTestSuite) TestName() {
	s.Equal("TypeScript Project", s.check.Name())
}

func (s *ProjectTestSuite) TestRun_NoTsConfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No tsconfig.json found")
}

func (s *ProjectTestSuite) TestRun_BasicTsConfig() {
	s.writeFile("tsconfig.json", `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "strict": true
  }
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "tsconfig.json")
	s.Contains(result.Reason, "ES2020")
	s.Contains(result.Reason, "strict mode")
}

func (s *ProjectTestSuite) TestRun_TsConfigWithPackageJson() {
	s.writeFile("tsconfig.json", `{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext"
  }
}`)
	s.writeFile("package.json", `{
  "name": "my-app",
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "TypeScript ^5.0.0")
}

func (s *ProjectTestSuite) TestRun_TsConfigBase() {
	s.writeFile("tsconfig.base.json", `{
  "compilerOptions": {
    "target": "ES2021"
  }
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *ProjectTestSuite) TestRun_ResultLanguage() {
	s.writeFile("tsconfig.json", `{}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
