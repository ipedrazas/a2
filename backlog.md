# a2 Checks Backlog

This document tracks proposed new checks for a2 to help assess **application maturity**. With the rise of coding agents producing large amounts of code, having a tool to quickly understand an application's production-readiness is essential.

## Vision

a2 should answer: "Is this application mature and production-ready?" by checking for:
- Security best practices
- Proper documentation
- Comprehensive testing
- Observability instrumentation
- Production-ready configuration
- Clean architecture patterns
- Solid development workflow

---

## Current Checks (52 total)

| Language | Checks |
|----------|--------|
| Go | module, build, tests, race, format, vet, coverage, deps, cyclomatic, logging |
| Python | project, build, tests, format, lint, type, coverage, deps, complexity, logging |
| Node.js | project, build, tests, format, lint, type, coverage, deps, logging |
| Java | project, build, tests, format, lint, coverage, deps, logging |
| Common | file_exists, dockerfile, ci, health, **secrets**, **api_docs**, **changelog**, **integration**, **metrics**, **errors**, **precommit**, **k8s**, **shutdown**, external |

---

## Proposed New Checks

### Priority Levels
- **P0 (High)**: Critical for production readiness assessment
- **P1 (Medium)**: Important for mature applications
- **P2 (Low)**: Nice-to-have for comprehensive assessment

---

## 1. Security Checks

### `common:secrets` [P0] ✅ IMPLEMENTED
**Detect hardcoded secrets, API keys, and tokens**

Mature applications never store secrets in code. They use environment variables, secret managers, or encrypted configs.

| Attribute | Value |
|-----------|-------|
| Order | 940 |
| Critical | No (Warn) |

**Detection approach:**
- Check for secret scanning config: `.gitleaks.toml`, `.secretlintrc`, `trufflehog.yml`
- Scan for common patterns: `API_KEY=`, `password=`, `secret=`, AWS keys (`AKIA...`), JWT tokens
- Skip `.env.example` template files

**Implementation:** `pkg/checks/common/secrets.go`

---

### `common:license` [P1]
**Verify dependency license compliance**

Mature applications track license obligations for legal compliance.

| Attribute | Value |
|-----------|-------|
| Order | 950 |
| Critical | No (Warn) |

**Detection approach:**
- Check for license audit configs: `.licensrc`, `.fossa.yml`, `license-checker.json`
- Check for tools in dependencies:
  - Go: `go-licenses`
  - Python: `pip-licenses`, `liccheck`
  - Node: `license-checker`, `license-compliance`

---

### `common:env` [P1]
**Validate environment variable handling**

Mature applications separate config from code and document required env vars.

| Attribute | Value |
|-----------|-------|
| Order | 945 |
| Critical | No (Warn) |

**Detection approach:**
- Check for `.env.example` or `.env.sample` (documents required vars)
- Verify `.env` is in `.gitignore`
- Check for dotenv library usage
- Warn if `.env` exists in repo (should be gitignored)

---

### `common:sast` [P1]
**Verify SAST tooling is configured**

Mature applications have automated security scanning in their pipeline.

| Attribute | Value |
|-----------|-------|
| Order | 955 |
| Critical | No (Warn) |

**Detection approach:**
- Check for configs: `.semgrep.yml`, `sonar-project.properties`, `.snyk`
- Check for CodeQL: `.github/workflows/codeql*.yml`
- Check for security step in CI config

---

### `go:gosec` / `python:bandit` / `node:security` [P1]
**Language-specific static security analysis**

| Attribute | Value |
|-----------|-------|
| Order | 260 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Run `gosec ./...` or check for gosec in CI
- Python: Run `bandit -r .` or check for bandit in CI/pyproject.toml
- Node: Run `npm audit` or check for audit in CI

---

## 2. Documentation Checks

### `common:api_docs` [P0] ✅ IMPLEMENTED
**Verify API documentation exists**

Mature APIs are well-documented; OpenAPI specs enable client generation and testing.

| Attribute | Value |
|-----------|-------|
| Order | 960 |
| Critical | No (Warn) |

