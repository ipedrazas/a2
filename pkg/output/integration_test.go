package output

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite tests output formatters end-to-end.
type IntegrationTestSuite struct {
	suite.Suite
}

// createTestSuiteResult creates a consistent SuiteResult for testing.
func createTestSuiteResult() runner.SuiteResult {
	return runner.SuiteResult{
		Results: []checker.Result{
			{
				Name:     "Go Build",
				ID:       "go:build",
				Passed:   true,
				Status:   checker.Pass,
				Message:  "Build successful",
				Language: checker.LangGo,
				Duration: 150 * time.Millisecond,
			},
			{
				Name:     "Go Tests",
				ID:       "go:tests",
				Passed:   true,
				Status:   checker.Pass,
				Message:  "All tests passed",
				Language: checker.LangGo,
				Duration: 2500 * time.Millisecond,
			},
			{
				Name:     "Go Format",
				ID:       "go:format",
				Passed:   false,
				Status:   checker.Warn,
				Message:  "3 files need formatting",
				Language: checker.LangGo,
				Duration: 50 * time.Millisecond,
			},
			{
				Name:     "Go Coverage",
				ID:       "go:coverage",
				Passed:   false,
				Status:   checker.Fail,
				Message:  "Coverage 45% below threshold 80%",
				Language: checker.LangGo,
				Duration: 3000 * time.Millisecond,
			},
			{
				Name:     "Version Info",
				ID:       "common:version",
				Passed:   true,
				Status:   checker.Info,
				Message:  "v1.0.0",
				Language: checker.LangCommon,
				Duration: 5 * time.Millisecond,
			},
		},
		Passed:        2,
		Warnings:      1,
		Failed:        1,
		Info:          1,
		Aborted:       false,
		TotalDuration: 5705 * time.Millisecond,
	}
}

// createTestDetectionResult creates a consistent DetectionResult for testing.
func createTestDetectionResult() language.DetectionResult {
	return language.DetectionResult{
		Languages: []checker.Language{checker.LangGo},
	}
}

// TestJSON_IntegrationOutput tests the JSON formatter produces valid JSON.
func (suite *IntegrationTestSuite) TestJSON_IntegrationOutput() {
	result := createTestSuiteResult()
	detected := createTestDetectionResult()

	output := captureStdout(func() {
		success, err := JSON(result, detected, VerbosityNormal)
		suite.NoError(err)
		suite.False(success) // Has failures
	})

	// Verify it's valid JSON
	var jsonOutput JSONOutput
	err := json.Unmarshal([]byte(output), &jsonOutput)
	suite.NoError(err, "JSON output should be valid JSON")

	// Verify structure
	suite.Equal([]string{"go"}, jsonOutput.Languages)
	suite.Len(jsonOutput.Results, 5)
	suite.Equal(4, jsonOutput.Summary.Total) // Excludes Info
	suite.Equal(2, jsonOutput.Summary.Passed)
	suite.Equal(1, jsonOutput.Summary.Warnings)
	suite.Equal(1, jsonOutput.Summary.Failed)
	suite.Equal(1, jsonOutput.Summary.Info)
	suite.Equal(50.0, jsonOutput.Summary.Score) // 2/4 = 50%
	suite.False(jsonOutput.Aborted)
	suite.False(jsonOutput.Success)

	// Verify results
	suite.Equal("go:build", jsonOutput.Results[0].ID)
	suite.Equal("pass", jsonOutput.Results[0].Status)
	suite.Equal("go:coverage", jsonOutput.Results[3].ID)
	suite.Equal("fail", jsonOutput.Results[3].Status)
	suite.Equal("common:version", jsonOutput.Results[4].ID)
	suite.Equal("info", jsonOutput.Results[4].Status)

	// Verify durations are captured
	suite.Equal(int64(150), jsonOutput.Results[0].DurationMs)
	suite.Equal(int64(2500), jsonOutput.Results[1].DurationMs)
	suite.Greater(jsonOutput.Summary.TotalDurationMs, int64(0))

	// Verify maturity
	suite.NotEmpty(jsonOutput.Maturity.Level)
	suite.NotEmpty(jsonOutput.Maturity.Description)
}

