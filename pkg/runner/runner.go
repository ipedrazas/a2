package runner

import (
	"context"
	"sync"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
)

// SuiteResult holds the results of running a suite of checks.
type SuiteResult struct {
	Results       []checker.Result // All check results
	Aborted       bool             // True if a Fail (veto) check stopped execution
	Passed        int              // Number of passed checks
	Warnings      int              // Number of warning checks
	Failed        int              // Number of failed checks
	Info          int              // Number of informational checks (excluded from score)
	TotalDuration time.Duration    // Total time taken to run all checks
}

// RunSuiteOptions configures how the suite is executed.
type RunSuiteOptions struct {
	Parallel bool          // Run checks in parallel (default: true)
	Timeout  time.Duration // Timeout for each individual check (0 = no timeout)
}

// RunSuite executes a suite of checks against the given path.
// Uses parallel execution by default for better performance.
// If a check returns Fail status, Aborted is set to true.
func RunSuite(path string, checks []checker.Checker) SuiteResult {
	return RunSuiteWithOptions(path, checks, RunSuiteOptions{Parallel: true})
}

// RunSuiteSequential executes checks sequentially, stopping on first critical failure.
// Use this when you need veto power behavior or have limited CPU resources.
func RunSuiteSequential(path string, checks []checker.Checker) SuiteResult {
	return RunSuiteWithOptions(path, checks, RunSuiteOptions{Parallel: false})
}

// RunSuiteWithOptions executes a suite of checks with the given options.
// In parallel mode: all checks run concurrently, Aborted is set if any critical check fails.
// In sequential mode: stops immediately on first critical failure (veto power).
func RunSuiteWithOptions(path string, checks []checker.Checker, opts RunSuiteOptions) SuiteResult {
	if opts.Parallel {
		return runParallel(path, checks, opts.Timeout)
	}
	return runSequential(path, checks, opts.Timeout)
}

// runParallel executes all checks concurrently.
func runParallel(path string, checks []checker.Checker, timeout time.Duration) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, len(checks)),
	}

	if len(checks) == 0 {
		return result
	}

	suiteStart := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(checks))

	for i, check := range checks {
		go func(idx int, c checker.Checker) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					result.Results[idx] = checker.Result{
						Name:    c.Name(),
						ID:      c.ID(),
						Passed:  false,
						Status:  checker.Fail,
						Message: formatPanicMessage(r),
					}
				}
			}()
			res := runCheckWithTimeout(path, c, timeout)
			result.Results[idx] = res
		}(i, check)
	}

	wg.Wait()
	result.TotalDuration = time.Since(suiteStart)

	// Count results and check for critical failures
	for _, res := range result.Results {
		switch res.Status {
		case checker.Pass:
			result.Passed++
		case checker.Warn:
			result.Warnings++
		case checker.Fail:
			result.Failed++
			if !res.Passed {
				result.Aborted = true
			}
		case checker.Info:
			result.Info++
		}
	}

	return result
}

// runSequential executes checks one by one, stopping on first critical failure.
func runSequential(path string, checks []checker.Checker, timeout time.Duration) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, 0, len(checks)),
	}

	suiteStart := time.Now()

	for _, check := range checks {
		res := runCheckWithTimeout(path, check, timeout)
		result.Results = append(result.Results, res)

		switch res.Status {
		case checker.Pass:
			result.Passed++
		case checker.Warn:
			result.Warnings++
		case checker.Fail:
			result.Failed++
		case checker.Info:
			result.Info++
		}

		// Veto Power Logic: stop if a critical check fails
		if !res.Passed && res.Status == checker.Fail {
			result.Aborted = true
			break
		}
	}

	result.TotalDuration = time.Since(suiteStart)
	return result
}

// runCheckWithTimeout executes a single check with an optional timeout.
// If timeout is 0, no timeout is applied.
func runCheckWithTimeout(path string, c checker.Checker, timeout time.Duration) checker.Result {
	start := time.Now()

	// If no timeout, run directly
	if timeout == 0 {
		res, err := c.Run(path)
		duration := time.Since(start)
		if err != nil {
			res = checker.Result{
				Name:    c.Name(),
				ID:      c.ID(),
				Passed:  false,
				Status:  checker.Fail,
				Message: "Internal error: " + err.Error(),
			}
		}
		res.Duration = duration
		return res
	}

	// Run with timeout using context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channel for result
	type checkResult struct {
		res checker.Result
		err error
	}
	resultCh := make(chan checkResult, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultCh <- checkResult{
					res: checker.Result{
						Name:    c.Name(),
						ID:      c.ID(),
						Passed:  false,
						Status:  checker.Fail,
						Message: formatPanicMessage(r),
					},
				}
			}
		}()
		res, err := c.Run(path)
		resultCh <- checkResult{res: res, err: err}
	}()

	select {
	case <-ctx.Done():
		// Timeout occurred
		duration := time.Since(start)
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Fail,
			Message:  "Check timed out after " + timeout.String(),
			Duration: duration,
		}
	case cr := <-resultCh:
		duration := time.Since(start)
		res := cr.res
		if cr.err != nil {
			res = checker.Result{
				Name:    c.Name(),
				ID:      c.ID(),
				Passed:  false,
				Status:  checker.Fail,
				Message: "Internal error: " + cr.err.Error(),
			}
		}
		res.Duration = duration
		return res
	}
}

// TotalChecks returns the total number of checks that were run.
func (s *SuiteResult) TotalChecks() int {
	return len(s.Results)
}

// ScoredChecks returns the number of checks that affect the maturity score.
// This excludes Info status checks which are informational only.
func (s *SuiteResult) ScoredChecks() int {
	return s.Passed + s.Warnings + s.Failed
}

// Success returns true if no checks failed (warnings are allowed).
func (s *SuiteResult) Success() bool {
	return s.Failed == 0
}

// formatPanicMessage creates a user-friendly message from a panic value.
func formatPanicMessage(r interface{}) string {
	switch v := r.(type) {
	case error:
		return "Check panicked: " + v.Error()
	case string:
		return "Check panicked: " + v
	default:
		return "Check panicked unexpectedly"
	}
}
