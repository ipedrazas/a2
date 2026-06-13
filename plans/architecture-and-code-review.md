# a2 — Architecture & Code Review, Improvements Backlog

_Review date: 2026-06-13_

## Implementation progress

- [x] **#1 Remove web client + server** (Section 5) — deleted `pkg/server/`,
  `cmd/server.go`, `ui/`, `docs/SERVER.md`, `docker-compose.yml`; dropped
  `gorilla/mux` + `google/uuid` via `go mod tidy`; simplified `Dockerfile`
  (removed ui-builder stage, EXPOSE 8080); removed `server:*` Taskfile tasks;
  scrubbed README / user-journeys server sections; retired the Firecracker
  roadmap item; CHANGELOG "Removed" entry added. `go build/vet/test ./...` green.
- [x] **#2 Unify parallel runner on bounded worker pool** (Section 6.3) —
  collapsed `runParallel` + `runParallelFailFast` into one bounded-pool
  `runParallel` (workers = `GOMAXPROCS`, overridable via
  `execution.concurrency`). Extracted `applyOptional` (Warn→Info) and `tally`
  helpers; the unbounded goroutine-per-check stampede is gone.
- [x] **#6 Default `--timeout` + `Errored`/`Skipped` status** (2.2.B/C) —
  `--timeout` now defaults to `3m` (0 still disables). Added `checker.Errored`
  (tool crash / panic / timeout — "a2 couldn't evaluate this") and
  `checker.Skipped` (cancelled by fail-fast); both are excluded from the score
  denominator and do not, on their own, fail the suite or abort the run.
  Surfaced in pretty/json/toon output and `maturity.Estimation`.
- [x] **#3 Add `--quick`/`--fast` run mode + `Speed` on `CheckMeta`** (6.1–6.2) —
  added `checker.Speed` (`SpeedFast` default / `SpeedSlow`); classified all 45
  build/test/coverage/race/vet/deps/deps_freshness/type/deadcode + clippy +
  duplication/sast checks as Slow. `--quick`/`--fast` applies an orthogonal
  `checks.FilterFast` after profile/target/skip selection. External checks opt
  into Slow via `speed: slow`. Mirrored in the GitHub Action (`quick` input) and
  a new `a2 check --quick` pre-commit hook. Slow checks already schedule first
  via their low `Order`, so longest-pole-first is satisfied.
- [ ] #4 Rename `a2 add` → `a2 init` + first-run nudge (7.1)
- [ ] #5 Surface description + `a2 explain` hint on failures (7.2)
- [ ] #7 `a2 list checks` shows descriptions by default (7.2)
- [ ] #8 Richer generated `.a2.yaml` template (7.1)
- [ ] #9 Move `os.Exit` out of `RunE`; return errors (3.2)
- [ ] #10 Extract repeated check boilerplate into `checkutil` helpers (3.3)
- [ ] #11 Data-driven security patterns (`//go:embed`) (3.3)
- [ ] #12 Score weighting for critical checks (4)

This document captures an architecture review, a code review, and a prioritized
list of improvements for `a2`. It also contains a concrete plan to **remove the
web client / server**, and proposals to address two recurring pieces of user
feedback:

1. _"Some checks take too long — give us a way to run a short set and a way to run everything."_
2. _"Users don't know how to create the config, or what each check actually does."_

---

## 1. Executive summary

`a2` is in good shape. The core abstraction (a `Checker` interface +
`CheckMeta` + a parallel runner + pluggable output formats) is clean, the check
library is large (80+ checks, 7 languages), and test coverage across the
`pkg/checks/*` tree is broad. The codebase is idiomatic Go with a sensible
package layout.

The highest-leverage moves right now are **subtractive and UX-focused**, not
new features:

- **Remove the web client + server.** It has not earned its keep, it adds a
  meaningful maintenance and dependency surface, and it pulls the project away
  from the terminal-first experience users actually want. (Section 5.)
