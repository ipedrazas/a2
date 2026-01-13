# A2 Improvements Tracker

This file tracks planned improvements for the A2 codebase.

## Completed

- [x] **Common Utilities Extraction** - Created `pkg/checkutil` package with shared utilities (`TruncateMessage`, `Pluralize`, `RunCommand`, etc.)
- [x] **Check Execution Timing Metrics** - Added `Duration` field to results and `TotalDuration` to suite, displayed in all output formats
- [x] **Result Builder Pattern** - Added `ResultBuilder` to `pkg/checkutil` with `Pass()`, `Fail()`, `Warn()`, `Info()` methods. Updated ALL checks across all languages (Go, Python, Node, TypeScript, Java, Rust, Swift, Common). Reduces 6-field struct construction to single method call.
- [x] **Exit Code Return** - Output formatters now return `(bool, error)` instead of calling `os.Exit(1)` directly. Exit code handling moved to `cmd/root.go` for proper graceful shutdown support.
- [x] **Context/Timeout Support** - Added `Timeout` field to `RunSuiteOptions`, `runCheckWithTimeout()` function with context-based timeout, and CLI flag `--timeout 30s`. Individual checks that exceed timeout fail with "Check timed out" message.
- [x] **Consistent Message Truncation** - Applied `checkutil.TruncateMessage()` to all Go checks (tests, build, format, vet, deps), Python build, and TypeScript build. Tool output is now consistently truncated to 200 chars.
- [x] **Recommendations Metadata** - Added `Suggestion` field to `CheckMeta` and populated it for all checks across all languages. `printRecommendations()` now dynamically looks up suggestions from check metadata via `checks.GetSuggestions()`.

## Medium Priority

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
