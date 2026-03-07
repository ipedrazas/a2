package profiles

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checks"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/validation"
)

// ValidateProfile validates a single profile against the check registry.
func ValidateProfile(p Profile, validCheckIDs map[string]bool, validIDList []string) validation.ValidationResult {
	result := validation.ValidationResult{
		Name:  p.Name,
		Valid: true,
	}

	// Track seen check IDs to detect duplicates
	seen := make(map[string]bool)

	for _, checkID := range p.Disabled {
		// Check for duplicates
		if seen[checkID] {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("duplicate disabled check: %s", checkID))
		}
		seen[checkID] = true

		// Skip wildcard patterns — they are valid pattern syntax, not literal check IDs
		if strings.Contains(checkID, "*") {
			continue
		}

		// Check if the check ID exists
		if !validCheckIDs[checkID] {
			result.Errors = append(result.Errors,
				fmt.Sprintf("unknown check ID: %s", checkID))
			result.Valid = false

			// Suggest similar IDs for typos
			similar := validation.FindSimilar(checkID, validIDList, 3)
			if len(similar) > 0 {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("did you mean: %v", similar))
			}
		}
	}

	// Warn if overriding a built-in profile and detect stale copies
	if builtIn, isBuiltIn := BuiltInProfiles[p.Name]; isBuiltIn && p.Source == SourceUser {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("overrides built-in profile: %s", p.Name))

		// Check for checks present in built-in but missing from user copy
		userSet := make(map[string]bool)
		for _, d := range p.Disabled {
			userSet[d] = true
		}
		var missing []string
		for _, d := range builtIn.Disabled {
			if !userSet[d] {
				missing = append(missing, d)
			}
		}
		if len(missing) > 0 {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("missing %d check(s) from latest built-in: %v (run 'a2 profiles init --force' to update)", len(missing), missing))
		}
	}

	return result
}

// ValidateAllUserProfiles validates all user-defined profiles.
// Returns a map of profile name to validation result.
func ValidateAllUserProfiles() map[string]validation.ValidationResult {
	results := make(map[string]validation.ValidationResult)

	// Get all valid check IDs from the registry
	cfg := config.DefaultConfig()
	allRegs := checks.GetAllCheckRegistrations(cfg)

	validCheckIDs := make(map[string]bool)
	var validIDList []string
	for _, reg := range allRegs {
		validCheckIDs[reg.Meta.ID] = true
		validIDList = append(validIDList, reg.Meta.ID)
	}

	// Load user profiles (reload to get fresh data)
	userProfiles, err := LoadUserProfiles()
	if err != nil {
		// Return a single error result
		results["_error"] = validation.ValidationResult{
			Name:   "_error",
			Valid:  false,
			Errors: []string{fmt.Sprintf("failed to load user profiles: %v", err)},
		}
		return results
	}

	// No user profiles to validate
	if len(userProfiles) == 0 {
		return results
	}

	// Validate each user profile
	for name, profile := range userProfiles {
		results[name] = ValidateProfile(profile, validCheckIDs, validIDList)
	}

	return results
}

// GetValidCheckIDs returns all valid check IDs from the registry.
func GetValidCheckIDs() (map[string]bool, []string) {
	cfg := config.DefaultConfig()
	allRegs := checks.GetAllCheckRegistrations(cfg)

	validCheckIDs := make(map[string]bool)
	var validIDList []string
	for _, reg := range allRegs {
		validCheckIDs[reg.Meta.ID] = true
		validIDList = append(validIDList, reg.Meta.ID)
	}

	return validCheckIDs, validIDList
}