**Detection approach:**
- Check for OpenAPI specs: `openapi.yaml`, `swagger.yaml`, `api.yaml`
- Check for doc directories: `docs/api/`, `api-docs/`
- Check for doc generators:
  - Go: `swaggo/swag`
  - Python: `drf-spectacular`, `flasgger`
  - Node: `swagger-jsdoc`, `@nestjs/swagger`

---

### `common:changelog` [P0] ✅ IMPLEMENTED
**Verify changelog exists**

Mature applications document changes for users and maintainers.

| Attribute | Value |
|-----------|-------|
| Order | 965 |
| Critical | No (Warn) |

**Detection approach:**
- Check for: `CHANGELOG.md`, `CHANGES.md`, `HISTORY.md`, `NEWS.md`
- Bonus: Check for Keep a Changelog format (`## [Unreleased]`, `### Added`)
- Check for release tools: `.releaserc`, `.goreleaser.yml`, `release.config.js`

**Implementation:** `pkg/checks/common/changelog.go`

---

### `common:contributing` [P1]
**Verify contribution guidelines exist**

Mature projects have clear contribution guidelines.

| Attribute | Value |
|-----------|-------|
| Order | 970 |
| Critical | No (Warn) |

**Detection approach:**
- Check for: `CONTRIBUTING.md`, `.github/CONTRIBUTING.md`
- Check for templates: `.github/PULL_REQUEST_TEMPLATE.md`, `.github/ISSUE_TEMPLATE/`

---

### `common:conduct` [P2]
**Verify code of conduct exists**

Mature community projects establish behavioral expectations.

| Attribute | Value |
|-----------|-------|
| Order | 975 |
| Critical | No (Warn) |

**Detection approach:**
- Check for: `CODE_OF_CONDUCT.md`, `.github/CODE_OF_CONDUCT.md`

---

### `go:godoc` / `python:docstrings` / `node:jsdoc` [P2]
**Validate exported functions have documentation**

Mature code is self-documenting; exported APIs should have clear docs.

| Attribute | Value |
|-----------|-------|
| Order | 270 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Parse AST, check exported functions have doc comments
- Python: Check for docstrings on public functions/classes
- Node/TS: Check for JSDoc on exported functions
- Set threshold (e.g., >70% documented)

---

## 3. Testing Depth Checks

### `common:integration` [P0] ✅ IMPLEMENTED
**Verify integration tests exist**

Mature applications test components working together, not just in isolation.

| Attribute | Value |
|-----------|-------|
| Order | 980 |
| Critical | No (Warn) |

**Detection approach:**
- Check for directories: `tests/integration/`, `integration_tests/`, `test/integration/`
- Check for naming patterns: `*_integration_test.go`, `test_integration_*.py`, `*.integration.test.ts`
- Check for test infrastructure: `docker-compose.test.yml`, `docker-compose.e2e.yml`

---

### `common:e2e` [P1]
**Verify end-to-end tests exist**

Mature applications validate the entire user flow.

| Attribute | Value |
|-----------|-------|
| Order | 985 |
| Critical | No (Warn) |

**Detection approach:**
- Check for e2e frameworks:
  - `cypress.config.js`, `playwright.config.ts`, `wdio.conf.js`
- Check for directories: `e2e/`, `cypress/`, `playwright/`
- Check for dependencies: `selenium`, `playwright`, `puppeteer`

---

### `common:fixtures` [P2]
**Verify test fixtures/factories exist**

Mature test suites use consistent, maintainable test data patterns.

| Attribute | Value |
|-----------|-------|
| Order | 990 |
| Critical | No (Warn) |

**Detection approach:**
- Check for directories: `fixtures/`, `testdata/`, `factories/`, `__fixtures__/`
- Check for factory libraries:
  - Go: custom factories or testify suites
  - Python: `factory_boy`
  - Node: `fishery`, `factory.ts`

---

### `common:mocking` [P2]
**Verify mocking infrastructure exists**

Mature applications mock external dependencies for reliable, fast tests.

| Attribute | Value |
|-----------|-------|
| Order | 995 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Check for `*_mock.go` files, `gomock`, `testify/mock`
- Python: Check for `unittest.mock`, `pytest-mock`, `responses`, `httpretty`
- Node: Check for `jest.mock`, `sinon`, `nock`, `msw`

