package runner

import (
	"errors"
	"testing"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// mockChecker is a simple mock implementation of the Checker interface for testing.
type mockChecker struct {
	id     string
	name   string
	result checker.Result
	err    error
}

func (m *mockChecker) ID() string {
	return m.id
}

func (m *mockChecker) Name() string {
	return m.name
}

func (m *mockChecker) Run(path string) (checker.Result, error) {
	return m.result, m.err
}

// slowMockChecker is a mock checker that adds a delay to simulate execution time.
type slowMockChecker struct {
	mockChecker
	delay time.Duration
}

func (m *slowMockChecker) Run(path string) (checker.Result, error) {
	time.Sleep(m.delay)
	return m.result, m.err
}

// panicMockChecker is a mock checker that panics when Run is called.
type panicMockChecker struct {
	id         string
	name       string
	panicValue interface{}
}

func (m *panicMockChecker) ID() string   { return m.id }
func (m *panicMockChecker) Name() string { return m.name }
func (m *panicMockChecker) Run(path string) (checker.Result, error) {
	panic(m.panicValue)
}

// RunnerTestSuite is the test suite for the runner package.
type RunnerTestSuite struct {
	suite.Suite
}

// SetupTest is called before each test method.
func (suite *RunnerTestSuite) SetupTest() {
	// Setup code if needed
}

// TestRunSuite_AllPass tests when all checks pass.
func (suite *RunnerTestSuite) TestRunSuite_AllPass() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
	}

	result := RunSuite("/test/path", checks)

	suite.Equal(2, result.TotalChecks())
	suite.Equal(2, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(0, result.Failed)
	suite.False(result.Aborted)
	suite.True(result.Success())
	suite.Len(result.Results, 2)
}

// TestRunSuite_WithWarnings tests when some checks have warnings.
func (suite *RunnerTestSuite) TestRunSuite_WithWarnings() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Warn,
				Message: "Warning message",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
	}

	result := RunSuite("/test/path", checks)

	suite.Equal(3, result.TotalChecks())
	suite.Equal(2, result.Passed)
	suite.Equal(1, result.Warnings)
	suite.Equal(0, result.Failed)
	suite.False(result.Aborted)
	suite.True(result.Success()) // Warnings don't count as failures
	suite.Len(result.Results, 3)
}

// TestRunSuite_WithFailure tests when a check fails in sequential mode (should abort).
func (suite *RunnerTestSuite) TestRunSuite_WithFailure() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Fail,
				Message: "Critical failure",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "This should not run",
			},
		},
	}

	// Use sequential mode to test abort behavior
	result := RunSuiteSequential("/test/path", checks)

	suite.Equal(2, result.TotalChecks()) // Only 2 checks ran (aborted after check2)
	suite.Equal(1, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)
	suite.False(result.Success())
	suite.Len(result.Results, 2)

	// Verify check3 was not executed
	suite.Equal("check1", result.Results[0].ID)
	suite.Equal("check2", result.Results[1].ID)
}

// TestRunSuite_InternalError tests when a check returns an error in sequential mode.
func (suite *RunnerTestSuite) TestRunSuite_InternalError() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			err:  errors.New("internal error occurred"),
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "This should not run",
			},
		},
	}

	// Use sequential mode to test abort behavior
	result := RunSuiteSequential("/test/path", checks)

	suite.Equal(2, result.TotalChecks()) // Only 2 checks ran (aborted after check2 error)
	suite.Equal(1, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)
	suite.False(result.Success())
	suite.Len(result.Results, 2)

	// Verify the error was converted to a Fail result
	suite.Equal("check1", result.Results[0].ID)
	suite.Equal("check2", result.Results[1].ID)
	suite.False(result.Results[1].Passed)
	suite.Equal(checker.Fail, result.Results[1].Status)
	suite.Contains(result.Results[1].Message, "Internal error")
}

