# Node.js Checks

This document describes all Node.js-specific checks available in A2.

## Overview

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

---

## node:project

Verifies that `package.json` exists and contains required fields.

**Status:**
- **Pass**: package.json with valid name and version
- **Warn**: Missing version field
- **Fail**: Missing package.json or missing name field

---

## node:build

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

---

## node:tests

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

---

## node:format

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

---

## node:lint

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

---

## node:type

Runs TypeScript compiler for type checking. Only activates for TypeScript projects.

**TypeScript project detection:**
- `tsconfig.json` or variants (`tsconfig.base.json`, etc.)
- TypeScript in devDependencies or dependencies

**Status:**
- **Pass**: No type errors or not a TypeScript project
- **Warn**: Type errors found

**Fix:** Run `npx tsc --noEmit`

---

## node:coverage

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

---

## node:deps

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

## node:logging

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

## Configuration Example

```yaml
language:
  node:
    package_manager: auto    # auto, npm, yarn, pnpm, bun
    test_runner: auto        # auto, jest, vitest, mocha, npm-test
    formatter: auto          # auto, prettier, biome
    linter: auto             # auto, eslint, biome, oxlint
    coverage_threshold: 80
```
