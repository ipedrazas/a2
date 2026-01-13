package typescriptcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/stretchr/testify/suite"
)

type LintTestSuite struct {
	suite.Suite
	tempDir string
	check   *LintCheck
}

func (s *LintTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "ts-lint-test-*")
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
	s.Equal("typescript:lint", s.check.ID())
}

func (s *LintTestSuite) TestName() {
	s.Equal("TypeScript Lint", s.check.Name())
}

func (s *LintTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No tsconfig.json found")
}

func (s *LintTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *LintTestSuite) TestDetectLinter_ESLint() {
	s.writeFile(".eslintrc.json", "{}")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("eslint", linter)
}

func (s *LintTestSuite) TestDetectLinter_Biome() {
	s.writeFile("biome.json", "{}")
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("biome", linter)
}

func (s *LintTestSuite) TestDetectLinter_None() {
	linter := s.check.detectLinter(s.tempDir)
	s.Equal("", linter)
}

func (s *LintTestSuite) TestParseESLintOutput_WithSummary() {
	output := `
/path/to/file.ts
  1:10  error  Something wrong  no-unused-vars
  5:3   warning  Another issue  prefer-const

✖ 2 problems (1 error, 1 warning)
`
	errors, warnings := parseESLintOutput(output)
	s.Equal(1, errors)
	s.Equal(1, warnings)
}

func (s *LintTestSuite) TestPluralize() {
	s.Equal("error", checkutil.Pluralize(1, "error", "errors"))
	s.Equal("errors", checkutil.Pluralize(2, "error", "errors"))
	s.Equal("errors", checkutil.Pluralize(0, "error", "errors"))
}

func TestLintTestSuite(t *testing.T) {
	suite.Run(t, new(LintTestSuite))
}
