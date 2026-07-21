package n8n

import "testing"

func node(nodeType string, params map[string]any) *Node {
	return &Node{Name: "test", Type: nodeType, Parameters: params}
}

func TestIsUntrustedTrigger(t *testing.T) {
	tests := []struct {
		nodeType string
		want     bool
	}{
		{"n8n-nodes-base.webhook", true},
		{"n8n-nodes-base.formTrigger", true},
		{"n8n-nodes-base.emailReadImap", true},
		{"@n8n/n8n-nodes-langchain.chatTrigger", true},
		{"n8n-nodes-base.rssFeedReadTrigger", true},
		{"n8n-nodes-base.scheduleTrigger", false},
		{"n8n-nodes-base.manualTrigger", false},
		{"n8n-nodes-base.set", false},
	}
	for _, tt := range tests {
		t.Run(tt.nodeType, func(t *testing.T) {
			if got := IsUntrustedTrigger(node(tt.nodeType, nil)); got != tt.want {
				t.Errorf("IsUntrustedTrigger(%s) = %v, want %v", tt.nodeType, got, tt.want)
			}
		})
	}
}

func TestClassifySideEffect(t *testing.T) {
	tests := []struct {
		name     string
		nodeType string
		params   map[string]any
		want     Effect
	}{
		{"email send", "n8n-nodes-base.emailSend", nil, EffectWrite},
		{"respond to webhook", "n8n-nodes-base.respondToWebhook", nil, EffectNone},
		{"http post", "n8n-nodes-base.httpRequest", map[string]any{"method": "POST"}, EffectWrite},
		{"http delete", "n8n-nodes-base.httpRequest", map[string]any{"method": "DELETE"}, EffectWrite},
		{"http default get", "n8n-nodes-base.httpRequest", map[string]any{}, EffectNone},
		{"http expression method", "n8n-nodes-base.httpRequest", map[string]any{"method": "={{ $json.m }}"}, EffectUnknown},
		{"hubspot create", "n8n-nodes-base.hubspot", map[string]any{"operation": "create"}, EffectWrite},
		{"postgres insert", "n8n-nodes-base.postgres", map[string]any{"operation": "insert"}, EffectWrite},
		{"sheets read", "n8n-nodes-base.googleSheets", map[string]any{"operation": "getAll"}, EffectNone},
		{"slack post", "n8n-nodes-base.slack", map[string]any{"operation": "postMessage"}, EffectWrite},
		{"op-dependent without op", "n8n-nodes-base.gmail", map[string]any{}, EffectUnknown},
		{"op-dependent odd op", "n8n-nodes-base.gmail", map[string]any{"operation": "frobnicate"}, EffectUnknown},
		{"set node", "n8n-nodes-base.set", nil, EffectNone},
		{"if node", "n8n-nodes-base.if", nil, EffectNone},
		{"unknown vendor node", "acme.futureWidget", nil, EffectUnknown},
		{"exec command", "n8n-nodes-base.executeCommand", nil, EffectWrite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ClassifySideEffect(node(tt.nodeType, tt.params))
			if got != tt.want {
				t.Errorf("ClassifySideEffect = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassifyToolCapability(t *testing.T) {
	tests := []struct {
		name     string
		nodeType string
		params   map[string]any
		want     Capability
	}{
		{"gmail send tool", "n8n-nodes-base.gmailTool", map[string]any{"operation": "send"}, CapWrite},
		{"gmail read tool", "n8n-nodes-base.gmailTool", map[string]any{"operation": "getAll"}, CapRead},
		{"gmail get tool", "n8n-nodes-base.gmailTool", map[string]any{"operation": "get"}, CapRead},
		{"sheets append tool", "n8n-nodes-base.googleSheetsTool", map[string]any{"operation": "append"}, CapWrite},
		{"calculator", "@n8n/n8n-nodes-langchain.toolCalculator", nil, CapRead},
		{"wikipedia", "@n8n/n8n-nodes-langchain.toolWikipedia", nil, CapRead},
		{"code tool", "@n8n/n8n-nodes-langchain.toolCode", nil, CapUnknown},
		{"workflow tool", "@n8n/n8n-nodes-langchain.toolWorkflow", nil, CapUnknown},
		{"http tool get", "@n8n/n8n-nodes-langchain.toolHttpRequest", map[string]any{}, CapRead},
		{"http tool post", "@n8n/n8n-nodes-langchain.toolHttpRequest", map[string]any{"method": "POST"}, CapWrite},
		{"tool without operation", "n8n-nodes-base.slackTool", map[string]any{}, CapUnknown},
		{"unrecognized tool", "acme.futureTool2000", nil, CapUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ClassifyToolCapability(node(tt.nodeType, tt.params))
			if got != tt.want {
				t.Errorf("ClassifyToolCapability = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExternalCallNode(t *testing.T) {
	tests := []struct {
		nodeType string
		want     bool
	}{
		{"n8n-nodes-base.httpRequest", true},
		{"n8n-nodes-base.postgres", true},
		{"n8n-nodes-base.hubspot", true},
		{"n8n-nodes-base.emailSend", true},
		{"n8n-nodes-base.set", false},
		{"n8n-nodes-base.if", false},
		{"n8n-nodes-base.code", false},
		{"n8n-nodes-base.respondToWebhook", false},
	}
	for _, tt := range tests {
		t.Run(tt.nodeType, func(t *testing.T) {
			if got := IsExternalCallNode(node(tt.nodeType, nil)); got != tt.want {
				t.Errorf("IsExternalCallNode(%s) = %v, want %v", tt.nodeType, got, tt.want)
			}
		})
	}
}

func TestCapabilityString(t *testing.T) {
	if CapRead.String() != "read" || CapWrite.String() != "write" || CapUnknown.String() != "unknown" {
		t.Error("Capability.String() mismatch")
	}
}
