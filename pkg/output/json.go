package output

import (
	"encoding/json"
	"os"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/maturity"
	"github.com/ipedrazas/a2/pkg/runner"
)

// JSONResult is the JSON-friendly version of a check result.
type JSONResult struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Passed     bool   `json:"passed"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Language   string `json:"language,omitempty"`
	DurationMs int64  `json:"duration_ms"` // Duration in milliseconds
}

// JSONOutput is the complete JSON output structure.
type JSONOutput struct {
	Languages []string     `json:"languages"`
	Results   []JSONResult `json:"results"`
	Summary   JSONSummary  `json:"summary"`
	Maturity  JSONMaturity `json:"maturity"`
	Aborted   bool         `json:"aborted"`
	Success   bool         `json:"success"`
}

// JSONSummary provides aggregate statistics.
type JSONSummary struct {
	Total           int     `json:"total"`
	Passed          int     `json:"passed"`
	Warnings        int     `json:"warnings"`
	Failed          int     `json:"failed"`
	Info            int     `json:"info"`
	Score           float64 `json:"score"`
	TotalDurationMs int64   `json:"total_duration_ms"` // Total duration in milliseconds
}

// JSONMaturity provides maturity assessment.
type JSONMaturity struct {
	Level       string   `json:"level"`
	Description string   `json:"description"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// JSON outputs the results as formatted JSON.
func JSON(result runner.SuiteResult, detected language.DetectionResult) error {
	// Convert languages to strings
	langs := make([]string, len(detected.Languages))
	for i, l := range detected.Languages {
		langs[i] = string(l)
	}

	// Calculate maturity estimation
	est := maturity.Estimate(result)

	output := JSONOutput{
		Languages: langs,
		Results:   make([]JSONResult, 0, len(result.Results)),
		Summary: JSONSummary{
			Total:           result.ScoredChecks(), // Excludes Info from total
			Passed:          result.Passed,
			Warnings:        result.Warnings,
			Failed:          result.Failed,
			Info:            result.Info,
			Score:           calculateScore(result),
			TotalDurationMs: result.TotalDuration.Milliseconds(),
		},
		Maturity: JSONMaturity{
			Level:       est.Level.String(),
			Description: est.Level.Description(),
			Suggestions: est.Suggestions,
		},
		Aborted: result.Aborted,
		Success: result.Success(),
	}

	for _, r := range result.Results {
		output.Results = append(output.Results, JSONResult{
			Name:       r.Name,
			ID:         r.ID,
			Passed:     r.Passed,
			Status:     statusToString(r.Status),
			Message:    r.Message,
			Language:   string(r.Language),
			DurationMs: r.Duration.Milliseconds(),
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		return err
	}

	if !result.Success() {
		os.Exit(1)
	}
	return nil
}

func statusToString(s checker.Status) string {
	switch s {
	case checker.Pass:
		return "pass"
	case checker.Warn:
		return "warn"
	case checker.Fail:
		return "fail"
	case checker.Info:
		return "info"
	default:
		return "unknown"
	}
}

func calculateScore(result runner.SuiteResult) float64 {
	// Use ScoredChecks to exclude Info from score calculation
	scoredTotal := result.ScoredChecks()
	if scoredTotal == 0 {
		return 100.0
	}
	return float64(result.Passed) / float64(scoredTotal) * 100
}
