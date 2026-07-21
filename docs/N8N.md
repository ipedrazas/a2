# Auditing n8n Workflows with a2

a2 can statically audit exported [n8n](https://n8n.io) workflows the same way
it audits source code: point it at a directory of workflow JSON and it reports
security, reliability, and design issues — deterministically, with no workflow
execution and no access to your n8n instance.

## What you need

A directory of n8n workflow exports. Export from the n8n editor (workflow menu
→ **Download**) or via the CLI/API. One `.json` file per workflow; a2 finds
them in the directory you point it at, or one level below.

```
workflows/
  lead-intake.json
  crm-sync.json
  inbox-assistant.json
```

## Running the audit

Run only the n8n checks with the `n8n` profile:

```bash
a2 check ./workflows --profile=n8n
```

a2 also auto-detects n8n exports, so `a2 check ./workflows` works too — but in a
mixed repository (workflows alongside application code) it will run the language
checks as well. The profile restricts the run to the n8n suite.

Machine-readable output for pipelines:

```bash
a2 check ./workflows --profile=n8n -f json    # or -f toon
```

The command exits non-zero when a critical check fails
(`n8n:agent-injection-path`, `n8n:webhook-auth`, `n8n:secrets`,
`n8n:error-workflow`), so it drops straight into CI.

## What it checks

| Check | Looks for |
|---|---|
| `n8n:secrets` | Hardcoded keys/tokens in parameters, Code nodes, pinned data |
| `n8n:agent-injection-path` | Untrusted input reaching AI agents with write tools |
| `n8n:webhook-auth` | Unauthenticated webhooks that can trigger writes |
| `n8n:agent-tools` | Per-agent read/write tool inventory |
| `n8n:pii-to-ai` | PII fields or full payloads sent to AI models |
| `n8n:error-workflow` | Active workflows with no error workflow |
| `n8n:silent-failures` | Failures swallowed inline with no error branch |
| `n8n:retry-config` | External calls with no (or default-only) retries |
| `n8n:timeouts` | HTTP nodes with no timeout while a caller waits |
| `n8n:code-node-bulk` | Oversized Code nodes; HTTP calls inside code |
| `n8n:pinned-data` | Test data pinned on active workflows |
| `n8n:naming` | Default/non-compliant workflow and node names |
| `n8n:inline-auth` | Literal auth headers instead of credential references |
| `n8n:duplication` | Node sequences copied across workflows |
| `n8n:in-git` | Exports not tracked in version control |

Run `a2 explain n8n:<check>` for the detection logic, and `a2 run n8n:<check>
./workflows` to see the exact findings for one check.

### Reading a finding

Findings name a location, in a stable `file · workflow · node — note` format.
The traversal checks name the full path, which is the point:

```
$ a2 run n8n:agent-injection-path ./workflows
✗ FAIL Agent Injection Path (n8n:agent-injection-path)
    email-agent.json · Inbox Assistant — untrusted input path:
    Email Trigger (IMAP) → AI Agent — agent holds write-capable tools:
    Send Email (gmailTool — operation send)
```

That reads: an inbound email flows into an AI agent that can send email, so
text an outsider puts in an email can make the agent send email on your behalf.

## Configuration

Add an `n8n:` section to `.a2.yaml` to tune thresholds and extend the
classification tables for node types a2 doesn't know yet:

```yaml
n8n:
  code_node_line_limit: 80          # flag Code nodes longer than this
  naming_regex: "^[a-z][a-z0-9 -]+$" # enforce a naming convention (empty = detect defaults)
  error_workflow_fail_ratio: 0.5    # fail below this error-workflow coverage
  duplication_min_sequence: 5       # minimum copied-sequence length to report

  # Extend the built-in classification tables for custom/community nodes:
  untrusted_triggers_extra:
    - "n8n-nodes-acme.customWebhook"
  side_effect_nodes_extra:
    - "n8n-nodes-acme.crmWriter"
  write_tools_extra:
    - "n8n-nodes-acme.crmTool:update"   # a specific write operation
    - "n8n-nodes-acme.dangerTool"        # every operation of this tool

checks:
  n8n:
    disabled:
      - n8n:naming                  # skip a check for this domain
```

## Guarantees

- **Deterministic.** The same exports always produce the same results. a2
  computes a SHA-256 over the bundle so a run is reproducible and disputable on
  evidence.
- **No secret leakage.** Secret values are redacted at detection; only
  locations (workflow, node, field path) appear in output.
- **Never a false pass.** An export a2 can't parse degrades the affected check
  to an **error** status, and a partially-parsed workflow caps an otherwise
  passing check at **warn** — a broken or truncated export can never look clean.
