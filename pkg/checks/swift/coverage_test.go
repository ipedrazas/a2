package swiftcheck

import (
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type CoverageTestSuite struct {
	suite.Suite
	tempDir string
	check   *CoverageCheck
}

func (s *CoverageTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "swift-coverage-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &CoverageCheck{Threshold: 80.0}
}

func (s *CoverageTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *CoverageTestSuite) TestID() {
	s.Equal("swift:coverage", s.check.ID())
}

func (s *CoverageTestSuite) TestName() {
	s.Equal("Swift Coverage", s.check.Name())
}

func (s *CoverageTestSuite) TestRun_NoPackageSwift() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Package.swift found")
}

func (s *CoverageTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangSwift, result.Language)
}

func (s *CoverageTestSuite) TestThresholdDefault() {
	check := &CoverageCheck{}
	s.Equal(0.0, check.Threshold)
}

func (s *CoverageTestSuite) TestThresholdSet() {
	check := &CoverageCheck{Threshold: 90.0}
	s.Equal(90.0, check.Threshold)
}

func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