// TestRunSuite_EmptyChecks tests with an empty check list.
func (suite *RunnerTestSuite) TestRunSuite_EmptyChecks() {
	checks := []checker.Checker{}
	result := RunSuite("/test/path", checks)

	suite.Equal(0, result.TotalChecks())
	suite.Equal(0, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(0, result.Failed)
	suite.False(result.Aborted)
	suite.True(result.Success())
	suite.Len(result.Results, 0)
}

// TestRunSuite_MixedResults tests a mix of pass, warn, and fail in sequential mode.
func (suite *RunnerTestSuite) TestRunSuite_MixedResults() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Warn,
				Message: "Warning",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  false,
				Status:  checker.Fail,
				Message: "Failed",
			},
		},
		&mockChecker{
			id:   "check4",
			name: "Check 4",
			result: checker.Result{
				Name:    "Check 4",
				ID:      "check4",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should not run",
			},
		},
	}

	// Use sequential mode to test abort behavior
	result := RunSuiteSequential("/test/path", checks)

	suite.Equal(3, result.TotalChecks()) // Aborted after check3
	suite.Equal(1, result.Passed)
	suite.Equal(1, result.Warnings)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)
	suite.False(result.Success())
}

// TestRunSuite_WarnDoesNotAbort tests that warnings don't cause abortion.
func (suite *RunnerTestSuite) TestRunSuite_WarnDoesNotAbort() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  false,
				Status:  checker.Warn,
				Message: "Warning",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Warn,
				Message: "Another warning",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should run",
			},
		},
	}

	result := RunSuite("/test/path", checks)

	suite.Equal(3, result.TotalChecks()) // All checks should run
	suite.Equal(1, result.Passed)
	suite.Equal(2, result.Warnings)
	suite.Equal(0, result.Failed)
	suite.False(result.Aborted)
	suite.True(result.Success())
}

// TestSuiteResult_TotalChecks tests the TotalChecks method.
func (suite *RunnerTestSuite) TestSuiteResult_TotalChecks() {
	result := SuiteResult{
		Results: []checker.Result{
			{ID: "check1"},
			{ID: "check2"},
			{ID: "check3"},
		},
	}

	suite.Equal(3, result.TotalChecks())
}

// TestSuiteResult_Success tests the Success method.
func (suite *RunnerTestSuite) TestSuiteResult_Success() {
	// Success with no failures
	result1 := SuiteResult{
		Failed: 0,
	}
	suite.True(result1.Success())

	// Success with warnings but no failures
	result2 := SuiteResult{
		Warnings: 2,
		Failed:   0,
	}
	suite.True(result2.Success())

	// Failure when there are failures
	result3 := SuiteResult{
		Failed: 1,
	}
	suite.False(result3.Success())
}

// TestRunSuite_ParallelRunsAllChecks tests that parallel mode runs all checks even with failures.
func (suite *RunnerTestSuite) TestRunSuite_ParallelRunsAllChecks() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Fail,
				Message: "Critical failure",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should run in parallel mode",
			},
		},
	}

	// Parallel mode should run all checks
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(3, result.TotalChecks()) // All 3 checks ran
	suite.Equal(2, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted) // Still marked as aborted due to critical failure
	suite.False(result.Success())

	// Verify all checks ran
	suite.Equal("check1", result.Results[0].ID)
	suite.Equal("check2", result.Results[1].ID)
	suite.Equal("check3", result.Results[2].ID)
}

// TestRunSuite_SequentialStopsOnFailure tests that sequential mode stops on first failure.
func (suite *RunnerTestSuite) TestRunSuite_SequentialStopsOnFailure() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Fail,
				Message: "Critical failure",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should NOT run in sequential mode",
			},
		},
	}

	// Sequential mode should stop on failure
	result := RunSuiteSequential("/test/path", checks)

	suite.Equal(2, result.TotalChecks()) // Only 2 checks ran (aborted after check2)
	suite.Equal(1, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)
	suite.False(result.Success())

	// Verify check3 was not executed
	suite.Len(result.Results, 2)
	suite.Equal("check1", result.Results[0].ID)
	suite.Equal("check2", result.Results[1].ID)
}

