package n8n

import (
	"slices"
	"sort"
	"strings"
)

// Bounds for path extraction so pathological graphs stay cheap.
const (
	maxPathsPerWorkflow = 64
	maxPathDepth        = 200
)

// DuplicationCluster is a node-type sequence shared by two or more
// workflows: a copy-pasted flow and a sub-workflow candidate.
type DuplicationCluster struct {
	Types       []string // the shared node-type sequence, in flow order
	Occurrences []ClusterOccurrence
	// ParamSimilarity is the average parameter-key overlap (0..1) between
	// the first two occurrences' aligned nodes — a rough "how identical is
	// the copy" signal.
	ParamSimilarity float64
}

// ClusterOccurrence is one instance of a shared sequence.
type ClusterOccurrence struct {
	File     string
	Workflow string
	Nodes    []string // node names, in flow order
}

// FindDuplicateSequences finds order-preserving node-type sequences of at
// least minLen shared across different workflow files. For each pair of
// workflows it reports the longest common flow subsequence (per path pair),
// merging identical sequences into one cluster. Same-type sequences inside
// a single workflow are not reported — the check is about cross-workflow
// copy-paste.
func FindDuplicateSequences(workflows []*WorkflowFile, minLen int) []DuplicationCluster {
	if minLen < 2 {
		minLen = 2
	}
	type pathInfo struct {
		names []string
		types []string
		nodes map[string]*Node
	}
	perFile := make([][]pathInfo, len(workflows))
	for i, wf := range workflows {
		g := BuildGraph(wf.Workflow)
		byName := map[string]*Node{}
		for j := range wf.Workflow.Nodes {
			byName[wf.Workflow.Nodes[j].Name] = &wf.Workflow.Nodes[j]
		}
		for _, names := range mainFlowPaths(wf.Workflow, g) {
			types := make([]string, len(names))
			for k, n := range names {
				types[k] = byName[n].Type
			}
			perFile[i] = append(perFile[i], pathInfo{names: names, types: types, nodes: byName})
		}
	}

	clusters := map[string]*DuplicationCluster{}
	occurrenceSeen := map[string]bool{}
	addOccurrence := func(key string, types []string, wf *WorkflowFile, names []string, similarity float64) {
		cluster, ok := clusters[key]
		if !ok {
			cluster = &DuplicationCluster{Types: types, ParamSimilarity: similarity}
			clusters[key] = cluster
		}
		occKey := key + "|" + wf.Path + "|" + strings.Join(names, "|")
		if occurrenceSeen[occKey] {
			return
		}
		occurrenceSeen[occKey] = true
		cluster.Occurrences = append(cluster.Occurrences, ClusterOccurrence{
			File:     wf.Path,
			Workflow: wf.Workflow.Name,
			Nodes:    names,
		})
	}

	for i := range workflows {
		for j := i + 1; j < len(workflows); j++ {
			for _, pa := range perFile[i] {
				for _, pb := range perFile[j] {
					length, endA, endB := longestCommonRun(pa.types, pb.types)
					if length < minLen {
						continue
					}
					types := pa.types[endA-length : endA]
					key := strings.Join(types, "→")
					namesA := pa.names[endA-length : endA]
					namesB := pb.names[endB-length : endB]
					similarity := paramKeySimilarity(pa.nodes, namesA, pb.nodes, namesB)
					addOccurrence(key, types, workflows[i], namesA, similarity)
					addOccurrence(key, types, workflows[j], namesB, similarity)
				}
			}
		}
	}

	var result []DuplicationCluster
	for _, cluster := range clusters {
		files := map[string]bool{}
		for _, occ := range cluster.Occurrences {
			files[occ.File] = true
		}
		if len(files) < 2 {
			continue
		}
		sort.Slice(cluster.Occurrences, func(a, b int) bool {
			if cluster.Occurrences[a].File != cluster.Occurrences[b].File {
				return cluster.Occurrences[a].File < cluster.Occurrences[b].File
			}
			return strings.Join(cluster.Occurrences[a].Nodes, "|") < strings.Join(cluster.Occurrences[b].Nodes, "|")
		})
		result = append(result, *cluster)
	}
	// Longest sequences first, then lexical, for deterministic output.
	sort.Slice(result, func(a, b int) bool {
		if len(result[a].Types) != len(result[b].Types) {
			return len(result[a].Types) > len(result[b].Types)
		}
		return strings.Join(result[a].Types, "→") < strings.Join(result[b].Types, "→")
	})
	return result
}

// mainFlowPaths linearizes a workflow into main-connection paths from each
// root (a node with no incoming main edge), depth-first with sorted edges
// for determinism. Branches produce one path per branch; cycles are cut.
func mainFlowPaths(wf *Workflow, g *Graph) [][]string {
	var roots []string
	for i := range wf.Nodes {
		name := wf.Nodes[i].Name
		if name == "" {
			continue
		}
		hasMainIn := false
		for _, e := range g.InEdges(name) {
			if e.Type == ConnMain {
				hasMainIn = true
				break
			}
		}
		if !hasMainIn {
			roots = append(roots, name)
		}
	}
	sort.Strings(roots)

	var paths [][]string
	var walk func(name string, current []string)
	walk = func(name string, current []string) {
		if len(paths) >= maxPathsPerWorkflow || len(current) >= maxPathDepth {
			paths = append(paths, append([]string(nil), current...))
			return
		}
		current = append(current, name)
		var next []string
		for _, e := range g.OutEdges(name) {
			if e.Type != ConnMain {
				continue
			}
			if !slices.Contains(current, e.To) {
				next = append(next, e.To)
			}
		}
		if len(next) == 0 {
			paths = append(paths, append([]string(nil), current...))
			return
		}
		for _, n := range next {
			if len(paths) >= maxPathsPerWorkflow {
				break
			}
			walk(n, current)
		}
	}
	for _, root := range roots {
		if len(paths) >= maxPathsPerWorkflow {
			break
		}
		walk(root, nil)
	}
	return paths
}

// longestCommonRun returns the length and end offsets (exclusive) of the
// longest contiguous run shared by a and b.
func longestCommonRun(a, b []string) (length, endA, endB int) {
	prev := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				current[j] = prev[j-1] + 1
				if current[j] > length {
					length, endA, endB = current[j], i, j
				}
			} else {
				current[j] = 0
			}
		}
		prev, current = current, prev
	}
	return length, endA, endB
}

// paramKeySimilarity averages the Jaccard overlap of top-level parameter
// keys across aligned node pairs — parameter *shape* similarity, no values.
func paramKeySimilarity(nodesA map[string]*Node, namesA []string, nodesB map[string]*Node, namesB []string) float64 {
	if len(namesA) == 0 || len(namesA) != len(namesB) {
		return 0
	}
	var total float64
	for i := range namesA {
		total += jaccardKeys(nodesA[namesA[i]], nodesB[namesB[i]])
	}
	return total / float64(len(namesA))
}

func jaccardKeys(a, b *Node) float64 {
	if a == nil || b == nil {
		return 0
	}
	if len(a.Parameters) == 0 && len(b.Parameters) == 0 {
		return 1
	}
	union := map[string]bool{}
	var intersection int
	for k := range a.Parameters {
		union[k] = true
	}
	for k := range b.Parameters {
		if union[k] {
			intersection++
		}
		union[k] = true
	}
	if len(union) == 0 {
		return 1
	}
	return float64(intersection) / float64(len(union))
}
