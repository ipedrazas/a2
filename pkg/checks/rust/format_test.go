package rustcheck

import (
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type FormatTestSuite struct {
	suite.Suite
	tempDir string
	check   *FormatCheck
}

func (s *FormatTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "rust-format-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &FormatCheck{}
}

func (s *FormatTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *FormatTestSuite) TestID() {
	s.Equal("rust:format", s.check.ID())
}

func (s *FormatTestSuite) TestName() {
	s.Equal("Rust Format", s.check.Name())
}

func (s *FormatTestSuite) TestRun_NoCargoToml() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No Cargo.toml found")
}

func (s *FormatTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangRust, result.Language)
}

func TestFormatTestSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}
