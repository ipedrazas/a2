Implementation Plan: Three New CLI Commands

 Overview

 Add three new commands to the a2 CLI:
 1. a2 run CHECK_ID - Run a specific check with full output
 2. a2 explain CHECK_ID - Show detailed explanation of a check
 3. a2 profiles validate / a2 targets validate - Validate user definitions

 ---
 Command 1: a2 run CHECK_ID [path]

 Run a specific check and display the full stdout/stderr from the underlying tool.

 Files to Modify

 pkg/checker/types.go - Add output field to Result:
 type Result struct {
     // ... existing fields ...
     RawOutput string // Full command output for verbose display
 }

 pkg/checkutil/util.go - Add method to ResultBuilder for setting raw output:
 func (b *ResultBuilder) WithOutput(r checker.Result, output string) checker.Result {
     r.RawOutput = output
     return r
 }

 Check files (update progressively to populate RawOutput):
 - pkg/checks/go/race.go
 - pkg/checks/go/build.go
 - pkg/checks/go/tests.go
 - Other checks as needed

 Files to Create

 cmd/run.go - New command:
 var runCmd = &cobra.Command{
     Use:   "run CHECK_ID [path]",
     Short: "Run a specific check with full output",
     Args:  cobra.RangeArgs(1, 2),
     RunE:  runSingleCheck,
 }

 Implementation:
 1. Parse check ID from args
 2. Load check registrations via checks.GetAllCheckRegistrations(cfg)
 3. Find matching check by ID
 4. Execute check using runner
 5. Display: status, message, duration, and full RawOutput

 Command Output Example

 $ a2 run go:race

 ✗ WARN Go Race Detection (2.4s)
     Race condition detected

 --- Output ---
 ==================
 WARNING: DATA RACE
 Read at 0x00c000120000 by goroutine 8:
   main.worker()
       /app/main.go:45 +0x3c
 ...
 ==================

 ---
 Command 2: a2 explain CHECK_ID and a2 list checks --explain

 Show detailed explanation of what a check does.

 Files to Modify

 pkg/checker/types.go - Add Description field:
 type CheckMeta struct {
     ID          string
     Name        string
     Description string  // NEW: Detailed explanation
     Languages   []Language
     Critical    bool
     Order       int
     Suggestion  string
 }

 All register.go files - Add descriptions to each check:
 - pkg/checks/go/register.go (10 checks)
 - pkg/checks/python/register.go
 - pkg/checks/node/register.go
 - pkg/checks/java/register.go
 - pkg/checks/rust/register.go
 - pkg/checks/typescript/register.go
 - pkg/checks/swift/register.go
 - pkg/checks/common/register.go (22 checks)

 Example:
 {
     Checker: &RaceCheck{},
     Meta: checker.CheckMeta{
         ID:          "go:race",
         Name:        "Go Race Detection",
         Description: "Runs tests with -race flag to detect data races. Data races occur when two goroutines access the same variable concurrently and at least one is a write.",
         Languages:   []checker.Language{checker.LangGo},
         Critical:    false,
         Order:       125,
         Suggestion:  "Fix race conditions detected by -race flag",
     },
 },

 cmd/list.go - Add --explain flag:
 var listExplain bool

 func init() {
     listChecksCmd.Flags().BoolVar(&listExplain, "explain", false,
         "Show detailed descriptions for each check")
 }

 Files to Create

 cmd/explain.go - New command:
 var explainCmd = &cobra.Command{
     Use:   "explain CHECK_ID",
     Short: "Show detailed explanation of a check",
     Args:  cobra.ExactArgs(1),
     Run:   runExplain,
 }

 Command Output Example

 $ a2 explain go:race

 Check ID:     go:race
 Name:         Go Race Detection
 Description:  Runs tests with -race flag to detect data races. Data races
               occur when two goroutines access the same variable concurrently
               and at least one is a write.
 Languages:    go
 Critical:     No
 Suggestion:   Fix race conditions detected by -race flag

 ---
 Command 3: a2 profiles validate and a2 targets validate

 Validate user-defined profiles and targets with comprehensive checks.

 Files to Create

 pkg/validation/validation.go - Shared validation types:
 package validation

 type ValidationResult struct {
     Valid    bool
     Errors   []string
     Warnings []string
 }

 func LevenshteinDistance(a, b string) int { ... }
 func FindSimilar(id string, validIDs []string, maxDistance int) []string { ... }

 pkg/profiles/validate.go - Profile validation:
 func ValidateProfile(p Profile, validCheckIDs map[string]bool) ValidationResult {
     // 1. Check for duplicate disabled check IDs
     // 2. Check if each disabled ID exists in registry
     // 3. Suggest similar IDs for typos (Levenshtein distance)
     // 4. Warn if overriding built-in profile
 }

 func ValidateAllUserProfiles() map[string]ValidationResult { ... }

 pkg/targets/validate.go - Target validation (same pattern).

 Files to Modify

 cmd/root.go - Add validate subcommands:
 var profilesValidateCmd = &cobra.Command{
     Use:   "validate",
     Short: "Validate user-defined profiles",
     Run:   runProfilesValidate,
 }

 var targetsValidateCmd = &cobra.Command{
     Use:   "validate",
     Short: "Validate user-defined targets",
     Run:   runTargetsValidate,
 }

 func init() {
     profilesCmd.AddCommand(profilesValidateCmd)
     targetsCmd.AddCommand(targetsValidateCmd)
 }

 Validation Checks

 1. YAML structure - Valid syntax, required fields
 2. Check ID existence - Each disabled ID must exist in registry
 3. Duplicates - Warn about duplicate disabled check IDs
 4. Typos - Suggest similar check IDs using Levenshtein distance
 5. Override warning - Warn if user profile/target overrides built-in

 Command Output Example

 $ a2 profiles validate

 Profile: my-api
   Status: INVALID
   ERROR: Unknown check ID: common:helth
   WARNING: Did you mean: common:health
   WARNING: Duplicate disabled check: common:secrets

 Profile: my-cli
   Status: VALID
   WARNING: Overrides built-in profile: cli

 ---
 Implementation Order

 Phase 1: explain command (lowest risk)

 1. Add Description field to CheckMeta in types.go
 2. Create cmd/explain.go
 3. Add --explain flag to cmd/list.go
 4. Add descriptions to all register.go files

 Phase 2: validate commands

 1. Create pkg/validation/validation.go
 2. Create pkg/profiles/validate.go
 3. Create pkg/targets/validate.go
 4. Add validate subcommands to cmd/root.go

 Phase 3: run command

 1. Add RawOutput field to Result in types.go
 2. Create cmd/run.go
 3. Update key checks to populate RawOutput (go:race, go:build, go:tests first)

 ---
 Files Summary

 New Files (7)

 - cmd/run.go
 - cmd/explain.go
 - pkg/validation/validation.go
 - pkg/profiles/validate.go
 - pkg/targets/validate.go

 Modified Files (12)

 - pkg/checker/types.go (add Description, RawOutput)
 - pkg/checkutil/util.go (add WithOutput helper)
 - cmd/root.go (register new commands)
 - cmd/list.go (add --explain flag)
 - pkg/checks/go/register.go (add descriptions)
 - pkg/checks/python/register.go (add descriptions)
 - pkg/checks/node/register.go (add descriptions)
 - pkg/checks/java/register.go (add descriptions)
 - pkg/checks/rust/register.go (add descriptions)
 - pkg/checks/typescript/register.go (add descriptions)
 - pkg/checks/swift/register.go (add descriptions)
 - pkg/checks/common/register.go (add descriptions)

 ---
 Verification

 Test the run command

 a2 run go:race
 a2 run go:build ./path/to/project
 a2 run common:dockerfile --format json

 Test the explain command

 a2 explain go:race
 a2 explain common:health
 a2 list checks --explain

 Test the validate commands

 # Create test profile with typo
 echo "disabled: [common:helth]" > ~/.config/a2/profiles/test.yaml
 a2 profiles validate

 # Create valid profile
 echo "disabled: [common:health]" > ~/.config/a2/profiles/test.yaml
 a2 profiles validate

 a2 targets validate

 Run existing tests

 task test