- **Give users control over check duration.** There is currently no concept of
  a "fast" vs "full" run, and the default parallel path spawns unbounded
  subprocesses, which can make runs *slower*. (Section 6.)
- **Fix config & check discoverability.** All the machinery exists
  (`a2 add`, `a2 explain`, `a2 list`, `a2 doctor`) but new users never find it.
  (Section 7.)

The rest of this document details each area with file references and concrete
steps.

---

## 2. Architecture review

### 2.1 What's good

- **Clean check abstraction.** `pkg/checker/types.go` defines a minimal
  `Checker` interface (`ID`, `Name`, `Run(path)`) plus rich `CheckMeta`
  (`Critical`, `Optional`, `Order`, `Description`, `Suggestion`). Adding a check
  is mechanical and low-risk.
- **Good separation of concerns.** Detection (`pkg/language`), assembly
  (`pkg/checks/registry.go`), execution (`pkg/runner`), scoring
  (`pkg/maturity`), and presentation (`pkg/output`) are cleanly split. Output
  formats (pretty/json/toon) are swappable behind a common shape.
- **Monorepo handling is thoughtful.** `pathResolvingChecker` and
  `multiPathChecker` (`pkg/checks/registry.go`) let language checks target
  configured `source_dir`s, and `scopePerLanguage` allows per-language disabling.
- **Graceful tool degradation.** Most checks detect missing external tools via
  `checkutil.ToolAvailable` and downgrade to an `Info` result rather than
  failing, so a missing `jscpd`/`gitleaks`/`semgrep` doesn't break a run.
- **Extensibility via external checks** (`external:` in config) and
  profiles/targets that are shippable as editable YAML.

### 2.2 Architectural concerns

**A. Two divergent parallel execution paths** (`pkg/runner/runner.go`)
- `runParallel` (non-fail-fast, the **default** path) launches **one goroutine
  per check with no concurrency bound** (`runner.go:80`). For 80+ checks that
  each shell out to heavy subprocesses (`go build`, `go test -race`, `npm`,
  `cargo`), this stampedes the machine — oversubscribing CPU and memory, which
  often makes a run *slower*, not faster. This directly feeds the "checks take
  too long" complaint.
- `runParallelFailFast` (`runner.go:133`) *does* use a bounded worker pool
  (`GOMAXPROCS`). The two paths should be unified on a single bounded-pool
  implementation. See Section 6.

**B. No default timeout.** `--timeout` defaults to `0` = no timeout
(`cmd/root.go:120`). A single hung tool (a wedged `npm install`, a network-bound
check) hangs the whole run indefinitely. A sane per-check default (e.g. 2–5 min)
plus the existing override would protect users.

**C. Internal/infra errors are conflated with check failures.** In
`runCheckWithTimeout` a tool error or panic becomes a `Fail`
(`runner.go:308-317, 380-389`). That means "the check couldn't run" (tool
crashed, env broken) is reported identically to "your code is bad," and can
veto the suite. A dedicated `Errored`/`Skipped` status would let the report
distinguish *"a2 couldn't evaluate this"* from *"this failed."*

**D. No check cost/speed model.** `CheckMeta` has `Critical`, `Optional`,
`Order` but no notion of how expensive a check is. There is no way to express
"run the cheap static checks only." This is the structural gap behind the
short-vs-full feedback (Section 6).

**E. The server is an architectural outlier.** `pkg/server` introduces an
entirely different runtime model (HTTP, job queue, workspace cloning, GitHub
webhooks, untrusted-code execution) and the only `gorilla/mux` + `google/uuid`
dependencies. It roughly doubles the conceptual surface of the project for a
feature users don't use. (Section 5.)

---

## 3. Code review — findings

Ordered roughly by impact. None are critical bugs; most are
maintainability/consistency items.

### 3.1 Runner

- **Unbounded goroutines in the default path** — see 2.2.A. _Fix: single
  worker-pool implementation used by all parallel modes._ (`runner.go:64-131`)
