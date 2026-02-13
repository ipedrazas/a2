package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/maturity"
	"github.com/ipedrazas/a2/pkg/output"
	"github.com/ipedrazas/a2/pkg/profiles"
	"github.com/ipedrazas/a2/pkg/runner"
	"github.com/ipedrazas/a2/pkg/targets"
)

// ProcessJob processes a single check job.
// This is the main worker function that clones the repo and runs checks.
func ProcessJob(ctx context.Context, job *Job, wm *WorkspaceManager) error {
	// Update workspace in job (create new workspace for this job)
	workspaceDir, err := wm.CreateWorkspace(job.ID)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	job.WorkspaceDir = workspaceDir

	// Parse GitHub URL
	repo, err := ParseGitHubURL(job.GitHubURL)
	if err != nil {
		return fmt.Errorf("invalid GitHub URL: %w", err)
	}

	// Clone repository
	if err := CloneRepository(repo, workspaceDir); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Clean up workspace after job completes
	if wm.cleanupAfter {
		defer func() {
			if err := wm.Cleanup(workspaceDir); err != nil {
				// Log but don't fail the job if cleanup fails
				fmt.Fprintf(os.Stderr, "Warning: failed to cleanup workspace %s: %v\n", workspaceDir, err)
			}
		}()
	}

	// Load configuration
	cfg, err := config.Load(workspaceDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Apply target if specified (maturity level)
	if job.Request.Target != "" {
		t, ok := targets.Get(job.Request.Target)
		if !ok {
			return fmt.Errorf("unknown target: %s", job.Request.Target)
		}
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, t.Disabled...)
	}

	// Apply profile if specified (application type)
	if job.Request.Profile != "" {
		p, ok := profiles.Get(job.Request.Profile)
		if !ok {
			return fmt.Errorf("unknown profile: %s", job.Request.Profile)
		}
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, p.Disabled...)
	}

	// Apply skip checks from request
	if len(job.Request.SkipChecks) > 0 {
		cfg.Checks.Disabled = append(cfg.Checks.Disabled, job.Request.SkipChecks...)
	}

	// Detect or use explicit languages
	var detected language.DetectionResult
	if len(job.Request.Languages) > 0 {
		// Convert string flags to Language types
		langs := make([]checker.Language, len(job.Request.Languages))
		for i, l := range job.Request.Languages {
			langs[i] = checker.Language(l)
		}
		detected = language.DetectWithOverride(workspaceDir, langs)
	} else if len(cfg.Language.Explicit) > 0 {
		// Use explicit languages from config
		langs := make([]checker.Language, len(cfg.Language.Explicit))
		for i, l := range cfg.Language.Explicit {
			langs[i] = checker.Language(l)
		}
		detected = language.DetectWithOverride(workspaceDir, langs)
	} else {
		// Auto-detect languages
		detected = language.DetectWithSourceDirs(workspaceDir, cfg.GetSourceDirs())
	}

	// Check if any language was detected
	if len(detected.Languages) == 0 {
		return fmt.Errorf("no supported language detected")
	}

	// Get the list of checks to run
	registrations := checks.GetChecks(cfg, detected)

	// Set timeout
	var timeout time.Duration
	if job.Request.TimeoutSecs > 0 {
		timeout = time.Duration(job.Request.TimeoutSecs) * time.Second
	}

	// Run the suite with configured execution options
	opts := runner.RunSuiteOptions{
		Parallel:   cfg.Execution.Parallel,
		Timeout:    timeout,
		OnProgress: nil, // No progress callback for server mode
	}

	result := runner.RunSuiteWithOptions(workspaceDir, registrations, opts)

	// Determine verbosity level based on request
	var verbosity output.VerbosityLevel
	if job.Request.Verbose {
		verbosity = output.VerbosityFailures // Show output for failures and warnings
	} else {
		verbosity = output.VerbosityLevel(0) // No raw output
	}

	// Convert result to JSON format
	jsonOutput := toJSONOutput(result, detected, verbosity)

	// Store result in job (this updates the in-memory job struct)
	job.Result = jsonOutput

	// Note: The job status will be marked as completed by the queue when this returns nil
	return nil
}

// toJSONOutput converts a SuiteResult to JSONOutput format.
// This is extracted from output/json.go to avoid writing to stdout.
func toJSONOutput(result runner.SuiteResult, detected language.DetectionResult, verbosity output.VerbosityLevel) *output.JSONOutput {
	// Convert languages to strings
	langs := make([]string, len(detected.Languages))
	for i, l := range detected.Languages {
		langs[i] = string(l)
	}

	// Calculate maturity estimation
	est := maturity.Estimate(result)

	jsonOutput := &output.JSONOutput{
		Languages: langs,
		Results:   make([]output.JSONResult, 0, len(result.Results)),
		Summary: output.JSONSummary{
			Total:           result.ScoredChecks(), // Excludes Info from total
			Passed:          result.Passed,
			Warnings:        result.Warnings,
			Failed:          result.Failed,
			Info:            result.Info,
			Score:           calculateScore(result),
			TotalDurationMs: result.TotalDuration.Milliseconds(),
		},
		Maturity: output.JSONMaturity{
			Level:       est.Level.String(),
			Description: est.Level.Description(),
			Suggestions: est.Suggestions,
		},
		Aborted: result.Aborted,
		Success: result.Success(),
	}

	for _, r := range result.Results {
		jsonResult := output.JSONResult{
			Name:       r.Name,
			ID:         r.ID,
			Passed:     r.Passed,
			Status:     statusToString(r.Status),
			Message:    r.Message,
			Reason:     r.Reason,
			Language:   string(r.Language),
			DurationMs: r.Duration.Milliseconds(),
		}

		// Include raw output based on verbosity level
		if r.RawOutput != "" {
			shouldInclude := verbosity == output.VerbosityAll ||
				(verbosity == output.VerbosityFailures && (r.Status == checker.Fail || r.Status == checker.Warn))
			if shouldInclude {
				jsonResult.RawOutput = r.RawOutput
			}
		}

		jsonOutput.Results = append(jsonOutput.Results, jsonResult)
	}

	return jsonOutput
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