// TestTOON_IntegrationOutput tests the TOON formatter produces valid TOON.
func (suite *IntegrationTestSuite) TestTOON_IntegrationOutput() {
	result := createTestSuiteResult()
	detected := createTestDetectionResult()

	output := captureStdout(func() {
		success, err := TOON(result, detected, VerbosityNormal)
		suite.NoError(err)
		suite.False(success) // Has failures
	})

	// Verify TOON structure
	suite.Contains(output, "languages[1]: go")
	suite.Contains(output, "results[5]{name,id,passed,status,message,language,duration_ms}:")
	suite.Contains(output, "summary:")
	suite.Contains(output, "total: 4")
	suite.Contains(output, "passed: 2")
	suite.Contains(output, "warnings: 1")
	suite.Contains(output, "failed: 1")
	suite.Contains(output, "info: 1")
	suite.Contains(output, "score: 50")
	suite.Contains(output, "maturity:")
	suite.Contains(output, "aborted: false")
	suite.Contains(output, "success: false")

	// Verify result rows (note: IDs with colons are quoted in TOON format)
	suite.Contains(output, "Go Build,\"go:build\",true,pass,Build successful,go,150")
	suite.Contains(output, "Go Tests,\"go:tests\",true,pass,All tests passed,go,2500")
}

// TestPretty_IntegrationOutput tests the Pretty formatter produces readable output.
func (suite *IntegrationTestSuite) TestPretty_IntegrationOutput() {
	result := createTestSuiteResult()
	detected := createTestDetectionResult()

	output := captureStdout(func() {
		success, err := Pretty(result, path, detected, VerbosityNormal)
		suite.NoError(err)
		suite.False(success) // Has failures
	})

	// Verify Pretty output contains expected elements
	suite.Contains(output, "Go Build")
	suite.Contains(output, "Go Tests")
	suite.Contains(output, "Go Format")
	suite.Contains(output, "Go Coverage")
	suite.Contains(output, "Version Info")

	// Verify summary information
	suite.Contains(output, "2") // Passed
	suite.Contains(output, "1") // Warnings
	suite.Contains(output, "1") // Failed

	// Verify it contains timing info
	suite.Contains(output, "ms")
}

// TestJSON_EmptyResults tests JSON output with empty results.
func (suite *IntegrationTestSuite) TestJSON_EmptyResults() {
	result := runner.SuiteResult{
		Results:       []checker.Result{},
		TotalDuration: 10 * time.Millisecond,
	}
	detected := language.DetectionResult{Languages: []checker.Language{}}

	output := captureStdout(func() {
		success, err := JSON(result, detected, VerbosityNormal)
		suite.NoError(err)
		suite.True(success) // No failures
	})

	var jsonOutput JSONOutput
	err := json.Unmarshal([]byte(output), &jsonOutput)
	suite.NoError(err)

	suite.Empty(jsonOutput.Results)
	suite.Empty(jsonOutput.Languages)
	suite.Equal(100.0, jsonOutput.Summary.Score)
	suite.True(jsonOutput.Success)
}

// TestTOON_EmptyResults tests TOON output with empty results.
func (suite *IntegrationTestSuite) TestTOON_EmptyResults() {
	result := runner.SuiteResult{
		Results:       []checker.Result{},
		TotalDuration: 10 * time.Millisecond,
	}
	detected := language.DetectionResult{Languages: []checker.Language{}}

	output := captureStdout(func() {
		success, err := TOON(result, detected, VerbosityNormal)
		suite.NoError(err)
		suite.True(success)
	})

	suite.Contains(output, "languages[0]:")
	suite.Contains(output, "results[0]{name,id,passed,status,message,language,duration_ms}:")
	suite.Contains(output, "score: 100")
	suite.Contains(output, "success: true")
}

// TestJSON_AllPass tests JSON output when all checks pass.
func (suite *IntegrationTestSuite) TestJSON_AllPass() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{
				Name:     "Check 1",
				ID:       "check1",
				Passed:   true,
				Status:   checker.Pass,
				Message:  "OK",
				Language: checker.LangGo,
				Duration: 100 * time.Millisecond,
			},
			{
				Name:     "Check 2",
				ID:       "check2",
				Passed:   true,
				Status:   checker.Pass,
				Message:  "OK",
				Language: checker.LangGo,
				Duration: 100 * time.Millisecond,
			},
		},
		Passed:        2,
		TotalDuration: 200 * time.Millisecond,
	}
	detected := createTestDetectionResult()

	output := captureStdout(func() {
		success, err := JSON(result, detected, VerbosityNormal)
		suite.NoError(err)
		suite.True(success)
	})

	var jsonOutput JSONOutput
	err := json.Unmarshal([]byte(output), &jsonOutput)
	suite.NoError(err)

	suite.Equal(100.0, jsonOutput.Summary.Score)
	suite.True(jsonOutput.Success)
	suite.False(jsonOutput.Aborted)
}

