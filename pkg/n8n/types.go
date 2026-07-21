// Package n8n parses exported n8n workflow JSON files into a defensive,
// audit-oriented model: a loader that discovers workflow files in a directory,
// typed workflow/node structures that tolerate format drift across n8n
// versions, a directed graph with transitive reachability over the
// connections object, and classification tables for trigger/side-effect/AI
// node types.
//
// Parsing is deliberately lossless and non-fatal: unknown node types are
// retained with their raw parameters, and structural surprises are recorded
// as parse issues rather than dropped, so checks can degrade to an error
// status instead of silently passing.
package n8n

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrNotWorkflow is returned by ParseWorkflow for valid JSON that does not
// have the n8n export shape (a top-level "nodes" array).
var ErrNotWorkflow = errors.New("not an n8n workflow export")

// Workflow is one parsed n8n workflow export.
type Workflow struct {
	Name        string
	ID          string
	Active      bool
	Nodes       []Node
	Connections Connections // nil when missing or unparseable (see ParseIssues)
	Settings    Settings
	PinData     map[string]json.RawMessage // node name -> pinned items; presence is what checks care about
	// ParseIssues records structural surprises hit while parsing (wrong-typed
	// fields, missing connections, unparseable nodes). A workflow with issues
	// is degraded: checks that depend on the affected structure must report
	// an error/warn status, never a pass.
	ParseIssues []string
}

// Degraded reports whether parsing hit structural surprises.
func (w *Workflow) Degraded() bool { return len(w.ParseIssues) > 0 }

// NodeByName returns the node with the given name, or nil.
func (w *Workflow) NodeByName(name string) *Node {
	for i := range w.Nodes {
		if w.Nodes[i].Name == name {
			return &w.Nodes[i]
		}
	}
	return nil
}

// Settings is the workflow-level settings block.
type Settings struct {
	ErrorWorkflow    string         // ID of the linked error workflow; empty = none
	ExecutionTimeout float64        // seconds; 0 = not set
	Raw              map[string]any // full settings object for keys we don't model
}

// Node is one node of a workflow. Unknown node types are retained as-is:
// Type carries whatever the export said and Parameters/Raw keep the full
// payload for classification and secret scanning.
type Node struct {
	ID               string
	Name             string
	Type             string
	TypeVersion      float64
	Disabled         bool
	Parameters       map[string]any
	Credentials      map[string]any // credential-store references (name/id), not secret material
	ContinueOnFail   bool           // legacy error setting
	OnError          string         // "", "stopWorkflow", "continueRegularOutput", "continueErrorOutput"
	RetryOnFail      bool
	MaxTries         int
	WaitBetweenTries int
	Raw              json.RawMessage
}

