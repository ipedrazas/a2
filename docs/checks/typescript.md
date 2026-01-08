# TypeScript Checks

This document describes all TypeScript-specific checks available in A2.

## Overview

| Check ID | Name | Critical | Order | Description |
|----------|------|----------|-------|-------------|
| `typescript:project` | TypeScript Project | Yes | 100 | Detects tsconfig.json, extracts compiler options and TypeScript version |
| `typescript:build` | TypeScript Build | Yes | 110 | Runs build script or `tsc --noEmit` to verify compilation |
| `typescript:tests` | TypeScript Tests | Yes | 120 | Runs tests with Jest, Vitest, or Mocha (auto-detects) |
| `typescript:format` | TypeScript Format | No | 200 | Checks code formatting with Prettier, Biome, or dprint |
| `typescript:lint` | TypeScript Lint | No | 210 | Runs ESLint, Biome, or oxlint for linting |
| `typescript:type` | TypeScript Type Check | Yes | 215 | Type checking with `tsc --noEmit` (critical check) |
| `typescript:coverage` | TypeScript Coverage | No | 220 | Detects coverage tools (Jest, Vitest, c8, nyc) |
| `typescript:deps` | TypeScript Vulnerabilities | No | 230 | Checks for vulnerabilities using npm/yarn/pnpm audit |
| `typescript:logging` | TypeScript Logging | No | 250 | Detects logging libraries (winston, pino, tslog), warns on console.log |

---

## typescript:project

Detects TypeScript project configuration and extracts compiler settings.

**Config files detected:**
- `tsconfig.json`
- `tsconfig.base.json`
- `tsconfig.build.json`

**Information extracted:**
- Target ECMAScript version
- Module system
- Strict mode enabled
- TypeScript version from package.json

**Status:**
- **Pass**: tsconfig.json found with valid configuration
- **Fail**: No tsconfig.json found

---

## typescript:build

Runs the build process to verify TypeScript compilation.

**Build detection:**
1. `build` script in package.json → runs via package manager
2. Falls back to `tsc --noEmit`

**Package managers supported:**
- npm, yarn, pnpm, bun

**Status:**
- **Pass**: Build completes without errors
- **Fail**: Compilation errors

---

## typescript:tests

Runs TypeScript tests using the detected test runner.

**Configuration:**
```yaml
language:
  typescript:
    test_runner: auto  # auto, jest, vitest, mocha
```

**Auto-detection:**
- `jest.config.*` → Jest
- `vitest.config.*` → Vitest
- `.mocharc.*` → Mocha
- Checks devDependencies
- Falls back to `npm test`

**Status:**
- **Pass**: All tests pass or no tests found
- **Fail**: Tests fail

---

## typescript:format

Checks code formatting using the detected formatter.

**Configuration:**
```yaml
language:
  typescript:
    formatter: auto  # auto, prettier, biome, dprint
```

**Formatters detected:**
- Prettier: `.prettierrc*`, `prettier.config.*`
- Biome: `biome.json`, `biome.jsonc`
- dprint: `dprint.json`, `.dprint.json`
- Also checks devDependencies

**Status:**
- **Pass**: All files formatted correctly
- **Warn**: Formatting issues found or no formatter configured

**Fix:** Run `npx prettier --write .` or `npx @biomejs/biome format --write .`

---

## typescript:lint

Runs linting using the detected linter.

**Configuration:**
```yaml
language:
  typescript:
    linter: auto  # auto, eslint, biome, oxlint
```

**Linters detected:**
- ESLint: `.eslintrc*`, `eslint.config.*`
- Biome: `biome.json`, `biome.jsonc`
- oxlint: `oxlint.json`, `.oxlintrc.json`
- Also checks devDependencies

**Status:**
- **Pass**: No linting issues
- **Warn**: Linting issues found or no linter configured

---

## typescript:type

Runs TypeScript compiler for type checking. This is a critical check.

**Command:** `tsc --noEmit`

**Status:**
- **Pass**: No type errors
- **Fail**: Type errors found

**Fix:** Run `npx tsc --noEmit` and fix reported errors.

---

## typescript:coverage

Detects and runs test coverage tools.

**Configuration:**
```yaml
language:
  typescript:
    coverage_threshold: 80  # Default: 80%
```

**Coverage tools detected:**
- Jest (`jest --coverage`)
- Vitest (`vitest run --coverage`)
- c8 (native V8 coverage)
- nyc (Istanbul)

**Status:**
- **Pass**: Coverage meets or exceeds threshold
- **Warn**: Coverage below threshold or tools not installed

---

## typescript:deps

Scans for security vulnerabilities in dependencies.

**Commands (based on package manager):**
- npm: `npm audit --json`
- yarn: `yarn audit --json`
- pnpm: `pnpm audit --json`
- bun: (skipped - no built-in audit)

**Status:**
- **Pass**: No vulnerabilities found
- **Warn**: Vulnerabilities detected

---

## typescript:logging

Checks for structured logging practices instead of console.log.

**Logging libraries detected (in package.json):**
- `winston`
- `pino`
- `bunyan`
- `log4js`
- `loglevel`
- `signale`
- `tslog`
- `@sentry/node`, `@sentry/react`
- `dd-trace` (Datadog)
- `newrelic`

**Anti-patterns detected (in non-test files):**
- `console.log`, `console.info`, `console.warn`, `console.error`, `console.debug`

**Directories skipped:**
- `node_modules`, `dist`, `build`, `.git`, `coverage`, `__tests__`

**Status:**
- **Pass**: Logging library found, no console.log calls
- **Warn**: console.log found or no logging library detected

**Fix:** Use `winston`, `pino`, or `tslog` for structured logging.

---

## Configuration Example

```yaml
language:
  typescript:
    package_manager: auto    # auto, npm, yarn, pnpm, bun
    test_runner: auto        # auto, jest, vitest, mocha
    formatter: auto          # auto, prettier, biome, dprint
    linter: auto             # auto, eslint, biome, oxlint
    coverage_threshold: 80
```
