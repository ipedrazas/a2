# A2 Checks Reference

This document describes all checks available in A2, organized by category.

## Table of Contents

- [Check Statuses](#check-statuses)
- [Language Checks](#language-checks)
- [Common Checks](#common-checks)
- [External Checks](#external-checks)
- [Configuration Reference](#configuration-reference)

---

## Check Statuses

Each check returns one of four statuses:

| Status | Symbol | Color | Affects Score | Description |
|--------|--------|-------|---------------|-------------|
| **Pass** | `✓` | Green | Yes | Check passed successfully |
| **Warn** | `!` | Yellow | Yes | Warning - something needs attention but isn't critical |
| **Fail** | `✗` | Red | Yes | Check failed - critical issue detected |
| **Info** | `i` | Cyan | No | Informational only - does not affect maturity score |

### Score Calculation

The maturity score is calculated as:

```
Score = (Passed / ScoredChecks) × 100%

Where ScoredChecks = Passed + Warnings + Failed
```

**Info checks are excluded from the score calculation.** This allows you to run checks that provide useful information without penalizing the project's maturity rating.

### Use Cases for Info Status

- Optional recommendations (e.g., "consider adding a CONTRIBUTING.md")
- Informational metrics (e.g., code statistics)
- Checks you want visibility into without enforcement

---

## Language Checks

Language-specific checks are documented in separate files:

| Language | Checks | Documentation |
|----------|--------|---------------|
| Go | 10 | [Go Checks](checks/go.md) |
| Python | 10 | [Python Checks](checks/python.md) |
| Node.js | 9 | [Node.js Checks](checks/node.md) |
| TypeScript | 9 | [TypeScript Checks](checks/typescript.md) |
| Java | 8 | [Java Checks](checks/java.md) |
| Rust | 8 | [Rust Checks](checks/rust.md) |

---

## Common Checks

These checks apply to all projects regardless of language.

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `file_exists` | Required Files | No | 900 | Checks for required files |
| `common:dockerfile` | Container Ready | No | 910 | Checks for Dockerfile/Containerfile |
| `common:ci` | CI Pipeline | No | 920 | Detects CI/CD configuration |
| `common:health` | Health Endpoint | No | 930 | Detects health check endpoints |
| `common:secrets` | Secrets Detection | No | 940 | Detects secret scanning config or hardcoded secrets |
| `common:env` | Environment Config | No | 945 | Validates environment variable handling practices |
| `common:license` | License Compliance | No | 950 | Verifies dependency license compliance tooling |
| `common:sast` | SAST Security Scanning | No | 955 | Verifies SAST tooling is configured |
| `common:api_docs` | API Documentation | No | 960 | Detects OpenAPI/Swagger specs and documentation generators |
| `common:changelog` | Changelog | No | 965 | Verifies changelog or release notes exist |
| `common:integration` | Integration Tests | No | 980 | Detects integration test directories, files, and E2E frameworks |
| `common:metrics` | Metrics Instrumentation | No | 1010 | Detects Prometheus, OpenTelemetry, and other metrics libraries |
| `common:errors` | Error Tracking | No | 1020 | Detects Sentry, Rollbar, Bugsnag, and other error tracking SDKs |
| `common:k8s` | Kubernetes Ready | No | 1030 | Detects Kubernetes manifests and deployment configurations |
| `common:shutdown` | Graceful Shutdown | No | 1035 | Detects signal handling and graceful shutdown configuration |
| `common:precommit` | Pre-commit Hooks | No | 1065 | Verifies pre-commit hooks are configured |
| `common:contributing` | Contributing Guidelines | No | 1070 | Detects CONTRIBUTING.md, PR templates, issue templates |
| `common:e2e` | E2E Testing | No | 1080 | Detects E2E testing frameworks (Cypress, Playwright, etc.) |
| `common:tracing` | Distributed Tracing | No | 1090 | Detects OpenTelemetry, Jaeger, Datadog tracing |
| `common:migrations` | Database Migrations | No | 1100 | Detects migration tools (Prisma, Alembic, Flyway, etc.) |
| `common:config_validation` | Config Validation | No | 1110 | Detects config validation (Pydantic, Zod, Viper, etc.) |
| `common:retry` | Retry/Resilience | No | 1120 | Detects retry libraries (tenacity, backoff, Resilience4j) |
| `common:editorconfig` | Editor Config | No | 1130 | Detects editor configuration files |

---

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

### common:dockerfile

Checks if the project is container-ready by looking for container configuration.

**Detection files:**
- `Dockerfile`
- `Containerfile`
- `dockerfile` (case-insensitive)
- `Dockerfile.*` (e.g., `Dockerfile.prod`)

**Bonus detection:**
- `.dockerignore` (reported in pass message)

**Status:**
- **Pass**: Dockerfile found
- **Warn**: No Dockerfile found

---

### common:ci

Detects CI/CD pipeline configuration in the project.

**Supported CI systems:**
- GitHub Actions: `.github/workflows/*.yml`, `.github/workflows/*.yaml`
- GitLab CI: `.gitlab-ci.yml`
- Jenkins: `Jenkinsfile`
- CircleCI: `.circleci/config.yml`
- Travis CI: `.travis.yml`
- Azure Pipelines: `azure-pipelines.yml`
- Bitbucket Pipelines: `bitbucket-pipelines.yml`

**Status:**
- **Pass**: CI configuration found (reports which CI system)
- **Warn**: No CI configuration found

---

### common:health

Detects health check endpoints in the codebase for production readiness.

**Endpoint patterns detected:**
- `/health`, `/healthz`, `/ready`, `/readiness`, `/live`, `/liveness`

**Function patterns detected:**
- `HealthCheck`, `healthCheck`, `health_check`

**Files scanned:**
- Go: `*.go`
- Python: `*.py`
- JavaScript/TypeScript: `*.js`, `*.ts`, `*.jsx`, `*.tsx`, `*.mjs`, `*.cjs`

**Status:**
- **Pass**: Health endpoint pattern found
- **Warn**: No health endpoint detected

---

### common:secrets

Detects secret scanning configuration or scans for hardcoded secrets in the codebase.

**Secret scanning tools detected:**
- Gitleaks: `.gitleaks.toml`, `.gitleaks.yaml`, `gitleaks.toml`
- TruffleHog: `.trufflehog.yml`, `trufflehog.yml`
- Secretlint: `.secretlintrc`, `.secretlintrc.json`
- detect-secrets: `.secrets.baseline`
- git-secrets: `.git-secrets`
- Pre-commit hooks with secret scanning

**Secret patterns detected (if no scanner configured):**
- AWS Access Keys (`AKIA...`)
- AWS Secret Keys
- Private Keys (RSA, DSA, EC, OPENSSH, PGP)
- GitHub Tokens (`ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_`)
- JWT Tokens
- Slack Tokens
- Stripe Keys (`sk_live_`)
- SendGrid Keys
- Database URLs with credentials
- Generic API keys and passwords

**Status:**
- **Pass**: Secret scanning tool configured
- **Warn**: Potential secrets found, or no secret scanning configured

**Recommendation:** Configure Gitleaks or similar tool for automated secret scanning.

---

### common:env

Validates environment variable handling practices to ensure configuration is properly separated from code.

**Environment template files detected:**
- `.env.example`, `.env.sample`, `.env.template`, `example.env`, `.env.local.example`

**Dotenv libraries detected:**
- Go: `github.com/joho/godotenv`, `github.com/caarlos0/env`, `github.com/kelseyhightower/envconfig`, `github.com/spf13/viper`
- Python: `python-dotenv`, `environs`, `pydantic-settings`, `django-environ`, `dynaconf`
- Node.js: `dotenv`, `dotenv-safe`, `dotenv-expand`, `env-cmd`, `cross-env`
- Java: Spring Boot (`application.properties`, `application.yml`), `dotenv-java`

**Status:**
- **Pass**: Environment template found, dotenv library detected, or `.env` properly gitignored
- **Warn**: `.env` exists but not gitignored, or no environment configuration found

**Recommendation:** Create `.env.example` to document required environment variables, and ensure `.env` is in `.gitignore`.

---

### common:license

Verifies that dependency license compliance tooling is configured to track license obligations.

**License audit config files detected:**
- `.licensrc`, `.licensrc.json`, `.licensrc.yaml`, `.licensrc.yml`
- `license-checker.json`, `.license-checker.json`

**FOSSA configuration detected:**
- `.fossa.yml`, `.fossa.yaml`, `fossa.yml`, `fossa.yaml`

**SPDX/SBOM files detected:**
- `spdx.json`, `spdx.yaml`, `sbom.spdx`, `sbom.spdx.json`

**Status:**
- **Pass**: License compliance tooling or SBOM generation found
- **Warn**: No license compliance tooling found

**Recommendation:** Configure FOSSA, go-licenses, license-checker, or generate SBOM with CycloneDX.

---

### common:sast

Verifies that SAST (Static Application Security Testing) tooling is configured for automated security scanning.

**Tools detected:**
- Semgrep: `.semgrep.yml`, `semgrep/` directories
- SonarQube/SonarCloud: `sonar-project.properties`
- Snyk: `.snyk`, `snyk.json`
- CodeQL: `.github/codeql/`, workflows with `github/codeql-action`
- Checkmarx, Veracode, Fortify, Coverity, Bearer, Horusec

**Language-specific security tools:**
- Go: gosec
- Python: Bandit, Safety
- Node.js: eslint-plugin-security, audit-ci
- Java: FindSecBugs, OWASP Dependency-Check

**Status:**
- **Pass**: SAST tooling found
- **Warn**: No SAST tooling found

**Recommendation:** Configure Semgrep, CodeQL, SonarQube, or Snyk for automated security scanning.

---

### common:api_docs

Detects API documentation configuration and files.

**OpenAPI/Swagger files detected:**
- `openapi.yaml`, `openapi.yml`, `openapi.json`
- `swagger.yaml`, `swagger.yml`, `swagger.json`
- Searched in root and `docs/`, `api/` directories

**Documentation generators detected:**
- Go: `swaggo/swag`, `go-swagger`, `grpc-gateway`
- Python: `drf-spectacular`, `drf-yasg`, `flasgger`, `FastAPI`
- Node.js: `swagger-jsdoc`, `swagger-ui-express`, `@nestjs/swagger`, `tsoa`

**Status:**
- **Pass**: API documentation or generator found
- **Warn**: No API documentation found

---

### common:changelog

Verifies that a changelog or release notes file exists.

**Changelog files detected:**
- `CHANGELOG.md`, `CHANGELOG.txt`, `CHANGELOG`
- `CHANGES.md`, `HISTORY.md`, `RELEASES.md`, `RELEASE_NOTES.md`

**Release tooling detected:**
- GoReleaser: `.goreleaser.yml`
- semantic-release: `.releaserc`, `release.config.js`
- release-please: `release-please-config.json`
- standard-version: `.versionrc`
- changesets: `.changeset/config.json`

**Status:**
- **Pass**: Changelog file found or release tooling configured
- **Warn**: No changelog or release tooling found

**Recommendation:** Create a `CHANGELOG.md` following [Keep a Changelog](https://keepachangelog.com) format.

---

### common:integration

Detects integration test directories, files, and E2E testing frameworks.

**Integration test directories detected:**
- `tests/integration/`, `test/integration/`
- `integration_tests/`, `integration-tests/`
- `tests/e2e/`, `test/e2e/`, `e2e/`

**Test infrastructure detected:**
- Docker Compose for tests: `docker-compose.test.yml`
- testcontainers library usage

**E2E testing frameworks detected:**
- Cypress, Playwright, WebdriverIO, Selenium

**Status:**
- **Pass**: Integration tests or E2E framework found
- **Warn**: No integration tests found

---

### common:metrics

Detects metrics instrumentation libraries and configuration.

**Metrics libraries detected:**
- Go: `prometheus/client_golang`, `opentelemetry-go`
- Python: `prometheus_client`, `opentelemetry`
- Node.js: `prom-client`, `@opentelemetry/sdk-metrics`
- Java: `micrometer`, `prometheus`, `opentelemetry`

**Configuration files detected:**
- `prometheus.yml`, `otel-collector-config.yaml`

**Status:**
- **Pass**: Metrics library or configuration found
- **Warn**: No metrics instrumentation found

---

### common:errors

Detects error tracking SDK configuration.

**Error tracking SDKs detected:**
- Sentry: `sentry-sdk`, `@sentry/node`, `getsentry/sentry-go`
- Rollbar: `rollbar` packages
- Bugsnag: `bugsnag` packages
- Honeybadger

**Configuration files detected:**
- `.sentryclirc`, `sentry.properties`
- `.rollbar`, `bugsnag.json`

**Status:**
- **Pass**: Error tracking SDK or configuration found
- **Warn**: No error tracking found

---

### common:k8s

Detects Kubernetes manifests, Helm charts, and deployment configurations.

**Helm charts detected:**
- `Chart.yaml`, `charts/*/Chart.yaml`, `helm/Chart.yaml`

**Kustomize configurations detected:**
- `kustomization.yaml`, `kustomization.yml`

**Kubernetes manifest directories detected:**
- `k8s/`, `kubernetes/`, `deploy/`, `manifests/`

**Alternative deployment tools detected:**
- Docker Compose: `docker-compose.yaml`, `compose.yaml`
- Skaffold: `skaffold.yaml`
- Tilt: `Tiltfile`

**Status:**
- **Pass**: Kubernetes manifests, Helm chart, or deployment config found
- **Warn**: No deployment configuration found

---

### common:shutdown

Detects graceful shutdown handling for proper process termination.

**Signal handling detected:**
- Go: `signal.Notify`, `syscall.SIGTERM`
- Python: `signal.signal`, `signal.SIGTERM`
- Node.js: `process.on('SIGTERM')`, `process.on('SIGINT')`
- Java: `Runtime.getRuntime().addShutdownHook`, `@PreDestroy`

**Kubernetes lifecycle hooks detected:**
- `terminationGracePeriodSeconds`, `preStop` hooks

**Status:**
- **Pass**: Signal handling or shutdown hooks found
- **Warn**: No graceful shutdown handling found

---

### common:precommit

Verifies that pre-commit hooks are configured.

**Pre-commit tools detected:**
- pre-commit (Python): `.pre-commit-config.yaml`
- Husky (Node.js): `.husky/` directory
- Lefthook: `lefthook.yml`
- Overcommit (Ruby): `.overcommit.yml`
- Native git hooks: `.git/hooks/`

**Related tooling detected:**
- commitlint: `commitlint.config.js`
- lint-staged: `lint-staged.config.js`

**Status:**
- **Pass**: Pre-commit hooks or related tooling configured
- **Warn**: No pre-commit hooks found

---

### common:contributing

Detects contribution guidelines and templates.

**Files detected:**
- `CONTRIBUTING.md`
- `.github/PULL_REQUEST_TEMPLATE.md`
- `.github/ISSUE_TEMPLATE/`
- `CODEOWNERS`

**Status:**
- **Pass**: Contributing guidelines or templates found
- **Warn**: No contributing documentation found

---

### common:e2e

Detects E2E (End-to-End) testing frameworks.

**E2E frameworks detected:**
- Cypress: `cypress.config.js`, `cypress.config.ts`, `cypress/`
- Playwright: `playwright.config.js`, `playwright.config.ts`
- WebdriverIO: `wdio.conf.js`, `wdio.conf.ts`
- TestCafe: `.testcaferc.json`, `testcafe/`
- Puppeteer: `puppeteer` in dependencies
- Selenium: `selenium-webdriver` in dependencies

**Status:**
- **Pass**: E2E testing framework found
- **Warn**: No E2E testing found

---

### common:tracing

Detects distributed tracing configuration.

**Tracing libraries detected:**
- OpenTelemetry: `@opentelemetry/*`, `opentelemetry-*`
- Jaeger: `jaeger-client`, `github.com/uber/jaeger-client-go`
- Datadog: `dd-trace`, `ddtrace`
- Zipkin: `zipkin`
- Sentry tracing

**Configuration files detected:**
- OpenTelemetry config files

**Status:**
- **Pass**: Tracing library or configuration found
- **Warn**: No distributed tracing found

---

### common:migrations

Detects database migration tools and configurations.

**Migration tools detected:**
- Prisma: `prisma/schema.prisma`, `prisma/migrations/`
- Alembic (Python): `alembic.ini`, `alembic/`
- Flyway: `flyway.conf`, `flyway/`
- golang-migrate: `migrate` in go.mod
- GORM migrations
- TypeORM: `typeorm` in package.json
- Sequelize: `sequelize` in package.json
- Knex: `knex` in package.json

**Status:**
- **Pass**: Migration tool found
- **Warn**: No migration tooling found

---

### common:config_validation

Detects configuration validation libraries.

**Validation libraries detected:**
- Go: `github.com/spf13/viper`, `github.com/go-playground/validator`
- Python: `pydantic`, `pydantic-settings`
- Node.js/TypeScript: `zod`, `joi`, `yup`, `class-validator`
- Rust: `serde`, `config`

**Status:**
- **Pass**: Config validation library found
- **Warn**: No config validation found

---

### common:retry

Detects retry and resilience libraries.

**Retry libraries detected:**
- Go: `cenkalti/backoff`, `avast/retry-go`, `sony/gobreaker`
- Python: `tenacity`, `backoff`, `retrying`
- Node.js: `async-retry`, `retry`, `cockatiel`
- Java: Resilience4j, Spring Retry

**Circuit breaker patterns detected:**
- Hystrix, Resilience4j, gobreaker

**Status:**
- **Pass**: Retry/resilience library found
- **Warn**: No retry handling found

---

### common:editorconfig

Detects editor configuration files.

**Files detected:**
- `.editorconfig`
- VS Code: `.vscode/settings.json`
- JetBrains: `.idea/`
- Vim: `.vim/`, `.vimrc`
- Dev Containers: `.devcontainer/`

**Status:**
- **Pass**: Editor configuration found
- **Warn**: No editor configuration found

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
    cyclomatic_threshold: 15

  python:
    package_manager: auto    # auto, pip, poetry, pipenv
    test_runner: auto        # auto, pytest, unittest
    formatter: auto          # auto, black, ruff
    linter: auto             # auto, ruff, flake8, pylint
    coverage_threshold: 80
    cyclomatic_threshold: 15

  node:
    package_manager: auto    # auto, npm, yarn, pnpm, bun
    test_runner: auto        # auto, jest, vitest, mocha, npm-test
    formatter: auto          # auto, prettier, biome
    linter: auto             # auto, eslint, biome, oxlint
    coverage_threshold: 80

  typescript:
    package_manager: auto    # auto, npm, yarn, pnpm, bun
    test_runner: auto        # auto, jest, vitest, mocha
    formatter: auto          # auto, prettier, biome, dprint
    linter: auto             # auto, eslint, biome, oxlint
    coverage_threshold: 80

  java:
    build_tool: auto         # auto, maven, gradle
    test_runner: auto        # auto, junit, testng
    coverage_threshold: 80

  rust:
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

---

## Summary

| Category | Total Checks |
|----------|-------------|
| Go | 10 |
| Python | 10 |
| Node.js | 9 |
| TypeScript | 9 |
| Java | 8 |
| Rust | 8 |
| Common | 23 |
| **Total** | **77** |

**Critical checks** stop execution in sequential mode when they fail.
**Non-critical checks** report warnings but allow other checks to continue.
