package rustcheck

import (
	"os"
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
	dir, err := os.MkdirTemp("", "rust-lint-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &LintCheck{}
}

func (s *LintTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *LintTestSuite) TestID() {
	s.Equal("rust:lint", s.check.ID())
}

func (s *LintTestSuite) TestName() {
	s.Equal("Rust Clippy", s.check.Name())
}

func (s *LintTestSuite) TestRun_NoCargoToml() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Cargo.toml found")
}

func (s *LintTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangRust, result.Language)
}

func TestLintTestSuite(t *testing.T) {
	suite.Run(t, new(LintTestSuite))
}
