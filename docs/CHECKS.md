# A2 Checks Reference

This document describes all checks available in A2, organized by language.

## Table of Contents

- [Go Checks](#go-checks)
- [Python Checks](#python-checks)
- [Node.js Checks](#nodejs-checks)
- [Common Checks](#common-checks)
- [External Checks](#external-checks)
- [Configuration Reference](#configuration-reference)

---

## Go Checks

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

### go:module

Verifies that `go.mod` exists and contains a valid Go version directive.

**Status:**
- **Pass**: go.mod exists with valid syntax and Go version
- **Warn**: go.mod exists but missing Go version directive
- **Fail**: go.mod not found or invalid syntax

### go:build

Runs `go build ./...` to verify the project compiles successfully.

**Status:**
- **Pass**: Build completes without errors
- **Fail**: Compilation errors

### go:tests

Runs `go test ./...` to execute all test packages.

**Status:**
- **Pass**: All tests pass or no test files found
- **Fail**: One or more tests fail

### go:race

Runs `go test -race -short ./...` to detect data races in concurrent code.

**Status:**
- **Pass**: No race conditions detected or no test files
- **Warn**: Race conditions detected or tests fail during race detection

### go:format

Runs `gofmt -l` to check if all Go files are properly formatted.

**Status:**
- **Pass**: All files properly formatted
- **Warn**: Unformatted files found

**Fix:** Run `gofmt -w .`

### go:vet

Runs `go vet ./...` to find suspicious constructs and potential bugs.

**Status:**
- **Pass**: No issues found
- **Warn**: Issues detected

### go:coverage

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

### go:deps

Scans for known vulnerabilities in Go dependencies using `govulncheck`.

**Requirements:** Install with `go install golang.org/x/vuln/cmd/govulncheck@latest`

**Status:**
- **Pass**: No vulnerabilities found or govulncheck not installed
- **Warn**: Vulnerabilities detected

---

## Python Checks

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `python:project` | Python Project | Yes | 100 | Verifies project config exists |
| `python:build` | Python Build | Yes | 110 | Validates package manager setup |
| `python:tests` | Python Tests | Yes | 120 | Runs tests with pytest/unittest |
| `python:format` | Python Format | No | 200 | Checks formatting with ruff/black |
| `python:lint` | Python Lint | No | 210 | Lints code with ruff/flake8/pylint |
| `python:type` | Python Type Check | No | 215 | Type checks with mypy |
| `python:coverage` | Python Coverage | No | 220 | Measures test coverage |
| `python:deps` | Python Vulnerabilities | No | 230 | Scans for vulnerabilities |

### python:project

Verifies that a Python project configuration file exists.

**Detection priority:**
1. `pyproject.toml` (preferred)
2. `setup.py` (legacy, warns)
3. `requirements.txt` (minimal, warns)

**Status:**
- **Pass**: pyproject.toml found
- **Warn**: setup.py or requirements.txt found
- **Fail**: No project configuration found

### python:build

Validates dependencies using the detected package manager.

**Configuration:**
```yaml
language:
  python:
    package_manager: auto  # auto, pip, poetry, pipenv
```

**Auto-detection:**
- `poetry.lock` → poetry (`poetry check`)
- `Pipfile.lock` or `Pipfile` → pipenv (`pipenv check`)
- Otherwise → pip

**Status:**
- **Pass**: Validation succeeds or tool not installed
- **Fail**: Validation fails

### python:tests

Runs Python tests using the detected test runner.

**Configuration:**
```yaml
language:
  python:
    test_runner: auto  # auto, pytest, unittest
```

**Auto-detection:**
- `pytest.ini`, `conftest.py`, or `[tool.pytest]` in pyproject.toml → pytest
- Otherwise → pytest (default)

**Status:**
- **Pass**: All tests pass or no tests found
- **Fail**: Tests fail

### python:format

Checks Python code formatting.

**Configuration:**
```yaml
language:
  python:
    formatter: auto  # auto, black, ruff
```

**Auto-detection:**
- `ruff.toml` or `.ruff.toml` → ruff
- `[tool.black]` in pyproject.toml → black
- Tries ruff first, falls back to black

**Status:**
- **Pass**: All files properly formatted
- **Warn**: Files need formatting

**Fix:** Run `ruff format .` or `black .`

### python:lint

Runs Python linting to check for code quality issues.

**Configuration:**
```yaml
language:
  python:
    linter: auto  # auto, ruff, flake8, pylint
```

**Auto-detection:**
- `ruff.toml`, `.ruff.toml`, or `[tool.ruff]` → ruff
- `.flake8` or `[tool.flake8]` → flake8
- `.pylintrc` or `[tool.pylint]` → pylint
- Tries ruff first, falls back to flake8

**Status:**
- **Pass**: No issues found
- **Warn**: Linting issues detected

### python:type

Runs mypy for static type checking. Only activates for typed Python projects.

**Typed project detection:**
- `mypy.ini` or `.mypy.ini`
- `py.typed` marker (PEP 561)
- `[mypy]` section in setup.cfg
- `[tool.mypy]` in pyproject.toml

**Status:**
- **Pass**: No type errors, not a typed project, or mypy not installed
- **Warn**: Type errors found

**Fix:** Run `mypy .`

### python:coverage

Measures test coverage using pytest-cov.

**Configuration:**
```yaml
language:
  python:
    coverage_threshold: 80  # Default: 80%
```

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold or pytest-cov not installed

### python:deps

Scans for security vulnerabilities in Python dependencies.

**Tools (tried in order):**
1. `pip-audit`
2. `safety`

**Status:**
- **Pass**: No vulnerabilities found or scanner not installed
- **Warn**: Vulnerabilities detected

---

## Node.js Checks

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `node:project` | Node Project | Yes | 100 | Verifies package.json exists |
| `node:build` | Node Build | Yes | 110 | Validates dependencies |
| `node:tests` | Node Tests | Yes | 120 | Runs tests |
| `node:format` | Node Format | No | 200 | Checks formatting |
| `node:lint` | Node Lint | No | 210 | Lints code |
| `node:type` | TypeScript Type Check | No | 215 | Type checks TypeScript |
| `node:coverage` | Node Coverage | No | 220 | Measures test coverage |
| `node:deps` | Node Vulnerabilities | No | 230 | Scans for vulnerabilities |

### node:project

Verifies that `package.json` exists and contains required fields.

**Status:**
- **Pass**: package.json with valid name and version
- **Warn**: Missing version field
- **Fail**: Missing package.json or missing name field

### node:build

Validates Node.js dependencies using the detected package manager.

**Configuration:**
```yaml
language:
  node:
    package_manager: auto  # auto, npm, yarn, pnpm, bun
```

**Auto-detection:**
- `pnpm-lock.yaml` → pnpm
- `yarn.lock` → yarn
- `bun.lockb` → bun
- `package-lock.json` → npm

**Commands:**
- npm: `npm ci --dry-run` or `npm install --dry-run`
- yarn: `yarn install --check-files`
- pnpm: `pnpm install --frozen-lockfile --dry-run`
- bun: `bun install --dry-run`

**Status:**
- **Pass**: Validation succeeds or package manager not installed
- **Fail**: Validation fails

### node:tests

Runs Node.js tests using the detected test runner.

**Configuration:**
```yaml
language:
  node:
    test_runner: auto  # auto, jest, vitest, mocha, npm-test
```

**Auto-detection:**
- `jest.config.*` → jest
- `vitest.config.*` → vitest
- `.mocharc.*` → mocha
- Checks devDependencies
- Falls back to `npm test`

**Status:**
- **Pass**: All tests pass or no test script defined
- **Fail**: Tests fail

### node:format

Checks code formatting using prettier or biome.

**Configuration:**
```yaml
language:
  node:
    formatter: auto  # auto, prettier, biome
```

**Auto-detection:**
- `.prettierrc*` or `prettier.config.*` → prettier
- `biome.json*` → biome
- Checks devDependencies

**Status:**
- **Pass**: All files formatted
- **Warn**: Files need formatting

**Fix:** Run `npx prettier --write .` or `npx @biomejs/biome format --write .`

### node:lint

Runs linting using eslint, biome, or oxlint.

**Configuration:**
```yaml
language:
  node:
    linter: auto  # auto, eslint, biome, oxlint
```

**Auto-detection:**
- `.eslintrc*` or `eslint.config.*` → eslint
- `biome.json*` → biome
- `oxlint.json` or `.oxlintrc.json` → oxlint
- Checks devDependencies

**Status:**
- **Pass**: No linting issues
- **Warn**: Linting issues found

### node:type

Runs TypeScript compiler for type checking. Only activates for TypeScript projects.

**TypeScript project detection:**
- `tsconfig.json` or variants (`tsconfig.base.json`, etc.)
- TypeScript in devDependencies or dependencies

**Status:**
- **Pass**: No type errors or not a TypeScript project
- **Warn**: Type errors found

**Fix:** Run `npx tsc --noEmit`

### node:coverage

Measures test coverage using jest, vitest, c8, or nyc.

**Configuration:**
```yaml
language:
  node:
    coverage_threshold: 80  # Default: 80%
```

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold or tools not installed

### node:deps

Scans for security vulnerabilities in Node.js dependencies.

**Commands (based on package manager):**
- npm: `npm audit --json`
- yarn: `yarn audit --json`
- pnpm: `pnpm audit --json`
- bun: (skipped - no built-in audit)

**Status:**
- **Pass**: No vulnerabilities found
- **Warn**: Vulnerabilities detected

---

## Common Checks

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `file_exists` | Required Files | No | 900 | Checks for required files |

### file_exists

Checks that required files exist in the project.

**Configuration:**
```yaml
files:
  required:
    - README.md
    - LICENSE
    - CONTRIBUTING.md
```

**Status:**
- **Pass**: All required files exist
- **Warn**: One or more files missing

---

## External Checks

External checks allow you to run custom commands as quality checks.

**Configuration:**
```yaml
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

  - id: secrets
    name: Secret Detection
    command: gitleaks
    args: ["detect", "--no-git", "--redact", "-v", "."]
    severity: fail
```

**Fields:**
- `id`: Unique identifier for the check
- `name`: Human-readable name
- `command`: Command to run (must be in PATH)
- `args`: Command arguments
- `severity`: `warn` or `fail` (determines if check is critical)

**Exit Code Handling:**
- `0`: Pass
- `1`: Warning (or Fail if severity=fail)
- `2+`: Fail

**JSON Output Support:**
Commands can output JSON with the format:
```json
{"message": "Custom message", "status": "pass|warn|fail"}
```

**Security:**
- Commands are validated via PATH lookup
- No shell interpretation (arguments passed directly)
- Shell metacharacters are rejected

---

## Configuration Reference

### Full .a2.yaml Example

```yaml
# Language configuration
language:
  explicit: []          # Override auto-detect: ["go", "python", "node"]
  auto_detect: true

  go:
    coverage_threshold: 80

  python:
    package_manager: auto    # auto, pip, poetry, pipenv
    test_runner: auto        # auto, pytest, unittest
    formatter: auto          # auto, black, ruff
    linter: auto             # auto, ruff, flake8, pylint
    coverage_threshold: 80

  node:
    package_manager: auto    # auto, npm, yarn, pnpm, bun
    test_runner: auto        # auto, jest, vitest, mocha, npm-test
    formatter: auto          # auto, prettier, biome
    linter: auto             # auto, eslint, biome, oxlint
    coverage_threshold: 80

# Required files
files:
  required:
    - README.md
    - LICENSE

# Disable specific checks
checks:
  disabled:
    - go:deps
    - python:deps

# Execution options
execution:
  parallel: true    # Run checks in parallel (default)

# External checks
external:
  - id: lint
    name: Golangci-lint
    command: golangci-lint
    args: ["run", "./..."]
    severity: warn
```

### Check ID Aliases (Backward Compatibility)

| Alias | Maps To |
|-------|---------|
| `go_mod` | `go:module` |
| `build` | `go:build` |
| `tests` | `go:tests` |
| `gofmt` | `go:format` |
| `govet` | `go:vet` |
| `coverage` | `go:coverage` |
| `deps` | `go:deps` |

---

## Summary

| Language | Total Checks | Critical | Non-Critical |
|----------|-------------|----------|--------------|
| Go | 8 | 3 | 5 |
| Python | 8 | 3 | 5 |
| Node.js | 8 | 3 | 5 |
| Common | 1+ | 0 | 1+ |
| **Total** | **25+** | **9** | **16+** |

**Critical checks** stop execution in sequential mode when they fail.
**Non-critical checks** report warnings but allow other checks to continue.
