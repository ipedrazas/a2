package rustcheck

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CoverageTestSuite struct {
	suite.Suite
	check *CoverageCheck
}

func (s *CoverageTestSuite) SetupTest() {
	s.check = &CoverageCheck{Threshold: 80.0}
}

func (s *CoverageTestSuite) TestID() {
	s.Equal("rust:coverage", s.check.ID())
}

func (s *CoverageTestSuite) TestName() {
	s.Equal("Rust Coverage", s.check.Name())
}

func (s *CoverageTestSuite) TestParseTarpaulinCoverage() {
	output := `running 10 tests
test result: ok. 10 passed; 0 failed

85.00% coverage, 170/200 lines covered`
	s.Equal(85.0, parseTarpaulinCoverage(output))
}

func (s *CoverageTestSuite) TestParseTarpaulinCoverage_Empty() {
	s.Equal(0.0, parseTarpaulinCoverage(""))
}

func (s *CoverageTestSuite) TestParseLlvmCovCoverage() {
	output := `Filename  Regions  Missed  Cover  Lines  Missed  Cover
---
src/lib.rs  50  10  80.00%  200  30  85.00%
---
TOTAL  50  10  80.00%  200  30  85.00%`
	// Matches TOTAL line: "TOTAL  50  10  80.00%"
	// The regex matches first 3 number groups then %
	cov := parseLlvmCovCoverage(output)
	s.InDelta(80.0, cov, 0.1)
}

func (s *CoverageTestSuite) TestParseLlvmCovCoverage_Empty() {
	s.Equal(0.0, parseLlvmCovCoverage(""))
}

func (s *CoverageTestSuite) TestParseLlvmCovCoverage_FallbackPercent() {
	// When no TOTAL line, falls back to last percentage in final lines
	output := `some output
coverage: 72.50%`
	s.Equal(72.5, parseLlvmCovCoverage(output))
}

func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
