# PLAN — n8n Audit Profile in a2 (Phase 0 output)

Recon result: **a2's abstractions fit the n8n domain — no escape hatch needed.** The check
interface is `Run(path string) (checker.Result, error)`, not hardwired to language toolchains;
the "language" mechanism is an open string enum that n8n can join as a domain; profiles,
config, list/explain, and output all work off metadata that n8n checks can populate like any
other. Frictions found are listed in §3 with resolutions — none require refactoring a2's core.

**CHK-04 decision (§3 of HANDOVER): asked Ivan 2026-07-21 — confirmed option (a), full check
in a2**, including the write-capable tool classification table. Nothing is held back privately.

---

## 1. Architecture as found (what n8n code plugs into)

| Concern | Where it lives | Concrete symbols |
|---|---|---|
| Check interface | `pkg/checker/types.go` | `checker.Checker` (`ID()`, `Name()`, `Run(path) (Result, error)`), `checker.Result`, `checker.CheckMeta`, `checker.CheckRegistration` |
| Result statuses | `pkg/checker/types.go:25` | `Pass`, `Warn`, `Fail`, `Info`, `Errored`, `Skipped` |
| Result helpers | `pkg/checkutil/util.go` | `NewResultBuilder`, `rb.Pass/Warn/Fail/Info`, `*WithOutput` variants, `RunCommand` |
| Error degradation | `pkg/runner/runner.go:286` | returning a non-nil `error` from `Run` → runner builds `Errored` result (`erroredResult`); panics recovered the same way. This is the "never false-pass" path. |
| Exit code | `pkg/runner/runner.go:362` | `SuiteResult.Success()` = `Failed == 0`. Any `Fail` status → exit 1. `Meta.Critical` does **not** gate exit; it double-weights the maturity score (`WeightedTally`, `criticalWeight = 2.0`). |
| Registration | per-domain `Register(cfg *config.Config) []checker.CheckRegistration` (pattern: `pkg/checks/security/register.go`) | wired in `pkg/checks/registry.go`: `getChecksForLanguage` (switch on language) and `GetAllCheckRegistrations` (drives `a2 list checks` + `a2 explain`) |
| Language detection | `pkg/language/detect.go` | `LanguageIndicators` (filename-based), `Detect`, `DetectWithSourceDirs`, `orderedLanguages` (duplicated in both Detect funcs) |
| Profiles | `pkg/profiles/profiles.go` | `BuiltInProfiles` map; a profile is a **disable list**, applied in `cmd/root.go:runCheck` by appending to `cfg.Checks.Disabled` |
| Config | `pkg/config/config.go` | `Config` struct (top-level sections e.g. `Security`), `knownLanguages` map (config.go:268) gates per-language `checks:` blocks, `MatchDisabled` supports `domain:*` wildcards |
| Output | `pkg/output/json.go`, `toon.go`, `pretty.go` | flat `Result` fields only: `Message`, `Reason`, `RawOutput` strings — no structured detail |
| explain/list | `cmd/explain.go`, `cmd/list.go` | driven entirely by `CheckMeta` (`Description`, `Detail`, `FixPrompt`, `Command`) — n8n checks get these for free once registered |

## 2. Integration approach: n8n is a "language" (domain)

- Add `checker.LangN8n Language = "n8n"` in `pkg/checker/types.go`.
- Checks carry `Languages: []Language{LangN8n}` and register via a new
  `pkg/checks/n8n.Register(cfg)`, added to the `getChecksForLanguage` switch and to
  `GetAllCheckRegistrations` in `pkg/checks/registry.go`.
- Detection: n8n export files are identified by **content shape** (JSON object with top-level
  `nodes` array + `connections` object), not filename — a2's `LanguageIndicators` is
  filename-based, so add a content-based detector `detectN8n(path)` in `pkg/language/detect.go`
  called from `Detect` / `DetectWithSourceDirs` (bounded: `*.json` in root + one level of
  subdirs, read-limit per file, skip `node_modules`/`.git`). This satisfies HANDOVER §2.3
  ("auto-detection adds the n8n domain the same way languages are detected") and clears the
  "no supported language detected" error in `cmd/root.go:251` for pure workflow directories.
