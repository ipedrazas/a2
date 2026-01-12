package swiftcheck

import (
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

type LintTestSuite struct {
	suite.Suite
	tempDir string
	check   *LintCheck
}

func (s *LintTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "swift-lint-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &LintCheck{Config: &config.SwiftLanguageConfig{}}
}

func (s *LintTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *LintTestSuite) TestID() {
	s.Equal("swift:lint", s.check.ID())
}

func (s *LintTestSuite) TestName() {
	s.Equal("Swift Lint", s.check.Name())
}

func (s *LintTestSuite) TestRun_NoPackageSwift() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Package.swift found")
}

func (s *LintTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangSwift, result.Language)
}

func (s *LintTestSuite) TestFormatIssueCount() {
	s.Equal("1", formatIssueCount(1))
	s.Equal("5", formatIssueCount(5))
	s.Equal("9", formatIssueCount(9))
	s.Equal("multiple", formatIssueCount(10))
	s.Equal("multiple", formatIssueCount(100))
}

func TestLintTestSuite(t *testing.T) {
	suite.Run(t, new(LintTestSuite))
}
