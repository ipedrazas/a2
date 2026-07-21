package n8ncheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Default thresholds for configurable checks.
const (
	defaultCodeNodeLineLimit      = 80
	defaultErrorWorkflowFailRatio = 0.5
)

// Register returns all n8n check registrations. All checks parse exported
// workflow JSON only — no network, no n8n API, no workflow execution — so
// every check is SpeedFast.
func Register(cfg *config.Config) []checker.CheckRegistration {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	lineLimit := cfg.N8n.CodeNodeLineLimit
	if lineLimit <= 0 {
		lineLimit = defaultCodeNodeLineLimit
	}
	failRatio := cfg.N8n.ErrorWorkflowFailRatio
	if failRatio <= 0 {
		failRatio = defaultErrorWorkflowFailRatio
	}

	n8nLang := []checker.Language{checker.LangN8n}

	return []checker.CheckRegistration{
		{
			Checker: &ErrorWorkflowCheck{FailRatio: failRatio},
			Meta: checker.CheckMeta{
				ID:          "n8n:error-workflow",
				Name:        "Error Workflow Coverage",
				Description: "Verifies active workflows have an error workflow configured so failures alert someone.",
				Languages:   n8nLang,
				Critical:    true,
				Order:       150,
				Suggestion:  "Create an error workflow (e.g. notify a channel) and set it in each workflow's settings",
				Detail: "Reads settings.errorWorkflow on every active workflow in the export. " +
					"Without one, a failed execution dies silently: no alert, no retry, no record anyone looked. " +
					"Reports the coverage ratio; fails below 50% coverage (configurable via n8n.error_workflow_fail_ratio), warns below 100%.",
				FixPrompt: "Create a single n8n error workflow that posts failures (workflow name, execution URL, error message) to the team's alert channel, then set it as the error workflow in the settings of every active workflow listed in the findings.",
			},
		},
		{
			Checker: &SilentFailuresCheck{},
			Meta: checker.CheckMeta{
				ID:          "n8n:silent-failures",
				Name:        "Silent Failure Settings",
				Description: "Finds nodes whose failures are swallowed inline instead of routed to an error branch.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       151,
				Suggestion:  "Route node failures to a connected error branch, or let them stop the workflow",
				Detail: "Flags nodes with continueOnFail / onError=continueRegularOutput (errors merge into the normal output and downstream nodes process garbage), " +
					"and nodes with onError=continueErrorOutput whose error output is not connected to anything. " +
					"Either way a failure leaves no trace.",
				FixPrompt: "For each listed node, decide the failure policy: remove continueOnFail so failures stop the workflow and trigger the error workflow, or switch to onError=continueErrorOutput and wire the error output to a handling branch (alert, dead-letter, compensation).",
			},
		},
		{
			Checker: &RetryCheck{},
			Meta: checker.CheckMeta{
				ID:          "n8n:retry-config",
				Name:        "Retry Configuration",
				Description: "Verifies external-call nodes (HTTP, DB, SaaS) have consciously configured retries.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       152,
				Suggestion:  "Enable retryOnFail on external-call nodes and tune maxTries/waitBetweenTries per dependency",
				Detail: "External calls fail transiently; nodes without retryOnFail turn every blip into a workflow failure. " +
					"Also flags estates where every retry sits exactly on the defaults (maxTries=3, wait=1000ms) — uniformity implies nobody chose a policy. " +
					"Emits the retry-setting distribution as detail.",
				FixPrompt: "Enable retryOnFail on each listed external-call node. Choose maxTries and waitBetweenTries per dependency (idempotent reads can retry aggressively; writes need care about duplicates — prefer fewer tries plus an error branch).",
			},
		},
		{
			Checker: &TimeoutsCheck{},
			Meta: checker.CheckMeta{
				ID:          "n8n:timeouts",
				Name:        "HTTP Timeout Discipline",
				Description: "Flags HTTP Request nodes without a timeout in workflows where a caller is waiting.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       153,
				Suggestion:  "Set options.timeout on HTTP Request nodes in webhook/form/chat-triggered workflows",
				Detail: "In webhook-, form-, or chat-triggered workflows a caller is synchronously waiting on the response. " +
					"An HTTP Request node without an explicit timeout can hang the whole path on one slow upstream, tying up the caller and the execution slot.",
				FixPrompt: "Set an explicit options.timeout (in ms) on each listed HTTP Request node, sized to the upstream's realistic latency budget, and decide what the workflow returns to the waiting caller when the timeout fires.",
			},
		},
		{
			Checker: &CodeBulkCheck{LineLimit: lineLimit},
			Meta: checker.CheckMeta{
				ID:          "n8n:code-node-bulk",
				Name:        "Code Node Bulk",
				Description: "Flags oversized Code nodes and Code nodes that embed HTTP calls.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       154,
				Suggestion:  "Split oversized Code nodes and move embedded HTTP calls into HTTP Request nodes",
				Detail: "Code/Function nodes beyond the line threshold (default 80, configurable via n8n.code_node_line_limit) are applications hiding inside a workflow: untested, unversioned, invisible to review. " +
					"Code nodes that call fetch/axios/urllib entangle logic with integration and bypass n8n's retry, timeout, and credential handling. Reports per-node line counts.",
				FixPrompt: "For each listed Code node: extract the HTTP calls into dedicated HTTP Request nodes (restoring retries, timeouts, and credential references), and split the remaining logic into smaller named Code nodes or a sub-workflow so each piece is reviewable.",
			},
		},
		{
			Checker: &PinnedDataCheck{},
			Meta: checker.CheckMeta{
				ID:          "n8n:pinned-data",
				Name:        "Pinned Data In Production",
				Description: "Flags pinned test data left on active workflows.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       155,
				Suggestion:  "Unpin test data before activating workflows",
				Detail:      "pinData freezes a node's output for testing. On an active workflow it is a test fixture left live — a marker that production and testing share an environment, and a trap when someone re-executes pinned branches with stale data.",
				FixPrompt:   "Unpin the data on each listed node (workflow editor → node → unpin). If the pinned samples are useful test cases, move them into a dedicated inactive test copy of the workflow instead.",
			},
		},
		{
			Checker: &NamingCheck{Regex: cfg.N8n.NamingRegex},
			Meta: checker.CheckMeta{
				ID:          "n8n:naming",
				Name:        "Naming Convention",
				Description: "Measures workflow and node name compliance against a convention, or detects untouched default names.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       156,
				Suggestion:  "Rename default-named workflows/nodes; set n8n.naming_regex to enforce a real convention",
				Detail: "With n8n.naming_regex configured, checks every workflow and node name against the convention and reports the compliance percentage. " +
					"Without one, it falls back to detecting untouched n8n defaults (\"My workflow\", \"Webhook1\", \"Code2\") — names nobody chose, in an estate nobody can navigate.",
				FixPrompt: "Rename each listed workflow/node to say what it does (e.g. \"HTTP Request1\" → \"Fetch CRM contacts\"). Then agree a team convention and add it to .a2.yaml as n8n.naming_regex so compliance is measured, not guessed.",
			},
		},
		{
			Checker: &InlineAuthCheck{},
			Meta: checker.CheckMeta{
				ID:          "n8n:inline-auth",
				Name:        "Inline Auth Vs Credential Store",
				Description: "Flags HTTP nodes carrying literal auth headers/tokens instead of credential references.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       157,
				Suggestion:  "Replace literal auth parameters with n8n credential references",
				Detail: "HTTP nodes with literal values in auth-looking header or query parameters (Authorization, api_key, token, ...) bypass the n8n credential store: the secret is visible to every editor, exported with every workflow, and rotated never. " +
					"A weaker signal than n8n:secrets (the value may be a placeholder), so it warns. Findings name the parameter, never the value.",
				FixPrompt: "For each listed node, create an n8n credential of the right type (header auth, OAuth2, API key), move the secret into it, switch the node's authentication to that credential, and rotate the exposed value — it has already been exported in plain text.",
			},
		},
		{
			Checker: &InGitCheck{},
			Meta: checker.CheckMeta{
				ID:          "n8n:in-git",
				Name:        "Workflows Versioned In Git",
				Description: "Verifies workflow exports are tracked in a git repository.",
				Languages:   n8nLang,
				Critical:    false,
				Order:       158,
				Command:     "git ls-files -- <workflow files>",
				Suggestion:  "Commit workflow exports to a git repository with regular re-exports",
				Detail: "Automation that only exists inside the live n8n instance has no history, no review, and no rollback. " +
					"Checks that the export directory is a git work tree and every workflow file is tracked. " +
					"If producing a versioned export for this audit was itself the hard part, that is the finding.",
				FixPrompt: "Initialize (or choose) a git repository for the workflow exports, commit the listed files, and set up a recurring export-and-commit job (n8n CLI or API) so the repository stays the reviewable source of truth.",
			},
		},
	}
}