- `--profile=n8n`: a2 profiles are subtractive, so "run exactly this suite" = a built-in
  profile in `BuiltInProfiles` disabling every other domain:
  `go:*, python:*, node:*, java:*, rust:*, typescript:*, swift:*, common:*, devops:*, security:*`.
  Without the profile, auto-detection runs n8n checks *alongside* other detected languages —
  the desired default in mixed repos.
- Hardcoded language lists to audit when adding `LangN8n`: `checker.AllLanguages()`,
  `orderedLanguages` (×2 in detect.go), `config.knownLanguages`, the error message at
  `cmd/root.go:251`, `--lang` flag help text, `a2 doctor`/`a2 add` if they enumerate languages.

## 3. Registry assumptions that don't fit a2 — flagged, with resolutions

1. **Output contract mismatch.** The registry wants structured
   `{check_id, status, affected_workflows[], detail}`; a2's `Result` is flat strings.
   *Resolution:* affected workflows rendered as deterministic, line-oriented text —
   one `<file> · <workflow name> · <node name>` line per finding — in `Reason` (summary count
   in `Message`, full listing in `RawOutput`). The private-repo mapper (Phase 6) parses these
   lines; the format is documented in `docs/N8N.md` and covered by golden tests so it's a
   stable contract. No a2 output-schema change.
2. **`workflow_bundle_hash` has no home in a2 output.** *Resolution:* `pkg/n8n` exposes
   `Bundle.Hash` (SHA-256 over sorted file contents) as public API; a2 does not emit it. The
   Phase 6 mapper recomputes it with the same documented algorithm over the same files.
   Keeps a2 generic.
3. **gitleaks reuse (CHK-01) rejected.** `common:secrets` (`pkg/checks/common/secrets.go`)
   shells to gitleaks over raw files and puts raw tool output into `RawOutput` — that scans
   the wrong surface (files, not extracted node parameters) and would leak matched values,
   violating the redaction constraint. *Resolution:* implement pattern + entropy detection in
   `pkg/n8n` over extracted parameter values / Code-node source, per the registry, with
   redaction at detection time. No new dependency.
4. **jscpd reuse (CHK-13) rejected.** `common:duplication` shells to jscpd over source text —
   wrong input model for node graphs, and adds a tool dependency. *Resolution:* n-gram over
   node-type sequences as HANDOVER already specifies.
5. **`not_applicable` has no exact a2 status.** *Resolution:* map to `Info`
   (excluded from score, visible in output) with a "not applicable: …" message.
   `Skipped` is reserved by a2 for fail-fast cancellation.
6. **Per-check parse cost.** Each a2 check runs independently (and in parallel), so 12+ n8n
   checks would each re-parse the bundle. *Resolution:* `pkg/n8n` provides a small
   mutex-guarded memoized loader (`LoadDirCached(path)`) keyed by absolute path. Determinism
   unaffected.

## 4. Check mapping (registry → a2)

All checks `SpeedFast` (static parsing, no builds/network). "Critical" = `Meta.Critical: true`
(score weighting + reporting); exit-code impact comes from returning `Fail` status.

