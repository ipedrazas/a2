# Checks Improvements: Run Tools by Default

## Overview

This document tracks improvements to checks so they run tools with default settings when installed, rather than just checking for config file presence.

**Reference implementation**: `common:secrets` (gitleaks) - already updated to run gitleaks with defaults when installed.

---

## User Control: `RunByDefault` ✅ IMPLEMENTED

### The Problem
Some users may not want checks to run tools automatically, even if installed:
- CI environments where tool execution should be explicit
- Projects where running certain tools is expensive/slow
- Cases where default rules don't match project requirements

### Solution (Implemented)

Added `RunByDefault` attribute to the Tool struct in `pkg/tools/tools.go`:

```go
type Tool struct {
    Name         string
    Description  string
    CheckCmd     []string
    Language     checker.Language
    CheckIDs     []string
    Required     bool
    RunByDefault bool   // If true, run with defaults when installed
    Install      InstallCommands
}
```

Users can override in `.a2.yaml`:
```yaml
tools:
  gitleaks:
    run_by_default: false  # Don't run even if installed
  semgrep:
    run_by_default: true   # Force run with defaults
```

### Implementation Status
- [x] Add `RunByDefault` field to Tool struct
- [x] Update config parsing to support tool overrides (`pkg/config/config.go`)
- [x] Add `ShouldRunByDefault()` helper function (`pkg/tools/tools.go`)
- [x] Default `RunByDefault: true` for security tools (gitleaks, semgrep, trivy, govulncheck, pip-audit, cargo-audit)
- [x] Default `RunByDefault: true` for fast linters (ruff, eslint, biome, swiftlint, swift-format, prettier, black)
- [x] Default `RunByDefault: false` for slow/heavy tools (cargo-tarpaulin)
- [x] Default `RunByDefault: false` for tools needing config (mypy, pytest)
- [x] Added tests for all new functionality

---

## Check Categories

### Category 1: Config-Only Checks (Currently Just Check Presence)

These checks only verify config exists but don't run the tool. **High priority for improvement.**

| Check ID | Tool | Status | Notes |
|----------|------|--------|-------|
| `common:secrets` | gitleaks | ✅ DONE | Runs gitleaks with defaults if installed |
| `common:sast` | semgrep, trivy | ⬜ TODO | Should run with defaults |
| `common:dockerfile` | trivy | ⬜ TODO | Should run trivy for Dockerfile scanning |
| `python:lint` | ruff, flake8, pylint | ⬜ TODO | Should run linter if installed |
| `python:format` | ruff, black | ⬜ TODO | Should run formatter check if installed |
| `python:deps` | pip-audit, safety | ⬜ TODO | Should run scanner if installed |
| `node:lint` | eslint, biome | ⬜ TODO | Should run linter if installed |
| `node:format` | prettier, biome | ⬜ TODO | Should run formatter if installed |
| `swift:lint` | swiftlint | ⬜ TODO | Should run if installed |
| `swift:format` | swift-format | ⬜ TODO | Should run if installed |

---

### Category 2: Already Run with Defaults

These checks already run tools when available without requiring config.

| Check ID | Tool | Default Behavior |
|----------|------|------------------|
| `go:build` | go | Runs `go build ./...` |
| `go:tests` | go | Runs `go test ./...` |
| `go:race` | go | Runs `go test -race -short ./...` |
| `go:format` | gofmt | Runs `gofmt -l .` |
| `go:vet` | go vet | Runs `go vet ./...` |
| `go:coverage` | go test | Runs `go test -cover ./...` |
| `go:deps` | govulncheck | Runs if installed, Info status if not |
| `rust:build` | cargo | Runs `cargo check` |
| `rust:tests` | cargo | Runs `cargo test` |
| `rust:format` | cargo fmt | Runs `cargo fmt --check` |
| `rust:lint` | cargo clippy | Runs `cargo clippy` |
| `python:build` | pip/poetry | Runs package manager check |
| `python:tests` | pytest | Runs `pytest` if installed |
| `node:build` | npm/yarn | Runs `npm install --dry-run` |
| `node:deps` | npm audit | Runs `npm audit` |
| `typescript:type` | tsc | Runs `npx tsc --noEmit` |
| `java:build` | maven/gradle | Runs build command |
| `swift:build` | swift | Runs `swift build` |
| `swift:tests` | swift | Runs `swift test` |

---

### Category 3: Require Config (Cannot Run with Defaults)

These checks need project-specific configuration to run meaningfully.

| Check ID | Why Config Needed |
|----------|-------------------|
| `go:cyclomatic` | Uses AST parsing, no external tool needed |
| `python:type` | Needs mypy config or py.typed marker |
| `python:complexity` | Needs radon but could use defaults |
| `rust:coverage` | cargo-tarpaulin slow, needs explicit opt-in |
| `rust:deps` | cargo-audit could run with defaults |
| `java:*` (most) | Needs Maven/Gradle plugin config |
| `common:external` | User-defined commands by design |

---

### Category 4: File/Pattern Presence Only (No Tool Needed)

These checks don't run external tools, they just check for files or patterns.

