package n8n

import "sort"

// Edge is one directed connection between two named nodes.
type Edge struct {
	From        string
	To          string
	Type        string // "main", "ai_tool", "ai_languageModel", ...
	OutputIndex int    // source output port; for "main", port 1+ can be the error output
}

// Graph is the directed connection graph of a single workflow, keyed by node
// name (the n8n export references nodes by name, not id).
type Graph struct {
	byName map[string]*Node
	out    map[string][]Edge
	in     map[string][]Edge
}

// ConnMain is the connection type data flows through; ConnAITool attaches a
// tool node to an agent (edge points tool -> agent in the export).
const (
	ConnMain   = "main"
	ConnAITool = "ai_tool"
)

// BuildGraph builds the connection graph for a workflow. Edges referencing
// nodes that don't exist in the nodes array are dropped. Disabled nodes are
// kept: n8n passes data through them, so they stay on the reachability path.
func BuildGraph(wf *Workflow) *Graph {
	g := &Graph{
		byName: make(map[string]*Node, len(wf.Nodes)),
		out:    map[string][]Edge{},
		in:     map[string][]Edge{},
	}
	for i := range wf.Nodes {
		if wf.Nodes[i].Name != "" {
			g.byName[wf.Nodes[i].Name] = &wf.Nodes[i]
		}
	}
	for from, byType := range wf.Connections {
		if g.byName[from] == nil {
			continue
		}
		for connType, ports := range byType {
			for portIdx, targets := range ports {
				for _, t := range targets {
					if t.Node == "" || g.byName[t.Node] == nil {
						continue
					}
					e := Edge{From: from, To: t.Node, Type: connType, OutputIndex: portIdx}
					g.out[from] = append(g.out[from], e)
					g.in[t.Node] = append(g.in[t.Node], e)
				}
			}
		}
	}
	// Deterministic edge order regardless of map iteration.
	for _, edges := range g.out {
		sortEdges(edges)
	}
	for _, edges := range g.in {
		sortEdges(edges)
	}
	return g
}

func sortEdges(edges []Edge) {
	sort.Slice(edges, func(i, j int) bool {
		a, b := edges[i], edges[j]
		if a.From != b.From {
			return a.From < b.From
		}
		if a.To != b.To {
			return a.To < b.To
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return a.OutputIndex < b.OutputIndex
	})
}

// Node returns the node with the given name, or nil.
func (g *Graph) Node(name string) *Node { return g.byName[name] }

// OutEdges returns the outgoing edges of a node.
func (g *Graph) OutEdges(name string) []Edge { return g.out[name] }

// InEdges returns the incoming edges of a node.
func (g *Graph) InEdges(name string) []Edge { return g.in[name] }

// Reachable returns the names of all nodes transitively reachable from start
// (excluding start itself), following only edges of the given connection
// types. No types means all types. "main" traversal includes every output
// port, so error-output branches are followed too. The result is sorted.
func (g *Graph) Reachable(start string, types ...string) []string {
	allowed := typeSet(types)
	seen := map[string]bool{start: true}
	queue := []string{start}
	var result []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, e := range g.out[current] {
			if allowed != nil && !allowed[e.Type] {
				continue
			}
			if seen[e.To] {
				continue
			}
			seen[e.To] = true
			result = append(result, e.To)
			queue = append(queue, e.To)
		}
	}
	sort.Strings(result)
	return result
}

// PathTo returns a shortest path from start to target (inclusive of both),
// following only edges of the given connection types (none = all). Returns
// nil when target is unreachable. Ties resolve deterministically because
// edges are sorted.
func (g *Graph) PathTo(start, target string, types ...string) []string {
	if start == target {
		return []string{start}
	}
	allowed := typeSet(types)
	prev := map[string]string{}
	seen := map[string]bool{start: true}
	queue := []string{start}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, e := range g.out[current] {
			if allowed != nil && !allowed[e.Type] {
				continue
			}
			if seen[e.To] {
				continue
			}
			seen[e.To] = true
			prev[e.To] = current
			if e.To == target {
				var path []string
				for at := target; at != ""; at = prev[at] {
					path = append([]string{at}, path...)
					if at == start {
						break
					}
				}
				return path
			}
			queue = append(queue, e.To)
		}
	}
	return nil
}

// AttachedTools returns the tool nodes wired into the given agent node via
// ai_tool connections (which point tool -> agent in the export), sorted by
// node name.
func (g *Graph) AttachedTools(agentName string) []*Node {
	var tools []*Node
	for _, e := range g.in[agentName] {
		if e.Type != ConnAITool {
			continue
		}
		if n := g.byName[e.From]; n != nil {
			tools = append(tools, n)
		}
	}
	sort.Slice(tools, func(i, j int) bool { return tools[i].Name < tools[j].Name })
	return tools
}

func typeSet(types []string) map[string]bool {
	if len(types) == 0 {
		return nil
	}
	s := make(map[string]bool, len(types))
	for _, t := range types {
		s[t] = true
	}
	return s
}