// TestRunSuiteWithOptions_Parallel tests RunSuiteWithOptions with parallel=true.
func (suite *RunnerTestSuite) TestRunSuiteWithOptions_Parallel() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Warn,
				Message: "Warning",
			},
		},
	}

	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(2, result.TotalChecks())
	suite.Equal(1, result.Passed)
	suite.Equal(1, result.Warnings)
	suite.False(result.Aborted)
	suite.True(result.Success())
}

// TestRunSuiteWithOptions_Sequential tests RunSuiteWithOptions with parallel=false.
func (suite *RunnerTestSuite) TestRunSuiteWithOptions_Sequential() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  false,
				Status:  checker.Warn,
				Message: "Warning",
			},
		},
	}

	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: false})

	suite.Equal(2, result.TotalChecks())
	suite.Equal(1, result.Passed)
	suite.Equal(1, result.Warnings)
	suite.False(result.Aborted)
	suite.True(result.Success())
}

// TestRunSuite_ParallelEmptyChecks tests parallel mode with empty check list.
func (suite *RunnerTestSuite) TestRunSuite_ParallelEmptyChecks() {
	checks := []checker.Checker{}
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(0, result.TotalChecks())
	suite.Equal(0, result.Passed)
	suite.False(result.Aborted)
	suite.True(result.Success())
}

// TestRunSuite_ParallelWithError tests parallel mode when a check returns an error.
func (suite *RunnerTestSuite) TestRunSuite_ParallelWithError() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "All good",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			err:  errors.New("internal error occurred"),
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should run in parallel mode",
			},
		},
	}

	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(3, result.TotalChecks()) // All 3 checks ran in parallel
	suite.Equal(2, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)
	suite.False(result.Success())

	// Verify the error was converted to a Fail result
	suite.False(result.Results[1].Passed)
	suite.Equal(checker.Fail, result.Results[1].Status)
	suite.Contains(result.Results[1].Message, "Internal error")
}

// TestRunSuite_WithInfo tests that Info status checks are counted separately and don't affect score.
func (suite *RunnerTestSuite) TestRunSuite_WithInfo() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  true,
				Status:  checker.Info,
				Message: "Informational only",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Passed",
			},
		},
	}

	result := RunSuite("/test/path", checks)

	suite.Equal(3, result.TotalChecks())  // All checks ran
	suite.Equal(2, result.ScoredChecks()) // Only 2 count toward score (excluding Info)
	suite.Equal(2, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(0, result.Failed)
	suite.Equal(1, result.Info)
	suite.False(result.Aborted)
	suite.True(result.Success())
}

// TestRunSuite_InfoDoesNotAbort tests that Info status doesn't cause abortion.
func (suite *RunnerTestSuite) TestRunSuite_InfoDoesNotAbort() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Info,
				Message: "Info 1",
			},
		},
		&mockChecker{
			id:   "check2",
			name: "Check 2",
			result: checker.Result{
				Name:    "Check 2",
				ID:      "check2",
				Passed:  true,
				Status:  checker.Info,
				Message: "Info 2",
			},
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should run",
			},
		},
	}

	result := RunSuiteSequential("/test/path", checks)

	suite.Equal(3, result.TotalChecks())  // All checks should run
	suite.Equal(1, result.ScoredChecks()) // Only 1 counts toward score
	suite.Equal(1, result.Passed)
	suite.Equal(0, result.Warnings)
	suite.Equal(0, result.Failed)
	suite.Equal(2, result.Info)
	suite.False(result.Aborted)
	suite.True(result.Success())
}

// TestSuiteResult_ScoredChecks tests the ScoredChecks method.
func (suite *RunnerTestSuite) TestSuiteResult_ScoredChecks() {
	result := SuiteResult{
		Results: []checker.Result{
			{ID: "check1", Status: checker.Pass},
			{ID: "check2", Status: checker.Warn},
			{ID: "check3", Status: checker.Fail},
			{ID: "check4", Status: checker.Info},
			{ID: "check5", Status: checker.Info},
		},
		Passed:   1,
		Warnings: 1,
		Failed:   1,
		Info:     2,
	}

	suite.Equal(5, result.TotalChecks())
	suite.Equal(3, result.ScoredChecks()) // Excludes Info
}