- **Duplicated Warn→Info conversion logic** appears three times
  (`runner.go:102-105, 180-183, 255-258`). Extract a helper
  `applyOptional(reg, res)`.
- **Dead/while-internal wrappers.** `runSuite` and `runSuiteSequential`
  (`runner.go:40-48`) are unexported and appear used only by tests. Either
  promote to the public API intentionally or drop them.
- **`SuiteResult.Success()` vs `Aborted`.** Success is `Failed == 0` but the
  abort/exit path is computed separately in the output layer. Worth a single
  source of truth for "did the suite pass."

### 3.2 cmd / CLI

- **`os.Exit` inside `RunE`** (`cmd/root.go:164, 281`, and the validate
  commands) bypasses Cobra's error handling and any deferred cleanup, and makes
  the commands hard to unit-test. Prefer returning errors and letting
  `Execute()` own the exit code (it already does at `root.go:146-148`).
- **`a2 add` is a poorly-chosen verb for "create a config."** Users reach for
  `a2 init` or `a2 config init`. This is a real cause of the "don't know how to
  create the config" feedback — the command exists but is unguessable.
  _Fix: rename to `a2 init` (keep `add` as a hidden alias)._ See Section 7.
- **`--format` / `--output` aliasing** (`root.go:114-115, 152-154`) is handled
  ad hoc. Minor, but a single normalization point would be cleaner.

### 3.3 Checks

- **Recurring boilerplate across checks.** Four patterns repeat across ~15+
  checks: (a) tool-availability guard, (b) "find first config file that
  exists," (c) "count findings in output via regex," (d) priority fallback
  (tool → config → regex). `checkutil` already absorbs some of this; extract the
  remaining three into shared helpers (`checkutil.FirstExisting(...)`,
  `checkutil.CountMatches(...)`, a small `FallbackChain` helper) to shrink each
  check and standardize behavior.
- **Large hand-rolled pattern libraries** in `pkg/checks/security/`
  (`filesystem.go` 537 lines, `network.go` 452, `obfuscation.go` 395). These are
  hard to review and prone to false positives. _Suggestion: move the patterns to
  embedded data (`//go:embed` YAML/JSON), keep the Go code as a thin scanner,
  and table-test the patterns._ This also lets users extend/override patterns
  via config.
- **Per-check metadata is consistently populated** (every check has
  `Description` + `Suggestion`) — good. The problem is purely that this content
  is not *surfaced* (Section 7).

### 3.4 Repo hygiene

- **Committed build artifacts.** `ui/dist/assets/*.js` and `*.css` are
  generated bundles checked into git. Regardless of the removal decision, build
  output should not be versioned. (Resolved by Section 5 anyway.)
- **`docker-compose.yml` and `docs/SERVER.md`** exist solely for the server.

---

## 4. Maturity & scoring (`pkg/maturity`)

The scoring model (`maturity.go`) is simple and defensible: `score =
passed / scoredChecks * 100`, with `Info` excluded so optional/ungraded checks
don't drag the number down. Levels (PoC / Development / Mature / Production)
key off failure counts + score. No change needed, but two notes:

- Once an `Errored`/`Skipped` status exists (3.1.C), make sure it is excluded
  from the denominator the same way `Info` is.
- Consider weighting `Critical` checks more heavily than cosmetic ones so the
  score reflects risk, not just a raw pass ratio.

---

## 5. Plan: remove the web client & server

**Decision:** Remove the React UI and the `a2 server` command. Users prefer the
terminal; the server adds disproportionate maintenance, dependency, and
security surface (untrusted-code execution, GitHub cloning, job queue). This
also retires the unbuilt **Firecracker execution** roadmap item
(`plans/roadmap.md`), which only existed to make the web service safe.

### 5.1 Delete outright (server/UI-only)

- `pkg/server/` — all files: `server.go`, `handlers.go`, `job.go`, `queue.go`,
  `github.go`, `workspace.go`, `static.go`, `types.go`, `auth_test.go`,
  `bodylimit_test.go`.
