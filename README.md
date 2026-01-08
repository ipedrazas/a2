# A2 - Application Analysis

A2 is a multi-language code quality checker. It auto-detects project language(s), runs a suite of checks, and provides a health score with recommendations that you can give to your Coding Agent to improve your application.

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

It's up to you to decide which checks make sense for you and your project. `a2` allows you to configure what and what not to run. For example, if you run `a2 check` in this repo, you will get a response like this one:

```
%> a2 check --profile=cli
A2 Analysis: a2

─────────────────────────────────────
Languages: go

✓ PASS Go Module
    Module: github.com/ipedrazas/a2 (Go 1.23)
✓ PASS Go Build
    Build successful
✓ PASS Go Tests
    All tests passed
✓ PASS Go Race Detection
    No data races detected
✓ PASS Go Format
    All Go files are properly formatted
✓ PASS Go Vet
    No issues found
✓ PASS Go Coverage
    Coverage: 57.2%
✓ PASS Go Vulnerabilities
    No known vulnerabilities found
✓ PASS Go Complexity
    No functions exceed complexity threshold (15)
✓ PASS Required Files
    All required files present
✓ PASS Container Ready
    Dockerfile found
✓ PASS CI Pipeline
    GitHub Actions configured
✓ PASS Secrets Detection
    Secret scanning configured: Gitleaks
✓ PASS Environment Config
    Environment config: .env in .gitignore
✓ PASS SAST Security Scanning
    SAST configured: gosec
✓ PASS Changelog
    CHANGELOG.md found
✓ PASS Contributing Guidelines
    Found: CONTRIBUTING.md
✓ PASS Editor Config
    .editorconfig configured

─────────────────────────────────────

STATUS: ✓ ALL CHECKS PASSED

Score: 23/23 checks passed (100%)

Maturity: Production-Ready
   All checks pass, ready for production deployment
```

As you can see `a2` will tell you the maturity level of the project based on the results of the checks. While it's true that you could use `Claude Code` or any other coding agent to do the same, I'd rather use a deterministic approach because saying that something is "Production-ready" depends on what you consider what production level is.

## Features

- **7 Languages Supported**: Go, Python, Node.js, TypeScript, Java, Rust (auto-detected or explicit)
- **70+ Built-in Checks**: Build, tests, coverage, formatting, linting, type checking, vulnerabilities, and more
- **Application Profiles**: CLI, API, Library, Desktop - tailor checks to your app type
- **Maturity Targets**: PoC or Production - control check strictness
- **Maturity Assessment**: Automatic scoring with Production-Ready, Mature, Development, or PoC levels
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

# Use application profile
a2 check --profile=cli      # Command-line tool
a2 check --profile=api      # Web service/API
a2 check --profile=library  # Reusable package
a2 check --profile=desktop  # Desktop application

# Use maturity target
a2 check --target=poc       # Minimal checks for early development
a2 check --target=production # All checks (default)

# Combine profile and target
a2 check --profile=cli --target=poc

# JSON output for CI/CD
a2 check --format json

# Skip specific checks
a2 check --skip=license,k8s

# List available options
a2 profiles  # List application profiles
a2 targets   # List maturity targets
```

## Application Profiles

Profiles define which checks are relevant for your application type. Use `--profile` to select one:

| Profile | Description | Skipped Checks |
|---------|-------------|----------------|
| `cli` | Command-line tools | health, k8s, metrics, api_docs, integration, shutdown, errors, e2e, tracing |
| `api` | Web services/APIs | e2e (uses integration tests instead) |
| `library` | Reusable packages | dockerfile, health, k8s, shutdown, metrics, errors, integration, tracing, e2e, api_docs |
| `desktop` | Desktop applications | health, k8s, api_docs, tracing, metrics, shutdown |

```bash
# List all profiles with details
a2 profiles
```

## Maturity Targets

Targets control the strictness level of checks based on your project stage. Use `--target` to select one:

| Target | Description | Effect |
|--------|-------------|--------|
| `poc` | Proof of Concept | Skips non-critical checks (license, sast, coverage, deps, etc.) |
| `production` | Production-ready (default) | All checks enabled |

```bash
# List all targets with details
a2 targets

