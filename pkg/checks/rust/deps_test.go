package rustcheck

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
	dir, err := os.MkdirTemp("", "rust-deps-test-*")
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
	s.Equal("rust:deps", s.check.ID())
}

func (s *DepsTestSuite) TestName() {
	s.Equal("Rust Vulnerabilities", s.check.Name())
}

func (s *DepsTestSuite) TestRun_NoCargoToml() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No Cargo.toml found")
}

func (s *DepsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangRust, result.Language)
}

func (s *DepsTestSuite) TestCheckCIForAudit_GitHubActions() {
	s.writeFile(".github/workflows/ci.yml", `
name: CI
jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - run: cargo audit
`)
	s.True(s.check.checkCIForAudit(s.tempDir))
}

func (s *DepsTestSuite) TestCheckCIForAudit_NoCI() {
	s.False(s.check.checkCIForAudit(s.tempDir))
}

func TestDepsTestSuite(t *testing.T) {
	suite.Run(t, new(DepsTestSuite))
}
