package n8n

import "strings"

// This file is the node classification data for the n8n checks: which
// trigger types carry untrusted input, which node types cause side effects,
// which are AI/agent nodes, and how attached tools are classified read vs
// write. It is deliberately data-over-code and kept in one file so extending
// coverage for new node types is a one-line change (and so the tables can be
// merged with user-supplied extensions from .a2.yaml).
//
// Node type names follow the n8n export format (n8n 1.x): built-in nodes are
// "n8n-nodes-base.<name>", LangChain/AI nodes are
// "@n8n/n8n-nodes-langchain.<name>".

// UntrustedTriggerTypes are trigger node types whose payload an outside
// party can influence: anyone can call a webhook, send an email, submit a
// form, or author an RSS entry. Workflows starting here process attacker
// -controllable input.
var UntrustedTriggerTypes = map[string]bool{
	"n8n-nodes-base.webhook":                 true,
	"n8n-nodes-base.formTrigger":             true,
	"n8n-nodes-base.emailReadImap":           true,
	"n8n-nodes-base.gmailTrigger":            true,
	"n8n-nodes-base.microsoftOutlookTrigger": true,
	"n8n-nodes-base.rssFeedReadTrigger":      true,
	"n8n-nodes-base.telegramTrigger":         true,
	"n8n-nodes-base.slackTrigger":            true,
	"n8n-nodes-base.whatsAppTrigger":         true,
	"n8n-nodes-base.typeformTrigger":         true,
	"n8n-nodes-base.jotFormTrigger":          true,
	"@n8n/n8n-nodes-langchain.chatTrigger":   true,
	"@n8n/n8n-nodes-langchain.mcpTrigger":    true,
	"n8n-nodes-base.emailReadImapV2":         true, // defensive: seen in some version drifts
	"n8n-nodes-base.microsoftTeamsTrigger":   true,
	"n8n-nodes-base.discordTrigger":          true,
	"n8n-nodes-base.twilioTrigger":           true,
	"n8n-nodes-base.postmarkTrigger":         true,
	"n8n-nodes-base.mailgunTrigger":          true,
	"n8n-nodes-base.surveyMonkeyTrigger":     true,
	"n8n-nodes-base.googleFormsTrigger":      true,
}

// alwaysSideEffect are node types that mutate external state regardless of
// their parameters (they exist to send or write).
var alwaysSideEffect = map[string]string{
	"n8n-nodes-base.emailSend":        "sends email",
	"n8n-nodes-base.twilio":           "sends SMS/voice",
	"n8n-nodes-base.pushover":         "sends notifications",
	"n8n-nodes-base.pushbullet":       "sends notifications",
	"n8n-nodes-base.respondToWebhook": "", // responds to the caller only: not a side effect
	"n8n-nodes-base.executeCommand":   "executes shell commands",
	"n8n-nodes-base.ssh":              "executes remote commands",
}

// operationDependent are node types whose effect depends on the selected
// resource/operation parameter (e.g. Gmail "send" vs "getAll"). The effect
// is classified through readOperations/writeOperations below.
var operationDependent = map[string]bool{
	"n8n-nodes-base.gmail":            true,
	"n8n-nodes-base.microsoftOutlook": true,
	"n8n-nodes-base.slack":            true,
	"n8n-nodes-base.telegram":         true,
	"n8n-nodes-base.discord":          true,
	"n8n-nodes-base.postgres":         true,
	"n8n-nodes-base.mySql":            true,
	"n8n-nodes-base.microsoftSql":     true,
	"n8n-nodes-base.mongoDb":          true,
	"n8n-nodes-base.redis":            true,
	"n8n-nodes-base.supabase":         true,
	"n8n-nodes-base.airtable":         true,
	"n8n-nodes-base.notion":           true,
	"n8n-nodes-base.googleSheets":     true,
	"n8n-nodes-base.hubspot":          true,
	"n8n-nodes-base.salesforce":       true,
	"n8n-nodes-base.pipedrive":        true,
	"n8n-nodes-base.zendesk":          true,
	"n8n-nodes-base.stripe":           true,
	"n8n-nodes-base.shopify":          true,
	"n8n-nodes-base.jira":             true,
	"n8n-nodes-base.github":           true,
	"n8n-nodes-base.gitlab":           true,
	"n8n-nodes-base.googleDrive":      true,
	"n8n-nodes-base.s3":               true,
	"n8n-nodes-base.awsS3":            true,
	"n8n-nodes-base.ftp":              true,
}