// TestRunSuite_DurationIsSet tests that Duration is set on individual check results.
func (suite *RunnerTestSuite) TestRunSuite_DurationIsSet() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 10 * time.Millisecond,
		},
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check2",
				name: "Check 2",
				result: checker.Result{
					Name:    "Check 2",
					ID:      "check2",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 20 * time.Millisecond,
		},
	}

	result := RunSuite("/test/path", checks)

	suite.Equal(2, result.TotalChecks())
	// Verify Duration is set on each result
	for i, r := range result.Results {
		suite.Greater(r.Duration.Nanoseconds(), int64(0), "Result %d should have non-zero duration", i)
	}
	// Verify durations are reasonable (at least as long as the delay)
	suite.GreaterOrEqual(result.Results[0].Duration, 10*time.Millisecond)
	suite.GreaterOrEqual(result.Results[1].Duration, 20*time.Millisecond)
}

// TestRunSuite_TotalDurationIsSet tests that TotalDuration is set on SuiteResult.
func (suite *RunnerTestSuite) TestRunSuite_TotalDurationIsSet() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 15 * time.Millisecond,
		},
	}

	result := RunSuite("/test/path", checks)

	// Verify TotalDuration is set and reasonable
	suite.Greater(result.TotalDuration.Nanoseconds(), int64(0), "TotalDuration should be non-zero")
	suite.GreaterOrEqual(result.TotalDuration, 15*time.Millisecond)
}

// TestRunSuite_SequentialTiming tests timing in sequential mode.
func (suite *RunnerTestSuite) TestRunSuite_SequentialTiming() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 10 * time.Millisecond,
		},
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check2",
				name: "Check 2",
				result: checker.Result{
					Name:    "Check 2",
					ID:      "check2",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 10 * time.Millisecond,
		},
	}

	result := RunSuiteSequential("/test/path", checks)

	// In sequential mode, TotalDuration should be at least the sum of delays
	suite.GreaterOrEqual(result.TotalDuration, 20*time.Millisecond)
	// Individual durations should be set
	suite.GreaterOrEqual(result.Results[0].Duration, 10*time.Millisecond)
	suite.GreaterOrEqual(result.Results[1].Duration, 10*time.Millisecond)
}

// TestRunSuite_ParallelTiming tests timing in parallel mode.
func (suite *RunnerTestSuite) TestRunSuite_ParallelTiming() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 20 * time.Millisecond,
		},
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check2",
				name: "Check 2",
				result: checker.Result{
					Name:    "Check 2",
					ID:      "check2",
					Passed:  true,
					Status:  checker.Pass,
					Message: "All good",
				},
			},
			delay: 20 * time.Millisecond,
		},
	}

	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	// In parallel mode, TotalDuration should be less than the sum of delays
	// (since checks run concurrently) but at least as long as the longest check
	suite.GreaterOrEqual(result.TotalDuration, 20*time.Millisecond)
	// Should be significantly less than 40ms (sequential would take 40ms+)
	suite.Less(result.TotalDuration, 35*time.Millisecond)
	// Individual durations should still be set
	suite.GreaterOrEqual(result.Results[0].Duration, 20*time.Millisecond)
	suite.GreaterOrEqual(result.Results[1].Duration, 20*time.Millisecond)
}

// TestRunSuite_EmptyChecksTiming tests timing with empty check list.
func (suite *RunnerTestSuite) TestRunSuite_EmptyChecksTiming() {
	checks := []checker.Checker{}

	result := RunSuite("/test/path", checks)

	// TotalDuration should be very small but not negative
	suite.GreaterOrEqual(result.TotalDuration.Nanoseconds(), int64(0))
}

// TestRunSuite_TimeoutNotExceeded tests that checks complete normally when timeout is not exceeded.
func (suite *RunnerTestSuite) TestRunSuite_TimeoutNotExceeded() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "Completed in time",
				},
			},
			delay: 10 * time.Millisecond,
		},
	}

	// Timeout is 100ms, check takes 10ms - should succeed
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{
		Parallel: true,
		Timeout:  100 * time.Millisecond,
	})

	suite.Equal(1, result.TotalChecks())
	suite.Equal(1, result.Passed)
	suite.Equal(0, result.Failed)
	suite.True(result.Success())
	suite.Equal("Completed in time", result.Results[0].Message)
}