// StringParam returns Parameters[key] when it is a string.
func (n *Node) StringParam(key string) (string, bool) {
	v, ok := n.Parameters[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// Target is one endpoint of a connection.
type Target struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

// OutputPorts is the per-connection-type structure: one slice of targets per
// output port (port 1 is the error output for nodes with
// onError=continueErrorOutput; If/Switch nodes use ports for branches).
type OutputPorts [][]Target

// UnmarshalJSON accepts both the canonical nested form
// [[{node,...}], [{node,...}]] and a flattened [{node,...}] form seen in
// some template exports, and tolerates null port entries.
func (p *OutputPorts) UnmarshalJSON(data []byte) error {
	var nested [][]Target
	if err := json.Unmarshal(data, &nested); err == nil {
		*p = nested
		return nil
	}
	var flat []Target
	if err := json.Unmarshal(data, &flat); err == nil {
		*p = OutputPorts{flat}
		return nil
	}
	return errors.New("unrecognized connection structure")
}

// Connections maps source node name -> connection type ("main", "ai_tool",
// "ai_languageModel", ...) -> output ports -> targets. Note the n8n export
// points AI sub-node connections from the tool/model node INTO the agent.
type Connections map[string]map[string]OutputPorts

// ParseWorkflow parses one exported workflow file. It returns ErrNotWorkflow
// for JSON without the export shape, an error for non-JSON input, and
// otherwise always a workflow — structural surprises inside a valid export
// are recorded in ParseIssues instead of failing the parse.
func ParseWorkflow(data []byte) (*Workflow, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	rawNodes, ok := top["nodes"]
	if !ok {
		return nil, ErrNotWorkflow
	}
	var nodeList []json.RawMessage
	if err := json.Unmarshal(rawNodes, &nodeList); err != nil {
		return nil, ErrNotWorkflow // "nodes" exists but is not an array: not the export shape
	}

	wf := &Workflow{}
	issue := func(format string, args ...any) {
		wf.ParseIssues = append(wf.ParseIssues, fmt.Sprintf(format, args...))
	}

	decodeString(top, "name", &wf.Name, issue)
	decodeString(top, "id", &wf.ID, issue)
	if raw, ok := top["active"]; ok {
		if err := json.Unmarshal(raw, &wf.Active); err != nil {
			issue("active: %v", err)
		}
	}

	for i, rawNode := range nodeList {
		wf.Nodes = append(wf.Nodes, parseNode(i, rawNode, issue))
	}

	if raw, ok := top["connections"]; ok {
		var conns Connections
		if err := json.Unmarshal(raw, &conns); err != nil {
			issue("connections: %v", err)
		} else {
			wf.Connections = conns
		}
	} else {
		issue("connections object missing (truncated or hand-edited export?)")
	}

	if raw, ok := top["settings"]; ok {
		wf.Settings = parseSettings(raw, issue)
	}
	if raw, ok := top["pinData"]; ok {
		if err := json.Unmarshal(raw, &wf.PinData); err != nil {
			issue("pinData: %v", err)
		}
	}
	return wf, nil
}

// parseNode decodes a single node, salvaging what it can on type surprises.
// encoding/json keeps decoding sibling fields after a type mismatch, so a
// wrong-typed field degrades only itself and gets recorded as an issue.
func parseNode(index int, raw json.RawMessage, issue func(string, ...any)) Node {
	var aux struct {
		ID               string         `json:"id"`
		Name             string         `json:"name"`
		Type             string         `json:"type"`
		TypeVersion      float64        `json:"typeVersion"`
		Disabled         bool           `json:"disabled"`
		Parameters       map[string]any `json:"parameters"`
		Credentials      map[string]any `json:"credentials"`
		ContinueOnFail   bool           `json:"continueOnFail"`
		OnError          string         `json:"onError"`
		RetryOnFail      bool           `json:"retryOnFail"`
		MaxTries         int            `json:"maxTries"`
		WaitBetweenTries int            `json:"waitBetweenTries"`
	}
	if err := json.Unmarshal(raw, &aux); err != nil {
		name := aux.Name
		if name == "" {
			name = fmt.Sprintf("node[%d]", index)
		}
		issue("node %q: %v", name, err)
	}
	return Node{
		ID:               aux.ID,
		Name:             aux.Name,
		Type:             aux.Type,
		TypeVersion:      aux.TypeVersion,
		Disabled:         aux.Disabled,
		Parameters:       aux.Parameters,
		Credentials:      aux.Credentials,
		ContinueOnFail:   aux.ContinueOnFail,
		OnError:          aux.OnError,
		RetryOnFail:      aux.RetryOnFail,
		MaxTries:         aux.MaxTries,
		WaitBetweenTries: aux.WaitBetweenTries,
		Raw:              raw,
	}
}

func parseSettings(raw json.RawMessage, issue func(string, ...any)) Settings {
	s := Settings{}
	if err := json.Unmarshal(raw, &s.Raw); err != nil {
		issue("settings: %v", err)
		return s
	}
	if v, ok := s.Raw["errorWorkflow"]; ok {
		if str, ok := v.(string); ok {
			s.ErrorWorkflow = str
		} else {
			issue("settings.errorWorkflow: expected string, got %T", v)
		}
	}
	if v, ok := s.Raw["executionTimeout"]; ok {
		if f, ok := v.(float64); ok && f > 0 {
			s.ExecutionTimeout = f
		}
	}
	return s
}

func decodeString(top map[string]json.RawMessage, key string, dst *string, issue func(string, ...any)) {
	raw, ok := top[key]
	if !ok {
		return
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		issue("%s: %v", key, err)
	}
}
