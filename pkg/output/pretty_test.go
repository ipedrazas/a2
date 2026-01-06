package output

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/stretchr/testify/suite"
)

// PrettyTestSuite is the test suite for the pretty output package.
type PrettyTestSuite struct {
	suite.Suite
}

// captureStdout captures stdout during function execution.
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// TestPrintResult_Pass tests printResult with a passing check.
func (suite *PrettyTestSuite) TestPrintResult_Pass() {
	result := checker.Result{
		Name:    "Test Check",
		ID:      "test",
		Passed:  true,
		Status:  checker.Pass,
		Message: "All good",
	}

	output := captureStdout(func() {
		printResult(result)
	})

	suite.Contains(output, "PASS")
	suite.Contains(output, "Test Check")
	suite.Contains(output, "All good")
}

// TestPrintResult_Warn tests printResult with a warning check.
func (suite *PrettyTestSuite) TestPrintResult_Warn() {
	result := checker.Result{
		Name:    "Warning Check",
		ID:      "warn",
		Passed:  false,
		Status:  checker.Warn,
		Message: "Something needs attention",
	}

	output := captureStdout(func() {
		printResult(result)
	})

	suite.Contains(output, "WARN")
	suite.Contains(output, "Warning Check")
	suite.Contains(output, "Something needs attention")
}

// TestPrintResult_Fail tests printResult with a failing check.
func (suite *PrettyTestSuite) TestPrintResult_Fail() {
	result := checker.Result{
		Name:    "Failed Check",
		ID:      "fail",
		Passed:  false,
		Status:  checker.Fail,
		Message: "Critical error",
	}

	output := captureStdout(func() {
		printResult(result)
	})

	suite.Contains(output, "FAIL")
	suite.Contains(output, "Failed Check")
	suite.Contains(output, "Critical error")
}

// TestPrintResult_NoMessage tests printResult with no message.
func (suite *PrettyTestSuite) TestPrintResult_NoMessage() {
	result := checker.Result{
		Name:   "Check Without Message",
		ID:     "nomsg",
		Passed: true,
		Status: checker.Pass,
	}

	output := captureStdout(func() {
		printResult(result)
	})

	suite.Contains(output, "PASS")
	suite.Contains(output, "Check Without Message")
}

// TestPrintStatus_AllPassed tests printStatus when all checks pass.
func (suite *PrettyTestSuite) TestPrintStatus_AllPassed() {
	result := runner.SuiteResult{
		Results:  []checker.Result{},
		Passed:   3,
		Warnings: 0,
		Failed:   0,
	}

	output := captureStdout(func() {
		printStatus(result)
	})

	suite.Contains(output, "ALL CHECKS PASSED")
}

// TestPrintStatus_HasWarnings tests printStatus when there are warnings.
func (suite *PrettyTestSuite) TestPrintStatus_HasWarnings() {
	result := runner.SuiteResult{
		Results:  []checker.Result{},
		Passed:   2,
		Warnings: 1,
		Failed:   0,
	}

	output := captureStdout(func() {
		printStatus(result)
	})

	suite.Contains(output, "NEEDS ATTENTION")
}

// TestPrintStatus_HasFailures tests printStatus when there are failures.
func (suite *PrettyTestSuite) TestPrintStatus_HasFailures() {
	result := runner.SuiteResult{
		Results:  []checker.Result{},
		Passed:   1,
		Warnings: 0,
		Failed:   1,
		Aborted:  false,
	}

	output := captureStdout(func() {
		printStatus(result)
	})

	suite.Contains(output, "FAILED")
}

// TestPrintStatus_Aborted tests printStatus when execution was aborted.
func (suite *PrettyTestSuite) TestPrintStatus_Aborted() {
	result := runner.SuiteResult{
		Results:  []checker.Result{},
		Passed:   0,
		Warnings: 0,
		Failed:   1,
		Aborted:  true,
	}

	output := captureStdout(func() {
		printStatus(result)
	})

	suite.Contains(output, "CRITICAL FAILURE")
	suite.Contains(output, "Aborted")
}

