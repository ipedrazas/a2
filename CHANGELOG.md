# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

## [0.1.0] - 2024-01-01

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
