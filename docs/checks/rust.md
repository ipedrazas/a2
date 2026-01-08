# Rust Checks

This document describes all Rust-specific checks available in A2.

## Overview

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `rust:project` | Rust Project | Yes | 100 | Detects Cargo.toml, extracts package info, detects workspaces |
| `rust:build` | Rust Build | Yes | 110 | Runs `cargo check` to verify compilation |
| `rust:tests` | Rust Tests | Yes | 120 | Runs `cargo test` to execute tests |
| `rust:format` | Rust Format | No | 200 | Checks code formatting with `cargo fmt --check` |
| `rust:lint` | Rust Clippy | No | 210 | Runs Clippy for linting and code quality (`cargo clippy`) |
| `rust:coverage` | Rust Coverage | No | 220 | Detects coverage tools (tarpaulin, llvm-cov) and reports |
| `rust:deps` | Rust Vulnerabilities | No | 230 | Checks for vulnerabilities using `cargo audit` or cargo-deny |
| `rust:logging` | Rust Logging | No | 250 | Detects logging crates (tracing, log, slog), warns on println! |

---

## rust:project

Detects Rust project configuration and extracts package information.

**Detection:**
- `Cargo.toml` - Required

**Information extracted:**
- Package name
- Version
- Workspace detection (`[workspace]` section)

**Status:**
- **Pass**: Cargo.toml found with valid configuration
- **Fail**: No Cargo.toml found

---

## rust:build

Runs `cargo check` to verify the project compiles without producing binaries.

**Command:** `cargo check`

This is faster than `cargo build` as it only checks for compilation errors.

**Status:**
- **Pass**: Check completes without errors
- **Fail**: Compilation errors

---

## rust:tests

Runs all tests using Cargo's built-in test runner.

**Command:** `cargo test`

**Status:**
- **Pass**: All tests pass
- **Warn**: No tests found
- **Fail**: Tests fail

---

## rust:format

Checks code formatting using rustfmt via Cargo.

**Command:** `cargo fmt --check`

**Configuration files:**
- `rustfmt.toml`
- `.rustfmt.toml`

**Status:**
- **Pass**: All files properly formatted
- **Warn**: Formatting issues found

**Fix:** Run `cargo fmt`

---

## rust:lint

Runs Clippy for linting and code quality analysis.

**Command:** `cargo clippy -- -D warnings`

Clippy catches common mistakes and suggests improvements.

**Status:**
- **Pass**: No Clippy warnings
- **Warn**: Clippy warnings or errors found

**Fix:** Run `cargo clippy --fix` for auto-fixable issues.

---

## rust:coverage

Detects and checks test coverage configuration.

**Configuration:**
```yaml
language:
  rust:
    coverage_threshold: 80  # Default: 80%
```

**Coverage tools detected:**
- **cargo-tarpaulin**: `tarpaulin.toml`, `.tarpaulin.toml`
- **cargo-llvm-cov**: Referenced in Cargo.toml or CI config
- **Coverage reporting**: Codecov, Coveralls in CI configs

**Report locations checked:**
- `target/tarpaulin/cobertura.xml`
- `target/llvm-cov/html/index.html`
- `coverage/cobertura.xml`
- `coverage.xml`
- `lcov.info`

**CI files checked:**
- `.github/workflows/ci.yml`, `.github/workflows/rust.yml`
- `.gitlab-ci.yml`

**Status:**
- **Pass**: Coverage meets threshold or tooling configured
- **Warn**: Coverage below threshold or no tooling found

**Recommendation:** Install cargo-tarpaulin (`cargo install cargo-tarpaulin`) or cargo-llvm-cov.

---

## rust:deps

Checks for security vulnerabilities in dependencies.

**Tools:**
- `cargo audit` - RustSec Advisory Database
- `cargo-deny` - License and security checking

**Status:**
- **Pass**: No vulnerabilities found or audit tool configured
- **Warn**: Vulnerabilities detected or no audit tool installed

**Recommendation:** Install cargo-audit (`cargo install cargo-audit`).

---

## rust:logging

Checks for structured logging practices instead of println! macro.

**Structured logging crates detected (preferred):**
- `tracing` - Modern async-aware instrumentation
- `tracing-subscriber` - Tracing subscribers
- `slog` - Structured, contextual logging
- `slog-json` - JSON output for slog

**Basic logging crates detected:**
- `log` - Standard logging facade
- `env_logger` - Environment-based logger
- `fern` - Flexible logging
- `flexi_logger` - Flexible logger
- `log4rs` - Log4j-like logging
- `pretty_env_logger` - Pretty env_logger

**Anti-patterns detected (in src/main.rs, src/lib.rs):**
- `println!` macro
- `print!` macro

Comments are excluded from detection.

**Status:**
- **Pass**: Logging crate found, no println! in source
- **Warn**: println! found or no logging crate detected

**Recommendation:** Use `tracing` for async applications or `log` with `env_logger` for simpler cases.

---

## Configuration Example

```yaml
language:
  rust:
    coverage_threshold: 80
```

## Recommended Rust Tooling

```bash
# Install Clippy (usually included with rustup)
rustup component add clippy

# Install rustfmt (usually included with rustup)
rustup component add rustfmt

# Install cargo-audit for vulnerability scanning
cargo install cargo-audit

# Install cargo-tarpaulin for code coverage
cargo install cargo-tarpaulin

# Or install cargo-llvm-cov (requires LLVM)
cargo install cargo-llvm-cov
```
