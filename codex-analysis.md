# A2 Analysis and Improvement Suggestions (Revised)

## Executive Summary
A2 is a strong, pragmatic quality gate for multi-language repos with clear CLI ergonomics, deterministic checks, and a useful maturity model. The biggest opportunity is **signal-to-noise**: reduce false positives, improve explainability, and tighten prioritization so repeated daily runs highlight the *most important work* first. The second opportunity is **output UX**: make it obvious what to fix next and why.

This revision incorporates feedback on the initial report.

## What’s Working Well
- **Clear product intent and narrative**: The README and PRD define a compelling goal—baseline quality for AI-generated code, with profiles/targets to keep expectations realistic. See: `/Users/ivan/Workspace/github.com/ipedrazas/a2/README.md`, `/Users/ivan/Workspace/github.com/ipedrazas/a2/prd.md`.
- **Consistent architecture**: Checks are structured, registered, and run through a single runner. The system is easy to extend and reason about: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/checks`, `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/runner/runner.go`.
- **Output flexibility**: Pretty, JSON, and TOON formats cover human, CI, and agent workflows. Verbosity options are practical: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/output/pretty.go`.
- **Config-driven behavior**: `.a2.yaml` supports profiles, targets, thresholds, and external checks. This is the right foundation for enterprise-grade usage: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/config/config.go`.
- **Server mode**: The HTTP server and UI provide a path toward team usage beyond local CLI: `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/server/server.go`.

## Pain Points Affecting Daily Use
1. **Output lacks prioritization**
   - Results are rendered uniformly. There is no “fix these first” summary. Daily runs should immediately answer *what matters most*.

2. **Explainability is not structured**
   - `Result` has `Message` and `RawOutput`, but no dedicated “Reason/Evidence” field to explain *why* a check failed.

3. **Allowlist and suppression capabilities are under-documented**
   - The tool already has an allowlist system for filesystem checks and skips many file categories, but the documentation doesn’t surface this clearly.

4. **Skip reasoning is invisible**
   - When a check is skipped by profile/target/config, output doesn’t say why. This makes scores and results feel arbitrary.

## Suggestions: Focused Improvements

### 1) Add a “Top Issues” summary block
**Why:** This is the highest ROI change for daily use.

**Implementation:**
- Rank by severity (Fail > Warn), criticality, and check order.
- Show 3–5 issues at the top of pretty output with suggestions.
- Use existing check metadata.

**Where:** `/Users/ivan/Workspace/github.com/ipedrazas/a2/pkg/output/pretty.go`

### 2) Show skip reasons in verbose output
**Why:** Users should know *why* a check is absent (profile/target/skip list).

**Implementation:**
- When `-v` is used, list skipped checks and their reason.
- Source reasons from profile/target/CLI skip flags.

### 3) Document allowlist and suppression
**Why:** It already exists but users don’t know it.

**Implementation:**
- Add a section in `docs/CHECKS.md` showing `security.filesystem.allow` examples.
- Mention wildcard usage and file:line / file:match formats.

### 4) Add a structured Reason/Evidence field
**Why:** Helps users judge false positives and prioritize fixes.

**Implementation:**
- Extend `checker.Result` with a `Reason` field.
- Standardize messages to use `Reason` for “why” and `Message` for “what.”

### 5) Inline suppression comments
**Why:** Faster and more ergonomic than editing `.a2.yaml`.

**Implementation:**
- Support `// a2:allow filesystem` in Go (extend to other langs later).

### 6) Optional fail-fast for parallel mode
**Why:** Not a correctness issue, but improves time-to-feedback.

**Implementation:**
- Add `--fail-fast` to cancel checks on first critical failure.
- Keep current parallel behavior as default.

### 7) Baseline mode for legacy repos
**Why:** Turns overwhelming first runs into incremental improvement.

**Implementation:**
- Store baseline failures and only report regressions.
- Can be opt-in via `a2 baseline` or config toggle.

## Notes on Priorities (Adjusted)
- **Detective vs verification weighting**: Not urgent; keep score simple until users complain.
- **Caching for filesystem scans**: Defer until real performance pain exists.
- **Server hardening**: Long-term concern for multi-tenant deployments.

## Suggested Roadmap

### Short-term (1–2 weeks)
1. Top Issues summary block in pretty output.
2. Skipped-checks visibility in `-v`.
3. Allowlist/suppression documentation.

### Medium-term (1–2 months)
4. `Reason/Evidence` field in results.
5. Inline suppression comments.
6. Optional `--fail-fast` for parallel mode.
7. Baseline mode for legacy repos.

### Long-term (3+ months)
8. Safe-variable inference for non-Go languages.
9. “Daily/weekly/CI” presets as aliases of profiles/targets.
10. Server quotas and rate limits if multi-tenant usage grows.

## Closing Note
A2’s core idea is strong. The biggest gains now are **trust and prioritization**. If the tool consistently tells users what matters most, with minimal noise and clear evidence, daily usage will naturally improve without increasing complexity.
