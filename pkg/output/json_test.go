package output

import (
	// "bytes"
	"encoding/json"
	// "os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/stretchr/testify/suite"
)

// JSONTestSuite is the test suite for the JSON output package.
type JSONTestSuite struct {
	suite.Suite
	// originalStdout *os.File
	// capturedOutput *bytes.Buffer
}

// SetupTest is called before each test method.
func (suite *JSONTestSuite) SetupTest() {
	// Note: We can't easily capture os.Stdout and test os.Exit in unit tests
	// So we'll test the helper functions and JSON structure creation
}

// TestStatusToString_Pass tests that statusToString converts Pass to "pass".
func (suite *JSONTestSuite) TestStatusToString_Pass() {
	result := statusToString(checker.Pass)
	suite.Equal("pass", result)
}

// TestStatusToString_Warn tests that statusToString converts Warn to "warn".
func (suite *JSONTestSuite) TestStatusToString_Warn() {
	result := statusToString(checker.Warn)
	suite.Equal("warn", result)
}

// TestStatusToString_Fail tests that statusToString converts Fail to "fail".
func (suite *JSONTestSuite) TestStatusToString_Fail() {
	result := statusToString(checker.Fail)
	suite.Equal("fail", result)
}

// TestStatusToString_Info tests that statusToString converts Info to "info".
func (suite *JSONTestSuite) TestStatusToString_Info() {
	result := statusToString(checker.Info)
	suite.Equal("info", result)
}

// TestStatusToString_Unknown tests that statusToString handles unknown status.
func (suite *JSONTestSuite) TestStatusToString_Unknown() {
	result := statusToString(checker.Status(99))
	suite.Equal("unknown", result)
}

// TestCalculateScore_EmptyResults tests that calculateScore returns 100.0 for empty results.
func (suite *JSONTestSuite) TestCalculateScore_EmptyResults() {
	result := runner.SuiteResult{
		Results: []checker.Result{},
		Passed:  0,
	}

	score := calculateScore(result)
	suite.Equal(100.0, score)
}

// TestCalculateScore_AllPass tests that calculateScore calculates correctly for all pass.
func (suite *JSONTestSuite) TestCalculateScore_AllPass() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Status: checker.Pass},
			{Status: checker.Pass},
			{Status: checker.Pass},
		},
		Passed: 3,
	}

	score := calculateScore(result)
	suite.Equal(100.0, score)
}

// TestCalculateScore_PartialPass tests that calculateScore calculates percentage correctly.
func (suite *JSONTestSuite) TestCalculateScore_PartialPass() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Status: checker.Pass},
			{Status: checker.Pass},
			{Status: checker.Warn},
			{Status: checker.Fail},
		},
		Passed:   2,
		Warnings: 1,
		Failed:   1,
	}

	score := calculateScore(result)
	suite.Equal(50.0, score) // 2/4 * 100 = 50%
}

// TestCalculateScore_NoPass tests that calculateScore handles zero passed.
func (suite *JSONTestSuite) TestCalculateScore_NoPass() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Status: checker.Warn},
			{Status: checker.Fail},
		},
		Passed:   0,
		Warnings: 1,
		Failed:   1,
	}

	score := calculateScore(result)
	suite.Equal(0.0, score)
}

// TestCalculateScore_MixedResults tests that calculateScore handles mixed results.
func (suite *JSONTestSuite) TestCalculateScore_MixedResults() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Status: checker.Pass},
			{Status: checker.Pass},
			{Status: checker.Pass},
			{Status: checker.Warn},
			{Status: checker.Warn},
		},
		Passed:   3,
		Warnings: 2,
	}

	score := calculateScore(result)
	suite.Equal(60.0, score) // 3/5 * 100 = 60%
}

