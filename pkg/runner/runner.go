package runner

import (
	"sync"

	"github.com/ipedrazas/a2/pkg/checker"
)

// SuiteResult holds the results of running a suite of checks.
type SuiteResult struct {
	Results  []checker.Result // All check results
	Aborted  bool             // True if a Fail (veto) check stopped execution
	Passed   int              // Number of passed checks
	Warnings int              // Number of warning checks
	Failed   int              // Number of failed checks
}

// RunSuiteOptions configures how the suite is executed.
type RunSuiteOptions struct {
	Parallel bool // Run checks in parallel (default: true)
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
		return runParallel(path, checks)
	}
	return runSequential(path, checks)
}

// runParallel executes all checks concurrently.
func runParallel(path string, checks []checker.Checker) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, len(checks)),
	}

	if len(checks) == 0 {
		return result
	}

	var wg sync.WaitGroup
	wg.Add(len(checks))

	for i, check := range checks {
		go func(idx int, c checker.Checker) {
			defer wg.Done()
			res, err := c.Run(path)
			if err != nil {
				res = checker.Result{
					Name:    c.Name(),
					ID:      c.ID(),
					Passed:  false,
					Status:  checker.Fail,
					Message: "Internal error: " + err.Error(),
				}
			}
			result.Results[idx] = res
		}(i, check)
	}

	wg.Wait()

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
		}
	}

	return result
}

// runSequential executes checks one by one, stopping on first critical failure.
func runSequential(path string, checks []checker.Checker) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, 0, len(checks)),
	}

	for _, check := range checks {
		res, err := check.Run(path)
		if err != nil {
			res = checker.Result{
				Name:    check.Name(),
				ID:      check.ID(),
				Passed:  false,
				Status:  checker.Fail,
				Message: "Internal error: " + err.Error(),
			}
		}

		result.Results = append(result.Results, res)

		switch res.Status {
		case checker.Pass:
			result.Passed++
		case checker.Warn:
			result.Warnings++
		case checker.Fail:
			result.Failed++
		}

		// Veto Power Logic: stop if a critical check fails
		if !res.Passed && res.Status == checker.Fail {
			result.Aborted = true
			break
		}
	}

	return result
}

// TotalChecks returns the total number of checks that were run.
func (s *SuiteResult) TotalChecks() int {
	return len(s.Results)
}

// Success returns true if no checks failed (warnings are allowed).
func (s *SuiteResult) Success() bool {
	return s.Failed == 0
}
