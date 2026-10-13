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
	Aborted       bool             // True if a Critical check failed (the veto)
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
	Parallel        bool          // Run checks in parallel (default: true)
	FailFast        bool          // Cancel remaining and in-flight checks on first critical failure (parallel mode only)
	Strict          bool          // Fail non-optional checks that errored or whose tool is not installed
	Timeout         time.Duration // Timeout for each individual check (0 = no timeout)
	Concurrency     int           // Max concurrent fast checks in parallel mode (0 = GOMAXPROCS)
	SlowConcurrency int           // Max concurrent slow checks in parallel mode (0 = DefaultSlowConcurrency)
	OnProgress      ProgressFunc  // Optional callback for progress updates
}

// DefaultSlowConcurrency is how many slow checks (builds, tests, coverage, race)
// run at once by default. Each of them already uses every core, so running
// more in parallel oversubscribes the CPU and makes timeouts more likely.
const DefaultSlowConcurrency = 2

// errFailFast is the cancellation cause recorded when fail-fast stops the run.
var errFailFast = errors.New("fail-fast")

// runSuite executes a suite of checks against the given path.
// Uses parallel execution by default for better performance.
func runSuite(path string, registrations []checker.CheckRegistration) SuiteResult {
	return RunSuiteWithOptions(context.Background(), path, registrations, RunSuiteOptions{Parallel: true})
}

// runSuiteSequential executes checks sequentially, stopping on first critical failure.
// Use this when you need veto power behavior or have limited CPU resources.
func runSuiteSequential(path string, registrations []checker.CheckRegistration) SuiteResult {
	return RunSuiteWithOptions(context.Background(), path, registrations, RunSuiteOptions{Parallel: false})
}

// RunSuiteWithOptions executes a suite of checks with the given options.
// In parallel mode: fast and slow checks run on separate bounded worker pools.
// With FailFast, the first critical failure cancels the remaining checks,
// including those already running.
// In sequential mode: stops immediately on first critical failure (veto power).
// Cancelling ctx (e.g. on Ctrl-C) stops running checks and skips the rest.
func RunSuiteWithOptions(ctx context.Context, path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	if opts.Parallel {
		return runParallel(ctx, path, registrations, opts)
	}
	return runSequential(ctx, path, registrations, opts)
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

// slowPoolSize returns the number of workers for slow checks. It never exceeds
// the fast pool's configured size, so a low `concurrency` also limits slow checks.
func slowPoolSize(opts RunSuiteOptions, n int) int {
	workers := opts.SlowConcurrency
	if workers < 1 {
		workers = DefaultSlowConcurrency
	}
	if limit := poolSize(opts, n); workers > limit {
		workers = limit
	}
	return workers
}

// applyOptional converts a Warn result to Info for optional checks so that
// optional checks never drag down the maturity score, and stamps the result
// with the check's Critical flag so scoring can weight it.
func applyOptional(meta checker.CheckMeta, res checker.Result) checker.Result {
	if meta.Optional && res.Status == checker.Warn {
		res.Status = checker.Info
		res.Passed = true
	}
	res.Critical = meta.Critical
	return res
}

// applyStrict turns results a2 could not trust into failures: checks that
// errored (crash, panic, timeout) and checks whose tool is not installed.
// Optional checks are exempt, since they never gate the suite.
func applyStrict(meta checker.CheckMeta, res checker.Result) checker.Result {
	if meta.Optional || (res.Status != checker.Errored && !res.ToolMissing) {
		return res
	}
	res.Status = checker.Fail
	res.Passed = false
	res.Reason = "Strict mode: " + res.Reason
	return res
}

// finalize applies the per-check result policies in order.
func finalize(meta checker.CheckMeta, res checker.Result, opts RunSuiteOptions) checker.Result {
	res = applyOptional(meta, res)
	if opts.Strict {
		res = applyStrict(meta, res)
	}
	return res
}

// isVeto reports whether res is a critical failure that stops the run.
func isVeto(res checker.Result) bool {
	return res.Critical && res.Status == checker.Fail && !res.Passed
}

// runParallel executes checks on two bounded worker pools: one for fast
// checks and a smaller one for slow checks. When opts.FailFast is set, the
// first critical failure cancels every check that has not finished yet.
func runParallel(ctx context.Context, path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, len(registrations)),
	}
	if len(registrations) == 0 {
		return result
	}

	suiteStart := time.Now()
	total := len(registrations)
	var completed int32

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	type task struct {
		idx int
		reg checker.CheckRegistration
	}

	reportProgress := func() {
		if opts.OnProgress != nil {
			done := int(atomic.AddInt32(&completed, 1))
			opts.OnProgress(done, total)
		}
	}

	var fast, slow []task
	for idx, reg := range registrations {
		if reg.Meta.Speed == checker.SpeedSlow {
			slow = append(slow, task{idx: idx, reg: reg})
		} else {
			fast = append(fast, task{idx: idx, reg: reg})
		}
	}

	var wg sync.WaitGroup
	runPool := func(tasks []task, workers int) {
		if len(tasks) == 0 {
			return
		}
		if workers > len(tasks) {
			workers = len(tasks)
		}
		queue := make(chan task)
		for range workers {
			wg.Go(func() {
				for t := range queue {
					res := finalize(t.reg.Meta, runCheck(ctx, path, t.reg.Checker, opts.Timeout), opts)
					result.Results[t.idx] = res
					if opts.FailFast && isVeto(res) {
						cancel(errFailFast)
					}
					reportProgress()
				}
			})
		}
		wg.Go(func() {
			defer close(queue)
			for _, t := range tasks {
				select {
				case queue <- t:
				case <-ctx.Done():
					result.Results[t.idx] = skippedResult(ctx, t.reg)
					reportProgress()
				}
			}
		})
	}

	runPool(fast, poolSize(opts, len(registrations)))
	runPool(slow, slowPoolSize(opts, len(registrations)))
	wg.Wait()

	result.TotalDuration = time.Since(suiteStart)
	tally(&result)
	return result
}