- `cmd/server.go` — the entire file (only defines `serverCmd` and helpers).
- `ui/` — the entire directory (React source, `ui/dist` build artifacts, all
  build config).
- `docs/SERVER.md`.
- `docker-compose.yml` (server-only).

### 5.2 Edit

- `Dockerfile` — drop the `ui-builder` stage, the `COPY --from=ui-builder ...
  ui` line, `EXPOSE 8080`, and simplify the entrypoint to the CLI
  (`ENTRYPOINT ["a2"]` / `CMD ["check"]`).
- `Taskfile.yaml` — remove the `server:*` tasks (`server:ui:build`,
  `server:build`, `server:run`, `server:dev`, `server:dev:ui`).
- `go.mod` — after deletion, `go mod tidy` should drop `github.com/gorilla/mux`
  and `github.com/google/uuid` (server-only; verify nothing else imports them).
- `README.md` / `docs/user-journeys.md` — remove the "Server Mode" sections and
  any web-UI references.
- `plans/roadmap.md` — remove the Firecracker section (or note it as retired).

### 5.3 Suggested sequencing

1. Delete `pkg/server/`, `cmd/server.go`, `ui/`, `docs/SERVER.md`,
   `docker-compose.yml`.
2. `go build ./... && go test ./...` — confirm nothing references the removed
   packages.
3. `go mod tidy`; confirm `gorilla/mux` and `uuid` are gone.
4. Edit `Dockerfile`, `Taskfile.yaml`, docs.
5. Update `CHANGELOG.md` with a clear "Removed: web UI and server mode" entry
   and a one-line migration note (CLI is now the only interface).

This is a low-risk change: the server has no inbound dependencies from the CLI
core — `cmd/root.go` and `main.go` never import `pkg/server`.

---

## 6. Feature: short vs full check runs (+ make runs faster)

Two complementary changes: a **declarative speed/tier model** so users can pick
a subset, and a **runner fix** so any run uses the machine sensibly.

### 6.1 Add a cost classification to `CheckMeta`

Add a field to `CheckMeta` (`pkg/checker/types.go`):

```go
// Speed indicates the relative cost of a check.
// Fast checks are static/IO-light; Slow checks spawn builds/tests/network.
type Speed int
const ( SpeedFast Speed = iota; SpeedSlow )
```

Classify each registration: static checks (format, lint, file existence,
secrets-regex, config validation, naming, editorconfig) → `Fast`; checks that
build/test/cover/race/scan-the-network → `Slow`. This is a one-time pass over
the `register.go` files.

> Optional refinement: instead of a binary, record an empirical baseline
> duration and let `a2` learn real timings from `Result.Duration` (already
> measured) cached under `~/.config/a2`. Start with the static binary
> classification — it's simpler and good enough.

### 6.2 Surface it as run modes

- `a2 check --quick` (alias `--fast`): run only `SpeedFast` checks. Target:
  sub-second-to-few-seconds feedback for the inner loop / pre-commit.
- `a2 check` (default): unchanged — run everything.
- Make this composable with profiles/targets (quick is an orthogonal filter,
  not a profile).
- Mirror it in the pre-commit hook and GitHub Action inputs (`scripts/`,
  `scripts/action.yml`) so the cheap set runs on commit and the full set runs in
  CI.

### 6.3 Fix the runner so every run is faster

- **Unify on the bounded worker pool** (currently only fail-fast uses it). Cap
  concurrency at `GOMAXPROCS` (or a configurable `execution.concurrency`).
  Removing the unbounded goroutine stampede in `runParallel` will, on its own,
  speed up large runs and reduce memory spikes. (`runner.go:64-131`)
- **Ship a sane default `--timeout`** (e.g. 3m) so one wedged tool can't hang
  the suite (2.2.B).
- Consider scheduling `SpeedSlow` checks first (longest-pole-first) so total
  wall-clock is bounded by the slowest check rather than by scheduling luck.

---

## 7. UX: config creation & "what does this check do?"

