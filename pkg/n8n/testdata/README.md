# n8n workflow fixtures

Hand-written, minimal n8n workflow exports for the `n8n:*` checks. One concern
per directory:

- `clean/` — a small estate that passes everything: header-authed webhook,
  error workflow configured, retries with non-default values, timeouts on
  HTTP nodes, credential-store references only, descriptive names.
- `vulnerable/` — the demo estate: unauthenticated webhook feeding a CRM
  write; IMAP trigger feeding an AI agent with a send-capable Gmail tool;
  fake secrets in HTTP headers; 1-of-3 active workflows with an error
  workflow; a 122-line Code node with an embedded `fetch()`; `pinData` on an
  active workflow; default-ish node and workflow names.
- `mitigated/` — the vulnerable shapes, fixed: HMAC/header-authed webhook,
  agent tools read-only, credentials referenced from the store, error
  workflows set, retries and timeouts configured.
- `malformed/` — exercises error-degradation: truncated JSON, unknown node
  types with a wrong-typed `typeVersion`, a missing `connections` object,
  and a `connections` value with the wrong JSON type. Parsing must degrade,
  never panic or silently pass.
- `duplication/` — two workflows sharing a copied 5-node sequence
  (httpRequest → code → if → set → slack) with near-identical parameters.

## Export format

Fixtures target the n8n 1.x export format (verified against the n8n docs and
template exports, July 2026): top-level `nodes` array + `connections` object
keyed by node *name*; AI sub-nodes (models, tools) connect *into* the agent
via `ai_languageModel`/`ai_tool` connection types; node-level error/retry
settings are `onError`/`continueOnFail`, `retryOnFail`, `maxTries`,
`waitBetweenTries`.

## Secrets

All credential-looking strings are deliberately fake (`sk-FAKE…`,
`AKIAFAKE…`) but still match real detection patterns. Tests assert these
strings never appear in check output — do not replace them with anything
real-looking.