# Combine with profiles
a2 check --profile=api --target=poc  # API in early development
```

## Maturity Assessment

A2 automatically assesses your project's maturity level based on check results:

| Level | Criteria |
|-------|----------|
| **Production-Ready** | 100% score, 0 failures, 0 warnings |
| **Mature** | ≥80% score, 0 failures |
| **Development** | ≥60% score, ≤2 failures |
| **Proof of Concept** | <60% score or >2 failures |

## Supported Languages

A2 supports 6 programming languages with auto-detection:

| Language | Indicator Files | Checks |
|----------|-----------------|--------|
| **Go** | `go.mod`, `go.sum` | 10 checks |
| **Python** | `pyproject.toml`, `setup.py`, `requirements.txt`, `Pipfile` | 10 checks |
| **Node.js** | `package.json`, `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml` | 9 checks |
| **TypeScript** | `tsconfig.json` | 9 checks |
| **Java** | `pom.xml`, `build.gradle`, `build.gradle.kts` | 8 checks |
| **Rust** | `Cargo.toml` | 8 checks |

Use `--lang` flag or `language.explicit` config to override auto-detection.

## Built-in Checks

### Go Checks (10)

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Go Module | `go:module` | Fail | go.mod exists and has valid Go version |
| Go Build | `go:build` | Fail | `go build ./...` succeeds |
| Go Tests | `go:tests` | Fail | `go test ./...` passes |
| Go Race Detection | `go:race` | Warn | No data races (`go test -race`) |
| Go Format | `go:format` | Warn | Code is properly formatted (`gofmt`) |
| Go Vet | `go:vet` | Warn | No `go vet` issues |
| Go Coverage | `go:coverage` | Warn | Coverage >= threshold (default 80%) |
| Go Vulnerabilities | `go:deps` | Warn | No known vulns (`govulncheck`) |
| Go Complexity | `go:cyclomatic` | Warn | No functions exceed complexity threshold |
| Go Logging | `go:logging` | Warn | Uses structured logging library |

### Python Checks (10)

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Python Project | `python:project` | Fail | pyproject.toml or setup.py exists |
| Python Build | `python:build` | Fail | Dependencies install successfully |
| Python Tests | `python:tests` | Fail | pytest or unittest passes |
| Python Format | `python:format` | Warn | Code formatted (black/ruff) |
| Python Lint | `python:lint` | Warn | No lint issues (ruff/flake8/pylint) |
| Python Type Check | `python:type` | Warn | No type errors (mypy/pyright) |
| Python Coverage | `python:coverage` | Warn | Coverage >= threshold (default 80%) |
| Python Vulnerabilities | `python:deps` | Warn | No known vulns (pip-audit/safety) |
| Python Complexity | `python:complexity` | Warn | No functions exceed complexity threshold |
| Python Logging | `python:logging` | Warn | Uses logging library, not print() |

### Node.js Checks (9)

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Node Project | `node:project` | Fail | package.json exists with name/version |
| Node Build | `node:build` | Fail | Dependencies install successfully |
| Node Tests | `node:tests` | Fail | Tests pass (jest/vitest/mocha) |
| Node Format | `node:format` | Warn | Code formatted (prettier/biome) |
| Node Lint | `node:lint` | Warn | No lint issues (eslint/biome/oxlint) |
| Node Type Check | `node:type` | Warn | No type errors (tsc --noEmit) |
| Node Coverage | `node:coverage` | Warn | Coverage >= threshold (default 80%) |
| Node Vulnerabilities | `node:deps` | Warn | No known vulns (npm/yarn/pnpm audit) |
| Node Logging | `node:logging` | Warn | Uses logging library, warns on console.log |

### TypeScript Checks (9)

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| TypeScript Project | `typescript:project` | Fail | tsconfig.json exists |
| TypeScript Build | `typescript:build` | Fail | Build script or `tsc --noEmit` succeeds |
| TypeScript Tests | `typescript:tests` | Fail | Tests pass (jest/vitest/mocha) |
| TypeScript Format | `typescript:format` | Warn | Code formatted (prettier/biome/dprint) |
| TypeScript Lint | `typescript:lint` | Warn | No lint issues (eslint/biome/oxlint) |
| TypeScript Type Check | `typescript:type` | Fail | No type errors (`tsc --noEmit`) |
| TypeScript Coverage | `typescript:coverage` | Warn | Coverage tools configured |
| TypeScript Vulnerabilities | `typescript:deps` | Warn | No known vulns (npm/yarn/pnpm audit) |
| TypeScript Logging | `typescript:logging` | Warn | Uses logging library (winston/pino/tslog) |

### Java Checks (8)

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Java Project | `java:project` | Fail | pom.xml or build.gradle exists |
| Java Build | `java:build` | Fail | Maven/Gradle build succeeds |
| Java Tests | `java:tests` | Fail | JUnit/TestNG tests pass |
| Java Format | `java:format` | Warn | Formatter configured (Spotless/google-java-format) |
| Java Lint | `java:lint` | Warn | Linter configured (Checkstyle/SpotBugs/PMD) |
| Java Coverage | `java:coverage` | Warn | JaCoCo coverage >= threshold |
| Java Vulnerabilities | `java:deps` | Warn | Dependency scanning configured |
| Java Logging | `java:logging` | Warn | Uses SLF4J/Logback/Log4j2, not System.out |

### Rust Checks (8)

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Rust Project | `rust:project` | Fail | Cargo.toml exists |
| Rust Build | `rust:build` | Fail | `cargo check` succeeds |
| Rust Tests | `rust:tests` | Fail | `cargo test` passes |
| Rust Format | `rust:format` | Warn | Code formatted (`cargo fmt --check`) |
| Rust Lint | `rust:lint` | Warn | No Clippy warnings (`cargo clippy`) |
| Rust Coverage | `rust:coverage` | Warn | Coverage tools configured (tarpaulin/llvm-cov) |
| Rust Vulnerabilities | `rust:deps` | Warn | No known vulns (`cargo audit`) |
| Rust Logging | `rust:logging` | Warn | Uses tracing/log crate, not println! |

### Common Checks (23)

These checks apply to all projects regardless of language:

| Check | ID | Severity | Description |
|-------|-----|----------|-------------|
| Required Files | `file_exists` | Warn | README.md, LICENSE exist |
| Dockerfile | `common:dockerfile` | Warn | Container configuration present |
| CI Pipeline | `common:ci` | Warn | CI/CD configuration present |
| Health Check | `common:health` | Warn | Health endpoint configured |
| Kubernetes | `common:k8s` | Warn | K8s manifests, Helm, or Kustomize |
| Pre-commit Hooks | `common:precommit` | Warn | Git hooks configured |
| Changelog | `common:changelog` | Warn | CHANGELOG.md present |
| Secrets Detection | `common:secrets` | Warn | Secret scanning configured |
| Environment Config | `common:env` | Warn | .env handling configured |
| License Compliance | `common:license` | Warn | License scanning configured |
| SAST | `common:sast` | Warn | Static analysis security testing |
| API Documentation | `common:api_docs` | Warn | OpenAPI/Swagger docs present |
| Integration Tests | `common:integration` | Warn | Integration test infrastructure |
| Metrics | `common:metrics` | Warn | Prometheus/OpenTelemetry metrics |
| Error Tracking | `common:errors` | Warn | Sentry/Rollbar/Bugsnag configured |
| Graceful Shutdown | `common:shutdown` | Warn | Signal handling implemented |
| Contributing | `common:contributing` | Warn | CONTRIBUTING.md, PR/issue templates |
| E2E Tests | `common:e2e` | Warn | Cypress/Playwright/Selenium configured |
| Distributed Tracing | `common:tracing` | Warn | OpenTelemetry/Jaeger/Datadog |
| Database Migrations | `common:migrations` | Warn | Migration tool configured |
| Config Validation | `common:config_validation` | Warn | Config validation library used |
| Retry Logic | `common:retry` | Warn | Retry/resilience library used |
| Editor Config | `common:editorconfig` | Warn | .editorconfig present |

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
language:
  python:
    package_manager: auto  # auto, pip, poetry, pipenv
    test_runner: auto      # auto, pytest, unittest
    formatter: auto        # auto, black, ruff
    linter: auto           # auto, ruff, flake8, pylint
    coverage_threshold: 75

files:
  required:
    - README.md
    - LICENSE
    - pyproject.toml

checks:
  disabled:
    - python:deps

external:
  - id: security
    name: Security Scan
    command: bandit
    args: ["-r", "src/"]
    severity: warn
```

