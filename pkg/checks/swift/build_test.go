package swiftcheck

import (
	"os"
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
	dir, err := os.MkdirTemp("", "swift-build-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &BuildCheck{}
}

func (s *BuildTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *BuildTestSuite) TestID() {
	s.Equal("swift:build", s.check.ID())
}

func (s *BuildTestSuite) TestName() {
	s.Equal("Swift Build", s.check.Name())
}

func (s *BuildTestSuite) TestRun_NoPackageSwift() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No Package.swift found")
}

func (s *BuildTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangSwift, result.Language)
}

func TestBuildTestSuite(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}
