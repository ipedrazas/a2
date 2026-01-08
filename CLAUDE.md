# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
task build              # Build binary to dist/a2
task test               # Run all tests
task test:unit          # Run unit tests only (excludes integration)
task fmt                # Format code with gofmt
task lint               # Run go vet and check formatting
task ci                 # Full CI: fmt, lint, test, build, then run a2 check
./dist/a2 check         # Run a2 on current directory after building
./dist/a2 check --lang python  # Run Python checks only
```

Single test: `go test -v -run TestName ./pkg/checks/...`

## Architecture

A2 is a multi-language code quality checker. It auto-detects project language(s) and runs appropriate checks.

### Core Flow

1. **Language Detection** (`pkg/language/detect.go`): Scans for indicator files (go.mod, pyproject.toml, etc.) to determine languages
2. **Check Registry** (`pkg/checks/registry.go`): Returns checks filtered by detected languages and config
3. **Runner** (`pkg/runner/runner.go`): Executes checks (parallel by default, sequential mode stops on critical failures)
4. **Output** (`pkg/output/`): Formats results as pretty terminal output or JSON

### Check Organization

Checks are organized by language in subdirectories:
- `pkg/checks/go/` - Go checks (go:module, go:build, go:tests, go:format, go:vet, go:coverage, go:deps)
- `pkg/checks/python/` - Python checks (python:project, python:build, python:tests, python:format, python:lint, python:coverage, python:deps)
- `pkg/checks/common/` - Language-agnostic (file_exists, external checks)

Each language directory has a `register.go` that returns `[]checker.CheckRegistration` with metadata (ID, order, critical flag).

### Key Types

```go
// pkg/checker/types.go
type Checker interface {
    ID() string
    Name() string
    Run(path string) (Result, error)
}

type Result struct {
    Name, ID, Message string
    Passed            bool
    Status            Status   // Pass, Warn, Fail
    Language          Language // go, python, common
}
```

### Veto System

Critical checks (e.g., build, tests) have `Critical: true` in their `CheckMeta`. In sequential mode, a Fail status stops execution immediately.

### Configuration (.a2.yaml)

```yaml
language:
  explicit: []           # Override auto-detect: ["go", "python"]
  python:
    test_runner: auto    # pytest, unittest
    formatter: auto      # black, ruff
    linter: auto         # ruff, flake8, pylint
    coverage_threshold: 80

checks:
  disabled:
    - python:deps        # Disable by check ID

files:
  required:
    - README.md
```

Check ID aliases exist for backward compatibility (e.g., "gofmt" maps to "go:format").

## Adding a New Check

1. Create `pkg/checks/<lang>/<check>.go` implementing `Checker` interface
2. Add to `register.go` in that language's directory with `CheckMeta` (ID, order, critical flag)
3. Write tests in `<check>_test.go`
4. Use `safepath` package for path operations

## Testing

Tests use `testify/suite`. Target coverage: >80%.

## Implementation Notes

1. All new checks should follow the existing pattern in `pkg/checks/`
2. Use `safepath` package for all file operations
3. Common checks go in `pkg/checks/common/`, language-specific in `pkg/checks/<lang>/`
4. Each check needs: implementation file, registration in `register.go`, tests