// readOperations and writeOperations classify the `operation` parameter of
// operation-dependent nodes and app tool nodes. Matching is by exact,
// case-insensitive operation name.
var readOperations = map[string]bool{
	"get": true, "getall": true, "getmany": true, "read": true,
	"search": true, "list": true, "lookup": true, "download": true,
	"select": true, "find": true, "findone": true, "aggregate": true,
	"gethistory": true, "getprofile": true, "getlabels": true,
}

var writeOperations = map[string]bool{
	"send": true, "sendandwait": true, "create": true, "update": true,
	"upsert": true, "delete": true, "deleteall": true, "append": true,
	"appendorupdate": true, "insert": true, "post": true, "postmessage": true,
	"execute": true, "executequery": true, "executecommand": true, "set": true,
	"move": true, "copy": true, "upload": true, "write": true, "publish": true,
	"reply": true, "addlabels": true, "removelabels": true, "invoke": true,
}

// AgentNodeTypes are the AI agent node types: LLM nodes that can be given
// tools and decide autonomously which to call.
var AgentNodeTypes = map[string]bool{
	"@n8n/n8n-nodes-langchain.agent":           true,
	"@n8n/n8n-nodes-langchain.openAiAssistant": true,
}

// AIProcessingNodeTypes are node types that feed their input into an LLM
// (agents included). Untrusted input reaching any of these is a prompt
// -injection path; reaching an agent with write-capable tools is the
// critical case.
var AIProcessingNodeTypes = map[string]bool{
	"@n8n/n8n-nodes-langchain.agent":                true,
	"@n8n/n8n-nodes-langchain.openAiAssistant":      true,
	"@n8n/n8n-nodes-langchain.chainLlm":             true,
	"@n8n/n8n-nodes-langchain.chainSummarization":   true,
	"@n8n/n8n-nodes-langchain.informationExtractor": true,
	"@n8n/n8n-nodes-langchain.textClassifier":       true,
	"@n8n/n8n-nodes-langchain.sentimentAnalysis":    true,
	"@n8n/n8n-nodes-langchain.openAi":               true,
	"n8n-nodes-base.openAi":                         true, // legacy OpenAI node
}

// readOnlyTools are LangChain tool node types that can only look things up.
var readOnlyTools = map[string]bool{
	"@n8n/n8n-nodes-langchain.toolCalculator":       true,
	"@n8n/n8n-nodes-langchain.toolWikipedia":        true,
	"@n8n/n8n-nodes-langchain.toolSerpApi":          true,
	"@n8n/n8n-nodes-langchain.toolWolframAlpha":     true,
	"@n8n/n8n-nodes-langchain.toolVectorStore":      true,
	"@n8n/n8n-nodes-langchain.retrieverVectorStore": true,
}

// arbitraryCodeTools run code or sub-workflows the classifier cannot see
// into; their capability is unknown, which checks must treat as
// potentially write-capable (never classify uncertainty as safe).
var arbitraryCodeTools = map[string]bool{
	"@n8n/n8n-nodes-langchain.toolCode":     true,
	"@n8n/n8n-nodes-langchain.toolWorkflow": true,
	"n8n-nodes-base.executeWorkflowTool":    true,
}

// Effect is the side-effect classification of a node.
type Effect int

const (
	// EffectNone: the node only reads or transforms data.
	EffectNone Effect = iota
	// EffectWrite: the node mutates external state (sends, writes, executes).
	EffectWrite
	// EffectUnknown: the classifier cannot tell (unrecognized type or
	// expression-driven operation). Checks must not treat this as safe.
	EffectUnknown
)

// Capability is the read/write classification of an agent tool.
type Capability int

const (
	CapRead Capability = iota
	CapWrite
	CapUnknown
)

func (c Capability) String() string {
	switch c {
	case CapRead:
		return "read"
	case CapWrite:
		return "write"
	default:
		return "unknown"
	}
}

// IsUntrustedTrigger reports whether the node is a trigger whose payload an
// outside party can influence.
func IsUntrustedTrigger(n *Node) bool { return UntrustedTriggerTypes[n.Type] }

// IsAgentNode reports whether the node is an AI agent (tool-capable LLM).
func IsAgentNode(n *Node) bool { return AgentNodeTypes[n.Type] }