// TestRunSuite_TimeoutExceeded tests that checks fail when timeout is exceeded.
func (suite *RunnerTestSuite) TestRunSuite_TimeoutExceeded() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "Should not see this",
				},
			},
			delay: 100 * time.Millisecond, // Takes 100ms
		},
	}

	// Timeout is 20ms, check takes 100ms - should timeout
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{
		Parallel: true,
		Timeout:  20 * time.Millisecond,
	})

	suite.Equal(1, result.TotalChecks())
	suite.Equal(0, result.Passed)
	suite.Equal(1, result.Failed)
	suite.False(result.Success())
	suite.Contains(result.Results[0].Message, "timed out")
	suite.Equal("check1", result.Results[0].ID)
	suite.Equal("Check 1", result.Results[0].Name)
}

// TestRunSuite_TimeoutSequential tests timeout in sequential mode.
func (suite *RunnerTestSuite) TestRunSuite_TimeoutSequential() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "Fast check",
				},
			},
			delay: 5 * time.Millisecond,
		},
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check2",
				name: "Check 2",
				result: checker.Result{
					Name:    "Check 2",
					ID:      "check2",
					Passed:  true,
					Status:  checker.Pass,
					Message: "Slow check",
				},
			},
			delay: 100 * time.Millisecond, // Too slow
		},
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check3",
				name: "Check 3",
				result: checker.Result{
					Name:    "Check 3",
					ID:      "check3",
					Passed:  true,
					Status:  checker.Pass,
					Message: "Should not run",
				},
			},
			delay: 5 * time.Millisecond,
		},
	}

	// Timeout is 20ms per check
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{
		Parallel: false,
		Timeout:  20 * time.Millisecond,
	})

	// First check passes, second times out (causing abort in sequential mode)
	suite.Equal(2, result.TotalChecks())
	suite.Equal(1, result.Passed)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted) // Should abort because timeout produces a Fail result
	suite.Contains(result.Results[1].Message, "timed out")
}

// TestRunSuite_TimeoutZeroMeansNoTimeout tests that timeout=0 means no timeout.
func (suite *RunnerTestSuite) TestRunSuite_TimeoutZeroMeansNoTimeout() {
	checks := []checker.Checker{
		&slowMockChecker{
			mockChecker: mockChecker{
				id:   "check1",
				name: "Check 1",
				result: checker.Result{
					Name:    "Check 1",
					ID:      "check1",
					Passed:  true,
					Status:  checker.Pass,
					Message: "Slow but should complete",
				},
			},
			delay: 50 * time.Millisecond,
		},
	}

	// Timeout is 0 (no timeout)
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{
		Parallel: true,
		Timeout:  0,
	})

	suite.Equal(1, result.TotalChecks())
	suite.Equal(1, result.Passed)
	suite.True(result.Success())
	suite.Equal("Slow but should complete", result.Results[0].Message)
}

// TestRunSuite_PanicRecoveryParallel tests that panics are recovered in parallel mode.
func (suite *RunnerTestSuite) TestRunSuite_PanicRecoveryParallel() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should complete",
			},
		},
		&panicMockChecker{
			id:         "check2",
			name:       "Panicking Check",
			panicValue: "something went wrong",
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should also complete",
			},
		},
	}

	// Parallel mode should recover from panic and continue
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(3, result.TotalChecks()) // All 3 checks should have results
	suite.Equal(2, result.Passed)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)
	suite.False(result.Success())

	// Verify the panicking check was converted to a Fail result
	suite.Equal("check2", result.Results[1].ID)
	suite.Equal("Panicking Check", result.Results[1].Name)
	suite.False(result.Results[1].Passed)
	suite.Equal(checker.Fail, result.Results[1].Status)
	suite.Contains(result.Results[1].Message, "panicked")
	suite.Contains(result.Results[1].Message, "something went wrong")

	// Verify other checks completed normally
	suite.True(result.Results[0].Passed)
	suite.True(result.Results[2].Passed)
}

