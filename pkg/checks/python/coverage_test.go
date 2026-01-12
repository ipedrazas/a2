package pythoncheck

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
	dir, err := os.MkdirTemp("", "python-coverage-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &CoverageCheck{}
}

func (s *CoverageTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *CoverageTestSuite) TestID() {
	s.Equal("python:coverage", s.check.ID())
}

func (s *CoverageTestSuite) TestName() {
	s.Equal("Python Coverage", s.check.Name())
}

func (s *CoverageTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func (s *CoverageTestSuite) TestParsePythonCoverage_Total() {
	output := `
Name                 Stmts   Miss  Cover
----------------------------------------
src/main.py             50     10    80%
src/utils.py            30      5    83%
----------------------------------------
TOTAL                   80     15    81%
`
	cov := parsePythonCoverage(output)
	s.Equal(81.0, cov)
}

func (s *CoverageTestSuite) TestParsePythonCoverage_Percentage() {
	output := "Coverage: 75.5%"
	cov := parsePythonCoverage(output)
	s.Equal(75.5, cov)
}

func (s *CoverageTestSuite) TestParsePythonCoverage_NoCoverage() {
	output := "no tests ran"
	cov := parsePythonCoverage(output)
	s.Equal(0.0, cov)
}

func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