// TestJSONOutput_Structure tests that JSON creates correct output structure.
func (suite *JSONTestSuite) TestJSONOutput_Structure() {
	// Create a test result
	result := runner.SuiteResult{
		Results: []checker.Result{
			{
				Name:   "Test Check 1",
				ID:     "test1",
				Passed: true,
				Status: checker.Pass,
				Reason: "All good",
			},
			{
				Name:   "Test Check 2",
				ID:     "test2",
				Passed: false,
				Status: checker.Warn,
				Reason: "Warning message",
			},
		},
		Passed:   1,
		Warnings: 1,
		Failed:   0,
		Aborted:  false,
	}

	// We can't easily test JSON() function directly due to os.Stdout and os.Exit
	// But we can test the structure creation logic
	output := JSONOutput{
		Results: make([]JSONResult, 0, len(result.Results)),
		Summary: JSONSummary{
			Total:    result.TotalChecks(),
			Passed:   result.Passed,
			Warnings: result.Warnings,
			Failed:   result.Failed,
			Score:    calculateScore(result),
		},
		Aborted: result.Aborted,
		Success: result.Success(),
	}

	for _, r := range result.Results {
		output.Results = append(output.Results, JSONResult{
			Name:   r.Name,
			ID:     r.ID,
			Passed: r.Passed,
			Status: statusToString(r.Status),
			Reason: r.Reason,
		})
	}

	// Verify structure
	suite.Equal(2, len(output.Results))
	suite.Equal(2, output.Summary.Total)
	suite.Equal(1, output.Summary.Passed)
	suite.Equal(1, output.Summary.Warnings)
	suite.Equal(0, output.Summary.Failed)
	suite.Equal(50.0, output.Summary.Score)
	suite.False(output.Aborted)
	suite.True(output.Success)

	// Verify first result
	suite.Equal("Test Check 1", output.Results[0].Name)
	suite.Equal("test1", output.Results[0].ID)
	suite.True(output.Results[0].Passed)
	suite.Equal("pass", output.Results[0].Status)
	suite.Equal("All good", output.Results[0].Reason)

	// Verify second result
	suite.Equal("Test Check 2", output.Results[1].Name)
	suite.Equal("test2", output.Results[1].ID)
	suite.False(output.Results[1].Passed)
	suite.Equal("warn", output.Results[1].Status)
	suite.Equal("Warning message", output.Results[1].Reason)
}

// TestJSONOutput_EmptyResults tests JSON output with empty results.
func (suite *JSONTestSuite) TestJSONOutput_EmptyResults() {
	result := runner.SuiteResult{
		Results: []checker.Result{},
		Passed:  0,
		Aborted: false,
	}

	output := JSONOutput{
		Results: make([]JSONResult, 0),
		Summary: JSONSummary{
			Total:    result.TotalChecks(),
			Passed:   result.Passed,
			Warnings: result.Warnings,
			Failed:   result.Failed,
			Score:    calculateScore(result),
		},
		Aborted: result.Aborted,
		Success: result.Success(),
	}

	suite.Empty(output.Results)
	suite.Equal(0, output.Summary.Total)
	suite.Equal(100.0, output.Summary.Score)
	suite.True(output.Success)
}

// TestJSONOutput_Aborted tests JSON output with aborted result.
func (suite *JSONTestSuite) TestJSONOutput_Aborted() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{
				Name:   "Critical Check",
				ID:     "critical",
				Passed: false,
				Status: checker.Fail,
				Reason: "Critical failure",
			},
		},
		Passed:  0,
		Failed:  1,
		Aborted: true,
	}

	output := JSONOutput{
		Results: make([]JSONResult, 0, len(result.Results)),
		Summary: JSONSummary{
			Total:    result.TotalChecks(),
			Passed:   result.Passed,
			Warnings: result.Warnings,
			Failed:   result.Failed,
			Score:    calculateScore(result),
		},
		Aborted: result.Aborted,
		Success: result.Success(),
	}

	suite.True(output.Aborted)
	suite.False(output.Success)
	suite.Equal(1, output.Summary.Failed)
	suite.Equal(0.0, output.Summary.Score)
}

// TestJSONOutput_Serialization tests that JSON output can be serialized.
func (suite *JSONTestSuite) TestJSONOutput_Serialization() {
	output := JSONOutput{
		Results: []JSONResult{
			{
				Name:   "Test",
				ID:     "test",
				Passed: true,
				Status: "pass",
				Reason: "Message",
			},
		},
		Summary: JSONSummary{
			Total:    1,
			Passed:   1,
			Warnings: 0,
			Failed:   0,
			Score:    100.0,
		},
		Aborted: false,
		Success: true,
	}

	// Test JSON serialization
	data, err := json.Marshal(output)
	suite.NoError(err)
	suite.NotEmpty(data)

	// Test JSON deserialization
	var decoded JSONOutput
	err = json.Unmarshal(data, &decoded)
	suite.NoError(err)
	suite.Equal(output.Success, decoded.Success)
	suite.Equal(output.Aborted, decoded.Aborted)
	suite.Equal(len(output.Results), len(decoded.Results))
}

// TestJSONTestSuite runs all the tests in the suite.
func TestJSONTestSuite(t *testing.T) {
	suite.Run(t, new(JSONTestSuite))
}
