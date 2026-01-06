# AGENTS.md - Guide for AI Agents

This document provides guidance for AI agents working with the A2 codebase. A2 is a code quality checker for Go projects that runs a suite of checks and provides a health score.

## Project Overview

A2 is a CLI tool that:
- Runs 8+ built-in checks on Go projects
- Supports custom external checks via binaries
- Provides both pretty (colored terminal) and JSON output
- Integrates with CI/CD (GitHub Actions, pre-commit hooks)
- Uses a veto system where critical checks stop execution on failure

## Architecture

### Directory Structure

```
a2/
├── cmd/              # CLI commands (Cobra-based)
│   └── root.go       # Main command entry point
├── pkg/
│   ├── checker/      # Core check interface and types
│   ├── checks/       # Built-in check implementations
│   ├── config/       # Configuration loading (.a2.yaml)
│   ├── output/       # Output formatters (pretty, JSON)
│   ├── runner/       # Check execution engine (parallel/sequential)
│   └── safepath/     # Path validation utilities
├── main.go           # Application entry point
└── .a2.yaml          # Example configuration file
```

### Key Components

#### 1. Check Interface (`pkg/checker/types.go`)

All checks implement the `Checker` interface:

```go
type Checker interface {
    ID() string      // Unique identifier (e.g., "build", "tests")
    Name() string    // Human-readable name (e.g., "Build")
    Run(path string) (Result, error)
}
```

**Status Levels:**
- `Pass`: Check passed, no issues
- `Warn`: Non-critical issue, execution continues
- `Fail`: Critical failure, stops execution (veto power)

#### 2. Check Registry (`pkg/checks/registry.go`)

- `GetChecks(cfg *config.Config)`: Returns enabled checks based on config
- Critical checks (Fail severity) are ordered first
- Disabled checks are filtered out
- External checks from config are appended

#### 3. Runner (`pkg/runner/runner.go`)

- **Parallel mode** (default): All checks run concurrently
- **Sequential mode**: Stops on first critical failure (veto power)
- Returns `SuiteResult` with aggregated results

#### 4. Configuration (`pkg/config/config.go`)

Loads `.a2.yaml` with:
- Coverage thresholds
- Required files
- Disabled checks
- External check definitions
- Execution options (parallel/sequential)

## Adding New Checks

### Built-in Check

1. **Create check file** in `pkg/checks/`:
   - `{check_name}.go` - Implementation
   - `{check_name}_test.go` - Tests

2. **Implement the Checker interface**:

```go
package checks

type MyCheck struct {
    // Configuration fields
}

func (c *MyCheck) ID() string {
    return "my_check"
}

func (c *MyCheck) Name() string {
    return "My Check"
}

func (c *MyCheck) Run(path string) (checker.Result, error) {
    // Implementation
    return checker.Result{
        Name:    c.Name(),
        ID:      c.ID(),
        Passed:  true,
        Status:  checker.Pass, // or Warn, Fail
        Message: "Check completed",
    }, nil
}
```

3. **Register in `registry.go`**:
   - Add to `allChecks` slice in `GetChecks()`
   - Place critical checks (Fail) before warning checks

4. **Write tests**:
   - Test both success and failure cases
   - Test with different path configurations
   - Follow existing test patterns

### External Check

External checks are configured in `.a2.yaml`:

```yaml
external:
  - id: my_tool
    name: My Tool
    command: my-tool
    args: ["--flag", "value"]
    severity: warn  # or fail
```

**Exit Code Protocol:**
- `0`: Pass
- `1`: Warning
- `2+`: Fail

**Output Format:**
- Plain text (displayed as-is)
- JSON: `{"message": "...", "status": "warn"}`

## Code Conventions

### Naming

- **Check IDs**: Use snake_case (e.g., `go_mod`, `file_exists`)
- **Check Names**: Use Title Case (e.g., "Go Module", "File Exists")
- **Files**: Match check ID (e.g., `go_mod.go`, `file_exists.go`)

### Error Handling

- Checks should return descriptive error messages
- Internal errors should set `Status: Fail` and `Passed: false`
- Use `safepath` utilities for path validation

### Testing

- **Test files**: Always include `_test.go` files
- **Coverage**: Maintain high test coverage
- **Test patterns**: Use table-driven tests where appropriate
- **Run tests**: `go test ./...`

### Status Determination

- **Pass**: Everything is correct
- **Warn**: Issue found but non-critical (e.g., low coverage, missing optional file)
- **Fail**: Critical issue (e.g., build fails, tests fail, missing go.mod)

## Common Tasks

### Running Checks Locally

```bash
# Run on current directory
go run . check

# Run on specific path
go run . check /path/to/project

# JSON output
go run . check --format json
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/checks/...
```

### Building

```bash
# Build binary
go build -o dist/a2

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o dist/a2-linux
```

### Adding Configuration Options

1. Add field to `pkg/config/config.go` struct
2. Update YAML parsing logic
3. Use in check implementation
4. Document in README.md

## Important Files

- **`cmd/root.go`**: CLI entry point, command definitions
- **`pkg/checks/registry.go`**: Check registration and filtering
- **`pkg/runner/runner.go`**: Execution engine (parallel/sequential)
- **`pkg/config/config.go`**: Configuration loading and defaults
- **`pkg/output/pretty.go`**: Colored terminal output
- **`pkg/output/json.go`**: JSON output for CI/CD

## Configuration File (.a2.yaml)

Located in project root. Structure:

```yaml
coverage:
  threshold: 80  # Coverage percentage threshold

files:
  required:
    - README.md
    - LICENSE

checks:
  disabled:
    - deps  # Disable specific checks

execution:
  parallel: true  # Run checks in parallel

external:
  - id: lint
    name: Golangci-lint
    command: golangci-lint
    args: ["run", "./..."]
    severity: warn
```

## Exit Codes

- `0`: All checks passed (warnings allowed)
- `1`: One or more critical checks failed

## Dependencies

Key dependencies:
- `github.com/spf13/cobra`: CLI framework
- `github.com/charmbracelet/lipgloss`: Terminal styling
- `gopkg.in/yaml.v3`: YAML parsing

## CI/CD Integration

- **GitHub Action**: `scripts/action.yml`
- **Pre-commit**: Supported via pre-commit hooks
- **Docker**: `Dockerfile` available for containerized runs

## Notes for Agents

1. **Always write tests** when adding new checks or modifying existing ones
2. **Maintain backward compatibility** with existing `.a2.yaml` configurations
3. **Follow Go conventions**: Use `go fmt` and `go vet` before committing
4. **Check severity**: Understand when to use Pass/Warn/Fail
5. **Parallel execution**: Be aware that checks run concurrently by default
6. **Veto system**: Critical checks (Fail) stop execution in sequential mode
7. **Path handling**: Use `safepath` utilities for safe path operations
8. **Error messages**: Provide clear, actionable error messages

## Testing Checklist

Before submitting changes:
- [ ] All tests pass: `go test ./...`
- [ ] Code formatted: `go fmt ./...`
- [ ] No vet issues: `go vet ./...`
- [ ] Coverage maintained or improved
- [ ] Documentation updated if needed
- [ ] Configuration changes tested with `.a2.yaml`