---

## 4. Observability Checks

### `common:metrics` [P0] ✅ IMPLEMENTED
**Verify metrics instrumentation exists**

Mature applications expose metrics for monitoring and alerting.

| Attribute | Value |
|-----------|-------|
| Order | 1010 |
| Critical | No (Warn) |

**Detection approach:**
- Check for metrics libraries:
  - Go: `prometheus/client_golang`, `go-metrics`
  - Python: `prometheus_client`, `statsd`
  - Node: `prom-client`, `hot-shots`
- Check for `/metrics` endpoint in code
- Check for dashboards: `grafana/*.json`

---

### `common:tracing` [P1]
**Verify distributed tracing is implemented**

Mature microservices need request tracing for debugging.

| Attribute | Value |
|-----------|-------|
| Order | 1015 |
| Critical | No (Warn) |

**Detection approach:**
- Check for tracing libraries:
  - Go: `opentelemetry-go`, `jaeger-client-go`
  - Python: `opentelemetry`, `ddtrace`
  - Node: `@opentelemetry/sdk-trace-node`, `dd-trace`
- Check for OTEL config: `otel-collector-config.yaml`

---

### `common:errors` [P0] ✅ IMPLEMENTED
**Verify error tracking is configured**

Mature applications capture and report errors to monitoring services.

| Attribute | Value |
|-----------|-------|
| Order | 1020 |
| Critical | No (Warn) |

**Detection approach:**
- Check for error tracking SDKs:
  - Go: `getsentry/sentry-go`, `rollbar/rollbar-go`
  - Python: `sentry-sdk`, `rollbar`, `bugsnag`
  - Node: `@sentry/node`, `rollbar`
- Check for config: `.sentryclirc`, `sentry.properties`

---

### `common:logging_config` [P2]
**Verify log aggregation is configured**

Mature applications ship logs to centralized systems.

| Attribute | Value |
|-----------|-------|
| Order | 1025 |
| Critical | No (Warn) |

**Detection approach:**
- Check for log shipper configs: `fluentd.conf`, `filebeat.yml`, `vector.toml`
- Check for structured log config in app settings

---

## 5. Production Readiness Checks

### `common:k8s` [P0] ✅ IMPLEMENTED
**Verify Kubernetes manifests exist**

Mature cloud-native applications have declarative deployment configurations.

| Attribute | Value |
|-----------|-------|
| Order | 1030 |
| Critical | No (Warn) |

**Detection approach:**
- Check for directories: `k8s/`, `kubernetes/`, `deploy/`, `charts/`
- Check for files: `deployment.yaml`, `service.yaml`, files with `kind: Deployment`
- Check for Helm: `Chart.yaml`, `values.yaml`
- Check for Kustomize: `kustomization.yaml`

**Implementation:** `pkg/checks/common/k8s.go`

---

### `common:shutdown` [P0] ✅ IMPLEMENTED
**Verify graceful shutdown handling**

Mature applications handle SIGTERM/SIGINT to complete in-flight requests.

| Attribute | Value |
|-----------|-------|
| Order | 1035 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Search for `signal.Notify`, `os.Signal`, `syscall.SIGTERM`
- Python: Search for `signal.signal`, `atexit.register`
- Node: Search for `process.on('SIGTERM')`, `process.on('SIGINT')`

---

### `common:migrations` [P1]
**Verify database migrations are managed**

Mature applications version their database schema changes.

| Attribute | Value |
|-----------|-------|
| Order | 1040 |
| Critical | No (Warn) |

**Detection approach:**
- Check for directories: `migrations/`, `db/migrations/`, `alembic/`
- Check for migration tools:
  - Go: `golang-migrate`, `goose`, `atlas`
  - Python: `alembic/`, `django/migrations/`
  - Node: `prisma/migrations/`, `knex/migrations/`

---

### `common:config_validation` [P1]
**Verify configuration is validated at startup**

Mature applications fail fast with clear errors when misconfigured.