| a2 check ID | Registry | Statuses emitted | Critical | Notes |
|---|---|---|---|---|
| `n8n:secrets` | CHK-01 | Fail (confirmed pattern), Warn (high-entropy candidate) | ✔ | own scanner, redacted (§3.3) |
| `n8n:inline-auth` | CHK-02 | Warn | | literal auth headers/tokens vs `credentials` reference |
| `n8n:webhook-auth` | CHK-03 | Fail (side-effectful reachable), Warn (read-only path) | ✔ | graph traversal |
| `n8n:agent-injection-path` | CHK-04 | Fail (write-capable tools), Warn (tool-less AI) | ✔ | **full check in OSS (confirmed)**; readable trigger→…→agent path in `Reason` |
| `n8n:agent-tools` | CHK-05 | Warn (write-capable attached), Pass with inventory | | inventory in `RawOutput` |
| `n8n:pii-to-ai` | CHK-06 | Warn | | field-name heuristics + `{{ $json }}` passthrough |
| `n8n:error-workflow` | CHK-07 | Fail (<50% coverage), Warn (<100%) | ✔ | coverage ratio in `Message` |
| `n8n:silent-failures` | CHK-08 | Warn | | `continueOnFail`/`onError` without error branch |
| `n8n:retry-config` | CHK-09 | Warn | | distribution in `RawOutput` |
| `n8n:timeouts` | CHK-10 | Warn | | scoped to webhook-triggered workflows |
| `n8n:code-node-bulk` | CHK-11 | Warn | | line limit configurable (default 80); embedded fetch/HTTP |
| `n8n:naming` | CHK-12 | Warn | | regex configurable; defaults-detection fallback |
| `n8n:duplication` | CHK-13 | Warn | | n-gram clusters (§3.4) |
| `n8n:pinned-data` | CHK-14 | Warn | | `pinData` on active workflows |
| `n8n:zombies` | CHK-15 | — | | **deferred** (needs n8n API), per HANDOVER |
| `n8n:in-git` | CHK-16 | Warn | | no existing a2 check covers "workflow files tracked in git"; implement via `checkutil.RunCommand` + `git ls-files` |

Every status downgrade on uncertainty: unparseable workflow relevant to a check → that check
returns an `error` → runner `Errored` (never silent pass); partially parseable bundles degrade
only the affected checks.

## 5. Config surface (`.a2.yaml`)

New top-level section (like `security:`), struct `N8nConfig` added to `config.Config`:

```yaml
n8n:
  code_node_line_limit: 80        # n8n:code-node-bulk
  naming_regex: ""                # n8n:naming; empty = defaults-detection mode
  error_workflow_fail_ratio: 0.5  # n8n:error-workflow fail threshold
  duplication_min_sequence: 5     # n8n:duplication
  untrusted_triggers_extra: []    # extend classification tables
  side_effect_nodes_extra: []
  write_tools_extra: []           # entries "nodeType:operation"
checks:
  n8n:                            # per-domain disable list (requires knownLanguages entry)
    disabled: []
```

## 6. Phases (restated against real symbols)

### Phase 1 — parsing layer: new package `pkg/n8n`
- `types.go` — `Workflow` (Name, ID, Active, Settings.ErrorWorkflow, PinData), `Node` (ID,
  Name, Type, TypeVersion, Parameters as `map[string]any` + raw `json.RawMessage` retained,
  Credentials, ContinueOnFail/OnError, RetryOnFail, MaxTries, WaitBetweenTries),
  `Connections`. Defensive unmarshal: unknown node types retained; structural surprises
  recorded as `ParseError`s on the bundle, not dropped.
- `loader.go` — `IsWorkflowFile([]byte) bool` (nodes+connections shape), `LoadDir(path)
  (*Bundle, error)`, `LoadDirCached`, `Bundle{Workflows, ParseErrors, Hash}` with SHA-256 over
  sorted file contents.
- `graph.go` — `BuildGraph(*Workflow) *Graph`, `Graph.Reachable(nodeID) []NodeRef`
  (transitive, main + error outputs), `Graph.PathTo(from, to) []NodeRef` for readable paths.
- `classify.go` — data tables only: `UntrustedTriggers`, `SideEffectNodes`, `AIAgentNodes`
  (+ how tools attach in the export), `WriteCapableTools` (nodeType+operation, e.g. Gmail
  `send` vs `getAll`). Merged with `_extra` config lists at check construction.
- `testdata/` fixtures per HANDOVER §6 (`clean/`, `vulnerable/`, `mitigated/`, `malformed/`,
  `duplication/`) + `testdata/README.md` recording the n8n version the export shape targets
  (verify current export format against n8n docs/source before authoring).
