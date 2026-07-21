package n8n

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

// SecretFinding locates secret material in a workflow WITHOUT carrying it:
// the matched value is redacted at detection time and never leaves the
// scanner — findings hold workflow, node, and field path only.
type SecretFinding struct {
	File      string // set when scanning a bundle
	Workflow  string
	Node      string // empty for pinData findings ("pinData.<node>" is in FieldPath)
	FieldPath string // e.g. "parameters.headerParameters.parameters[0].value"
	Rule      string // which detector matched, e.g. "openai-style key (sk-...)"
	Confirmed bool   // true: known key format; false: high-entropy candidate
}

// secretPatterns are known key formats. A match is a confirmed finding.
var secretPatterns = []struct {
	rule    string
	pattern *regexp.Regexp
}{
	{"openai-style key (sk-...)", regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{20,}`)},
	{"github token (gh?_...)", regexp.MustCompile(`\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36}\b|\bgithub_pat_[A-Za-z0-9_]{22,}`)},
	{"slack token (xox...)", regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{10,}`)},
	{"aws access key id (AKIA/ASIA...)", regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`)},
	{"google api key (AIza...)", regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}`)},
	{"jwt", regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{5,}`)},
	{"private key block", regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
}

// authFieldPattern marks parameter names/paths that typically carry
// credentials; a high-entropy literal in one of these is a candidate finding.
var authFieldPattern = regexp.MustCompile(`(?i)(auth|token|secret|api[-_]?key|access[-_]?key|password|bearer|credential)`)

// Entropy-candidate thresholds: shorter or lower-entropy strings are too
// noisy to report even as candidates.
const (
	entropyMinLength = 16
	entropyThreshold = 3.5 // bits per character
)

// authModeFields hold n8n's auth *method* enums ("headerAuth",
// "genericCredentialType", ...), not secret material — excluded from
// entropy candidates.
var authModeFields = map[string]bool{
	"authentication":     true,
	"genericAuthType":    true,
	"nodeCredentialType": true,
	"authType":           true,
}

// digitOrSymbol distinguishes key material from prose/enums: real keys
// virtually always mix in digits or symbols.
var digitOrSymbol = regexp.MustCompile(`[0-9_\-+/=]`)

// ScanBundleSecrets scans every workflow in the bundle (including degraded
// ones — whatever parsed is scannable) and returns findings sorted by
// file/node/field path.
func ScanBundleSecrets(b *Bundle) []SecretFinding {
	var findings []SecretFinding
	for _, wf := range b.Workflows {
		for _, f := range ScanWorkflowSecrets(wf.Workflow) {
			f.File = wf.Path
			findings = append(findings, f)
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		a, b := findings[i], findings[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Node != b.Node {
			return a.Node < b.Node
		}
		return a.FieldPath < b.FieldPath
	})
	return findings
}

// ScanWorkflowSecrets scans node parameters (including Code node source,
// which is a parameter) and pinned data of one workflow.
func ScanWorkflowSecrets(wf *Workflow) []SecretFinding {
	var findings []SecretFinding
	emit := func(node, path, label, value string) {
		findings = append(findings, scanValue(wf.Name, node, path, label, value)...)
	}
	for i := range wf.Nodes {
		node := &wf.Nodes[i]
		walkStrings("parameters", "", node.Parameters, func(path, label, value string) {
			emit(node.Name, path, label, value)
		})
	}
	nodeNames := make([]string, 0, len(wf.PinData))
	for name := range wf.PinData {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)
	for _, name := range nodeNames {
		var pinned any
		if err := json.Unmarshal(wf.PinData[name], &pinned); err != nil {
			continue
		}
		walkStrings("pinData."+name, "", pinned, func(path, label, value string) {
			emit("", path, label, value)
		})
	}
	return findings
}

// scanValue runs the detectors on one string value. The value itself is
// dropped here: only rule name and location survive into findings.
func scanValue(workflow, node, path, label, value string) []SecretFinding {
	var findings []SecretFinding
	for _, sp := range secretPatterns {
		if sp.pattern.MatchString(value) {
			findings = append(findings, SecretFinding{
				Workflow:  workflow,
				Node:      node,
				FieldPath: path,
				Rule:      sp.rule,
				Confirmed: true,
			})
		}
	}
	if len(findings) > 0 {
		return findings
	}
	// Entropy candidates: literal (not an n8n expression) values in
	// auth-adjacent fields that look random enough to be real key material.
	if strings.HasPrefix(value, "=") {
		return nil
	}
	if authModeFields[label] {
		return nil
	}
	if !authFieldPattern.MatchString(label) && !authFieldPattern.MatchString(path) {
		return nil
	}
	if len(value) >= entropyMinLength && shannonEntropy(value) >= entropyThreshold && digitOrSymbol.MatchString(value) {
		findings = append(findings, SecretFinding{
			Workflow:  workflow,
			Node:      node,
			FieldPath: path,
			Rule:      "high-entropy value in auth field",
			Confirmed: false,
		})
	}
	return findings
}

// walkStrings visits every string in a decoded JSON structure with its
// path. For n8n's {name, value} parameter pairs, the sibling "name" becomes
// the value's label so auth-adjacency works on what the parameter is called,
// not just where it sits. Map keys are visited in sorted order for
// deterministic output.
func walkStrings(path, label string, v any, fn func(path, label, value string)) {
	switch t := v.(type) {
	case string:
		fn(path, label, t)
	case map[string]any:
		pairName, _ := t["name"].(string)
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			childLabel := k
			if k == "value" && pairName != "" {
				childLabel = pairName
			}
			walkStrings(path+"."+k, childLabel, t[k], fn)
		}
	case []any:
		for i, e := range t {
			walkStrings(fmt.Sprintf("%s[%d]", path, i), label, e, fn)
		}
	}
}

// shannonEntropy returns the entropy of s in bits per character.
func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}
	freq := map[rune]int{}
	var total int
	for _, r := range s {
		freq[r]++
		total++
	}
	var entropy float64
	for _, count := range freq {
		p := float64(count) / float64(total)
		entropy -= p * math.Log2(p)
	}
	return entropy
}
