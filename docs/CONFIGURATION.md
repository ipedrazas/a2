# A2 Configuration Reference

A2 is configured via a `.a2.yaml` file in your project root. You can generate one with `a2 init` (use `-i` for interactive mode).

## Table of Contents

- [Generating Configuration](#generating-configuration)
- [Configuration Examples](#configuration-examples)
- [Source Directories](#source-directories)
- [External Checks](#external-checks)
- [Wildcard Patterns](#wildcard-patterns)
- [Per-Language Disabling](#per-language-disabling)

---

## Generating Configuration

### Interactive Mode

Run with `-i` flag for guided prompts:

```bash
a2 init -i
```

This will prompt you for:
1. **Application profile** (cli, api, library, desktop)
2. **Maturity target** (poc, production)
3. **Languages** (auto-detected, with option to modify)
4. **Required files** (README.md, LICENSE, etc.)
5. **Coverage threshold** (default: 80%)

Shows a preview before writing the file.

### Non-Interactive Mode

Pass options directly via flags:

```bash
# Basic usage
a2 init --profile cli --target poc

# With language and coverage
a2 init --lang go,python --coverage 90

# Custom required files
a2 init --files README.md,LICENSE,CHANGELOG.md

# Overwrite existing config
a2 init --profile api --force

# Custom output path
a2 init --output custom-config.yaml
```

### Available Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-i, --interactive` | Run in interactive mode | `false` |
| `--profile` | Application profile (cli, api, library, desktop) | - |
| `--target` | Maturity target (poc, production) | - |
| `--lang` | Languages (go, python, node, java, rust, typescript, swift) | auto-detect |
| `--files` | Required files (comma-separated) | README.md,LICENSE |
| `--coverage` | Coverage threshold (0-100) | 80 |
| `-o, --output` | Output file path | .a2.yaml |
| `-f, --force` | Overwrite existing file | `false` |

---

## Configuration Examples

### Go Project

```yaml
# Language settings
language:
  go:
    coverage_threshold: 80

# Required files to check
files:
  required:
    - README.md
    - LICENSE
    - CONTRIBUTING.md

# Disable specific checks
checks:
  disabled:
    - go:deps  # Skip vulnerability check
    - "*:tests"  # Skip all test checks across all languages
    - "node:*"   # Skip all Node.js checks
    - "*:logging"  # Skip logging checks for all languages

# Execution options
execution:
  parallel: true  # Run checks concurrently (default)

# Custom external checks
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
```

### Python Project

```yaml
language:
  python:
    package_manager: auto  # auto, pip, poetry, pipenv
    test_runner: auto      # auto, pytest, unittest
    formatter: auto        # auto, black, ruff
    linter: auto           # auto, ruff, flake8, pylint
    coverage_threshold: 75

files:
  required:
    - README.md
    - LICENSE
    - pyproject.toml

checks:
  disabled:
    - python:deps

external:
  - id: security
    name: Security Scan
    command: bandit
    args: ["-r", "src/"]
    severity: warn
```

### Node.js/TypeScript Project

```yaml
language:
  node:
    package_manager: auto  # auto, npm, yarn, pnpm, bun
    test_runner: auto      # auto, jest, vitest, mocha
    formatter: auto        # auto, prettier, biome
    linter: auto           # auto, eslint, biome, oxlint
    coverage_threshold: 80
  typescript:
    coverage_threshold: 80

files:
  required:
    - README.md
    - LICENSE
    - package.json

checks:
  disabled:
    - node:deps
```

### Multi-Language Project (Monorepo)

In a monorepo, each language usually lives in one or more subdirectories rather
than at the repo root. List every component directory under that language's
`source_dir` (see [Source Directories](#source-directories) for the full
syntax). For example, three Go components and two TypeScript components:

```yaml
# Explicit language selection (overrides auto-detect)
language:
  explicit:
    - go
    - typescript
  go:
    source_dir:           # a list — one entry per Go component
      - services/api
      - services/worker
      - cmd/cli
    coverage_threshold: 80
  typescript:
    source_dir:           # a list — one entry per TypeScript component
      - web/app
      - web/admin
    coverage_threshold: 75

files:
  required:
    - README.md
```

Each listed directory is treated as an independent root for that language's
checks (build, tests, coverage, format, …).

---

## Source Directories

By default A2 runs each language's checks from the repository root. When a
language's code lives in one or more subdirectories — common in monorepos — set
`source_dir` under that language so checks run against the right subtree(s).

`source_dir` is **not** limited to a single path. It accepts three forms:

### 1. Single path

```yaml
language:
  go:
    source_dir: backend
```

### 2. List of paths

Use a list when a language has multiple components. Each path becomes an
independent root for that language's checks.

```yaml
language:
  go:
    source_dir:
      - services/api
      - services/worker
      - cmd/cli
  typescript:
    source_dir:
      - web/app
      - web/admin
```

The inline form `source_dir: [services/api, services/worker]` is equivalent.

### 3. List of objects (per-directory profile / threshold)

Use the object form when individual components need different rules. Each entry
is a `{path, profile, coverage_threshold}` object, where `profile` and
`coverage_threshold` are optional.

```yaml
language:
  go:
    coverage_threshold: 80        # language-wide default
    source_dir:
      - path: services/api
        profile: api              # api component: keep all operational checks
      - path: services/worker
        profile: api
        coverage_threshold: 60    # override just for this directory
      - path: pkg/shared
        profile: library          # library component: skip server/devops checks
```

| Field | Required | Description |
|-------|----------|-------------|
| `path` | yes | Subdirectory containing the component's code |
| `profile` | no | Name of a profile whose `disabled` list is applied to **this directory only** |
| `coverage_threshold` | no | Per-directory coverage threshold, overriding the language-wide value |

**Profiles** are the same ones used by `a2 init --profile`. Built-in profiles are
`cli`, `api`, `library`, and `desktop`; you can also define your own (run
`a2 profiles` to list them). Attaching a profile to a directory disables that profile's checks
for that directory without affecting the others — so a `library` package and an
`api` service can live in the same language block with different rules.

> Each language's `source_dir` is independent: Go can use the object form while
> TypeScript uses a plain list. Within a *single* `source_dir`, however, all
> entries must use the same form — either all bare paths or all `{path: ...}`
> objects, not a mix of the two.

You can pair `source_dir` with [Per-Language Disabling](#per-language-disabling)
so that repo-wide `common:*` / `devops:*` checks are re-scoped to the correct
subtree per language.

---

## External Checks

A2 supports external check binaries. The protocol is simple:

- **Exit code 0**: Pass
- **Exit code 1**: Warning
- **Exit code 2+**: Fail

Output can be plain text or JSON:

```json
{
  "message": "Found 3 issues",
  "status": "warn"
}
```

Define external checks in `.a2.yaml`:

```yaml
external:
  - id: lint
    name: Golangci-lint
    command: golangci-lint
    args: ["run", "./..."]
    severity: warn
```

---

## Wildcard Patterns

Skip patterns support wildcards for flexible check filtering. Use wildcards in `.a2.yaml` or with the `--skip` CLI flag.

**Important:** In `.a2.yaml`, patterns containing `*` must be quoted (e.g. `"*:logging"` or `'*:logging'`). Unquoted `*` is interpreted by YAML as an alias and will cause a parse error.

### Pattern Syntax

| Pattern | Description | Example Matches |
|---------|-------------|-----------------|
| `*:suffix` | Match any check ending with `:suffix` | `*:tests` matches `go:tests`, `python:tests`, `node:tests` |
| `prefix:*` | Match any check starting with `prefix:` | `node:*` matches `node:deps`, `node:logging`, `node:tests` |
| `*:*` | Match any check with a colon (all checks) | `*:*` matches all `lang:check` format checks |
| `*` | Match everything | `*` matches all check IDs |
| Exact | Exact string match | `go:tests` matches only `go:tests` |

### Examples

**Skip all test checks:**
```yaml
checks:
  disabled:
    - "*:tests"  # Skips go:tests, python:tests, node:tests, etc.
```

**Skip all Node.js checks:**
```yaml
checks:
  disabled:
    - "node:*"  # Skips node:deps, node:logging, node:tests, etc.
```

**Skip specific check types across all languages:**
```yaml
checks:
  disabled:
    - "*:logging"    # Skip logging checks
    - "*:tests"      # Skip test checks
    - "*:deps"       # Skip dependency vulnerability checks
```

**Combine wildcards with exact matches:**
```yaml
checks:
  disabled:
    - "common:*"     # Skip all common checks
    - "node:*"       # Skip all Node.js checks
    - "go:race"      # Also skip specific Go race detection
```

**CLI flag examples:**
```bash
# Skip all test checks
a2 check --skip="*:tests"

# Skip all Node.js checks
a2 check --skip="node:*"

# Skip multiple patterns
a2 check --skip="*:tests" --skip="*:license"
```

### Notes

- Wildcards work with both `.a2.yaml` configuration and `--skip` CLI flags
- Aliases (e.g., `tests` -> `go:tests`) continue to work alongside wildcards
- Invalid patterns with multiple wildcards (e.g., `*:*:*`) are treated as non-matching

---

## Per-Language Disabling

The top-level `checks.disabled` list applies to **all** languages. In a
monorepo you often want a check to apply to one language but not another — for
example, requiring metrics instrumentation in a Go backend but not in a
TypeScript frontend.

Add a language-keyed block under `checks:` with its own `disabled:` list. The
same wildcard and alias rules apply.

```yaml
checks:
  disabled:            # applies to every language
    - go:logging
    - "devops:*"
  typescript:
    disabled:          # applies only when TypeScript is active
      - common:metrics
  go:
    disabled:
      - common:tracing
```

Valid language keys are: `go`, `python`, `node`, `java`, `rust`, `typescript`,
`swift`. An unknown key (e.g. a typo like `typscript`) is rejected with an error.

### How language-agnostic checks are scoped

Most `common:*`, `devops:*`, and `security:*` checks normally run **once at the
repo root**. When you disable such a check for only *some* of the detected
languages, A2 re-scopes it instead of dropping it:

- It runs against the `source_dir`(s) of the languages that **still** want it.
- It is skipped for the languages where it is disabled.
- If it is disabled for **every** detected language, it is removed entirely.
- If no per-language list mentions it, it keeps running once at the root.

This means per-language disabling of common checks is most useful alongside
per-language `source_dir` configuration, so each check runs against the right
subtree:

```yaml
language:
  go:
    source_dir: backend
  typescript:
    source_dir: frontend

checks:
  typescript:
    disabled:
      - common:metrics   # metrics required in backend/ (go), skipped for frontend/ (ts)
```

With the above, `common:metrics` runs only in `backend/` (where it finds the Go
Prometheus dependency and passes) and is not evaluated against `frontend/`.

> Language-specific checks (e.g. `go:logging`, `typescript:tests`) are simply
> removed for that language when listed in its block.
