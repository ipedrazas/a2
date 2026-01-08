# A2 Checks Reference

This document describes all checks available in A2, organized by language.

## Table of Contents

- [Go Checks](#go-checks)
- [Python Checks](#python-checks)
- [Node.js Checks](#nodejs-checks)
- [Java Checks](#java-checks)
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
| `go:cyclomatic` | Go Complexity | No | 240 | Analyzes cyclomatic complexity of functions |
| `go:logging` | Go Logging | No | 250 | Detects structured logging vs fmt.Print |

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

### go:cyclomatic

Analyzes cyclomatic complexity of Go functions using the go/ast package.

**Configuration:**
```yaml
language:
  go:
    cyclomatic_threshold: 15  # Default: 15
```

**Complexity Calculation:**
- Base complexity: 1
- +1 for each: `if`, `for`, `range`, `case` (non-default), `&&`, `||`

**Status:**
- **Pass**: No functions exceed threshold
- **Warn**: Functions exceed complexity threshold

**Fix:** Break complex functions into smaller, focused functions.

### go:logging

Checks for proper structured logging practices instead of fmt.Print statements.

**Structured loggers detected:**
- `log/slog` (Go 1.21+)
- `go.uber.org/zap`
- `github.com/rs/zerolog`
- `github.com/sirupsen/logrus`

**Anti-patterns detected (in non-test files):**
- `fmt.Print`, `fmt.Println`, `fmt.Printf`

**Status:**
- **Pass**: Uses structured logging, no fmt.Print statements
- **Warn**: Uses fmt.Print for logging or no structured logger detected

**Fix:** Use `log/slog` or other structured logging libraries.

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
| `python:complexity` | Python Complexity | No | 240 | Analyzes cyclomatic complexity with radon |
| `python:logging` | Python Logging | No | 250 | Detects structured logging vs print() |

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

### python:complexity

Analyzes cyclomatic complexity of Python functions using radon.

**Requirements:** Install with `pip install radon`

**Configuration:**
```yaml
language:
  python:
    cyclomatic_threshold: 15  # Default: 15
```

**Complexity Grades:**
- A (1-5): Low complexity
- B (6-10): Moderate complexity
- C (11-20): High complexity
- D/E/F (21+): Very high complexity

**Status:**
- **Pass**: No functions exceed threshold or radon not installed
- **Warn**: Functions exceed complexity threshold

**Fix:** Break complex functions into smaller, focused functions.

### python:logging

Checks for proper logging practices instead of print() statements.

**Logging modules detected:**
- `logging` (standard library)
- `structlog`
- `loguru`

**Anti-patterns detected (in non-test files):**
- `print()` statements

**Status:**
- **Pass**: Uses logging module, no print() statements
- **Warn**: Uses print() for logging or no logging module detected

**Fix:** Use the `logging` module or `structlog`/`loguru`.

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
| `node:logging` | Node Logging | No | 250 | Detects structured logging vs console.log |

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

### node:logging

Checks for proper structured logging practices instead of console.log.

**Logging libraries detected (in package.json):**
- `winston`
- `pino`
- `bunyan`
- `log4js`
- `loglevel`
- `signale`

**Anti-patterns detected (in non-test files):**
- `console.log`, `console.error`, `console.warn`, `console.info`, `console.debug`

**Status:**
- **Pass**: Uses structured logging, no console.* statements
- **Warn**: Uses console.log for logging or no structured logger detected

**Fix:** Use `winston`, `pino`, or another structured logging library.

---

## Java Checks

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `java:project` | Java Project | Yes | 100 | Detects Maven (pom.xml) or Gradle (build.gradle) projects |
| `java:build` | Java Build | Yes | 110 | Compiles with Maven or Gradle (auto-detects wrapper scripts) |
| `java:tests` | Java Tests | Yes | 120 | Runs JUnit/TestNG tests |
| `java:format` | Java Format | No | 200 | Detects Spotless, google-java-format, EditorConfig, IDE formatters |
| `java:lint` | Java Lint | No | 210 | Detects Checkstyle, SpotBugs, PMD, Error Prone, SonarQube |
| `java:coverage` | Java Coverage | No | 220 | Checks JaCoCo coverage reports against threshold |
| `java:deps` | Java Dependencies | No | 230 | Detects OWASP Dependency-Check, Snyk, Dependabot, Renovate |
| `java:logging` | Java Logging | No | 250 | Detects SLF4J, Logback, Log4j2; warns on System.out.println |

### java:project

Detects Maven or Gradle project configuration.

**Build tool detection priority:**
1. `build.gradle` or `build.gradle.kts` → Gradle (Groovy or Kotlin DSL)
2. `pom.xml` → Maven
3. `gradlew` → Gradle wrapper
4. `mvnw` → Maven wrapper

When both Maven and Gradle are present, Gradle is preferred as the more modern build system.

**Wrapper detection:**
- `mvnw` (Maven Wrapper)
- `gradlew` (Gradle Wrapper)

**Status:**
- **Pass**: pom.xml or build.gradle found
- **Fail**: No Java project configuration found

### java:build

Compiles the Java project using the detected build tool.

**Configuration:**
```yaml
language:
  java:
    build_tool: auto  # auto, maven, gradle
```

**Build commands:**
- Maven: `./mvnw compile -q -DskipTests` (or `mvn` if no wrapper)
- Gradle: `./gradlew compileJava -q --no-daemon` (or `gradle` if no wrapper)

Wrapper scripts are preferred when available for reproducible builds.

**Status:**
- **Pass**: Compilation succeeds
- **Fail**: Compilation fails or no build tool found

### java:tests

Runs tests using the detected build tool and test runner.

**Configuration:**
```yaml
language:
  java:
    test_runner: auto  # auto, junit, testng
```

**Test commands:**
- Maven: `./mvnw test -q` (or `mvn test -q`)
- Gradle: `./gradlew test --no-daemon` (or `gradle test`)

**Status:**
- **Pass**: All tests pass
- **Warn**: No tests found
- **Fail**: Tests fail

### java:format

Detects code formatting configuration for Java projects.

**Formatters detected:**
- google-java-format: `google-java-format` in pom.xml/build.gradle, or `google-java-format.jar`
- Spotless: `spotless` plugin in pom.xml or build.gradle
- EditorConfig: `.editorconfig` with Java settings (`[*.java]`)
- IntelliJ IDEA: `.idea/codeStyles/` directory
- Eclipse: `.settings/org.eclipse.jdt.core.prefs`

**Status:**
- **Pass**: Formatter configuration found
- **Warn**: No formatter configuration found

**Recommendation:** Configure Spotless or google-java-format for consistent code formatting.

### java:lint

Detects static analysis tools configured for the project.

**Tools detected:**
- Checkstyle: `checkstyle.xml`, `checkstyle` plugin in pom.xml/build.gradle
- SpotBugs: `spotbugs.xml`, `spotbugsInclude.xml`, `spotbugs` plugin
- PMD: `pmd.xml`, `ruleset.xml`, `pmd` plugin
- Error Prone: `errorprone` in pom.xml/build.gradle
- SonarQube: `sonar-project.properties`, `sonarqube` plugin

**Status:**
- **Pass**: One or more linting tools configured
- **Warn**: No linting tools found

**Recommendation:** Configure Checkstyle, SpotBugs, or PMD for static analysis.

### java:coverage

Checks test coverage using JaCoCo reports.

**Configuration:**
```yaml
language:
  java:
    coverage_threshold: 80  # Default: 80%
```

**Report locations:**
- Maven: `target/site/jacoco/jacoco.xml`
- Gradle: `build/reports/jacoco/test/jacocoTestReport.xml`

Parses LINE or INSTRUCTION counter types from the XML report.

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold or JaCoCo not configured/no report found

**Recommendation:** Add JaCoCo plugin to pom.xml or build.gradle.

### java:deps

Detects dependency vulnerability scanning configuration.

**Tools detected:**
- OWASP Dependency-Check: `dependency-check` plugin in pom.xml/build.gradle
- Snyk: `.snyk` file, `snyk` in CI config
- Dependabot: `.github/dependabot.yml` with maven/gradle ecosystem
- Renovate: `renovate.json`, `renovate.json5`, `.renovaterc`

**Maven plugins checked:**
- `org.owasp:dependency-check-maven`
- `snyk-maven-plugin`

**Gradle plugins checked:**
- `org.owasp.dependencycheck`
- `io.snyk.gradle.plugin`

**Status:**
- **Pass**: Dependency scanning tool configured
- **Warn**: No dependency scanning found

**Recommendation:** Configure OWASP Dependency-Check or Snyk for vulnerability scanning.

### java:logging

Checks for structured logging practices instead of System.out.println.

**Logging frameworks detected:**
- SLF4J: `slf4j` dependency in pom.xml/build.gradle
- Logback: `logback.xml`, `logback-spring.xml`, or `logback` dependency
- Log4j2: `log4j2.xml`, `log4j2.yaml`, `log4j2.properties`, or `log4j` dependency

**Anti-patterns detected (in src/main/java/):**
- `System.out.print` statements
- `System.err.print` statements

Comments are excluded from detection.

**Status:**
- **Pass**: Structured logging library found, no System.out.println
- **Warn**: System.out.println found, or no structured logging detected

**Recommendation:** Use SLF4J with Logback or Log4j2 for structured logging.

---

## Common Checks

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `file_exists` | Required Files | No | 900 | Checks for required files |
| `common:dockerfile` | Container Ready | No | 910 | Checks for Dockerfile/Containerfile |
| `common:ci` | CI Pipeline | No | 920 | Detects CI/CD configuration |
| `common:health` | Health Endpoint | No | 930 | Detects health check endpoints |
| `common:secrets` | Secrets Detection | No | 940 | Detects secret scanning config or hardcoded secrets |
| `common:env` | Environment Config | No | 945 | Validates environment variable handling practices |
| `common:api_docs` | API Documentation | No | 960 | Detects OpenAPI/Swagger specs and documentation generators |
| `common:changelog` | Changelog | No | 965 | Verifies changelog or release notes exist |
| `common:integration` | Integration Tests | No | 980 | Detects integration test directories, files, and E2E frameworks |
| `common:metrics` | Metrics Instrumentation | No | 1010 | Detects Prometheus, OpenTelemetry, and other metrics libraries |
| `common:errors` | Error Tracking | No | 1020 | Detects Sentry, Rollbar, Bugsnag, and other error tracking SDKs |
| `common:k8s` | Kubernetes Ready | No | 1030 | Detects Kubernetes manifests and deployment configurations |
| `common:shutdown` | Graceful Shutdown | No | 1035 | Detects signal handling and graceful shutdown configuration |
| `common:precommit` | Pre-commit Hooks | No | 1065 | Verifies pre-commit hooks are configured |

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

**Files scanned:**
- Code files: `*.go`, `*.py`, `*.js`, `*.ts`, `*.java`, `*.rb`, `*.php`, `*.cs`
- Config files: `*.yaml`, `*.yml`, `*.json`, `*.xml`, `*.config`
- Shell scripts: `*.sh`, `*.bash`
- Environment files: `.env` (but not `.env.example`)

**Directories skipped:**
- `node_modules/`, `vendor/`, `.git/`, `__pycache__/`, `.venv/`, `venv/`, `dist/`, `build/`

**Status:**
- **Pass**: Secret scanning tool configured
- **Warn**: Potential secrets found, or no secret scanning configured

**Recommendation:** Configure Gitleaks or similar tool for automated secret scanning.

### common:env

Validates environment variable handling practices to ensure configuration is properly separated from code.

**Environment template files detected:**
- `.env.example`, `.env.sample`, `.env.template`, `example.env`, `.env.local.example`

**Dotenv libraries detected:**
- Go: `github.com/joho/godotenv`, `github.com/caarlos0/env`, `github.com/kelseyhightower/envconfig`, `github.com/spf13/viper`
- Python: `python-dotenv`, `environs`, `pydantic-settings`, `django-environ`, `dynaconf`
- Node.js: `dotenv`, `dotenv-safe`, `dotenv-expand`, `env-cmd`, `cross-env`
- Java: Spring Boot (`application.properties`, `application.yml`), `dotenv-java`

**Gitignore detection:**
- Checks if `.env` is listed in `.gitignore`
- Recognizes patterns: `.env`, `.env*`, `.env.*`, `*.env`, `.env.local`

**Issues detected:**
- `.env` file exists but not in `.gitignore` (potential secret exposure)

**Status:**
- **Pass**: Environment template found, dotenv library detected, or `.env` properly gitignored
- **Warn**: `.env` exists but not gitignored, or no environment configuration found

**Recommendation:** Create `.env.example` to document required environment variables, and ensure `.env` is in `.gitignore`.

### common:api_docs

Detects API documentation configuration and files.

**OpenAPI/Swagger files detected:**
- `openapi.yaml`, `openapi.yml`, `openapi.json`
- `swagger.yaml`, `swagger.yml`, `swagger.json`
- `api.yaml`, `api.yml`, `api.json`
- Searched in root and `docs/`, `api/` directories

**Documentation directories detected:**
- `docs/api/`, `api-docs/`, `api/docs/`, `swagger/`, `swagger-ui/`

**Documentation generators detected:**
- Go: `swaggo/swag`, `go-swagger`, `grpc-gateway`
- Python: `drf-spectacular`, `drf-yasg`, `flasgger`, `flask-restx`, `FastAPI`
- Node.js: `swagger-jsdoc`, `swagger-ui-express`, `@nestjs/swagger`, `tsoa`

**Also detected:**
- GraphQL schemas: `schema.graphql`, `schema.gql`
- Protocol Buffers: `*.proto` files

**Status:**
- **Pass**: API documentation or generator found
- **Warn**: No API documentation found

**Recommendation:** Add OpenAPI/Swagger specification or use a documentation generator.

### common:changelog

Verifies that a changelog or release notes file exists, and detects release tooling configuration.

**Changelog files detected:**
- `CHANGELOG.md`, `CHANGELOG.txt`, `CHANGELOG`
- `CHANGES.md`, `CHANGES.txt`, `CHANGES`
- `HISTORY.md`, `HISTORY.txt`, `HISTORY`
- `NEWS.md`, `NEWS.txt`, `NEWS`
- `RELEASES.md`, `RELEASE_NOTES.md`

**Release tooling detected:**
- GoReleaser: `.goreleaser.yml`, `.goreleaser.yaml`
- semantic-release: `.releaserc`, `.releaserc.json`, `release.config.js`
- release-please: `release-please-config.json`, `.release-please-manifest.json`
- standard-version: `.versionrc`, `.versionrc.json`
- changesets: `.changeset/config.json`

**Changelog format detection:**
- Keep a Changelog: Detects `## [Unreleased]`, `### Added`, `### Changed`, `### Fixed`, etc.
- Conventional Changelog: Detects `### Features`, `### Bug Fixes`, `feat:`, `fix:`
- Plain markdown or text

**Status:**
- **Pass**: Changelog file found or release tooling configured
- **Warn**: No changelog or release tooling found

**Recommendation:** Create a `CHANGELOG.md` following [Keep a Changelog](https://keepachangelog.com) format.

### common:integration

Detects integration test directories, files, and E2E testing frameworks.

**Integration test directories detected:**
- `tests/integration/`, `test/integration/`
- `integration_tests/`, `integration-tests/`
- `tests/e2e/`, `test/e2e/`, `e2e/`, `e2e-tests/`

**Integration test file patterns:**
- Go: `*_integration_test.go`, `*_integ_test.go`
- Python: `test_integration_*.py`, `test_*_integration.py`
- Node.js/TypeScript: `*.integration.test.ts`, `*.integration.spec.ts`, `*.e2e.test.ts`

**Test infrastructure detected:**
- Docker Compose for tests: `docker-compose.test.yml`, `docker-compose.e2e.yml`
- testcontainers library usage

**E2E testing frameworks detected:**
- Cypress: `cypress.config.js`, `cypress.config.ts`
- Playwright: `playwright.config.js`, `playwright.config.ts`
- WebdriverIO: `wdio.conf.js`, `wdio.conf.ts`
- Selenium: `selenium-webdriver` in dependencies

**Status:**
- **Pass**: Integration tests or E2E framework found
- **Warn**: No integration tests found

**Recommendation:** Add integration tests in `tests/integration/` directory.

### common:metrics

Detects metrics instrumentation libraries and configuration.

**Metrics libraries detected:**
- Go: `prometheus/client_golang`, `go-metrics`, `opentelemetry-go`, `datadog/dd-trace-go`
- Python: `prometheus_client`, `statsd`, `opentelemetry`, `ddtrace`
- Node.js: `prom-client`, `hot-shots`, `@opentelemetry/sdk-metrics`, `dd-trace`
- Java: `micrometer`, `prometheus`, `dropwizard-metrics`, `opentelemetry`

**Configuration files detected:**
- `prometheus.yml`, `prometheus.yaml`
- `otel-collector-config.yaml`
- `metrics.yml`, `metrics.yaml`

**Also detected:**
- Grafana dashboards in `grafana/` or `dashboards/` directories
- `/metrics` endpoint patterns in source code

**Status:**
- **Pass**: Metrics library or configuration found
- **Warn**: No metrics instrumentation found

**Recommendation:** Add Prometheus client or OpenTelemetry for metrics instrumentation.

### common:errors

Detects error tracking SDK configuration.

**Error tracking SDKs detected:**
- Go: `getsentry/sentry-go`, `rollbar/rollbar-go`, `bugsnag/bugsnag-go`
- Python: `sentry-sdk`, `rollbar`, `bugsnag`, `honeybadger`
- Node.js: `@sentry/node`, `rollbar`, `@bugsnag/js`, `@honeybadger-io/js`
- Java: `sentry`, `rollbar`, `bugsnag`

**Configuration files detected:**
- `.sentryclirc`, `sentry.properties`
- `.rollbar`, `bugsnag.json`

**Environment variables detected (in .env.example):**
- `SENTRY_DSN`, `ROLLBAR_TOKEN`, `BUGSNAG_API_KEY`

**CI integration detected:**
- Sentry release uploads in GitHub Actions

**Status:**
- **Pass**: Error tracking SDK or configuration found
- **Warn**: No error tracking found

**Recommendation:** Add Sentry, Rollbar, or Bugsnag for error tracking.

### common:k8s

Detects Kubernetes manifests, Helm charts, and other deployment configurations for cloud-native applications.

**Helm charts detected:**
- `Chart.yaml` in root directory
- `charts/*/Chart.yaml` (chart repositories)
- `helm/Chart.yaml` or `helm/*/Chart.yaml`

**Kustomize configurations detected:**
- `kustomization.yaml`, `kustomization.yml`, `Kustomization`
- Searched in root and common directories: `k8s/`, `kubernetes/`, `deploy/`, `base/`, `overlays/`

**Kubernetes manifest directories detected:**
- `k8s/`, `kubernetes/`, `deploy/`, `deployment/`, `deployments/`, `manifests/`, `.k8s/`
- Scans for YAML files containing `apiVersion:` and `kind: Deployment`, `kind: Service`, etc.

**Standalone manifest files detected (in root):**
- `deployment.yaml`, `service.yaml`, `ingress.yaml`, `configmap.yaml`
- `secret.yaml`, `pod.yaml`, `statefulset.yaml`, `daemonset.yaml`
- `cronjob.yaml`, `job.yaml` (and `.yml` variants)

**Alternative deployment tools detected:**
- Docker Compose: `docker-compose.yaml`, `compose.yaml`, `docker-compose.prod.yaml`
- Skaffold: `skaffold.yaml`, `skaffold.yml`
- Tilt: `Tiltfile`

**Kubernetes resource types detected:**
- Deployment, Service, ConfigMap, Secret, Ingress, Pod
- StatefulSet, DaemonSet, Job, CronJob
- PersistentVolume, PersistentVolumeClaim, Namespace
- ServiceAccount, Role, RoleBinding, ClusterRole, ClusterRoleBinding
- HorizontalPodAutoscaler, NetworkPolicy

**Status:**
- **Pass**: Kubernetes manifests, Helm chart, Kustomize, or alternative deployment config found
- **Warn**: No deployment configuration found

**Recommendation:** Configure Kubernetes manifests, Helm charts, or Docker Compose for reproducible deployments.

### common:shutdown

Detects graceful shutdown handling for proper process termination.

**Go signal handling detected:**
- `signal.Notify`, `os.Signal`, `syscall.SIGTERM`, `syscall.SIGINT`
- `server.Shutdown`, `context.WithCancel`
- Graceful shutdown packages: `oklog/run`, `errgroup`

**Python signal handling detected:**
- `signal.signal`, `signal.SIGTERM`, `signal.SIGINT`
- `atexit.register`
- `asyncio` signal handlers

**Node.js signal handling detected:**
- `process.on('SIGTERM')`, `process.on('SIGINT')`
- `process.on('beforeExit')`, `process.on('exit')`
- Graceful shutdown packages: `http-terminator`, `stoppable`, `@godaddy/terminus`, `lightship`

**Java shutdown hooks detected:**
- `Runtime.getRuntime().addShutdownHook`
- `@PreDestroy` annotation
- `DisposableBean`, `SmartLifecycle` Spring interfaces

**Kubernetes lifecycle hooks detected:**
- `terminationGracePeriodSeconds` in manifests
- `preStop` lifecycle hooks

**Status:**
- **Pass**: Signal handling or shutdown hooks found
- **Warn**: No graceful shutdown handling found

**Recommendation:** Handle SIGTERM/SIGINT signals to complete in-flight requests before shutdown.

### common:precommit

Verifies that pre-commit hooks are configured to automate quality checks before code is committed.

**Pre-commit tools detected:**
- pre-commit (Python): `.pre-commit-config.yaml`, `.pre-commit-config.yml`
- Husky (Node.js): `.husky/` directory with hook files, or `husky` in package.json
- Lefthook: `lefthook.yml`, `lefthook.yaml`, `.lefthook.yml`
- Overcommit (Ruby): `.overcommit.yml`
- Native git hooks: `.git/hooks/` with executable hook files

**Related tooling detected:**
- commitlint: `commitlint.config.js`, `.commitlintrc`, `.commitlintrc.json`, or `commitlint` in package.json
- lint-staged: `lint-staged.config.js`, `.lintstagedrc`, or `lint-staged` in package.json

**Hook files checked:**
- `pre-commit`, `pre-push`, `commit-msg`, `prepare-commit-msg`

**Status:**
- **Pass**: Pre-commit hooks or related tooling configured
- **Warn**: No pre-commit hooks found

**Recommendation:** Configure pre-commit, Husky, or Lefthook for automated quality checks.

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

  java:
    build_tool: auto         # auto, maven, gradle
    test_runner: auto        # auto, junit, testng
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
| Go | 10 | 3 | 7 |
| Python | 10 | 3 | 7 |
| Node.js | 9 | 3 | 6 |
| Java | 8 | 3 | 5 |
| Common | 14+ | 0 | 14+ |
| **Total** | **51+** | **12** | **39+** |

**Critical checks** stop execution in sequential mode when they fail.
**Non-critical checks** report warnings but allow other checks to continue.
