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
	Errored       int              // Number of checks a2 could not evaluate (excluded from score)
	Skipped       int              // Number of checks that were not run (excluded from score)
	TotalDuration time.Duration    // Total time taken to run all checks
}

// RunSuiteOptions configures how the suite is executed.
type RunSuiteOptions struct {
	Parallel    bool          // Run checks in parallel (default: true)
	FailFast    bool          // Cancel remaining checks on first critical failure (parallel mode only)
	Timeout     time.Duration // Timeout for each individual check (0 = no timeout)
	Concurrency int           // Max concurrent checks in parallel mode (0 = GOMAXPROCS)
	OnProgress  ProgressFunc  // Optional callback for progress updates
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
// In parallel mode: checks run on a bounded worker pool. With FailFast, the first
// critical failure cancels the remaining (not-yet-started) checks.
// In sequential mode: stops immediately on first critical failure (veto power).
func RunSuiteWithOptions(path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	if opts.Parallel {
		return runParallel(path, registrations, opts)
	}
	return runSequential(path, registrations, opts)
}

// poolSize returns the number of workers to use for a parallel run, bounded by
// the number of registrations and the configured concurrency (GOMAXPROCS default).
func poolSize(opts RunSuiteOptions, n int) int {
	workers := opts.Concurrency
	if workers < 1 {
		workers = runtime.GOMAXPROCS(0)
	}
	if workers < 1 {
		workers = 1
	}
	if workers > n {
		workers = n
	}
	return workers
}

// applyOptional converts a Warn result to Info for optional checks so that
// optional checks never drag down the maturity score.
func applyOptional(meta checker.CheckMeta, res checker.Result) checker.Result {
	if meta.Optional && res.Status == checker.Warn {
		res.Status = checker.Info
		res.Passed = true
	}
	return res
}

// runParallel executes all checks on a bounded worker pool. When opts.FailFast
// is set, the first critical failure cancels checks that have not started yet.
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type task struct {
		idx int
		reg checker.CheckRegistration
	}

	workers := poolSize(opts, len(registrations))
	tasks := make(chan task)
	var wg sync.WaitGroup
	wg.Add(workers)

	reportProgress := func() {
		if opts.OnProgress != nil {
			done := int(atomic.AddInt32(&completed, 1))
			opts.OnProgress(done, total)
		}
	}

	for range workers {
		go func() {
			defer wg.Done()
			for t := range tasks {
				if opts.FailFast && ctx.Err() != nil {
					result.Results[t.idx] = skippedResult(t.reg)
					reportProgress()
					continue
				}

				res := applyOptional(t.reg.Meta, runCheckWithTimeout(ctx, path, t.reg.Checker, opts.Timeout))
				result.Results[t.idx] = res

				if opts.FailFast && t.reg.Meta.Critical && res.Status == checker.Fail && !res.Passed {
					cancel()
				}
				reportProgress()
			}
		}()
	}

	for idx, reg := range registrations {
		if opts.FailFast && ctx.Err() != nil {
			result.Results[idx] = skippedResult(reg)
			reportProgress()
			continue
		}
		select {
		case tasks <- task{idx: idx, reg: reg}:
		case <-ctx.Done():
			result.Results[idx] = skippedResult(reg)
			reportProgress()
		}
	}
	close(tasks)
	wg.Wait()

	result.TotalDuration = time.Since(suiteStart)
	tally(&result)
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
		res := applyOptional(reg.Meta, runCheckWithTimeout(context.Background(), path, reg.Checker, opts.Timeout))
		result.Results = append(result.Results, res)
		completed++

		if opts.OnProgress != nil {
			opts.OnProgress(completed, total)
		}

		// Veto Power Logic: stop if a critical check fails
		if !res.Passed && res.Status == checker.Fail {
			break
		}
	}

	result.TotalDuration = time.Since(suiteStart)
	tally(&result)
	return result
}

// tally counts each result by status and sets Aborted when a critical check failed.
func tally(result *SuiteResult) {
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
		case checker.Errored:
			result.Errored++
		case checker.Skipped:
			result.Skipped++
		}
	}
}

// runCheckWithTimeout executes a single check with an optional timeout and cancellation.
// If timeout is 0, no timeout is applied.
func runCheckWithTimeout(ctx context.Context, path string, c checker.Checker, timeout time.Duration) checker.Result {
	start := time.Now()

	// If no timeout, run directly (recovering from panics so a single check
	// cannot crash the whole run).
	if timeout == 0 {
		if ctx.Err() != nil {
			return skippedCheckResult(c)
		}
		return runCheckDirect(path, c, start)
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
				resultCh <- checkResult{res: erroredResult(c, formatPanicMessage(r))}
			}
		}()
		res, err := c.Run(path)
		resultCh <- checkResult{res: res, err: err}
	}()

	select {
	case <-ctx.Done():
		duration := time.Since(start)
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			res := erroredResult(c, "Check timed out after "+timeout.String())
			res.Duration = duration
			return res
		}
		res := skippedCheckResult(c)
		res.Duration = duration
		return res
	case cr := <-resultCh:
		duration := time.Since(start)
		res := cr.res
		if cr.err != nil {
			res = erroredResult(c, "Internal error: "+cr.err.Error())
		}
		res.Duration = duration
		return res
	}
}

// runCheckDirect runs a check synchronously, recovering from panics.
func runCheckDirect(path string, c checker.Checker, start time.Time) (res checker.Result) {
	defer func() {
		if r := recover(); r != nil {
			res = erroredResult(c, formatPanicMessage(r))
			res.Duration = time.Since(start)
		}
	}()
	res, err := c.Run(path)
	if err != nil {
		res = erroredResult(c, "Internal error: "+err.Error())
	}
	res.Duration = time.Since(start)
	return res
}

// erroredResult builds a result for a check that a2 could not evaluate.
func erroredResult(c checker.Checker, msg string) checker.Result {
	return checker.Result{
		Name:    c.Name(),
		ID:      c.ID(),
		Passed:  false,
		Status:  checker.Errored,
		Message: msg,
		Reason:  msg,
	}
}

// skippedCheckResult builds a result for a check that was cancelled before running.
func skippedCheckResult(c checker.Checker) checker.Result {
	return checker.Result{
		Name:    c.Name(),
		ID:      c.ID(),
		Passed:  true,
		Status:  checker.Skipped,
		Message: "Skipped",
		Reason:  "Skipped due to fail-fast",
	}
}

func skippedResult(reg checker.CheckRegistration) checker.Result {
	return skippedCheckResult(reg.Checker)
}

// TotalChecks returns the total number of checks that were run.
func (s *SuiteResult) TotalChecks() int {
	return len(s.Results)
}

// ScoredChecks returns the number of checks that affect the maturity score.
// This excludes Info, Errored, and Skipped status checks which do not score.
func (s *SuiteResult) ScoredChecks() int {
	return s.Passed + s.Warnings + s.Failed
}

// Success returns true if no checks failed (warnings are allowed).
// Errored and Skipped checks do not, on their own, fail the suite — they are
// surfaced separately so infrastructure problems are not mistaken for code defects.
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
