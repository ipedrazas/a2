package godogs

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

// TestFeatures runs the godog test suite
func TestFeatures(t *testing.T) {
	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
		Paths:  []string{"features"},
	}

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		t.Fatalf("failed to run feature suite: %d", status)
	}
}

// InitializeScenario registers all step definitions
func InitializeScenario(ctx *godog.ScenarioContext) {
	// Installation and setup steps
	ctx.Step(`^I have Go installed on my system$`, iHaveGoInstalled)
	ctx.Step(`^I have an existing Go project directory$`, iHaveExistingProject)
	ctx.Step(`^I install A2 using "([^"]*)"$`, iInstallA2)
	ctx.Step(`^I verify the installation with "([^"]*)"$`, iVerifyInstallation)
	ctx.Step(`^I navigate to my project directory$`, iNavigateToProject)

	// Running checks
	ctx.Step(`^I run "([^"]*)"$`, iRunCommand)
	ctx.Step(`^I run "([^"]*)" in directory "([^"]*)"$`, iRunCommandInDirectory)
	ctx.Step(`^I run "([^"]*)" with output format "([^"]*)"$`, iRunWithOutputFormat)

	// A2 behavior
	ctx.Step(`^A2 should auto-detect the project language$`, a2ShouldAutoDetectLanguage)
	ctx.Step(`^A2 should run appropriate checks in parallel$`, a2ShouldRunChecksInParallel)
	ctx.Step(`^A2 should display results with color coding$`, a2ShouldDisplayResultsWithColor)
	ctx.Step(`^I should see a maturity score$`, iShouldSeeMaturityScore)
	ctx.Step(`^I should receive clear suggestions for improvement$`, iShouldReceiveSuggestions)

	// Results interpretation
	ctx.Step(`^I view the output$`, iViewOutput)
	ctx.Step(`^I should see green checkmarks for passed checks$`, iShouldSeeGreenCheckmarks)
	ctx.Step(`^I should see red X marks for critical failures$`, iShouldSeeRedXMarks)
	ctx.Step(`^I should see yellow warnings for recommendations$`, iShouldSeeYellowWarnings)
	ctx.Step(`^I should see blue info for optional tools not installed$`, iShouldSeeBlueInfo)

	// Fixing issues
	ctx.Step(`^A2 detected issues in my code$`, a2DetectedIssues)
	ctx.Step(`^I received actionable suggestions$`, iReceivedSuggestions)
	ctx.Step(`^I fix the identified issues$`, iFixIssues)
	ctx.Step(`^my maturity score should improve$`, maturityScoreShouldImprove)
	ctx.Step(`^I should be able to commit with confidence$`, iShouldCommitWithConfidence)

	// Missing tools
	ctx.Step(`^some required tools are not installed$`, someToolsNotInstalled)
	ctx.Step(`^A2 should display Info status for missing tools$`, a2ShouldDisplayInfoStatus)
	ctx.Step(`^A2 should continue running available checks$`, a2ShouldContinueRunning)
	ctx.Step(`^A2 should suggest how to install missing tools$`, a2ShouldSuggestToolInstallation)
	ctx.Step(`^the check should complete successfully$`, checkShouldComplete)

	// Configuration
	ctx.Step(`^I have A2 installed$`, iHaveA2Installed)
	ctx.Step(`^I have a new API project$`, iHaveNewAPIProject)
	ctx.Step(`^I run "([^"]*)" in interactive mode$`, iRunInInteractiveMode)
	ctx.Step(`^I select "([^"]*)" as the application type$`, iSelectApplicationType)
	ctx.Step(`^I select "([^"]*)" as the maturity level$`, iSelectMaturityLevel)
	ctx.Step(`^I select "([^"]*)" for language detection$`, iSelectLanguageDetection)
	ctx.Step(`^I choose not to enable external checks$`, iChooseNoExternalChecks)

	// Configuration file
	ctx.Step(`^A2 should create "([^"]*)" with sensible defaults$`, a2CreatesConfig)
	ctx.Step(`^the configuration should include API profile$`, configIncludesAPIProfile)
	ctx.Step(`^the configuration should include production target$`, configIncludesProductionTarget)
	ctx.Step(`^E2E tests should be disabled in the configuration$`, e2eTestsDisabled)

	// Customization
	ctx.Step(`^I have created an initial "([^"]*)" configuration$`, iHaveInitialConfig)
	ctx.Step(`^my team maintains high standards$`, teamMaintainsHighStandards)
	ctx.Step(`^I edit the configuration file$`, iEditConfig)
	ctx.Step(`^I set ([^"]+)\.([^"]*) to (\d+)$`, iSetConfigValue)
	ctx.Step(`^I save the file$`, iSaveFile)
	ctx.Step(`^A2 should use the stricter thresholds$`, a2UsesStricterThresholds)
	ctx.Step(`^checks should enforce the new standards$`, checksEnforceNewStandards)

	// Required files
	ctx.Step(`^I have a basic "([^"]*)" configuration$`, iHaveBasicConfig)
	ctx.Step(`^I add required documentation files to the configuration$`, iAddRequiredDocs)
	ctx.Step(`^I include "([^"]*)"$`, iIncludeRequiredFile)
	ctx.Step(`^A2 should verify these files exist$`, a2VerifiesFilesExist)
	ctx.Step(`^A2 should fail if required files are missing$`, a2FailsOnMissingFiles)

	// Disabling checks
	ctx.Step(`^I have a Go-only project$`, iHaveGoOnlyProject)
	ctx.Step(`^I disable "([^"]*)" checks$`, iDisableChecks)
	ctx.Step(`^A2 should skip these checks$`, a2SkipsChecks)
	ctx.Step(`^the check execution should be faster$`, checksFaster)
	ctx.Step(`^the results should only show relevant checks$`, resultsShowOnlyRelevant)

	// Team adoption
	ctx.Step(`^I have configured A2 for my project$`, iHaveConfiguredA2)
	ctx.Step(`^I commit the configuration file$`, iCommitConfig)
	ctx.Step(`^I push to the repository$`, iPushToRepo)
	ctx.Step(`^I communicate the setup to my team$`, iCommunicateToTeam)
	ctx.Step(`^team members should be able to run "([^"]*)"$`, teamCanRun)
	ctx.Step(`^everyone should see the same quality standards$`, everyoneSeesSameStandards)
	ctx.Step(`^the configuration should be version controlled$`, configVersionControlled)

	// AI-assisted development
	ctx.Step(`^I use an AI assistant to generate code$`, iUseAIGenerateCode)
	ctx.Step(`^I have not reviewed the code yet$`, iHaveNotReviewedCode)
	ctx.Step(`^A2 should detect any build failures$`, a2DetectsBuildFailures)
	ctx.Step(`^A2 should identify missing tests$`, a2IdentifiesMissingTests)
	ctx.Step(`^A2 should flag formatting issues$`, a2FlagsFormatIssues)
	ctx.Step(`^A2 should check for security vulnerabilities$`, a2ChecksSecurity)
	ctx.Step(`^I should receive actionable feedback$`, iReceiveActionableFeedback)

	// Iterative improvement
	ctx.Step(`^I received specific failure messages$`, iReceivedFailureMessages)
	ctx.Step(`^I fix the build issues$`, iFixBuildIssues)
	ctx.Step(`^I should see progress indicators$`, iSeeProgressIndicators)
	ctx.Step(`^the number of failures should decrease$`, failuresShouldDecrease)
	ctx.Step(`^I should continue fixing until all critical issues pass$`, iContinueFixing)

	// Quick validation
	ctx.Step(`^I have been working on a feature for (\d+) hours$`, iWorkOnFeature)
	ctx.Step(`^I want to push my changes$`, iWantToPush)
	ctx.Step(`^I should receive a minimal token format output$`, iReceiveTokenFormat)
	ctx.Step(`^I should see the status of all checks$`, iSeeAllCheckStatus)
	ctx.Step(`^I should identify any remaining issues$`, iIdentifyRemainingIssues)
	ctx.Step(`^I address the warnings$`, iAddressWarnings)
	ctx.Step(`^I run "([^"]*)" one final time$`, iRunFinalTime)
	ctx.Step(`^all checks should pass$`, allChecksPass)
	ctx.Step(`^my maturity score should be (\d+)%$`, maturityScoreIs)
	ctx.Step(`^I can push with confidence$`, iCanPush)

	// Getting help
	ctx.Step(`^A2 reported a failing test check$`, a2ReportedFailingTest)
	ctx.Step(`^I don't know what's required$`, iDontKnowRequirements)
	ctx.Step(`^I should see the check description$`, iSeeCheckDescription)
	ctx.Step(`^I should see the required tool command$`, iSeeToolCommand)
	ctx.Step(`^I should see the requirements for passing$`, iSeeRequirements)
	ctx.Step(`^I should receive specific suggestions for fixing$`, iReceiveFixSuggestions)

	// Incremental development
	ctx.Step(`^I am implementing a new feature step by step$`, iImplementingStepByStep)
	ctx.Step(`^I run "([^"]*)" after each small change$`, iRunAfterEachChange)
	ctx.Step(`^I should receive fast feedback$`, iReceiveFastFeedback)
	ctx.Step(`^I should catch issues early$`, iCatchIssuesEarly)
	ctx.Step(`^the feedback loop should be less than (\d+) minutes$`, feedbackLoopLessThan)
	ctx.Step(`^my development velocity should be maintained$`, velocityMaintained)

	// Before/After hooks
	ctx.BeforeScenario(func(*godog.Scenario) {
		// Reset state before each scenario
		ResetState()
	})

	ctx.AfterScenario(func(*godog.Scenario, error) {
		// Cleanup after each scenario
		cleanup()
	})
}

func cleanup() {
	// Clean up temporary files, directories, etc.
}
