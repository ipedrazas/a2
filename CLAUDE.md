# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

This project uses [Task](https://taskfile.dev) as the task runner. Key commands:

```bash
task build              # Build binary to dist/a2
task test               # Run all tests
task test:unit          # Run unit tests only (excludes integration)
task test:coverage      # Run tests with coverage report
task fmt                # Format code with gofmt
task lint               # Run go vet and format check
task ci                 # Run full CI pipeline (fmt, lint, test, build, check)
task dependencies       # Install all external tool dependencies
```

Run single test file or function:
```bash
go test -v ./pkg/checks/go/... -run TestBuildCheck
go test -v ./pkg/runner/... -run TestRunSuite
```

## Architecture Overview

A2 is a code quality checker that runs configurable checks against repositories and provides a health/maturity score.

### Core Flow

1. **CLI** (`cmd/root.go`): Parses flags, loads config, detects languages, runs checks, outputs results
2. **Language Detection** (`pkg/language/`): Auto-detects project languages from indicator files (go.mod, package.json, etc.)
3. **Check Registry** (`pkg/checks/registry.go`): Aggregates checks from all language packages and common checks
4. **Runner** (`pkg/runner/`): Executes checks in parallel or sequential mode with optional timeout
5. **Output** (`pkg/output/`): Formats results as pretty terminal output, JSON, or TOON (minimal token format)

### Package Structure

- `pkg/checker/types.go` - Core types: `Checker` interface, `Result`, `Status` (Pass/Warn/Fail/Info), `CheckMeta`
- `pkg/checks/{go,python,node,java,rust,typescript,swift}/` - Language-specific checks, each has a `register.go`
- `pkg/checks/common/` - Language-agnostic checks (dockerfile, CI, health, k8s, etc.)
- `pkg/checkutil/` - Shared utilities: `ResultBuilder` for creating results, `RunCommand` for exec
- `pkg/config/` - Configuration loading from `.a2.yaml`
- `pkg/profiles/` - Application profiles (cli, api, library, desktop) that skip irrelevant checks
- `pkg/targets/` - Maturity targets (poc, production) controlling check strictness
- `pkg/maturity/` - Maturity level assessment based on check results

### Implementing a New Check

1. Create check struct implementing `checker.Checker` interface:
   ```go
   type MyCheck struct{}
   func (c *MyCheck) ID() string { return "lang:mycheck" }
   func (c *MyCheck) Name() string { return "My Check" }
   func (c *MyCheck) Run(path string) (checker.Result, error) { ... }
   ```

2. Use `checkutil.ResultBuilder` for consistent result creation:
   ```go
   rb := checkutil.NewResultBuilder(c, checker.LangGo)
   return rb.Pass("All good"), nil
   return rb.Warn("Minor issue found"), nil
   return rb.Fail("Critical failure"), nil
   return rb.ToolNotInstalled("mytool", "go install mytool@latest"), nil
   ```

3. Register in the language's `register.go` with `CheckMeta` (ID, Name, Critical flag, Order, Suggestion)

### Check Severity Levels

- **Pass**: Check passed successfully
- **Warn**: Issue found but not critical (doesn't stop execution)
- **Fail**: Critical failure, sets `Aborted=true` in suite results
- **Info**: Informational only, doesn't affect maturity score (used for missing optional tools)

### External Checks

Custom checks can be defined in `.a2.yaml`:
```yaml
external:
  - id: lint
    command: golangci-lint
    args: ["run", "./..."]
    severity: warn  # or fail
```

### Wildcard Support for Skip List

Skip patterns support wildcards for flexible check filtering:
- `*:tests` - skips all test checks (go:tests, python:tests, etc.)
- `node:*` - skips all Node.js checks
- `common:*` - skips all common checks
- `*:*` - skips ALL checks (use with caution)

**Important:** In `.a2.yaml`, patterns that contain `*` must be quoted (e.g. `"*:logging"` or `'*:logging'`). Unquoted `*` is interpreted by YAML as an alias and will cause a parse error.

Examples:
```yaml
checks:
  disabled:
    - "*:tests"      # Skip all test checks across all languages (quotes required)
    - "node:*"       # Skip all Node.js checks
    - "*:logging"     # Skip logging checks for all languages
    - "go:deps"      # Skip specific Go vulnerability check (exact match)
```

Exit codes: 0=Pass, 1=Warn, 2+=Fail

### Testing Patterns

Tests use testify/suite. Example pattern:
```go
type MyCheckTestSuite struct {
    suite.Suite
}

func (s *MyCheckTestSuite) TestMyCheck_Success() {
    // test implementation
}

func TestMyCheckTestSuite(t *testing.T) {
    suite.Run(t, new(MyCheckTestSuite))
}
```
