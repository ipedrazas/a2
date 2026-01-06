# A2 - Application Analysis

A2 is a code quality checker for Go projects. It runs a suite of checks against your repository and provides a health score with actionable recommendations.

## Features

- **8 Built-in Checks**: Module validation, build, tests, coverage, formatting, vet, file existence, vulnerabilities
- **Veto System**: Critical checks (build, tests) stop execution on failure
- **Pretty Output**: Colored terminal output with recommendations
- **JSON Output**: Machine-readable format for CI/CD integration
- **Configurable**: `.a2.yaml` for thresholds, disabled checks, and custom checks
- **Extensible**: Add your own checks via external binaries
- **CI/CD Ready**: GitHub Action and pre-commit hook support

## Installation

```bash
go install github.com/ipedrazas/a2@latest
```

## Usage

```bash
# Run checks on current directory
a2 check

# Run checks on specific path
a2 check /path/to/project

# JSON output for CI/CD
a2 check --format json
```

### Sample Output

```
A2 Analysis: myproject
─────────────────────────────────────

✓ PASS Go Module
    Module: github.com/user/myproject (Go 1.23)
✓ PASS Build
    Build successful
✓ PASS Unit Tests
    All tests passed
✓ PASS Required Files
    All required files present
✓ PASS Code Formatting
    All Go files are properly formatted
✓ PASS Go Vet
    No issues found
! WARN Test Coverage
    Coverage 65.0% is below threshold 80.0%
✓ PASS Vulnerabilities
    No known vulnerabilities found

─────────────────────────────────────

STATUS: ⚠ NEEDS ATTENTION

Score: 7/8 checks passed (88%)

Recommendations:
→ Add more tests to improve coverage
```

## Built-in Checks

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Go Module | `go_mod` | Fail | go.mod exists and has valid Go version |
| Build | `build` | Fail | `go build ./...` succeeds |
| Unit Tests | `tests` | Fail | `go test ./...` passes |
| Required Files | `file_exists` | Warn | README.md, LICENSE exist |
| Code Formatting | `gofmt` | Warn | Code is properly formatted |
| Go Vet | `govet` | Warn | No `go vet` issues |
| Test Coverage | `coverage` | Warn | Coverage >= threshold (default 80%) |
| Vulnerabilities | `deps` | Warn | No known vulns (requires govulncheck) |

**Severity Levels:**
- **Fail**: Critical check - stops execution immediately (veto power)
- **Warn**: Non-critical - reported but doesn't stop execution
- **Pass**: Check passed

## Configuration

Create a `.a2.yaml` file in your project root:

```yaml
# Coverage threshold (default: 80)
coverage:
  threshold: 70

# Required files to check
files:
  required:
    - README.md
    - LICENSE
    - CONTRIBUTING.md

# Disable specific checks
checks:
  disabled:
    - deps  # Skip vulnerability check

# Custom external checks
external:
  - id: lint
    name: Golangci-lint
    command: golangci-lint
    args: ["run", "./..."]
    severity: warn

  - id: security
    name: Security Scan
    command: gosec
    args: ["./..."]
    severity: fail
```

## External Checks

A2 supports external check binaries. The protocol is simple:

- **Exit code 0**: Pass
- **Exit code 1**: Warning
- **Exit code 2+**: Fail

Output can be plain text or JSON:

```json
{
  "message": "Found 3 issues",
  "status": "warn"
}
```

## GitHub Action

```yaml
name: A2 Analysis
on: [push, pull_request]

jobs:
  a2:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run A2
        uses: ipedrazas/a2@v1
        with:
          path: '.'
          format: 'pretty'
          fail-on-warning: 'false'
```

### Inputs

| Input | Description | Default |
|-------|-------------|---------|
| `path` | Directory to analyze | `.` |
| `format` | Output format (pretty/json) | `pretty` |
| `fail-on-warning` | Fail if warnings exist | `false` |

### Outputs

| Output | Description |
|--------|-------------|
| `score` | Percentage of checks passed |
| `passed` | Number of passed checks |
| `total` | Total checks run |
| `success` | Whether critical checks passed |

## Pre-commit Hook

Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/ipedrazas/a2
    rev: v1.0.0
    hooks:
      - id: a2
```

## Docker

```bash
# Build image
docker build -t a2 .

# Run checks
docker run -v $(pwd):/workspace a2 check
```

## Exit Codes

- `0`: All checks passed (warnings allowed)
- `1`: One or more critical checks failed

## License

MIT License - see [LICENSE](LICENSE)
