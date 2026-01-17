# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Verbose output flag** (`-v`, `-vv`) for `a2 check` command:
  - `-v` shows command output for failed and warning checks only
  - `-vv` shows command output for ALL checks (including passed)
  - Works with all output formats (pretty, JSON, TOON)
  - Useful for debugging why checks failed without needing `a2 run`
  - Example: `a2 check -v` or `a2 check -vv --format json`

## [0.5.1] - 2026-01-14

### Added
- **Run command** (`a2 run CHECK_ID`) to run a specific check with full output:
  - Shows complete stdout/stderr from the underlying tool (not truncated)
  - Useful for debugging why a check failed
  - Supports `--format json` for machine-readable output
  - Example: `a2 run go:race`, `a2 run common:secrets`
- **Explain command** (`a2 explain CHECK_ID`) to show detailed check information:
  - Displays: ID, Name, Description, Languages, Critical status, Suggestion
  - Example: `a2 explain go:race`
- **List checks with descriptions** (`a2 list checks --explain`):
  - Shows detailed descriptions inline with each check
  - Helps understand what each check does without running it
- **Validate commands** for user-defined profiles and targets:
  - `a2 profiles validate` - Validates all user profiles in ~/.config/a2/profiles/
  - `a2 targets validate` - Validates all user targets in ~/.config/a2/targets/
  - Checks for unknown check IDs with typo suggestions (using Levenshtein distance)
  - Warns about duplicate disabled checks
  - Warns when overriding built-in profiles/targets
- **Check descriptions**: Added detailed descriptions to all 80+ checks for use with `explain` command
- Added progress feedback to a2 check that shows "Running checks... (X/Y completed)" during execution.
  - Pretty format: Shows progress on stderr, e.g., "Running checks... (5/25 completed)"
  - JSON/TOON formats: No progress output (stdout remains clean for piping)
  - Progress updates as each check completes (works with both parallel and sequential modes)

### Changed
- **Check output now shows check ID**: Each result line now displays the check ID at the end
  - Example: `✓ PASS Go Build (1.9s) - go:build`
  - Helps users learn check IDs for use with `--skip`, `a2 run`, or `a2 explain`

## [0.5.0] - 2026-01-13

### Added
- **Doctor command** (`a2 doctor`) to check system for required external tools:
  - Detects OS, architecture, and available package managers
  - Shows installed vs missing tools for detected languages
  - Provides platform-specific install commands (brew, apt, go install, cargo, npm, pip)
  - Use `--all` to show tools for all languages
  - Supports `--format json` and `--format toon` for coding agents
- **List checks command** (`a2 list checks`) to display all available checks:
  - Shows all 84 checks grouped by language (Go, Python, Node.js, TypeScript, Java, Rust, Swift, Common)
  - Displays check IDs for use with `--skip` flag
  - Marks critical checks that cause failure
  - Alias: `a2 list c`

## [0.4.0] - 2026-01-12

### Added
- **Swift language support** with 8 checks:
  - `swift:project` - Detects Package.swift, extracts package info and dependencies
  - `swift:build` - Runs `swift build` to verify compilation
  - `swift:tests` - Runs `swift test` to execute tests
  - `swift:format` - Checks code formatting with swift-format or SwiftLint
  - `swift:lint` - Runs SwiftLint for linting and code quality
  - `swift:coverage` - Detects coverage tools (llvm-cov) and reports
  - `swift:deps` - Checks for dependency scanning configuration
  - `swift:logging` - Detects logging (OSLog, swift-log), warns on print()
- **TOON output format** (`--format toon`) for coding agents:
  - Token-Oriented Object Notation - minimal token usage for LLM consumption
  - Line-oriented format with indentation-based structure
  - Efficient encoding of arrays and tabular data
  - Spec: https://github.com/toon-format/spec
- **Config generator command** (`a2 add`) for creating `.a2.yaml` configuration files:
  - Interactive mode (`-i` flag) with guided prompts for profile, target, languages, files, and coverage
  - Non-interactive mode with CLI flags for scripted/automated config generation
  - Auto-detects project languages and shows them for confirmation in interactive mode
  - Validates profile and target names against available options
  - Generates minimal YAML (only non-default values) with helpful comments
  - Supports `--force` flag to overwrite existing config files
  - Supports custom output path via `--output` flag

