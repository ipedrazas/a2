package runner

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
)

// ProgressFunc is called when a check completes.
// completed is the number of checks finished so far, total is the total number of checks.
type ProgressFunc func(completed, total int)

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
	Parallel   bool          // Run checks in parallel (default: true)
	FailFast   bool          // Cancel remaining checks on first critical failure (parallel mode only)
	Timeout    time.Duration // Timeout for each individual check (0 = no timeout)
	OnProgress ProgressFunc  // Optional callback for progress updates
}

// runSuite executes a suite of checks against the given path.
// Uses parallel execution by default for better performance.
// If a check returns Fail status, Aborted is set to true.
func runSuite(path string, registrations []checker.CheckRegistration) SuiteResult {
	return RunSuiteWithOptions(path, registrations, RunSuiteOptions{Parallel: true})
}

// runSuiteSequential executes checks sequentially, stopping on first critical failure.
// Use this when you need veto power behavior or have limited CPU resources.
func runSuiteSequential(path string, registrations []checker.CheckRegistration) SuiteResult {
	return RunSuiteWithOptions(path, registrations, RunSuiteOptions{Parallel: false})
}

// RunSuiteWithOptions executes a suite of checks with the given options.
// In parallel mode: all checks run concurrently, Aborted is set if any critical check fails.
// In sequential mode: stops immediately on first critical failure (veto power).
func RunSuiteWithOptions(path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	if opts.Parallel {
		if opts.FailFast {
			return runParallelFailFast(path, registrations, opts)
		}
		return runParallel(path, registrations, opts)
	}
	return runSequential(path, registrations, opts)
}

// runParallel executes all checks concurrently.
func runParallel(path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, len(registrations)),
	}

	if len(registrations) == 0 {
		return result
	}

	suiteStart := time.Now()
	total := len(registrations)
	var completed int32

	var wg sync.WaitGroup
	wg.Add(len(registrations))

	for i, reg := range registrations {
		go func(idx int, registration checker.CheckRegistration) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					result.Results[idx] = checker.Result{
						Name:    registration.Checker.Name(),
						ID:      registration.Checker.ID(),
						Passed:  false,
						Status:  checker.Fail,
						Message: formatPanicMessage(r),
						Reason:  formatPanicMessage(r),
					}
				}
				// Update progress after check completes
				if opts.OnProgress != nil {
					done := int(atomic.AddInt32(&completed, 1))
					opts.OnProgress(done, total)
				}
			}()
			res := runCheckWithTimeout(context.Background(), path, registration.Checker, opts.Timeout)
			// Convert Warn to Info for optional checks
			if registration.Meta.Optional && res.Status == checker.Warn {
				res.Status = checker.Info
				res.Passed = true
			}
			result.Results[idx] = res
		}(i, reg)
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

func runParallelFailFast(path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, len(registrations)),
	}
	if len(registrations) == 0 {
		return result
	}

	suiteStart := time.Now()
	total := len(registrations)
	var completed int32

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type task struct {
		idx int
		reg checker.CheckRegistration
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(registrations) {
		workers = len(registrations)
	}

	tasks := make(chan task)
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for t := range tasks {
				if ctx.Err() != nil {
					result.Results[t.idx] = cancelledResult(t.reg)
					if opts.OnProgress != nil {
						done := int(atomic.AddInt32(&completed, 1))
						opts.OnProgress(done, total)
					}
					continue
				}

				res := runCheckWithTimeout(ctx, path, t.reg.Checker, opts.Timeout)
				// Convert Warn to Info for optional checks
				if t.reg.Meta.Optional && res.Status == checker.Warn {
					res.Status = checker.Info
					res.Passed = true
				}
				result.Results[t.idx] = res

				if t.reg.Meta.Critical && res.Status == checker.Fail && !res.Passed {
					cancel()
				}

				if opts.OnProgress != nil {
					done := int(atomic.AddInt32(&completed, 1))
					opts.OnProgress(done, total)
				}
			}
		}()
	}

	for idx, reg := range registrations {
		if ctx.Err() != nil {
			result.Results[idx] = cancelledResult(reg)
			if opts.OnProgress != nil {
				done := int(atomic.AddInt32(&completed, 1))
				opts.OnProgress(done, total)
			}
			continue
		}
		select {
		case tasks <- task{idx: idx, reg: reg}:
		case <-ctx.Done():
			result.Results[idx] = cancelledResult(reg)
			if opts.OnProgress != nil {
				done := int(atomic.AddInt32(&completed, 1))
				opts.OnProgress(done, total)
			}
		}
	}
	close(tasks)
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
func runSequential(path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, 0, len(registrations)),
	}

	suiteStart := time.Now()
	total := len(registrations)
	completed := 0

	for _, reg := range registrations {
		res := runCheckWithTimeout(context.Background(), path, reg.Checker, opts.Timeout)
		// Convert Warn to Info for optional checks
		if reg.Meta.Optional && res.Status == checker.Warn {
			res.Status = checker.Info
			res.Passed = true
		}
		result.Results = append(result.Results, res)
		completed++

		// Call progress callback
		if opts.OnProgress != nil {
			opts.OnProgress(completed, total)
		}

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

// runCheckWithTimeout executes a single check with an optional timeout and cancellation.
// If timeout is 0, no timeout is applied.
func runCheckWithTimeout(ctx context.Context, path string, c checker.Checker, timeout time.Duration) checker.Result {
	start := time.Now()

	// If no timeout, run directly
	if timeout == 0 {
		if ctx.Err() != nil {
			return checker.Result{
				Name:    c.Name(),
				ID:      c.ID(),
				Passed:  true,
				Status:  checker.Info,
				Message: "Cancelled",
				Reason:  "Cancelled due to fail-fast",
			}
		}
		res, err := c.Run(path)
		duration := time.Since(start)
		if err != nil {
			res = checker.Result{
				Name:    c.Name(),
				ID:      c.ID(),
				Passed:  false,
				Status:  checker.Fail,
				Message: "Internal error: " + err.Error(),
				Reason:  "Internal error: " + err.Error(),
			}
		}
		res.Duration = duration
		return res
	}

	// Run with timeout using context
	ctx, cancel := context.WithTimeout(ctx, timeout)
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
						Reason:  formatPanicMessage(r),
					},
				}
			}
		}()
		res, err := c.Run(path)
		resultCh <- checkResult{res: res, err: err}
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// Timeout occurred
			duration := time.Since(start)
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Fail,
				Message:  "Check timed out after " + timeout.String(),
				Reason:   "Check timed out after " + timeout.String(),
				Duration: duration,
			}
		}
		duration := time.Since(start)
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Info,
			Message:  "Cancelled",
			Reason:   "Cancelled due to fail-fast",
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
				Reason:  "Internal error: " + cr.err.Error(),
			}
		}
		res.Duration = duration
		return res
	}
}

func cancelledResult(reg checker.CheckRegistration) checker.Result {
	return checker.Result{
		Name:    reg.Checker.Name(),
		ID:      reg.Checker.ID(),
		Passed:  true,
		Status:  checker.Info,
		Message: "Cancelled",
		Reason:  "Cancelled due to fail-fast",
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
func formatPanicMessage(r any) string {
	switch v := r.(type) {
	case error:
		return "Check panicked: " + v.Error()
	case string:
		return "Check panicked: " + v
	default:
		return "Check panicked unexpectedly"
	}
}
