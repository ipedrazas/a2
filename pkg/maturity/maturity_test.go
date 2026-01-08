package maturity

import (
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/stretchr/testify/suite"
)

type MaturityTestSuite struct {
	suite.Suite
}

func (s *MaturityTestSuite) makeResult(passed, warnings, failed int) runner.SuiteResult {
	results := make([]checker.Result, passed+warnings+failed)
	idx := 0

	for i := 0; i < passed; i++ {
		results[idx] = checker.Result{Status: checker.Pass, Passed: true}
		idx++
	}
	for i := 0; i < warnings; i++ {
		results[idx] = checker.Result{Status: checker.Warn, Passed: true}
		idx++
	}
	for i := 0; i < failed; i++ {
		results[idx] = checker.Result{Status: checker.Fail, Passed: false}
		idx++
	}

	return runner.SuiteResult{
		Results:  results,
		Passed:   passed,
		Warnings: warnings,
		Failed:   failed,
	}
}

func (s *MaturityTestSuite) TestLevel_String() {
	s.Equal("Production-Ready", ProductionReady.String())
	s.Equal("Mature", Mature.String())
	s.Equal("Development", Development.String())
	s.Equal("Proof of Concept", PoC.String())
}

func (s *MaturityTestSuite) TestLevel_Description() {
	s.Contains(ProductionReady.Description(), "production deployment")
	s.Contains(Mature.Description(), "minor improvements")
	s.Contains(Development.Description(), "quality improvements")
	s.Contains(PoC.Description(), "Early stage")
}

func (s *MaturityTestSuite) TestEstimate_EmptyResults() {
	result := runner.SuiteResult{}
	est := Estimate(result)

	s.Equal(PoC, est.Level)
	s.Equal(float64(0), est.Score)
	s.Equal(0, est.Total)
}

func (s *MaturityTestSuite) TestEstimate_ProductionReady() {
	// 100% pass, 0 warnings, 0 failures
	result := s.makeResult(10, 0, 0)
	est := Estimate(result)

	s.Equal(ProductionReady, est.Level)
	s.Equal(float64(100), est.Score)
	s.Equal(10, est.Passed)
	s.Equal(0, est.Warnings)
	s.Equal(0, est.Failed)
	s.Empty(est.Suggestions)
}

func (s *MaturityTestSuite) TestEstimate_MatureWithWarnings() {
	// 80%+ pass, 0 failures, some warnings
	result := s.makeResult(8, 2, 0)
	est := Estimate(result)

	s.Equal(Mature, est.Level)
	s.Equal(float64(80), est.Score)
	s.NotEmpty(est.Suggestions)
	s.Contains(est.Suggestions[0], "warnings")
}

func (s *MaturityTestSuite) TestEstimate_MatureNoWarnings() {
	// 80%+ pass, 0 failures, 0 warnings (but not 100%)
	// This shouldn't happen in practice since pass+warn+fail = total
	// but the logic handles it
	result := s.makeResult(9, 0, 0)
	est := Estimate(result)

	// With 100% pass and 0 warnings, it's Production-Ready
	s.Equal(ProductionReady, est.Level)
}

func (s *MaturityTestSuite) TestEstimate_Development() {
	// 60-80% score with 1-2 failures
	result := s.makeResult(6, 2, 2)
	est := Estimate(result)

	s.Equal(Development, est.Level)
	s.Equal(float64(60), est.Score)
	s.NotEmpty(est.Suggestions)
}

func (s *MaturityTestSuite) TestEstimate_DevelopmentEdgeCase() {
	// Exactly 60% with 2 failures
	result := s.makeResult(6, 0, 4)
	est := Estimate(result)

	// 60% but 4 failures = PoC
	s.Equal(PoC, est.Level)
}

func (s *MaturityTestSuite) TestEstimate_PoC_LowScore() {
	// Below 60% score
	result := s.makeResult(3, 2, 5)
	est := Estimate(result)

	s.Equal(PoC, est.Level)
	s.Equal(float64(30), est.Score)
	s.NotEmpty(est.Suggestions)
	s.Contains(est.Suggestions[0], "critical checks")
}

func (s *MaturityTestSuite) TestEstimate_PoC_ManyFailures() {
	// More than 2 failures even with high score
	result := s.makeResult(7, 0, 3)
	est := Estimate(result)

	s.Equal(PoC, est.Level)
	s.Len(est.Suggestions, 2) // Both "critical checks" and "many failing"
}

func (s *MaturityTestSuite) TestEstimate_ScoreCalculation() {
	testCases := []struct {
		passed   int
		warnings int
		failed   int
		expected float64
	}{
		{10, 0, 0, 100},
		{5, 0, 5, 50},
		{3, 4, 3, 30},
		{0, 5, 5, 0},
		{8, 1, 1, 80},
	}

	for _, tc := range testCases {
		result := s.makeResult(tc.passed, tc.warnings, tc.failed)
		est := Estimate(result)
		s.Equal(tc.expected, est.Score, "passed=%d, warnings=%d, failed=%d",
			tc.passed, tc.warnings, tc.failed)
	}
}

func (s *MaturityTestSuite) TestEstimate_FieldsPopulated() {
	result := s.makeResult(5, 3, 2)
	est := Estimate(result)

	s.Equal(5, est.Passed)
	s.Equal(3, est.Warnings)
	s.Equal(2, est.Failed)
	s.Equal(10, est.Total)
}

func TestMaturityTestSuite(t *testing.T) {
	suite.Run(t, new(MaturityTestSuite))
}