### Changed
- **No-language detection now exits with error**: When no supported language is detected, a2 now exits with an error message listing supported languages instead of silently falling back to Go checks. Use `--lang` to explicitly specify the language.

## [0.3.1] - 2026-01-09

### Added
- **Monorepo support** via `source_dir` configuration per language:
  - Configure subdirectories for each language (e.g., `rust.source_dir: src-tauri`)
  - Language detection checks configured subdirectories for indicator files
  - Language-specific checks run in the configured source directory
  - Enables multi-language projects like Tauri (React + Rust) to be properly analyzed
- **Configurable profiles and targets**: User-defined profiles and targets stored in `~/.config/a2/`:
  - `~/.config/a2/profiles/` - Custom application profiles (override or extend built-ins)
  - `~/.config/a2/targets/` - Custom maturity targets
  - New commands: `a2 profiles init` and `a2 targets init` to bootstrap directories
  - User definitions override built-in profiles/targets with the same name
  - Profiles/targets show source indicator (built-in/user) in list output
- **Info status**: New check status for informational checks that don't affect maturity score
  - Checks with Info status are executed and displayed but excluded from score calculation
  - Useful for optional recommendations (e.g., "license not found" without penalty)
  - Output shows Info results with cyan color and `i` symbol in terminal
  - JSON output includes `info` count in summary
- **Application profiles** (`--profile` flag) for different application types:
  - `cli` - Command-line tools (skips health, k8s, metrics, api_docs, etc.)
  - `api` - Web services/APIs (all operational checks enabled)
  - `library` - Reusable packages (skips deployment and runtime checks)
  - `desktop` - Desktop applications (skips server-related checks)
- **Maturity targets** (`--target` flag) for project stage:
  - `poc` - Proof of Concept (minimal checks for early development)
  - `production` - All checks enabled (default)
- **New commands**:
  - `a2 profiles` - List available application profiles
  - `a2 targets` - List available maturity targets
- **P1 common checks** for application maturity:
  - `common:contributing` - Detects CONTRIBUTING.md, PR templates, issue templates, CODEOWNERS
  - `common:e2e` - Detects E2E testing (Cypress, Playwright, WebdriverIO, Puppeteer, etc.)
  - `common:tracing` - Detects distributed tracing (OpenTelemetry, Jaeger, Datadog, etc.)
  - `common:migrations` - Detects database migrations (Prisma, Alembic, Flyway, golang-migrate, etc.)
  - `common:config_validation` - Detects config validation (Pydantic, Zod, Viper, validator, etc.)
  - `common:retry` - Detects retry/resilience libraries (tenacity, backoff, Resilience4j, etc.)
  - `common:editorconfig` - Detects editor configuration (.editorconfig, VS Code, JetBrains)
- **TypeScript language support** with 9 checks:
  - `typescript:project` - Detects tsconfig.json, extracts compiler options and TypeScript version
  - `typescript:build` - Runs build script or `tsc --noEmit` to verify compilation
  - `typescript:tests` - Runs tests with Jest, Vitest, or Mocha (auto-detects)
  - `typescript:format` - Checks code formatting with Prettier, Biome, or dprint
  - `typescript:lint` - Runs ESLint, Biome, or oxlint for linting
  - `typescript:type` - Type checking with `tsc --noEmit` (critical check)
  - `typescript:coverage` - Detects coverage tools (Jest, Vitest, c8, nyc)
  - `typescript:deps` - Checks for vulnerabilities using npm/yarn/pnpm audit
  - `typescript:logging` - Detects logging libraries (winston, pino, tslog), warns on console.log
- **Rust language support** with 8 checks:
  - `rust:project` - Detects Cargo.toml, extracts package info, detects workspaces
  - `rust:build` - Runs `cargo check` to verify compilation
  - `rust:tests` - Runs `cargo test` to execute tests
  - `rust:format` - Checks code formatting with `cargo fmt --check`
  - `rust:lint` - Runs Clippy for linting and code quality (`cargo clippy`)
  - `rust:coverage` - Detects coverage tools (tarpaulin, llvm-cov) and reports
  - `rust:deps` - Checks for vulnerabilities using `cargo audit` or cargo-deny
  - `rust:logging` - Detects logging crates (tracing, log, slog), warns on println!

