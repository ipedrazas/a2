# Code Quality Metrics Analysis

## Current Coverage

| Category | Checks | Status |
|---|---|---|
| Buildability | `*:build` per language | Covered |
| Testing | `*:tests`, `common:integration`, `common:e2e` | Covered |
| Coverage | `*:coverage` with configurable thresholds | Covered |
| Formatting | `*:format`, `common:editorconfig` | Covered |
| Linting | `*:lint` | Covered |
| Type safety | `*:type` (Python, Node, TS) | Covered |
| Complexity | `go:cyclomatic`, `python:complexity` | Partial |
| Security | `common:secrets`, `common:sast`, `security:*` | Covered |
| Dependencies | `*:deps` | Covered |
| Operational readiness | `common:health`, `common:metrics`, `common:tracing`, `common:shutdown`, `common:errors` | Covered |

## Missing Metrics

### Priority 1 — High impact, implement first

- [x] **Code duplication / DRY analysis**
  - Universal, language-agnostic, high signal
  - Tools: `jscpd`, `CPD` (PMD), `flay` (Ruby), SonarQube duplication engine
  - High duplication is one of the strongest predictors of maintenance burden
  - Suggested check ID: `common:duplication`

- [x] **Dependency freshness / staleness**
  - `*:deps` validates dependency management exists but doesn't check how outdated dependencies are
  - Tools: `npm outdated`, `pip list --outdated`, `go list -m -u all`, `cargo outdated`
  - Stale deps accumulate security risk and upgrade pain
  - Suggested check IDs: `go:deps_freshness`, `python:deps_freshness`, `node:deps_freshness`, etc.

- [ ] **Mutation testing** (test quality)
  - Coverage can be 90% with weak assertions; mutation testing catches that
  - Tools: `go-mutesting` (Go), `mutmut` (Python), `Stryker` (JS/TS), `PIT` (Java)
  - Suggested check IDs: `*:mutation`

- [x] **Dead code / unused exports**
  - Reduces cognitive load and attack surface
  - Tools: `deadcode` (Go), `vulture` (Python), `ts-prune` (TS), `unused`
  - Suggested check IDs: `*:deadcode`

### Priority 2 — Valuable additions

- [ ] **Cognitive complexity** (distinct from cyclomatic)
  - Cyclomatic counts branches; cognitive penalizes nesting depth and flow-breaking constructs
  - Better captures readability than cyclomatic complexity alone
  - SonarQube defines the metric; some linters support it (e.g., `gocognit`, `cognitive_complexity` flake8 plugin)
  - Currently only have cyclomatic for Go and Python
  - Suggested check IDs: `*:cognitive_complexity`

- [ ] **Documentation coverage**
  - Public API functions/types having docstrings
  - Tools: `interrogate` (Python), `golint` doc checks, `typedoc` coverage, `javadoc` warnings
  - Suggested check IDs: `*:doc_coverage`

- [ ] **Binary / artifact size**
  - Especially relevant for CLI, desktop, and mobile profiles
  - Unexpected size growth signals accidentally bundled dependencies
  - Suggested check ID: `common:artifact_size`

### Priority 3 — Nice to have

- [ ] **Build reproducibility**
  - Can two builds from the same commit produce identical artifacts?
  - Checks: lockfile presence (partially covered), pinned dependencies, deterministic build flags
  - Suggested check ID: `common:reproducibility`

- [ ] **API compatibility / breaking changes**
  - Tools: `go-apidiff`, `api-extractor` (TS), `bump2version`
  - Most relevant for library profile
  - Suggested check IDs: `*:api_compat`

- [ ] **Commit / branch hygiene**
  - Conventional commits validation, branch naming, PR size limits
  - Large PRs correlate with defect introduction
  - Suggested check ID: `common:commit_hygiene`

- [ ] **Test-to-code ratio**
  - Simple heuristic: ratio of test lines to source lines
  - Suggested check ID: `common:test_ratio`

- [ ] **Flaky test detection**
  - Repeated test runs to surface non-determinism
  - Suggested check IDs: `*:flaky_tests`

## Implementation Notes

- All new checks should implement the `Checker` interface in `pkg/checker/types.go`
- Language-specific checks go in `pkg/checks/{language}/`
- Common checks go in `pkg/checks/common/`
- Update profiles in `pkg/profiles/profiles.go` (e.g., `common:artifact_size` irrelevant for library profile)
- Update targets in `pkg/targets/targets.go` (e.g., mutation testing only expected at production target)
