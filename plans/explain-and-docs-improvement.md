# Plan: Make check failures actionable (explain + docs + security allowlists)

## Motivation

Customer feedback: when `a2 check -v` reports a failing check (e.g. `security:obfuscation`),
there is no clear path to **what** the problem is, **where** it is, or **how** to fix it.
`a2 explain <id>` only prints static metadata, and the docs don't cover most checks.

The customer's four expectations, which this plan must satisfy:
1. What the problem is.
2. Where it is.
3. How to fix it (a command, or an AI prompt if no command exists).
4. A way of finding out — via docs or `a2 explain`.

## Root-cause findings

1. **The "where" is computed then hidden at `-v`.** Scan checks already produce `file:line`
   findings in `result.Message` (`obfuscation.go:374` `formatFindings`), but `pretty.go:219`
   only prints `Message` when `verbosity > 1` (`-vv`). The customer ran `-v` (verbosity = 1)
   and saw nothing useful. `a2 run security:obfuscation` shows it, but nothing points there.
2. **`explain` is metadata-only.** `Meta.Command` is set on only ~11 of ~95 checks, and scan
   checks have no command and no pointer to `a2 run <id>` where the findings live.
3. **Docs gaps in `docs/CHECKS.md`:** security documents only `filesystem` (missing
   `obfuscation`, `shell_injection`, `network`); devops documents only `k8s` (missing
   `ansible`, `helm`, `pulumi`, `terraform`); there is no `docs/checks/swift.md` (8 checks).
4. **No suppression for Critical security scans.** Only `security:filesystem` has an
   allowlist. `obfuscation` / `shell_injection` / `network` are `Critical` and hard-block;
   a false positive forces users to disable the whole check.

## Decisions (confirmed with user)

- **Rich `explain` via new `CheckMeta` metadata** — extend `CheckMeta`, populate for all
  checks, make `explain` self-contained (the single source of truth).
- **Add a false-positive allowlist** to `obfuscation`, `shell_injection`, and `network`
  (parity with `filesystem`) in this effort.

---

## Phase 1 — Surface *where* at `-v` (highest impact, smallest change)

- `pkg/output/pretty.go:219`: print `r.Message` (and `r.Reason`) for **Fail/Warn** checks at
  `VerbosityFailures` (`-v`), not only `-vv`. Keep `-vv` showing messages for all checks.
- Confirm `toon.go` / `json.go` verbosity gating matches and already carry `Message`/`Reason`.
- Tests: `pkg/output/pretty_test.go` — assert a failing check's `Message` appears at `-v`.

## Phase 2 — Extend `CheckMeta` and make `explain` self-contained

### 2a. New `CheckMeta` fields (`pkg/checker/types.go`)
```go
Detail    string // Long-form: what it detects, what patterns/conditions trigger it, why it matters
FixPrompt string // Copy-paste remediation prompt for an AI, for checks with no fix command
```
- Keep existing `Description` (one-line), `Suggestion` (short fix), `Command` (template).
- `Detail` answers "what it does" + "why important"; `FixPrompt` answers "how to fix" when
  there is no command.

### 2b. `explain` rendering (`cmd/explain.go`)
Render, in order: ID, Name, Description, **Detail**, Languages, Critical, Suggestion,
**Command** (or, for scan checks, `→ To see exact files & lines:  a2 run <id>`),
**Fix prompt** (`FixPrompt`), Speed, Docs link. Each new field prints only when non-empty.

### 2c. Populate metadata for all checks
- Fill `Command` for every check that shells out (~95 register entries / check files).
- Fill `Detail` for every check (what triggers it + why it matters).
- Fill `FixPrompt` for checks with no fix command (security scans, secrets, obfuscation).
- Files: all `pkg/checks/*/register.go` and relevant check files.
- Mitigation for volume: do it language-package by language-package; the docs-coverage test
  (Phase 4) doubles as a checklist.

### 2d. Tests (`cmd/explain_test.go`)
- A command-based check (`go:vet`) renders `Command`.
- A scan check (`security:obfuscation`) renders `Detail`, the `a2 run <id>` pointer, and
  `FixPrompt`.

## Phase 3 — Allowlists for the remaining security scans

- Generalise the `filesystem` allowlist mechanism (`security:filesystem` already supports
  `file:line`, `file:match` w/ wildcards, `file`).
- Add `Allow []string` config under `security.obfuscation`, `security.shell_injection`,
  `security.network` (mirror `security.filesystem.allow` in `pkg/config`).
- Wire allowlists into the three checks (`obfuscation.go`, `shell_injection.go`,
  `network.go`) so matching findings are suppressed before Pass/Fail is decided.
- Tests: each check suppresses an allowlisted finding; non-matching findings still fail.
- Files: `pkg/config/*`, `pkg/checks/security/{obfuscation,shell_injection,network}.go`
  + `register.go` (pass config through like `FileSystemCheck{Allowlist: ...}`).

## Phase 4 — Document every check (what / what it does / why it matters)

- `docs/CHECKS.md`: add `security:obfuscation`, `security:shell_injection`, `security:network`
  (with trigger patterns + allowlist usage) and devops `ansible`, `helm`, `pulumi`,
  `terraform`, matching the existing `### id` format.
- Create `docs/checks/swift.md` (8 checks) in the `go.md` table+sections format.
- Audit `docs/checks/*.md` so each check answers what / what-it-does / why-it-matters.
- **Add `docs_coverage_test.go`**: iterate `GetAllCheckRegistrations` and fail if any check ID
  lacks a doc heading (in `CHECKS.md` or the per-language file). Prevents future drift.

## Phase 5 — Verify & wrap up

- Reproduce the customer session: `a2 check -v` → `a2 explain security:obfuscation` →
  `a2 run security:obfuscation`, and confirm all four expectations are answerable.
- `go test ./...` green; run `a2 check` on this repo.
- Update `CHANGELOG.md`.

## Effort / risk

| Phase | Size | Risk | Notes |
|-------|------|------|-------|
| 1 | S | Low | One file; biggest UX win |
| 2 | L | Low | Wide but mechanical; touches every register.go |
| 3 | M | Med | New config surface; mirror existing filesystem allowlist |
| 4 | L | Low | Docs volume; coverage test guards drift |
| 5 | S | Low | Verification + changelog |

## Suggested order

1 → 2a/2b (shape) → 3 (so security checks have allowlist before docs describe it) →
2c (bulk metadata) → 4 (docs) → 5.
