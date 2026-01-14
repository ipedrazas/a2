// Package validation provides utilities for validating profiles and targets.
package validation

// ValidationResult holds the outcome of a validation.
type ValidationResult struct {
	Name     string   // Name of the profile/target being validated
	Valid    bool     // Whether validation passed
	Errors   []string // Critical errors that make the configuration invalid
	Warnings []string // Non-critical warnings about potential issues
}

// LevenshteinDistance calculates the edit distance between two strings.
// Used to suggest similar check IDs when a typo is detected.
func LevenshteinDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Create a matrix to store distances
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := range dp[0] {
		dp[0][j] = j
	}

	// Fill in the matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,      // deletion
				dp[i][j-1]+1,      // insertion
				dp[i-1][j-1]+cost, // substitution
			)
		}
	}

	return dp[m][n]
}

// FindSimilar finds check IDs similar to the given one within maxDistance.
// Returns up to 3 similar IDs sorted by distance.
func FindSimilar(checkID string, validIDs []string, maxDistance int) []string {
	type candidate struct {
		id       string
		distance int
	}

	var candidates []candidate
	for _, id := range validIDs {
		dist := LevenshteinDistance(checkID, id)
		if dist > 0 && dist <= maxDistance {
			candidates = append(candidates, candidate{id: id, distance: dist})
		}
	}

	// Sort by distance
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].distance < candidates[i].distance {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	// Return up to 3 suggestions
	var result []string
	for i := 0; i < len(candidates) && i < 3; i++ {
		result = append(result, candidates[i].id)
	}

	return result
}
