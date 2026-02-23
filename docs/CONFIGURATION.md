# A2 Configuration Reference

A2 is configured via a `.a2.yaml` file in your project root. You can generate one with `a2 add` (use `-i` for interactive mode).

## Table of Contents

- [Generating Configuration](#generating-configuration)
- [Configuration Examples](#configuration-examples)
- [External Checks](#external-checks)
- [Wildcard Patterns](#wildcard-patterns)

---

## Generating Configuration

### Interactive Mode

Run with `-i` flag for guided prompts:

```bash
a2 add -i
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
a2 add --profile cli --target poc

# With language and coverage
a2 add --lang go,python --coverage 90

# Custom required files
a2 add --files README.md,LICENSE,CHANGELOG.md

# Overwrite existing config
a2 add --profile api --force

# Custom output path
a2 add --output custom-config.yaml
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

```yaml
# Explicit language selection (overrides auto-detect)
language:
  explicit:
    - go
    - python
    - typescript
  go:
    coverage_threshold: 80
  python:
    coverage_threshold: 70
    linter: ruff
  typescript:
    coverage_threshold: 75

files:
  required:
    - README.md
```

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