// TestOutputConsistency_ScoreCalculation verifies score is consistent across formats.
func (suite *IntegrationTestSuite) TestOutputConsistency_ScoreCalculation() {
	result := createTestSuiteResult()
	detected := createTestDetectionResult()

	// Get JSON output
	jsonStr := captureStdout(func() {
		_, err := JSON(result, detected, VerbosityNormal)
		suite.NoError(err)
	})
	var jsonOutput JSONOutput
	err := json.Unmarshal([]byte(jsonStr), &jsonOutput)
	suite.NoError(err)

	// Get TOON output
	toonStr := captureStdout(func() {
		_, err := TOON(result, detected, VerbosityNormal)
		suite.NoError(err)
	})

	// Both should have the same score
	expectedScore := 50.0 // 2 passed / 4 total (excluding Info)
	suite.Equal(expectedScore, jsonOutput.Summary.Score)
	suite.Contains(toonStr, "score: 50")
}

// TestOutputConsistency_SuccessFlag verifies success flag is consistent.
func (suite *IntegrationTestSuite) TestOutputConsistency_SuccessFlag() {
	// Test with failures
	resultFail := createTestSuiteResult()
	detected := createTestDetectionResult()

	jsonFail := captureStdout(func() {
		success, err := JSON(resultFail, detected, VerbosityNormal)
		suite.NoError(err)
		suite.False(success)
	})
	var jsonOutputFail JSONOutput
	err := json.Unmarshal([]byte(jsonFail), &jsonOutputFail)
	suite.NoError(err)
	suite.False(jsonOutputFail.Success)

	toonFail := captureStdout(func() {
		success, err := TOON(resultFail, detected, VerbosityNormal)
		suite.NoError(err)
		suite.False(success)
	})
	suite.Contains(toonFail, "success: false")

	// Test without failures
	resultPass := runner.SuiteResult{
		Results: []checker.Result{
			{Name: "Test", ID: "test", Passed: true, Status: checker.Pass, Message: "OK"},
		},
		Passed: 1,
	}

	jsonPass := captureStdout(func() {
		success, err := JSON(resultPass, detected, VerbosityNormal)
		suite.NoError(err)
		suite.True(success)
	})
	var jsonOutputPass JSONOutput
	err = json.Unmarshal([]byte(jsonPass), &jsonOutputPass)
	suite.NoError(err)
	suite.True(jsonOutputPass.Success)

	toonPass := captureStdout(func() {
		success, err := TOON(resultPass, detected, VerbosityNormal)
		suite.NoError(err)
		suite.True(success)
	})
	suite.Contains(toonPass, "success: true")
}

// TestJSON_SpecialCharactersInMessage tests that special characters are handled.
func (suite *IntegrationTestSuite) TestJSON_SpecialCharactersInMessage() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{
				Name:    "Test",
				ID:      "test",
				Passed:  false,
				Status:  checker.Fail,
				Message: "Error: \"file not found\" at path/to/file",
			},
		},
		Failed: 1,
	}
	detected := createTestDetectionResult()

	output := captureStdout(func() {
		_, err := JSON(result, detected, VerbosityNormal)
		suite.NoError(err)
	})

	// Should be valid JSON even with quotes in message
	var jsonOutput JSONOutput
	err := json.Unmarshal([]byte(output), &jsonOutput)
	suite.NoError(err)
	suite.Contains(jsonOutput.Results[0].Message, "file not found")
}

// TestTOON_SpecialCharactersInMessage tests TOON handles special characters.
func (suite *IntegrationTestSuite) TestTOON_SpecialCharactersInMessage() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{
				Name:    "Test",
				ID:      "test",
				Passed:  false,
				Status:  checker.Fail,
				Message: "Error: has,comma and \"quotes\"",
			},
		},
		Failed: 1,
	}
	detected := createTestDetectionResult()

	output := captureStdout(func() {
		_, err := TOON(result, detected, VerbosityNormal)
		suite.NoError(err)
	})

	// Message with comma should be quoted in TOON
	suite.Contains(output, `"Error: has,comma and \"quotes\""`)
}

// TestJSON_MultiplLanguages tests JSON output with multiple languages.
func (suite *IntegrationTestSuite) TestJSON_MultipleLanguages() {
	result := runner.SuiteResult{
		Results: []checker.Result{
			{Name: "Go Test", ID: "go:test", Passed: true, Status: checker.Pass, Language: checker.LangGo},
			{Name: "Python Test", ID: "python:test", Passed: true, Status: checker.Pass, Language: checker.LangPython},
		},
		Passed: 2,
	}
	detected := language.DetectionResult{
		Languages: []checker.Language{checker.LangGo, checker.LangPython},
	}

	output := captureStdout(func() {
		_, err := JSON(result, detected, VerbosityNormal)
		suite.NoError(err)
	})

	var jsonOutput JSONOutput
	err := json.Unmarshal([]byte(output), &jsonOutput)
	suite.NoError(err)
	suite.ElementsMatch([]string{"go", "python"}, jsonOutput.Languages)
}

// path is used for Pretty formatter tests
var path = "."

// TestIntegrationTestSuite runs all integration tests.
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