- Acceptance: parse + reachability unit tests green; malformed fixtures produce `ParseError`s,
  never panic or silent pass.

### Phase 2 — simple checks: new package `pkg/checks/n8n`
- `register.go` (pattern: `pkg/checks/security/register.go`) + one file per check, each a
  small struct holding its thresholds from `config.N8nConfig`.
- Checks: `error_workflow.go`, `silent_failures.go`, `retry.go`, `timeouts.go`,
  `code_bulk.go`, `pinned_data.go`, `naming.go`, `inline_auth.go`, `in_git.go`.
- Core wiring in the same commit: `checker.LangN8n`, `pkg/checks/registry.go` (both call
  sites), `config.N8nConfig` + `knownLanguages`, `pkg/language` detection, §2 hardcoded-list
  audit.
- Acceptance: table-driven tests over fixtures; `a2 list checks` shows all;
  `a2 explain n8n:error-workflow` returns `Detail`/`FixPrompt` text derived from the registry.

### Phase 3 — secrets & duplication
- `pkg/n8n/secrets.go` scanner (patterns: `sk-`, `ghp_`, `xoxb-`, `AKIA`, JWT shape,
  auth-adjacent literal headers; entropy for candidates) returning findings as
  workflow/node/field-path **only** — matched value never leaves the scanner. Check wrapper
  `pkg/checks/n8n/secrets.go`.
- `pkg/n8n/duplication.go` — order-preserving node-type n-grams across workflows (min length
  configurable, default 5), parameter similarity as tiebreak, cluster output.
- Acceptance: vulnerable fixture fails `n8n:secrets` AND a test greps the fake secret string
  out of the full `-f json` serialized output (absent); duplication fixture reports one cluster.

### Phase 4 — traversal checks
- `webhook_auth.go`, `agent_injection.go`, `agent_tools.go`, `pii_to_ai.go` in
  `pkg/checks/n8n`, built on `Graph.Reachable`/`PathTo` + classification tables.
- Path rendering (`Reason`): `Webhook "lead-intake" → Set → AI Agent "triage" [tools:
  gmail.send (write), sheets.append (write)]` — treat this string's quality as a feature.
- Acceptance: vulnerable fixture → Fail on both traversal checks with correct named paths;
  mitigated fixture passes.

### Phase 5 — profile, docs, polish
- `BuiltInProfiles["n8n"]` in `pkg/profiles/profiles.go` (§2 disable list); profile validation
  tests updated.
- Docs: n8n sections in `docs/CHECKS.md` + `docs/checks/n8n.md` (house style: per-language md
  files exist), new `docs/N8N.md` (usage + the stable finding-line format from §3.1), README
  feature list, CHANGELOG entry, version bump per repo convention.
- Verify `-f json`/`-f toon` structure unchanged; exit 1 on Fail; `a2 check` self-check green;
  `task ci` + pre-commit before every commit.
- Acceptance: `a2 check --profile=n8n ./pkg/n8n/testdata/vulnerable` exits 1 naming the
  injection path (the sales demo); `.../clean` and `.../mitigated` pass.

### Phase 6 — schema mapper (private repo; not in a2)
- Consumes `a2 check --profile=n8n -f json`; owns `n8n:*` ↔ `CHK-XX`, `maps_to_questions`,
  `triggers_gate`, recomputes `workflow_bundle_hash` (§3.2), sets `tool_version`; validates
  against `maturity-assessment.schema.json` in tests.

## 7. Constraints checklist (HANDOVER §5 → mechanism)

- Determinism → no shelling out except `git` in `n8n:in-git`; sorted iteration everywhere;
  hash-keyed loader cache.
- No secret values in output → redaction inside the scanner + serialization-absence test.
- Never false-pass → `error` return → runner `Errored`; classification uncertainty → `Warn`.
- a2 stays generic → no Tavon/question/gate vocabulary anywhere; check `Detail` text explains
  each n8n risk on its own terms.
- Dependencies → none new planned (stdlib + existing go.mod only).