| Attribute | Value |
|-----------|-------|
| Order | 1045 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Check for `go-playground/validator`, `envconfig`
- Python: Check for `pydantic` Settings, `dynaconf`
- Node: Check for `joi`, `zod`, `convict`

---

### `common:retry` [P1]
**Verify retry logic exists for external calls**

Mature applications handle transient failures with retries and backoff.

| Attribute | Value |
|-----------|-------|
| Order | 1050 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Check for `cenkalti/backoff`, `avast/retry-go`, `go-retryablehttp`
- Python: Check for `tenacity`, `backoff`, `urllib3.util.retry`
- Node: Check for `async-retry`, `axios-retry`, `p-retry`

---

## 6. Architecture Quality Checks

### `go:errcheck` / `python:exceptions` / `node:errors` [P1]
**Validate proper error handling patterns**

Mature applications handle errors consistently.

| Attribute | Value |
|-----------|-------|
| Order | 280 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Check for ignored errors (err not checked)
- Python: Check for bare `except:` or `except Exception:` without re-raise
- Node: Check for empty catch blocks, unhandled promise rejections

---

### `common:di` [P2]
**Verify dependency injection patterns**

Mature applications use DI for testability and loose coupling.

| Attribute | Value |
|-----------|-------|
| Order | 1055 |
| Critical | No (Warn) |

**Detection approach:**
- Go: Check for `wire`, `dig`, `fx` imports
- Python: Check for `dependency-injector`, `fastapi.Depends`
- Node: Check for `inversify`, `tsyringe`, NestJS `@Injectable()`

---

### `common:architecture` [P2]
**Verify clean architecture structure**

Mature applications separate concerns into clear layers.

| Attribute | Value |
|-----------|-------|
| Order | 1060 |
| Critical | No (Warn) |

**Detection approach:**
- Check for layered directories:
  - Clean: `domain/`, `usecases/`, `interfaces/`, `infrastructure/`
  - Layered: `handlers/`, `services/`, `repositories/`, `models/`
  - Go standard: `cmd/`, `internal/`, `pkg/`

---

## 7. Development Workflow Checks

### `common:precommit` [P0] ✅ IMPLEMENTED
**Verify pre-commit hooks are configured**

Mature teams automate quality checks before code is committed.

| Attribute | Value |
|-----------|-------|
| Order | 1065 |
| Critical | No (Warn) |

**Detection approach:**
- Check for: `.pre-commit-config.yaml`, `.husky/`, `lefthook.yml`
- Check for husky config in `package.json`
- Check for commitlint: `commitlint.config.js`, `.commitlintrc`

**Implementation:** `pkg/checks/common/precommit.go`

---

### `common:editorconfig` [P1]
**Verify editor configuration exists**

Mature teams ensure consistent code style across all editors.

| Attribute | Value |
|-----------|-------|
| Order | 1070 |
| Critical | No (Warn) |

**Detection approach:**
- Check for: `.editorconfig`
- Check for IDE configs: `.vscode/settings.json`, `.idea/codeStyles/`

---

## Summary

| Priority | Count | Implemented | Categories |
|----------|-------|-------------|------------|
| P0 (High) | 10 | 9 | ✅ secrets, ✅ api_docs, ✅ changelog, ✅ integration, ✅ metrics, ✅ errors, ✅ k8s, ✅ shutdown, ✅ precommit |
| P1 (Medium) | 13 | 0 | license, env, sast, security, contributing, e2e, tracing, migrations, config_validation, retry, errcheck, editorconfig |
| P2 (Low) | 5 | 0 | conduct, godoc, fixtures, mocking, logging_config, di, architecture |

**Total: 28 new checks (9 implemented, 19 remaining)**

---

## Implementation Notes

1. All new checks should follow the existing pattern in `pkg/checks/`
2. Use `safepath` package for all file operations
3. Common checks go in `pkg/checks/common/`, language-specific in `pkg/checks/<lang>/`
4. Each check needs: implementation file, registration in `register.go`, tests
5. Consider adding a new config section for maturity-related thresholds
6. Update `docs/CHECKS.md` with the new check documentation
7. Update `backlog.md` marking the check as ✅ IMPLEMENTED
8. Update `register_test.go` to include the new check in test assertions
