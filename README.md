# A2 - Application Analysis

A2 is a multi-language code quality checker. It auto-detects project language(s), runs a suite of checks, and provides a health score with actionable recommendations.

Because of the amount of new code and projects created with the rise of `vibecoding`, I needed a way to assess the level of maturity of any project.

`a2` should help you to understand if a project is production-ready or not, or in which level of the application life cycle it is: PoC, alpha, beta, prod ready?

Checks are configurable and they cover big themes like:

- Security best practices
- Proper documentation
- Comprehensive testing
- Observability instrumentation
- Production-ready configuration
- Clean architecture patterns
- Solid development workflow

It'a up to you to decide which checks make sense for you and your project. `a2` allows you to configure what and what not to run. For example, if you run `a2 check` in this repo, you will get a response like this one:

```
a2 check
A2 Analysis: a2
               
─────────────────────────────────────
Languages: go

✓ PASS Go Module
    Module: github.com/ipedrazas/a2 (Go 1.25.5)
✓ PASS Go Build
    Build successful
✓ PASS Go Tests
    No test files found
✓ PASS Go Race Detection
    No test files to check for races
✓ PASS Go Format
    All Go files are properly formatted
✓ PASS Go Vet
    No issues found
✓ PASS Go Coverage
    Coverage: 60.0%
✓ PASS Go Vulnerabilities
    No known vulnerabilities found
✓ PASS Go Complexity
    No functions exceed complexity threshold (15)
✓ PASS Required Files
    All required files present
✓ PASS Container Ready
    Dockerfile found (consider adding .dockerignore)
✓ PASS CI Pipeline
    GitHub Actions, Taskfile configured
✓ PASS Secrets Detection
    Secret scanning configured: Gitleaks, pre-commit hook
✓ PASS Environment Config
    Environment config: .env in .gitignore
✓ PASS SAST Security Scanning
    SAST configured: gosec, gosec (CI), Bandit (CI), Semgrep (CI)
✓ PASS Changelog
    CHANGELOG.md found (Keep a Changelog format)
✓ PASS Security Scan
    Results:        
                    
                    
    Summary:        
      Gosec  : dev  
      Files  : 72   
      Lines  : 11034
      Nosec  : 7    
      Issues : 0    
✓ PASS Golangci-lint
✓ PASS Secret Detection
    ○                                                   
        │╲                                              
        │ ○                                             
        ○ ░                                             
        ░    gitleaks                                   
                                                        
    2:35AM INF scanned ~1549632 bytes (1.55 MB) in 874ms
    2:35AM INF no leaks found                           
✓ PASS Error Check
✓ PASS Static Analysis

─────────────────────────────────────

STATUS: ✓ ALL CHECKS PASSED

Score: 21/21 checks passed (100%)

Maturity: Production-Ready
   All checks pass, ready for production deployment
```

As you can see `a2` will tell you the maturity level of the project based on the results of the checks. While it's true that you could use `Claude Code` or any other coding agent to do the same, I'd rather use a deterministic approach because saying that something is "Production-ready" depends on what you consider what production level is, so asking an Agent becomes tricky because you have to consistently inform the agent what do you consider "Production" or "Alpha" o...

You can still achieve this by writing a prompt, but I like tinkering with clis 🤷‍♂️

## Features

- **Multi-Language Support**: Go, Python, Java, and Node.js (auto-detected or explicit)
- **50+ Built-in Checks**: Build, tests, coverage, formatting, linting, type checking, vulnerabilities
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

# Skip specific checks
a2 check --skip=license,k8s
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
| Go Race Detection | `go:race` | Warn | No data races (`go test -race`) |
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
| Python Type Check | `python:type` | Warn | No type errors (mypy) |
| Python Coverage | `python:coverage` | Warn | Coverage >= threshold (default 80%) |
| Python Vulnerabilities | `python:deps` | Warn | No known vulns (pip-audit/safety) |

### Node.js Checks

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Node Project | `node:project` | Fail | package.json exists with name/version |
| Node Build | `node:build` | Fail | Dependencies install successfully |
| Node Tests | `node:tests` | Fail | Tests pass (jest/vitest/mocha/npm test) |
| Node Format | `node:format` | Warn | Code formatted (prettier/biome) |
| Node Lint | `node:lint` | Warn | No lint issues (eslint/biome/oxlint) |
| TypeScript Type Check | `node:type` | Warn | No type errors (tsc --noEmit) |
| Node Coverage | `node:coverage` | Warn | Coverage >= threshold (default 80%) |
| Node Vulnerabilities | `node:deps` | Warn | No known vulns (npm/yarn/pnpm audit) |

### Common Checks

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Required Files | `file_exists` | Warn | README.md, LICENSE exist |

> **See [docs/CHECKS.md](docs/CHECKS.md) for detailed documentation** on all checks, including configuration options, auto-detection logic, and pass/warn/fail conditions.

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
  - id: security
    name: Security Scan
    command: bandit
    args: ["-r", "src/"]
    severity: warn
```

### Example: Node.js Project

```yaml
# Language settings
language:
  node:
    package_manager: auto  # auto, npm, yarn, pnpm, bun
    test_runner: auto      # auto, jest, vitest, mocha, npm-test
    formatter: auto        # auto, prettier, biome
    linter: auto           # auto, eslint, biome, oxlint
    coverage_threshold: 80

# Required files
files:
  required:
    - README.md
    - LICENSE
    - package.json

# Disable specific checks
checks:
  disabled:
    - node:deps  # Skip vulnerability scan
```

### Multi-Language Project (Monorepo)

```yaml
# Explicit language selection (overrides auto-detect)
language:
  explicit:
    - go
    - python
    - node
  go:
    coverage_threshold: 80
  python:
    coverage_threshold: 70
    linter: ruff
    formatter: ruff
  node:
    coverage_threshold: 75
    linter: eslint
    formatter: prettier

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
| Node.js | `package.json`, `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`, `bun.lockb` |

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