Everything needed already exists — it's a discoverability problem. Concrete
fixes, cheapest first:

### 7.1 Make config creation guessable

- **Rename `a2 add` → `a2 init`** (keep `add` as a hidden alias for
  compatibility). `init` is the universal verb for scaffolding.
  (`cmd/add.go:30-44`)
- **First-run nudge.** When `a2 check` runs with **no `.a2.yaml` present**,
  print a one-line banner: _"No .a2.yaml found — running with defaults. Run
  `a2 init` to customize thresholds and checks."_ (Hook in `cmd/root.go`
  `runCheck` after `config.Load`.)
- **Nudge on "no language detected"** (`root.go:225`): point at `a2 init`/
  `--lang` explicitly (already partly there — make it mention `a2 init`).
- **Generate a richer template.** The generator currently emits mostly
  commented-out examples (`pkg/config/generator.go`). Emit a small but *active*
  annotated config so users see a working baseline, with inline comments
  explaining each block.

### 7.2 Make checks self-explaining in the normal workflow

- **Show the suggestion + an explain hint on failures.** In pretty output, for
  each failed check append: _"→ why: <Description> · `a2 explain <id>` for
  details."_ The data is already on `CheckMeta`; it's just not rendered in the
  default view. (`pkg/output/pretty.go`)
- **`a2 list checks` should show descriptions by default** (or add a top-level
  `a2 checks` that defaults to `--explain`). Requiring `--explain` to see what a
  check does is backwards. (`cmd/list.go`)
- **Cross-link the docs.** `docs/CHECKS.md` and `docs/checks/<lang>.md` are
  thorough; have `a2 explain <id>` print the doc anchor/URL so users can go
  deeper.

### 7.3 Onboarding doc

- Add a short "Getting Started" section to `README.md` that is literally:
  `a2 init` → `a2 check` → `a2 doctor` (for missing tools) → `a2 explain <id>`.
  The detailed `docs/user-journeys.md` assumes prior knowledge; new users need
  the 4-command happy path up front.

---

## 8. Prioritized backlog

| # | Item | Effort | Impact | Section |
|---|------|--------|--------|---------|
| 1 | Remove web client + server (+ `go mod tidy`, Dockerfile, Taskfile, docs) | M | High | 5 |
| 2 | Unify parallel runner on bounded worker pool | S | High | 6.3 |
| 3 | Add `--quick`/`--fast` run mode + `Speed` on `CheckMeta` | M | High | 6.1–6.2 |
| 4 | Rename `a2 add` → `a2 init` + first-run nudge | S | High | 7.1 |
| 5 | Surface description + `a2 explain` hint on failures | S | High | 7.2 |
| 6 | Default `--timeout` + `Errored`/`Skipped` status | S | Med | 2.2.B/C |
| 7 | `a2 list checks` shows descriptions by default | S | Med | 7.2 |
| 8 | Richer generated `.a2.yaml` template | S | Med | 7.1 |
| 9 | Move `os.Exit` out of `RunE`; return errors | S | Med | 3.2 |
| 10 | Extract repeated check boilerplate into `checkutil` helpers | M | Med | 3.3 |
| 11 | Data-driven security patterns (`//go:embed`) | M | Med | 3.3 |
| 12 | Score weighting for critical checks | S | Low | 4 |

_Effort: S ≈ <½ day, M ≈ 1–2 days. Impact is on user-perceived value._

### Suggested order of execution

1. **#1 (remove server)** — clears the deck, shrinks the dependency/maintenance
   surface, and aligns the project with the terminal-first reality.
2. **#2 + #6 (runner)** — makes *every* run faster and safer; prerequisite mood
   for #3.
3. **#3 (quick mode)** — directly answers the "checks take too long" feedback.
4. **#4 + #5 + #7 + #8 (UX)** — directly answers the "don't know how to
   configure / what checks do" feedback. These are small and high-impact.
5. **#9–#12** — internal quality, do as you touch the surrounding code.