// runSequential executes checks one by one, stopping on first critical failure.
func runSequential(ctx context.Context, path string, registrations []checker.CheckRegistration, opts RunSuiteOptions) SuiteResult {
	result := SuiteResult{
		Results: make([]checker.Result, 0, len(registrations)),
	}

	suiteStart := time.Now()
	total := len(registrations)
	completed := 0

	for _, reg := range registrations {
		res := finalize(reg.Meta, runCheck(ctx, path, reg.Checker, opts.Timeout), opts)
		result.Results = append(result.Results, res)
		completed++

		if opts.OnProgress != nil {
			opts.OnProgress(completed, total)
		}

		// Veto Power Logic: stop if a critical check fails
		if isVeto(res) {
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
			if isVeto(res) {
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

// runCheck executes a single check with an optional timeout, recovering from
// panics so a single check cannot crash the whole run. If timeout is 0, no
// timeout is applied. The check's context is cancelled when it times out or
// ctx is cancelled, which kills any command it started; runCheck returns
// immediately even if the check itself ignores cancellation.
func runCheck(ctx context.Context, path string, c checker.Checker, timeout time.Duration) checker.Result {
	start := time.Now()
	if ctx.Err() != nil {
		return skippedCheckResult(ctx, c)
	}

	checkCtx, cancel := context.WithCancel(ctx)
	if timeout > 0 {
		checkCtx, cancel = context.WithTimeout(ctx, timeout)
	}
	defer cancel()

	resultCh := make(chan checker.Result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultCh <- erroredResult(c, formatPanicMessage(r))
			}
		}()
		res, err := c.Run(checkCtx, path)
		if err != nil {
			res = erroredResult(c, "Internal error: "+err.Error())
		}
		resultCh <- res
	}()

	var res checker.Result
	select {
	case res = <-resultCh:
		// The check finished before any cancellation: its result stands.
		if checkCtx.Err() == nil {
			res.Duration = time.Since(start)
			return res
		}
	case <-checkCtx.Done():
	}

	// The check was cancelled. Any command it ran was killed, so whatever it
	// returned is a misleading failure: report the cancellation instead.
	switch {
	case ctx.Err() != nil:
		res = skippedCheckResult(ctx, c)
	case errors.Is(checkCtx.Err(), context.DeadlineExceeded):
		res = erroredResult(c, "Check timed out after "+timeout.String())
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

// skippedCheckResult builds a result for a check that was cancelled before it
// finished, either by fail-fast or because the run was interrupted.
func skippedCheckResult(ctx context.Context, c checker.Checker) checker.Result {
	reason := "Skipped: run was cancelled"
	if errors.Is(context.Cause(ctx), errFailFast) {
		reason = "Skipped due to fail-fast"
	}
	return checker.Result{
		Name:    c.Name(),
		ID:      c.ID(),
		Passed:  true,
		Status:  checker.Skipped,
		Message: "Skipped",
		Reason:  reason,
	}
}

func skippedResult(ctx context.Context, reg checker.CheckRegistration) checker.Result {
	return skippedCheckResult(ctx, reg.Checker)
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

// criticalWeight is how much more a Critical check counts toward the maturity
// score than an ordinary check, so the score reflects risk rather than a flat
// pass ratio.
const criticalWeight = 2.0

// WeightedTally returns the weighted passed and scored totals, where Critical
// checks count criticalWeight times as much as ordinary checks. Info, Errored,
// and Skipped results are excluded, matching ScoredChecks. When no Critical
// checks are present this reduces exactly to the raw passed/scored ratio.
func (s *SuiteResult) WeightedTally() (passed, scored float64) {
	for _, r := range s.Results {
		w := 1.0
		if r.Critical {
			w = criticalWeight
		}
		switch r.Status {
		case checker.Pass:
			passed += w
			scored += w
		case checker.Warn, checker.Fail:
			scored += w
		}
	}
	return passed, scored
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