| Check ID | What It Checks |
|----------|----------------|
| `*:project` | go.mod, package.json, Cargo.toml, etc. |
| `*:logging` | Code patterns for structured logging |
| `common:ci` | CI config files exist |
| `common:dockerfile` | Dockerfile exists |
| `common:health` | Health endpoint patterns in code |
| `common:env` | .env.example exists |
| `common:license` | License files/configs |
| `common:api_docs` | OpenAPI/Swagger files |
| `common:changelog` | CHANGELOG.md exists |
| `common:k8s` | Kubernetes manifests |
| `common:precommit` | .pre-commit-config.yaml |
| `common:contributing` | CONTRIBUTING.md |
| `common:editorconfig` | .editorconfig |
| `common:e2e` | E2E test patterns |
| `common:tracing` | Tracing config patterns |
| `common:metrics` | Metrics patterns |
| `common:errors` | Error tracking patterns |
| `common:shutdown` | Graceful shutdown patterns |
| `common:migrations` | Migration files |
| `common:config_validation` | Config validation patterns |
| `common:retry` | Retry logic patterns |
| `common:integration` | Integration test files |

---

## Implementation Plan

### Phase 1: Security Tools (High Impact)
Priority: These protect against real security issues

1. [x] `common:secrets` - gitleaks ✅ DONE
2. [ ] `common:sast` - semgrep with defaults
3. [ ] `common:dockerfile` - trivy for Dockerfile scanning
4. [ ] `rust:deps` - cargo-audit with defaults
5. [ ] `python:deps` - pip-audit with defaults

### Phase 2: Linting Tools (High Value)
Priority: Fast, good defaults, immediate feedback

1. [ ] `python:lint` - ruff/flake8 with defaults
2. [ ] `python:format` - ruff/black with defaults
3. [ ] `node:lint` - eslint/biome with defaults
4. [ ] `node:format` - prettier/biome with defaults
5. [ ] `swift:lint` - swiftlint with defaults
6. [ ] `swift:format` - swift-format with defaults

### Phase 3: Heavy Tools (Opt-in)
Priority: Slow but valuable, should be opt-in

1. [ ] `rust:coverage` - cargo-tarpaulin (slow, keep as opt-in)
2. [ ] `python:complexity` - radon (fast, could default)

---

## Tool Registry Updates Needed

Add these tools to `pkg/tools/tools.go`:

```go
// Already in registry:
// - govulncheck, gocyclo (Go)
// - pytest, ruff, black, mypy, pip-audit, radon (Python)
// - eslint, prettier, biome (Node)
// - cargo-audit, cargo-tarpaulin (Rust)
// - swiftlint, swift-format (Swift)
// - gitleaks, semgrep, trivy (Common)

// Need to add:
// - flake8 (Python linter)
// - pylint (Python linter)
// - safety (Python security)
// - oxlint (Node linter - fast alternative)
```

---

## Check Implementation Pattern

When updating a check to run with defaults:

```go
func (c *MyCheck) Run(path string) (checker.Result, error) {
    rb := checkutil.NewResultBuilder(c, checker.LangXxx)

    // 1. Check if tool is installed
    toolInstalled := checkutil.ToolAvailable("mytool")
    configFile := c.findConfig(path)

    // 2. Priority 1: Tool installed + config → run with config
    // 3. Priority 2: Tool installed + no config → run with defaults
    if toolInstalled {
        return c.runTool(path, configFile, rb)
    }

    // 4. Priority 3: Config exists but no tool → trust CI runs it
    if configFile != "" {
        return rb.Pass("mytool configured (run manually or in CI)"), nil
    }

    // 5. Priority 4: No tool, no config → warn or fallback
    return rb.Warn("No mytool configured (consider installing)"), nil
}
```

---

## Tracking Progress

### Completed
- [x] `common:secrets` - gitleaks runs with defaults
- [x] `RunByDefault` attribute added to Tool struct
- [x] Config override support (`tools:` section in `.a2.yaml`)
- [x] Helper function `tools.ShouldRunByDefault()`
- [x] Tests for all new functionality

### To Do - Check Implementations
These checks need to be updated to use `tools.ShouldRunByDefault()`:

- [ ] `common:sast` - semgrep
- [ ] `common:dockerfile` - trivy
- [ ] `python:lint` - ruff/flake8
- [ ] `python:format` - ruff/black
- [ ] `python:deps` - pip-audit
- [ ] `node:lint` - eslint/biome
- [ ] `node:format` - prettier/biome
- [ ] `swift:lint` - swiftlint
- [ ] `swift:format` - swift-format
- [ ] `rust:deps` - cargo-audit

---

## Output Message Examples

| Scenario | Status | Message |
|----------|--------|---------|
| Tool ran, no issues (with config) | Pass | "mytool: No issues (using .mytool.yaml)" |
| Tool ran, no issues (defaults) | Pass | "mytool: No issues (default rules)" |
| Tool ran, found issues | Warn | "mytool: 3 issues found" |
| Tool not installed, config exists | Pass | "mytool configured (install to run locally)" |
| Tool not installed, no config | Warn/Info | "No mytool configured (consider adding)" |
