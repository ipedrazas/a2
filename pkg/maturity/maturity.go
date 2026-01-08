// Package maturity provides project maturity estimation based on check results.
package maturity

import (
	"github.com/ipedrazas/a2/pkg/runner"
)

// Level represents the maturity level of a project.
type Level int

const (
	// PoC indicates early stage development, focus on core functionality.
	PoC Level = iota
	// Development indicates core functionality works but quality improvements needed.
	Development
	// Mature indicates most checks pass with minor improvements recommended.
	Mature
	// ProductionReady indicates all checks pass, ready for production deployment.
	ProductionReady
)

// String returns the human-readable name of the maturity level.
func (l Level) String() string {
	switch l {
	case ProductionReady:
		return "Production-Ready"
	case Mature:
		return "Mature"
	case Development:
		return "Development"
	default:
		return "Proof of Concept"
	}
}

// Description returns a detailed description of the maturity level.
func (l Level) Description() string {
	switch l {
	case ProductionReady:
		return "All checks pass, ready for production deployment"
	case Mature:
		return "Most checks pass, minor improvements recommended"
	case Development:
		return "Core functionality works, quality improvements needed"
	default:
		return "Early stage, focus on core functionality first"
	}
}

// Estimation contains the maturity assessment results.
type Estimation struct {
	Level       Level    // The assessed maturity level
	Score       float64  // Score from 0-100 (percentage of passed checks)
	Passed      int      // Number of passed checks
	Warnings    int      // Number of warnings
	Failed      int      // Number of failed checks
	Total       int      // Total number of checks run
	Suggestions []string // Recommendations for improvement
}

// Estimate analyzes check results and returns a maturity estimation.
func Estimate(result runner.SuiteResult) Estimation {
	total := result.TotalChecks()
	if total == 0 {
		return Estimation{Level: PoC, Score: 0}
	}

	score := float64(result.Passed) / float64(total) * 100

	est := Estimation{
		Score:    score,
		Passed:   result.Passed,
		Warnings: result.Warnings,
		Failed:   result.Failed,
		Total:    total,
	}

	// Determine level based on score and failure count
	switch {
	case result.Failed == 0 && result.Warnings == 0 && score == 100:
		est.Level = ProductionReady
	case result.Failed == 0 && score >= 80:
		est.Level = Mature
		if result.Warnings > 0 {
			est.Suggestions = append(est.Suggestions,
				"Address warnings to reach production-ready status")
		}
	case result.Failed <= 2 && score >= 60:
		est.Level = Development
		est.Suggestions = append(est.Suggestions,
			"Fix failing checks to improve maturity")
		if result.Warnings > 0 {
			est.Suggestions = append(est.Suggestions,
				"Review and address warnings")
		}
	default:
		est.Level = PoC
		est.Suggestions = append(est.Suggestions,
			"Focus on critical checks first (build, tests)")
		if result.Failed > 2 {
			est.Suggestions = append(est.Suggestions,
				"Many checks failing - prioritize fixing build and test failures")
		}
	}

	return est
}