### Changed
- **BREAKING**: Refactored profiles system into profiles and targets:
  - Old `--profile=poc` is now `--target=poc`
  - Old `--profile=production` is now `--target=production`
  - New `--profile` flag now selects application type (cli, api, library, desktop)
  - Both flags can be combined: `a2 check --profile=cli --target=poc`

## [0.2.0] - 2026-01-08

### Added
- **Skip flags**: `--skip` flag to exclude checks by ID (e.g., `--skip=license,k8s`)
- **Built-in profiles**: `--profile` flag for predefined check sets:
  - `poc` - Minimal checks for proof of concept / early development
  - `library` - Focus on code quality, skip deployment checks
  - `production` - All checks enabled (default behavior)
- **Profiles command**: `a2 profiles` lists available profiles
- **Maturity estimation**: Automatic assessment of project maturity level based on check results:
  - Production-Ready: All checks pass (100% score, 0 warnings, 0 failures)
  - Mature: Most checks pass (≥80% score, 0 failures)
  - Development: Core functionality works (≥60% score, ≤2 failures)
  - Proof of Concept: Early stage (<60% score or >2 failures)
- **Java language support** with 8 checks:
  - `java:project` - Detects Maven (pom.xml) or Gradle (build.gradle) projects
  - `java:build` - Compiles with Maven or Gradle (auto-detects wrapper scripts)
  - `java:tests` - Runs JUnit/TestNG tests
  - `java:format` - Detects Spotless, google-java-format, EditorConfig, IDE formatters
  - `java:lint` - Detects Checkstyle, SpotBugs, PMD, Error Prone, SonarQube
  - `java:coverage` - Checks JaCoCo coverage reports against threshold
  - `java:deps` - Detects OWASP Dependency-Check, Snyk, Dependabot, Renovate
  - `java:logging` - Detects SLF4J, Logback, Log4j2; warns on System.out.println
- `common:k8s` check - Detects Kubernetes manifests, Helm charts, Kustomize, Docker Compose, Skaffold, and Tilt configurations
- `common:precommit` check - Detects pre-commit hooks (pre-commit, Husky, Lefthook, Overcommit, commitlint, lint-staged)
- `common:changelog` check - Detects changelog files and release tooling (GoReleaser, semantic-release, release-please)
- `common:secrets` check - Detects secret scanning tools and scans for hardcoded secrets
- `common:api_docs` check - Detects API documentation (OpenAPI/Swagger specs, documentation generators)
- `common:integration` check - Detects integration tests (directories, test files, E2E frameworks, testcontainers)
- `common:metrics` check - Detects metrics instrumentation (Prometheus, OpenTelemetry, Datadog, etc.)
- `common:errors` check - Detects error tracking configuration (Sentry, Rollbar, Bugsnag, etc.)
- `common:shutdown` check - Detects graceful shutdown handling (signal handlers, K8s lifecycle hooks)
- `common:env` check - Validates environment variable handling (.env.example, dotenv libraries, .gitignore)
- `common:license` check - Verifies dependency license compliance (FOSSA, go-licenses, SPDX, CycloneDX)
- `common:sast` check - Verifies SAST tooling is configured (Semgrep, CodeQL, SonarQube, Snyk, etc.)
- `backlog.md` documenting proposed maturity checks for future implementation

### Changed
- Updated documentation in `docs/CHECKS.md` with new checks

### Tests
- Added comprehensive test suites for all Java checks (build, tests, format, lint, coverage, deps, logging)
- Added test suites for new common checks (api_docs, integration, metrics, errors, shutdown)

## [0.1.0] - 2026-01-01

### Added
- Initial release
- Go checks: module, build, tests, race, format, vet, coverage, deps, cyclomatic, logging
- Python checks: project, build, tests, format, lint, type, coverage, deps, complexity, logging
- Node.js checks: project, build, tests, format, lint, type, coverage, deps, logging
- Common checks: file_exists, dockerfile, ci, health, external
- Multi-language auto-detection
- Parallel and sequential execution modes
- Pretty terminal output and JSON output formats
- External check support via `.a2.yaml` configuration
