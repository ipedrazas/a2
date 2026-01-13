# A2 Improvements Tracker

This file tracks planned improvements for the A2 codebase.

## Completed

- [x] **Common Utilities Extraction** - Created `pkg/checkutil` package with shared utilities (`TruncateMessage`, `Pluralize`, `RunCommand`, etc.)
- [x] **Check Execution Timing Metrics** - Added `Duration` field to results and `TotalDuration` to suite, displayed in all output formats
- [x] **Result Builder Pattern** - Added `ResultBuilder` to `pkg/checkutil` with `Pass()`, `Fail()`, `Warn()`, `Info()` methods. Updated ALL checks across all languages (Go, Python, Node, TypeScript, Java, Rust, Swift, Common). Reduces 6-field struct construction to single method call.

## High Priority

- [ ] **Exit Code Return**
  - Location: `pkg/output/pretty.go`, `json.go`, `toon.go`
  - Issue: All output formatters call `os.Exit(1)` directly, preventing graceful shutdown or signal handling
  - Solution: Return an exit code from output functions and let `cmd/root.go` handle the exit

## Medium Priority

- [ ] **Context/Timeout Support**
  - Location: `pkg/runner/runner.go`, `pkg/checker/types.go`
  - Issue: No timeout mechanism - slow/hanging checks can block the entire suite
  - Solution: Add `context.Context` with timeout to `Run()` method; add CLI flag `--timeout 30s`

- [ ] **Consistent Message Truncation**
  - Location: `pkg/checks/` (various)
  - Issue: `checkutil.TruncateMessage()` exists but isn't used consistently across checks
  - Solution: Apply truncation consistently to all checks that report tool output

- [ ] **Recommendations Metadata**
  - Location: `pkg/output/pretty.go:241-285`, `pkg/checker/types.go`
  - Issue: Hardcoded recommendations only cover subset of check IDs (missing many language checks)
  - Solution: Add `Suggestion` field to `CheckMeta` for auto-population of recommendations

## Low Priority

- [ ] **Tool-Not-Found Status Standardization**
  - Location: `pkg/checks/python/*.go`, `pkg/checks/node/*.go`, `pkg/checks/go/*.go`
  - Issue: Inconsistent handling when tools aren't installed - some return `Pass`, others different messages
  - Solution: Standardize using `Info` status and consistent messaging via `checkutil.ToolNotFoundError()`

- [ ] **Output Formatter Integration Tests**
  - Location: `pkg/output/`
  - Issue: No end-to-end tests verifying output consistency across pretty/JSON/TOON formats
  - Solution: Add integration tests with actual runner results

- [ ] **Panic Recovery in Parallel Mode**
  - Location: `pkg/runner/runner.go:49-105`
  - Issue: No protection against panics in concurrent check execution
  - Solution: Add `recover()` in goroutines to prevent one check crash from killing the suite

## Implementation Notes

When implementing improvements:
1. Update this file to mark items as completed
2. Add tests for new functionality
3. Run `task ci` to verify all tests pass
4. Update CHANGELOG.md if the change is user-facing
