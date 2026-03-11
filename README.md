# a2 — Validate and verify every repository change

Deterministic repository verification for teams, CI, and coding agents.

`a2` checks whether a repository meets a defined baseline for build health, tests, coverage, security, formatting, and release readiness. It runs 80+ deterministic checks across multiple languages and returns a clear pass/fail result.

```bash
go install github.com/ipedrazas/a2@latest

a2 check
a2 check --profile=api
a2 check -f toon
```

## Why use `a2`?

Teams use prompts, `.cursorrules`, and agent skills to describe intent. But prompts are not deterministic.

`a2` turns repository intent into repeatable checks.

- Verify repo readiness in one command
- Standardize quality gates across projects
- Feed machine-readable failures to coding agents for automatic fixes

Prompts guide implementation. `a2` verifies the result.

## What `a2` does

`a2` runs deterministic checks for things like:

- buildability
- tests
- coverage
- formatting and linting
- type checking
- vulnerabilities
- repo and delivery hygiene

It auto-detects supported languages or lets you choose them explicitly.

Example of running `a2` against the a2 repo:

![screenshot](docs/assets/a2-screenshot.png)

## Example

Check a repository:

```bash
a2 check --profile=api
```

Export agent-friendly output:

```bash
a2 check -f toon > a2.out
```

Use that output in Cursor, Claude Code, or another coding agent, apply fixes, and rerun until the repo passes.

## When to use `a2`

Use `a2` when you want to:

- verify a pull request did not break the repo baseline
- give coding agents an objective list of issues to fix
- enforce different standards for PoC and production repos
- standardize checks across polyglot repositories

## How `a2` differs

| Tool | Role |
|------|------|
| Prompts / `.cursorrules` / skills | Describe how code should be written |
| Linters / tests | Check one part of a repo |
| CI scripts | Run repo-specific automation |
| `a2` | Verify the repository baseline end-to-end |

## Features

- **8 supported languages**: Go, Python, Node.js, TypeScript, Java, Rust, Swift
- **80+ built-in checks**
- **Profiles**: CLI, API, Library, Desktop
- **Targets**: PoC or Production
- **Maturity scoring**
- **JSON and TOON output**
- **Custom checks via external binaries**
- **GitHub Action, Docker, and pre-commit support**

## Common commands

```bash
a2 check
a2 check -v
a2 check --lang go,python
a2 check --profile=library
a2 check --target=poc
a2 check --skip="*:tests"
a2 list checks
a2 run rust:coverage agent
a2 explain go:race
a2 doctor
```

## Configuration

Create a `.a2.yaml` file to set thresholds, disable checks, and add external checks:

```yaml
language:
  go:
    coverage_threshold: 80

checks:
  disabled:
    - go:deps

external:
  - id: lint
    command: golangci-lint
    args: ["run", "./..."]
    severity: warn
```

More details:
- [Checks reference](docs/CHECKS.md)
- [Configuration](docs/CONFIGURATION.md)
- [Server mode](docs/SERVER.md)

## CI/CD

### GitHub Action

```yaml
name: A2 Analysis
on: [push, pull_request]

jobs:
  a2:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run A2
        uses: ipedrazas/a2@v1
        with:
          format: "pretty"
          profile: "api"
```

### Docker

```bash
docker run -v $(pwd):/workspace ipedrazas/a2 check
```

## Exit codes

- `0` — all critical checks passed
- `1` — one or more critical checks failed

## License

MIT — see [LICENSE](LICENSE)
