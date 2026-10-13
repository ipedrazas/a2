package runner

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// funcChecker runs an arbitrary function as a check.
type funcChecker struct {
	id string
	fn func(ctx context.Context) (checker.Result, error)
}

func (f *funcChecker) ID() string   { return f.id }
func (f *funcChecker) Name() string { return f.id }
func (f *funcChecker) Run(ctx context.Context, _ string) (checker.Result, error) {
	return f.fn(ctx)
}

func resultWith(id string, status checker.Status, passed bool) func(context.Context) (checker.Result, error) {
	return func(context.Context) (checker.Result, error) {
		return checker.Result{ID: id, Name: id, Status: status, Passed: passed}, nil
	}
}

func reg(c checker.Checker, meta checker.CheckMeta) checker.CheckRegistration {
	meta.ID = c.ID()
	meta.Name = c.Name()
	return checker.CheckRegistration{Checker: c, Meta: meta}
}

func TestSequential_NonCriticalFailDoesNotStopRun(t *testing.T) {
	regs := []checker.CheckRegistration{
		reg(&funcChecker{id: "a", fn: resultWith("a", checker.Fail, false)}, checker.CheckMeta{}),
		reg(&funcChecker{id: "b", fn: resultWith("b", checker.Pass, true)}, checker.CheckMeta{}),
	}
	result := RunSuiteWithOptions(context.Background(), ".", regs, RunSuiteOptions{Parallel: false})

	if len(result.Results) != 2 {
		t.Fatalf("expected both checks to run, got %d results", len(result.Results))
	}
	if result.Aborted {
		t.Error("a non-critical failure must not abort the run")
	}
	if result.Success() {
		t.Error("suite with a failure must not succeed")
	}
}

func TestStrict_FailsErroredAndMissingTool(t *testing.T) {
	missing := func(context.Context) (checker.Result, error) {
		return checker.Result{ID: "tool", Status: checker.Info, Passed: true, ToolMissing: true}, nil
	}
	regs := []checker.CheckRegistration{
		reg(&funcChecker{id: "errored", fn: resultWith("errored", checker.Errored, false)}, checker.CheckMeta{}),
		reg(&funcChecker{id: "tool", fn: missing}, checker.CheckMeta{}),
		reg(&funcChecker{id: "optional", fn: missing}, checker.CheckMeta{Optional: true}),
	}

	lenient := RunSuiteWithOptions(context.Background(), ".", regs, RunSuiteOptions{Parallel: true})
	if !lenient.Success() {
		t.Fatal("without strict, errored and missing-tool checks must not fail the suite")
	}

	strict := RunSuiteWithOptions(context.Background(), ".", regs, RunSuiteOptions{Parallel: true, Strict: true})
	if strict.Success() {
		t.Fatal("strict mode must fail the suite")
	}
	if strict.Failed != 2 {
		t.Errorf("expected 2 strict failures, got %d", strict.Failed)
	}
	if got := strict.Results[2].Status; got != checker.Info {
		t.Errorf("optional check must be exempt from strict mode, got %s", got)
	}
}

func TestTimeout_KillsRunningCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses the sleep command")
	}
	sleeper := &funcChecker{id: "sleep", fn: func(ctx context.Context) (checker.Result, error) {
		res := checkutil.RunCommand(ctx, ".", "sleep", "30")
		return checker.Result{ID: "sleep", Status: checker.Fail, Passed: false, Reason: res.Output()}, nil
	}}

	start := time.Now()
	result := RunSuiteWithOptions(context.Background(), ".", []checker.CheckRegistration{reg(sleeper, checker.CheckMeta{})},
		RunSuiteOptions{Parallel: true, Timeout: 100 * time.Millisecond})

	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("timed-out check took %s; the command was not killed", elapsed)
	}
	if got := result.Results[0].Status; got != checker.Errored {
		t.Errorf("expected timed-out check to be Errored, got %s", got)
	}
}

func TestFailFast_CancelsInFlightChecks(t *testing.T) {
	blocker := &funcChecker{id: "blocker", fn: func(ctx context.Context) (checker.Result, error) {
		<-ctx.Done()
		return checker.Result{ID: "blocker", Status: checker.Fail}, nil
	}}
	veto := &funcChecker{id: "veto", fn: func(context.Context) (checker.Result, error) {
		time.Sleep(20 * time.Millisecond) // let the blocker start first
		return checker.Result{ID: "veto", Status: checker.Fail, Passed: false}, nil
	}}
	regs := []checker.CheckRegistration{
		reg(blocker, checker.CheckMeta{}),
		reg(veto, checker.CheckMeta{Critical: true}),
	}

	done := make(chan SuiteResult, 1)
	go func() {
		done <- RunSuiteWithOptions(context.Background(), ".", regs,
			RunSuiteOptions{Parallel: true, FailFast: true, Concurrency: 2})
	}()

	select {
	case result := <-done:
		if got := result.Results[0]; got.Status != checker.Skipped || got.Reason != "Skipped due to fail-fast" {
			t.Errorf("expected in-flight check to be skipped by fail-fast, got %s %q", got.Status, got.Reason)
		}
		if !result.Aborted {
			t.Error("critical failure must abort the run")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("fail-fast did not cancel the in-flight check")
	}
}

func TestCancelledContext_SkipsChecks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	regs := []checker.CheckRegistration{
		reg(&funcChecker{id: "a", fn: resultWith("a", checker.Pass, true)}, checker.CheckMeta{}),
	}
	result := RunSuiteWithOptions(ctx, ".", regs, RunSuiteOptions{Parallel: false})

	if got := result.Results[0]; got.Status != checker.Skipped || got.Reason != "Skipped: run was cancelled" {
		t.Errorf("expected skipped result for cancelled run, got %s %q", got.Status, got.Reason)
	}
}

func TestParallel_SlowChecksAreBounded(t *testing.T) {
	var running, peak int32
	slow := func(id string) checker.CheckRegistration {
		return reg(&funcChecker{id: id, fn: func(context.Context) (checker.Result, error) {
			n := atomic.AddInt32(&running, 1)
			for {
				p := atomic.LoadInt32(&peak)
				if n <= p || atomic.CompareAndSwapInt32(&peak, p, n) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&running, -1)
			return checker.Result{ID: id, Status: checker.Pass, Passed: true}, nil
		}}, checker.CheckMeta{Speed: checker.SpeedSlow})
	}
	regs := []checker.CheckRegistration{slow("s1"), slow("s2"), slow("s3"), slow("s4"), slow("s5")}

	result := RunSuiteWithOptions(context.Background(), ".", regs, RunSuiteOptions{Parallel: true, Concurrency: 8})

	if result.Passed != 5 {
		t.Fatalf("expected all slow checks to pass, got %d", result.Passed)
	}
	if peak > DefaultSlowConcurrency {
		t.Errorf("expected at most %d slow checks at once, saw %d", DefaultSlowConcurrency, peak)
	}
}

func TestSlowPoolSize(t *testing.T) {
	tests := []struct {
		name string
		opts RunSuiteOptions
		want int
	}{
		{"default", RunSuiteOptions{Concurrency: 8}, DefaultSlowConcurrency},
		{"configured", RunSuiteOptions{Concurrency: 8, SlowConcurrency: 3}, 3},
		{"capped by concurrency", RunSuiteOptions{Concurrency: 1, SlowConcurrency: 4}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slowPoolSize(tt.opts, 10); got != tt.want {
				t.Errorf("slowPoolSize() = %d, want %d", got, tt.want)
			}
		})
	}
}
