# Go Checks

This document describes all Go-specific checks available in A2.

## Overview

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `go:module` | Go Module | Yes | 100 | Verifies go.mod exists with valid Go version |
| `go:build` | Go Build | Yes | 110 | Compiles the project with `go build ./...` |
| `go:tests` | Go Tests | Yes | 120 | Runs tests with `go test ./...` |
| `go:race` | Go Race Detection | No | 125 | Detects data races with `go test -race` |
| `go:format` | Go Format | No | 200 | Checks formatting with `gofmt` |
| `go:vet` | Go Vet | No | 210 | Finds suspicious code with `go vet` |
| `go:coverage` | Go Coverage | No | 220 | Measures test coverage against threshold |
| `go:deps` | Go Vulnerabilities | No | 230 | Scans for vulnerabilities with `govulncheck` |
| `go:cyclomatic` | Go Complexity | No | 240 | Analyzes cyclomatic complexity of functions |
| `go:logging` | Go Logging | No | 250 | Detects structured logging vs fmt.Print |

---

## go:module

Verifies that `go.mod` exists and contains a valid Go version directive.

**Status:**
- **Pass**: go.mod exists with valid syntax and Go version
- **Warn**: go.mod exists but missing Go version directive
- **Fail**: go.mod not found or invalid syntax

---

## go:build

Runs `go build ./...` to verify the project compiles successfully.

**Status:**
- **Pass**: Build completes without errors
- **Fail**: Compilation errors

---

## go:tests

Runs `go test ./...` to execute all test packages.

**Status:**
- **Pass**: All tests pass or no test files found
- **Fail**: One or more tests fail

---

## go:race

Runs `go test -race -short ./...` to detect data races in concurrent code.

**Status:**
- **Pass**: No race conditions detected or no test files
- **Warn**: Race conditions detected or tests fail during race detection

---

## go:format

Runs `gofmt -l` to check if all Go files are properly formatted.

**Status:**
- **Pass**: All files properly formatted
- **Warn**: Unformatted files found

**Fix:** Run `gofmt -w .`

---

## go:vet

Runs `go vet ./...` to find suspicious constructs and potential bugs.

**Status:**
- **Pass**: No issues found
- **Warn**: Issues detected

---

## go:coverage

Runs `go test -cover ./...` and compares coverage percentage against the configured threshold.

**Configuration:**
```yaml
language:
  go:
    coverage_threshold: 80  # Default: 80%
```

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold

---

## go:deps

Scans for known vulnerabilities in Go dependencies using `govulncheck`.

**Requirements:** Install with `go install golang.org/x/vuln/cmd/govulncheck@latest`

**Status:**
- **Pass**: No vulnerabilities found or govulncheck not installed
- **Warn**: Vulnerabilities detected

---

## go:cyclomatic

Analyzes cyclomatic complexity of Go functions using the go/ast package.

**Configuration:**
```yaml
language:
  go:
    cyclomatic_threshold: 15  # Default: 15
```

**Complexity Calculation:**
- Base complexity: 1
- +1 for each: `if`, `for`, `range`, `case` (non-default), `&&`, `||`

**Status:**
- **Pass**: No functions exceed threshold
- **Warn**: Functions exceed complexity threshold

**Fix:** Break complex functions into smaller, focused functions.

---

## go:logging

Checks for proper structured logging practices instead of fmt.Print statements.

**Structured loggers detected:**
- `log/slog` (Go 1.21+)
- `go.uber.org/zap`
- `github.com/rs/zerolog`
- `github.com/sirupsen/logrus`

**Anti-patterns detected (in non-test files):**
- `fmt.Print`, `fmt.Println`, `fmt.Printf`

**Status:**
- **Pass**: Uses structured logging, no fmt.Print statements
- **Warn**: Uses fmt.Print for logging or no structured logger detected

**Fix:** Use `log/slog` or other structured logging libraries.

---

## Configuration Example

```yaml
language:
  go:
    coverage_threshold: 80
    cyclomatic_threshold: 15
```

## Check ID Aliases (Backward Compatibility)

| Alias | Maps To |
|-------|---------|
| `go_mod` | `go:module` |
| `build` | `go:build` |
| `tests` | `go:tests` |
| `gofmt` | `go:format` |
| `govet` | `go:vet` |
| `coverage` | `go:coverage` |
| `deps` | `go:deps` |