// TestRunSuite_PanicRecoveryWithError tests panic with error value.
func (suite *RunnerTestSuite) TestRunSuite_PanicRecoveryWithError() {
	checks := []checker.Checker{
		&panicMockChecker{
			id:         "check1",
			name:       "Error Panicking Check",
			panicValue: errors.New("error panic value"),
		},
	}

	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(1, result.TotalChecks())
	suite.Equal(0, result.Passed)
	suite.Equal(1, result.Failed)
	suite.Contains(result.Results[0].Message, "panicked")
	suite.Contains(result.Results[0].Message, "error panic value")
}

// TestRunSuite_PanicRecoveryWithUnknownValue tests panic with non-string/non-error value.
func (suite *RunnerTestSuite) TestRunSuite_PanicRecoveryWithUnknownValue() {
	checks := []checker.Checker{
		&panicMockChecker{
			id:         "check1",
			name:       "Unknown Panic Check",
			panicValue: 12345, // integer panic value
		},
	}

	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{Parallel: true})

	suite.Equal(1, result.TotalChecks())
	suite.Equal(0, result.Passed)
	suite.Equal(1, result.Failed)
	suite.Contains(result.Results[0].Message, "panicked unexpectedly")
}

// TestRunSuite_PanicRecoveryWithTimeout tests panic recovery when timeout is set.
func (suite *RunnerTestSuite) TestRunSuite_PanicRecoveryWithTimeout() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should complete",
			},
		},
		&panicMockChecker{
			id:         "check2",
			name:       "Panicking Check",
			panicValue: "timeout test panic",
		},
	}

	// Panic recovery should work with timeout enabled
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{
		Parallel: true,
		Timeout:  1 * time.Second,
	})

	suite.Equal(2, result.TotalChecks())
	suite.Equal(1, result.Passed)
	suite.Equal(1, result.Failed)

	// Verify the panicking check was recovered
	suite.Contains(result.Results[1].Message, "panicked")
	suite.Contains(result.Results[1].Message, "timeout test panic")
}

// TestRunSuite_PanicRecoverySequential tests panic recovery in sequential mode.
func (suite *RunnerTestSuite) TestRunSuite_PanicRecoverySequential() {
	checks := []checker.Checker{
		&mockChecker{
			id:   "check1",
			name: "Check 1",
			result: checker.Result{
				Name:    "Check 1",
				ID:      "check1",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should complete first",
			},
		},
		&panicMockChecker{
			id:         "check2",
			name:       "Panicking Check",
			panicValue: "sequential panic",
		},
		&mockChecker{
			id:   "check3",
			name: "Check 3",
			result: checker.Result{
				Name:    "Check 3",
				ID:      "check3",
				Passed:  true,
				Status:  checker.Pass,
				Message: "Should not run (aborted)",
			},
		},
	}

	// Sequential mode should recover from panic but abort
	result := RunSuiteWithOptions("/test/path", checks, RunSuiteOptions{
		Parallel: false,
		Timeout:  1 * time.Second,
	})

	// Panic recovery converts to Fail, which aborts in sequential mode
	suite.Equal(2, result.TotalChecks()) // Only 2 ran (aborted after panic)
	suite.Equal(1, result.Passed)
	suite.Equal(1, result.Failed)
	suite.True(result.Aborted)

	suite.Contains(result.Results[1].Message, "panicked")
}

// TestFormatPanicMessage tests the formatPanicMessage helper function.
func (suite *RunnerTestSuite) TestFormatPanicMessage() {
	// Test with string value
	suite.Equal("Check panicked: test string", formatPanicMessage("test string"))

	// Test with error value
	suite.Equal("Check panicked: test error", formatPanicMessage(errors.New("test error")))

	// Test with other value
	suite.Equal("Check panicked unexpectedly", formatPanicMessage(123))
	suite.Equal("Check panicked unexpectedly", formatPanicMessage(nil))
	suite.Equal("Check panicked unexpectedly", formatPanicMessage(struct{}{}))
}

// TestRunnerTestSuite runs all the tests in the suite.
func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerTestSuite))
}