### Example: Node.js/TypeScript Project

```yaml
language:
  node:
    package_manager: auto  # auto, npm, yarn, pnpm, bun
    test_runner: auto      # auto, jest, vitest, mocha
    formatter: auto        # auto, prettier, biome
    linter: auto           # auto, eslint, biome, oxlint
    coverage_threshold: 80
  typescript:
    coverage_threshold: 80

files:
  required:
    - README.md
    - LICENSE
    - package.json

checks:
  disabled:
    - node:deps
```

### Multi-Language Project (Monorepo)

```yaml
# Explicit language selection (overrides auto-detect)
language:
  explicit:
    - go
    - python
    - typescript
  go:
    coverage_threshold: 80
  python:
    coverage_threshold: 70
    linter: ruff
  typescript:
    coverage_threshold: 75

files:
  required:
    - README.md
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
          profile: 'api'
          target: 'production'
          fail-on-warning: 'false'
```

### Inputs

| Input | Description | Default |
|-------|-------------|---------|
| `path` | Directory to analyze | `.` |
| `format` | Output format (pretty/json) | `pretty` |
| `profile` | Application profile (cli/api/library/desktop) | - |
| `target` | Maturity target (poc/production) | `production` |
| `fail-on-warning` | Fail if warnings exist | `false` |

### Outputs

| Output | Description |
|--------|-------------|
| `score` | Percentage of checks passed |
| `passed` | Number of passed checks |
| `total` | Total checks run |
| `success` | Whether critical checks passed |
| `maturity` | Maturity level assessment |

## Pre-commit Hook

Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/ipedrazas/a2
    rev: v1.0.0
    hooks:
      - id: a2
        args: ['--profile=library']  # Optional: specify profile
```

## Docker

```bash
# Build image
docker build -t a2 .

# Run checks
docker run -v $(pwd):/workspace a2 check

# Run with profile
docker run -v $(pwd):/workspace a2 check --profile=api
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
