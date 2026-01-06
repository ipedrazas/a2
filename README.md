# A2 - Application Analysis

A2 is a multi-language code quality checker. It auto-detects project language(s), runs a suite of checks, and provides a health score with actionable recommendations.

## Features

- **Multi-Language Support**: Go and Python (auto-detected or explicit)
- **15+ Built-in Checks**: Build, tests, coverage, formatting, linting, vulnerabilities
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
# Run checks on current directory (auto-detects language)
a2 check

# Run checks on specific path
a2 check /path/to/project

# Explicit language selection
a2 check --lang python
a2 check --lang go,python

# JSON output for CI/CD
a2 check --format json
```

### Sample Output

```
A2 Analysis: myproject
Detected: go
─────────────────────────────────────

✓ PASS Go Module
    Module: github.com/user/myproject (Go 1.23)
✓ PASS Go Build
    Build successful
✓ PASS Go Tests
    All tests passed
✓ PASS Required Files
    All required files present
✓ PASS Go Format
    All Go files are properly formatted
✓ PASS Go Vet
    No issues found
! WARN Go Coverage
    Coverage 65.0% is below threshold 80.0%
✓ PASS Go Vulnerabilities
    No known vulnerabilities found

─────────────────────────────────────

STATUS: ⚠ NEEDS ATTENTION

Score: 7/8 checks passed (88%)

Recommendations:
→ Add more tests to improve coverage
```

## Built-in Checks

### Go Checks

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Go Module | `go:module` | Fail | go.mod exists and has valid Go version |
| Go Build | `go:build` | Fail | `go build ./...` succeeds |
| Go Tests | `go:tests` | Fail | `go test ./...` passes |
| Go Format | `go:format` | Warn | Code is properly formatted |
| Go Vet | `go:vet` | Warn | No `go vet` issues |
| Go Coverage | `go:coverage` | Warn | Coverage >= threshold (default 80%) |
| Go Vulnerabilities | `go:deps` | Warn | No known vulns (requires govulncheck) |

### Python Checks

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Python Project | `python:project` | Fail | pyproject.toml or setup.py exists |
| Python Build | `python:build` | Fail | Dependencies install successfully |
| Python Tests | `python:tests` | Fail | pytest or unittest passes |
| Python Format | `python:format` | Warn | Code formatted (ruff/black) |
| Python Lint | `python:lint` | Warn | No lint issues (ruff/flake8/pylint) |
| Python Coverage | `python:coverage` | Warn | Coverage >= threshold (default 80%) |
| Python Vulnerabilities | `python:deps` | Warn | No known vulns (pip-audit/safety) |

### Common Checks

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Required Files | `file_exists` | Warn | README.md, LICENSE exist |

**Severity Levels:**
- **Fail**: Critical check - stops execution immediately (veto power)
- **Warn**: Non-critical - reported but doesn't stop execution
- **Pass**: Check passed

## Configuration

Create a `.a2.yaml` file in your project root.

### Example: Go Project

```yaml
# Language settings
language:
  go:
    coverage_threshold: 80

# Coverage threshold (legacy, also works)
coverage:
  threshold: 80

# Required files to check
files:
  required:
    - README.md
    - LICENSE
    - CONTRIBUTING.md

# Disable specific checks
checks:
  disabled:
    - go:deps  # Skip vulnerability check

# Execution options
execution:
  parallel: true  # Run checks concurrently (default)

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

### Example: Python Project

```yaml
# Language settings
language:
  python:
    package_manager: auto  # auto, pip, poetry, pipenv
    test_runner: auto      # auto, pytest, unittest
    formatter: auto        # auto, black, ruff
    linter: auto           # auto, ruff, flake8, pylint
    coverage_threshold: 75

# Required files
files:
  required:
    - README.md
    - LICENSE
    - pyproject.toml

# Disable specific checks
checks:
  disabled:
    - python:deps  # Skip vulnerability scan

# Custom external checks
external:
  - id: typecheck
    name: Type Check
    command: mypy
    args: ["src/"]
    severity: warn

  - id: security
    name: Security Scan
    command: bandit
    args: ["-r", "src/"]
    severity: warn
```

### Multi-Language Project (Monorepo)

```yaml
# Explicit language selection (overrides auto-detect)
language:
  explicit:
    - go
    - python
  go:
    coverage_threshold: 80
  python:
    coverage_threshold: 70
    linter: ruff
    formatter: ruff

files:
  required:
    - README.md
```

## Language Detection

A2 auto-detects languages based on indicator files:

| Language | Indicator Files |
|----------|----------------|
| Go | `go.mod`, `go.sum` |
| Python | `pyproject.toml`, `setup.py`, `requirements.txt`, `Pipfile`, `poetry.lock` |

Use `--lang` flag or `language.explicit` config to override auto-detection.

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

## Backward Compatibility

Old check IDs are aliased to new language-prefixed IDs:

| Old ID | New ID |
|--------|--------|
| `go_mod` | `go:module` |
| `build` | `go:build` |
| `tests` | `go:tests` |
| `gofmt` | `go:format` |
| `govet` | `go:vet` |
| `coverage` | `go:coverage` |
| `deps` | `go:deps` |

## License

MIT License - see [LICENSE](LICENSE)