// TestPrintScore_AllPassed tests printScore with 100% pass rate.
func (suite *PrettyTestSuite) TestPrintScore_AllPassed() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Status: checker.Pass},
			{Status: checker.Pass},
			{Status: checker.Pass},
		},
		Passed: 3,
	}

	output := captureStdout(func() {
		printScore(result)
	})

	suite.Contains(output, "3/3")
	suite.Contains(output, "100%")
}

// TestPrintScore_Partial tests printScore with partial pass rate.
func (suite *PrettyTestSuite) TestPrintScore_Partial() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Status: checker.Pass},
			{Status: checker.Warn},
			{Status: checker.Pass},
			{Status: checker.Fail},
		},
		Passed: 2,
	}

	output := captureStdout(func() {
		printScore(result)
	})

	suite.Contains(output, "2/4")
	suite.Contains(output, "50%")
}

// TestPrintScore_Empty tests printScore with empty results.
func (suite *PrettyTestSuite) TestPrintScore_Empty() {
	result := runner.SuiteResult{
		Results: []checker.Result{},
		Passed:  0,
	}

	output := captureStdout(func() {
		printScore(result)
	})

	suite.Contains(output, "0/0")
}

// TestPrintRecommendations_Coverage tests recommendations for coverage failure.
func (suite *PrettyTestSuite) TestPrintRecommendations_Coverage() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:coverage", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "Add more tests")
}

// TestPrintRecommendations_Gofmt tests recommendations for gofmt failure.
func (suite *PrettyTestSuite) TestPrintRecommendations_Gofmt() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:format", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "gofmt -w")
}

// TestPrintRecommendations_GoVet tests recommendations for govet failure.
func (suite *PrettyTestSuite) TestPrintRecommendations_GoVet() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:vet", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "go vet")
}

// TestPrintRecommendations_FileExists tests recommendations for missing files.
func (suite *PrettyTestSuite) TestPrintRecommendations_FileExists() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "file_exists", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "documentation files")
}

// TestPrintRecommendations_Tests tests recommendations for test failures.
func (suite *PrettyTestSuite) TestPrintRecommendations_Tests() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:tests", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "Fix failing tests")
}

// TestPrintRecommendations_Build tests recommendations for build failures.
func (suite *PrettyTestSuite) TestPrintRecommendations_Build() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:build", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "Fix build errors")
}

// TestPrintRecommendations_Deps tests recommendations for vulnerability findings.
func (suite *PrettyTestSuite) TestPrintRecommendations_Deps() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:deps", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "Update dependencies")
}

// TestPrintRecommendations_NoFailures tests no recommendations when all pass.
func (suite *PrettyTestSuite) TestPrintRecommendations_NoFailures() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:coverage", Passed: true},
			{ID: "go:format", Passed: true},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	// Should not contain "Recommendations" when all pass
	suite.NotContains(output, "Recommendations")
}

// TestPrintRecommendations_Multiple tests multiple recommendations.
func (suite *PrettyTestSuite) TestPrintRecommendations_Multiple() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "go:coverage", Passed: false},
			{ID: "go:format", Passed: false},
			{ID: "go:tests", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	suite.Contains(output, "Recommendations")
	suite.Contains(output, "Add more tests")
	suite.Contains(output, "gofmt")
	suite.Contains(output, "Fix failing tests")
}

// TestPrintRecommendations_UnknownCheck tests unknown check ID doesn't produce recommendation.
func (suite *PrettyTestSuite) TestPrintRecommendations_UnknownCheck() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{ID: "unknown_check", Passed: false},
		},
	}

	output := captureStdout(func() {
		printRecommendations(result)
	})

	// Unknown check IDs don't generate recommendations
	suite.NotContains(output, "Recommendations")
}

// TestPrettyTestSuite runs all the tests in the suite.
func TestPrettyTestSuite(t *testing.T) {
	suite.Run(t, new(PrettyTestSuite))
}