// IsAIProcessingNode reports whether the node feeds its input to an LLM.
func IsAIProcessingNode(n *Node) bool { return AIProcessingNodeTypes[n.Type] }

// IsHTTPRequestNode reports whether the node is the generic HTTP Request node.
func IsHTTPRequestNode(n *Node) bool { return n.Type == "n8n-nodes-base.httpRequest" }

// IsExternalCallNode reports whether the node calls an external system
// (relevant for retry/timeout discipline): HTTP, databases, SaaS
// integrations — anything operation-dependent or always side-effectful.
func IsExternalCallNode(n *Node) bool {
	if IsHTTPRequestNode(n) || operationDependent[n.Type] {
		return true
	}
	reason, known := alwaysSideEffect[n.Type]
	return known && reason != ""
}

// ClassifySideEffect classifies whether executing the node mutates external
// state, with a short human-readable reason.
func ClassifySideEffect(n *Node) (Effect, string) {
	if reason, ok := alwaysSideEffect[n.Type]; ok {
		if reason == "" {
			return EffectNone, ""
		}
		return EffectWrite, reason
	}
	if IsHTTPRequestNode(n) {
		return classifyHTTPMethod(n)
	}
	if operationDependent[n.Type] {
		return classifyOperation(n)
	}
	if strings.HasPrefix(n.Type, "n8n-nodes-base.") || strings.HasPrefix(n.Type, "@n8n/n8n-nodes-langchain.") {
		// A known-namespace node not in any table: assume core/transform
		// behaviour (Set, If, Code, Merge, ...). Code nodes can call out via
		// fetch; that is covered by the code-node check, not this classifier.
		return EffectNone, ""
	}
	return EffectUnknown, "unrecognized node type " + n.Type
}

// ClassifyToolCapability classifies an agent-attached tool node as read-only
// or write-capable, with a short reason. App nodes used as tools (type
// suffix "Tool", e.g. gmailTool) are classified by their operation
// parameter, mirroring the underlying node.
func ClassifyToolCapability(n *Node) (Capability, string) {
	switch {
	case readOnlyTools[n.Type]:
		return CapRead, "lookup-only tool"
	case arbitraryCodeTools[n.Type]:
		return CapUnknown, "runs arbitrary code/sub-workflow"
	case n.Type == "@n8n/n8n-nodes-langchain.toolHttpRequest":
		effect, reason := classifyHTTPMethod(n)
		return effectToCapability(effect, reason, "HTTP tool")
	case strings.HasSuffix(n.Type, "Tool"):
		effect, reason := classifyOperation(n)
		return effectToCapability(effect, reason, "app tool")
	default:
		return CapUnknown, "unrecognized tool type " + n.Type
	}
}

func effectToCapability(e Effect, reason, fallback string) (Capability, string) {
	if reason == "" {
		reason = fallback
	}
	switch e {
	case EffectNone:
		return CapRead, reason
	case EffectWrite:
		return CapWrite, reason
	default:
		return CapUnknown, reason
	}
}

// classifyOperation classifies a node by its `operation` parameter.
func classifyOperation(n *Node) (Effect, string) {
	op, ok := n.StringParam("operation")
	if !ok {
		// Many nodes default their operation when the parameter is absent
		// (often to a create/send-style default). Absent = cannot tell.
		return EffectUnknown, "operation not set"
	}
	if strings.HasPrefix(op, "=") {
		return EffectUnknown, "operation is an expression"
	}
	key := strings.ToLower(op)
	if writeOperations[key] {
		return EffectWrite, "operation " + op
	}
	if readOperations[key] {
		return EffectNone, "operation " + op
	}
	return EffectUnknown, "unrecognized operation " + op
}

// classifyHTTPMethod classifies an HTTP node by its method parameter
// (mutating methods are writes; n8n defaults to GET when absent).
func classifyHTTPMethod(n *Node) (Effect, string) {
	method, ok := n.StringParam("method")
	if !ok {
		method = "GET"
	}
	if strings.HasPrefix(method, "=") {
		return EffectUnknown, "HTTP method is an expression"
	}
	switch strings.ToUpper(method) {
	case "GET", "HEAD", "OPTIONS":
		return EffectNone, "HTTP " + strings.ToUpper(method)
	case "POST", "PUT", "PATCH", "DELETE":
		return EffectWrite, "HTTP " + strings.ToUpper(method)
	default:
		return EffectUnknown, "unrecognized HTTP method " + method
	}
}
