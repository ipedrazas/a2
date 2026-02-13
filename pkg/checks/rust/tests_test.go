package rustcheck

import (
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type TestsTestSuite struct {
	suite.Suite
	tempDir string
	check   *TestsCheck
}

func (s *TestsTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "rust-tests-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &TestsCheck{}
}

func (s *TestsTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *TestsTestSuite) TestID() {
	s.Equal("rust:tests", s.check.ID())
}

func (s *TestsTestSuite) TestName() {
	s.Equal("Rust Tests", s.check.Name())
}

func (s *TestsTestSuite) TestRun_NoCargoToml() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No Cargo.toml found")
}

func (s *TestsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangRust, result.Language)
}

func TestTestsTestSuite(t *testing.T) {
	suite.Run(t, new(TestsTestSuite))
}
