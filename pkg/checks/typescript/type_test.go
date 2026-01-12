package typescriptcheck

import (
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type TypeTestSuite struct {
	suite.Suite
	tempDir string
	check   *TypeCheck
}

func (s *TypeTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "ts-type-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &TypeCheck{}
}

func (s *TypeTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *TypeTestSuite) TestID() {
	s.Equal("typescript:type", s.check.ID())
}

func (s *TypeTestSuite) TestName() {
	s.Equal("TypeScript Type Check", s.check.Name())
}

func (s *TypeTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No tsconfig.json found")
}

func (s *TypeTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *TypeTestSuite) TestCountTypeErrors_FromSummary() {
	output := `
src/index.ts(5,10): error TS2322: Type 'string' is not assignable to type 'number'.
src/utils.ts(12,3): error TS2304: Cannot find name 'foo'.

Found 2 errors.
`
	count := countTypeErrors(output)
	s.Equal(2, count)
}

func (s *TypeTestSuite) TestCountTypeErrors_SingleError() {
	output := `
src/index.ts(5,10): error TS2322: Type 'string' is not assignable to type 'number'.

Found 1 error.
`
	count := countTypeErrors(output)
	s.Equal(1, count)
}

func (s *TypeTestSuite) TestCountTypeErrors_FallbackCounting() {
	output := `
src/index.ts(5,10): error TS2322: Type 'string' is not assignable to type 'number'.
src/utils.ts(12,3): error TS2304: Cannot find name 'foo'.
`
	count := countTypeErrors(output)
	s.Equal(2, count)
}

func (s *TypeTestSuite) TestCountTypeErrors_NoErrors() {
	output := "No errors found"
	count := countTypeErrors(output)
	s.Equal(0, count)
}

func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}
