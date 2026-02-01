package godogs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	// A2 Detection and Analysis Steps
	ctx.Step(`^A(\d+) detected a failing test$`, aDetectedAFailingTest)
	ctx.Step(`^A(\d+) detected issues in a PR$`, aDetectedIssuesInAPR)
	ctx.Step(`^A(\d+) detected issues with my code$`, aDetectedIssuesWithMyCode)
	ctx.Step(`^A(\d+) detects multiple critical issues$`, aDetectsMultipleCriticalIssues)
	ctx.Step(`^A(\d+) detects security issues \(e\.g\. leaked API key\)$`, aDetectsSecurityIssuesEgLeakedAPIKey)
	ctx.Step(`^A(\d+) should analyze cyclomatic complexity$`, aShouldAnalyzeCyclomaticComplexity)
	ctx.Step(`^A(\d+) should auto-detect the programming language$`, aShouldAutodetectTheProgrammingLanguage)
	ctx.Step(`^A(\d+) should check test coverage$`, aShouldCheckTestCoverage)
	ctx.Step(`^A(\d+) should display full output from the tool$`, aShouldDisplayFullOutputFromTheTool)
	ctx.Step(`^A(\d+) should display results in the terminal$`, aShouldDisplayResultsInTheTerminal)
	ctx.Step(`^A(\d+) should run all applicable checks$`, aShouldRunAllApplicableChecks)
	ctx.Step(`^A(\d+) should run only Go checks$`, aShouldRunOnlyGoChecks)
	ctx.Step(`^A(\d+) should run only the race detector check$`, aShouldRunOnlyTheRaceDetectorCheck)
	ctx.Step(`^A(\d+) should validate the build status$`, aShouldValidateTheBuildStatus)

	// CI/CD Integration Steps
	ctx.Step(`^A(\d+) should comment on the PR with a quality report$`, aShouldCommentOnThePRWithAQualityReport)
	ctx.Step(`^A(\d+) is running in CI for all projects$`, aIsRunningInCIForAllProjects)
	ctx.Step(`^A(\d+) detects security issues \(e\.g\., leaked API key\)$`, aDetectsSecurityIssuesEgLeakedAPIKey)
	ctx.Step(`^GitHub Actions should automatically run "([^"]*)"$`, gitHubActionsShouldAutomaticallyRun)
	ctx.Step(`^I add a GitHub Action script$`, iAddAGitHubActionScript)
	ctx.Step(`^I configure it to fail on exit code (\d+) \((\w+)\)$`, iConfigureItToFailOnExitCodeFailures)
	ctx.Step(`^I configure it to parse results\.json$`, iConfigureItToParseResultsjson)
	ctx.Step(`^I configure it to run on pull requests and pushes$`, iConfigureItToRunOnPullRequestsAndPushes)
	ctx.Step(`^I enable branch protection rules$`, iEnableBranchProtectionRules)
	ctx.Step(`^I install A(\d+) in the workflow$`, iInstallAInTheWorkflow)
	ctx.Step(`^I set it to comment on pull requests$`, iSetItToCommentOnPullRequests)
	ctx.Step(`^I update the workflow to check exit codes$`, iUpdateTheWorkflowToCheckExitCodes)
	ctx.Step(`^I upload the results as artifacts$`, iUploadTheResultsAsArtifacts)
	ctx.Step(`^CI should use different configs per branch$`, cIShouldUseDifferentConfigsPerBranch)
	ctx.Step(`^the check status should appear on GitHub$`, theCheckStatusShouldAppearOnGitHub)
	ctx.Step(`^the workflow should run on every PR$`, theWorkflowShouldRunOnEveryPR)

	// Pull Request Review Steps
	ctx.Step(`^a contributor submits a new PR$`, aContributorSubmitsANewPR)
	ctx.Step(`^a first-time contributor submits a PR$`, aFirsttimeContributorSubmitsAPR)
	ctx.Step(`^I can categorize green PRs \((\d+)%\+\) to review first$`, iCanCategorizeGreenPRsToReviewFirst)
	ctx.Step(`^I can categorize red PRs \(<(\d+)%\) to request fixes$`, iCanCategorizeRedPRsToRequestFixes)
	ctx.Step(`^I can categorize yellow PRs \((\d+)-(\d+)%\) to review second$`, iCanCategorizeYellowPRsToReviewSecond)
	ctx.Step(`^I can focus my time on high-quality PRs$`, iCanFocusMyTimeOnHighqualityPRs)
	ctx.Step(`^I filter PRs by A(\d+) quality score \((\d+)\)$`, iFilterPRsByAQualityScore)
	ctx.Step(`^I have (\d+) PRs to review$`, iHavePRsToReview)
	ctx.Step(`^I receive (\d+)\+ PRs per week$`, iReceivePRsPerWeek)
	ctx.Step(`^I review high-quality PRs first$`, iReviewHighqualityPRsFirst)
	ctx.Step(`^I should save ([\d.]+) hours per week$`, iShouldSaveHoursPerWeek)
	ctx.Step(`^I should save (\d+)\+ minutes of manual review time$`, iShouldSaveMinutesOfManualReviewTime)
	ctx.Step(`^low-quality PRs should be caught automatically$`, lowqualityPRsShouldBeCaughtAutomatically)
	ctx.Step(`^my review time should be significantly reduced$`, myReviewTimeShouldBeSignificantlyReduced)
	ctx.Step(`^the PR should be ready for review$`, thePRShouldBeReadyForReview)

	// AI-Assisted Development Steps
	ctx.Step(`^I ask the AI to fix the specific issue$`, iAskTheAIToFixTheSpecificIssue)
	ctx.Step(`^I asked AI to refactor a large module$`, iAskedAIToRefactorALargeModule)
	ctx.Step(`^I am an AI agent processing A(\d+) results$`, iAmAnAIAgentProcessingAResults)
	ctx.Step(`^I am using AI to develop features quickly$`, iAmUsingAIToDevelopFeaturesQuickly)
	ctx.Step(`^I didn't review the code carefully$`, iDidntReviewTheCodeCarefully)
	ctx.Step(`^I follow the validate-fix-recheck pattern$`, iFollowTheValidatefixrecheckPattern)
	ctx.Step(`^I relied too heavily on AI suggestions$`, iReliedTooHeavilyOnAISuggestions)
	ctx.Step(`^I run "([^"]*)" after each AI iteration$`, iRunAfterEachAIIteration)
	ctx.Step(`^I run "([^"]*)" before reviewing the code$`, iRunBeforeReviewingTheCode)
	ctx.Step(`^I run "([^"]*)" immediately after generation$`, iRunImmediatelyAfterGeneration)
	ctx.Step(`^I should catch issues before commit$`, iShouldCatchIssuesBeforeCommit)
	ctx.Step(`^I should establish a new workflow to generate code with AI$`, iShouldEstablishANewWorkflowToGenerateCodeWithAI)
	ctx.Step(`^I should maintain high quality at high velocity$`, iShouldMaintainHighQualityAtHighVelocity)
	ctx.Step(`^I should re-check before commit$`, iShouldRecheckBeforeCommit)
	ctx.Step(`^I should review and fix issues$`, iShouldReviewAndFixIssues)
	ctx.Step(`^I should run "([^"]*)" immediately after generation$`, iShouldRunImmediatelyAfterGeneration)
	ctx.Step(`^I should iterate with AI if needed$`, iShouldIterateWithAIIfNeeded)
	ctx.Step(`^my development cycle should be fast$`, myDevelopmentCycleShouldBeFast)
	ctx.Step(`^A(\d+) should prevent "([^"]*)" from entering the codebase$`, aShouldPreventFromEnteringTheCodebase)
	ctx.Step(`^the code appears to be AI-generated$`, theCodeAppearsToBeAIgenerated)
	ctx.Step(`^the quality of my AI-assisted code should improve$`, theQualityOfMyAIassistedCodeShouldImprove)

	// Configuration and Profile Steps
	ctx.Step(`^I create a custom profile in "([^"]*)"$`, iCreateACustomProfileIn)
	ctx.Step(`^I create improvement phases in "([^"]*)"$`, iCreateImprovementPhasesIn)
	ctx.Step(`^I define external checks in "([^"]*)"$`, iDefineExternalChecksIn)
	ctx.Step(`^I select "([^"]*)" profile$`, iSelectProfile)
	ctx.Step(`^I select "([^"]*)" target$`, iSelectTarget)
	ctx.Step(`^I set severity_mode to "([^"]*)"$`, iSetSeverity_modeTo)
	ctx.Step(`^I set target to "([^"]*)"$`, iSetTargetTo)
	ctx.Step(`^A(\d+) should disable API-specific checks$`, aShouldDisableAPIspecificChecks)
	ctx.Step(`^A(\d+) should skip all other language checks$`, aShouldSkipAllOtherLanguageChecks)
	ctx.Step(`^A(\d+) should skip coverage checks$`, aShouldSkipCoverageChecks)
	ctx.Step(`^A(\d+) should skip health endpoint checks$`, aShouldSkipHealthEndpointChecks)
	ctx.Step(`^A(\d+) should skip Kubernetes checks$`, aShouldSkipKubernetesChecks)
	ctx.Step(`^A(\d+) should skip license checks$`, aShouldSkipLicenseChecks)
	ctx.Step(`^A(\d+) should skip security scans$`, aShouldSkipSecurityScans)
	ctx.Step(`^I disable cloud-native checks \(health, metrics, tracing\)$`, iDisableCloudnativeChecksHealthMetricsTracing)
	ctx.Step(`^I disable container checks \(Dockerfile, K8s\)$`, iDisableContainerChecksDockerfileKS)

	// Project Setup and Context Steps
	ctx.Step(`^I am a DevOps engineer$`, iAmADevOpsEngineer)
	ctx.Step(`^I am building a CLI application$`, iAmBuildingACLIApplication)
	ctx.Step(`^I am building a healthcare application$`, iAmBuildingAHealthcareApplication)
	ctx.Step(`^I am debugging a check issue$`, iAmDebuggingACheckIssue)
	ctx.Step(`^I am evaluating an open source project$`, iAmEvaluatingAnOpenSourceProject)
	ctx.Step(`^I am in a project directory$`, iAmInAProjectDirectory)
	ctx.Step(`^I am in early development \(PoC phase\)$`, iAmInEarlyDevelopmentPoCPhase)
	ctx.Step(`^I have a GitHub repository$`, iHaveAGitHubRepository)
	ctx.Step(`^I have a legacy monolith application$`, iHaveALegacyMonolithApplication)
	ctx.Step(`^I have a legacy project at (\d+)% maturity$`, iHaveALegacyProjectAtMaturity)
	ctx.Step(`^I have a multi-language project$`, iHaveAMultilanguageProject)
	ctx.Step(`^I have a slow project$`, iHaveASlowProject)
	ctx.Step(`^I have different requirements per branch$`, iHaveDifferentRequirementsPerBranch)
	ctx.Step(`^I have limited time$`, iHaveLimitedTime)
	ctx.Step(`^I have not configured A(\d+)$`, iHaveNotConfiguredA)
	ctx.Step(`^I cannot easily add tests or containerization$`, iCannotEasilyAddTestsOrContainerization)

	// Quality Assessment Steps
	ctx.Step(`^I analyze the maturity scores$`, iAnalyzeTheMaturityScores)
	ctx.Step(`^I can assess technical debt risk$`, iCanAssessTechnicalDebtRisk)
	ctx.Step(`^I can calculate the team average$`, iCanCalculateTheTeamAverage)
	ctx.Step(`^I can identify best performers$`, iCanIdentifyBestPerformers)
	ctx.Step(`^I can identify best practices$`, iCanIdentifyBestPractices)
	ctx.Step(`^I can identify projects needing attention$`, iCanIdentifyProjectsNeedingAttention)
	ctx.Step(`^I can identify projects needing attention \(<(\d+)%\)$`, iCanIdentifyProjectsNeedingAttentionWithThreshold)
	ctx.Step(`^I can identify top performers \((\d+)%\+\)$`, iCanIdentifyTopPerformers)
	ctx.Step(`^I can make informed integration decisions$`, iCanMakeInformedIntegrationDecisions)
	ctx.Step(`^I can make quick adoption decisions$`, iCanMakeQuickAdoptionDecisions)
	ctx.Step(`^I can set realistic quality targets$`, iCanSetRealisticQualityTargets)
	ctx.Step(`^I have maturity scores for all projects$`, iHaveMaturityScoresForAllProjects)
	ctx.Step(`^I should have clear visibility into quality$`, iShouldHaveClearVisibilityIntoQuality)
	ctx.Step(`^I should see if the maturity score increased$`, iShouldSeeIfTheMaturityScoreIncreased)
	ctx.Step(`^my team maintains (\d+) microservices$`, myTeamMaintainsMicroservices)

	// Project Management and Planning Steps
	ctx.Step(`^a project scores (\d+)% with multiple failures$`, aProjectScoresWithMultipleFailures)
	ctx.Step(`^I access the team dashboard$`, iAccessTheTeamDashboard)
	ctx.Step(`^I am viewing the team dashboard$`, iAmViewingTheTeamDashboard)
	ctx.Step(`^I can assign resources to the most critical issues$`, iCanAssignResourcesToTheMostCriticalIssues)
	ctx.Step(`^I can assign senior developers to critical tasks$`, iCanAssignSeniorDevelopersToCriticalTasks)
	ctx.Step(`^I can assign work to team members$`, iCanAssignWorkToTeamMembers)
	ctx.Step(`^I can create an actionable improvement plan$`, iCanCreateAnActionableImprovementPlan)
	ctx.Step(`^I can create an executive summary$`, iCanCreateAnExecutiveSummary)
	ctx.Step(`^I can create a gradual improvement roadmap$`, iCanCreateAGradualImprovementRoadmap)
	ctx.Step(`^I can estimate remediation investment$`, iCanEstimateRemediationInvestment)
	ctx.Step(`^I can identify which check is causing issues$`, iCanIdentifyWhichCheckIsCausingIssues)
	ctx.Step(`^I can set clear milestones$`, iCanSetClearMilestones)
	ctx.Step(`^I can share it with stakeholders$`, iCanShareItWithStakeholders)
	ctx.Step(`^I can share results with my team$`, iCanShareResultsWithMyTeam)
	ctx.Step(`^I can track progress weekly$`, iCanTrackProgressWeekly)
	ctx.Step(`^I click on the order-service project$`, iClickOnTheOrderserviceProject)
	ctx.Step(`^I create a (\d+)-week improvement plan$`, iCreateAWeekImprovementPlan)
	ctx.Step(`^I create GitHub issues for low scores$`, iCreateGitHubIssuesForLowScores)
	ctx.Step(`^I examine the detailed results$`, iExamineTheDetailedResults)
	ctx.Step(`^I filter for failed checks$`, iFilterForFailedChecks)
	ctx.Step(`^I identified analytics-platform at (\d+)% maturity$`, iIdentifiedAnalyticsplatformAtMaturity)
	ctx.Step(`^I need to understand the specific issues$`, iNeedToUnderstandTheSpecificIssues)
	ctx.Step(`^I should see a score for each project$`, iShouldSeeAScoreForEachProject)
	ctx.Step(`^I should see alerts for score drops$`, iShouldSeeAlertsForScoreDrops)
	ctx.Step(`^I should see starting scores for all projects$`, iShouldSeeStartingScoresForAllProjects)
	ctx.Step(`^I want to assess the current state$`, iWantToAssessTheCurrentState)

	// Data Visualization and Reporting Steps
	ctx.Step(`^I choose the last (\d+) days$`, iChooseTheLastDays)
	ctx.Step(`^I collect results in JSON format$`, iCollectResultsInJSONFormat)
	ctx.Step(`^I compare the score with the baseline$`, iCompareTheScoreWithTheBaseline)
	ctx.Step(`^I convert JSON to HTML with charts$`, iConvertJSONToHTMLWithCharts)
	ctx.Step(`^I create a visualization of the data$`, iCreateAVisualizationOfTheData)
	ctx.Step(`^I create a script to run A(\d+) on all projects$`, iCreateAScriptToRunAOnAllProjects)
	ctx.Step(`^I extract score, level, failures, warnings$`, iExtractScoreLevelFailuresWarnings)
	ctx.Step(`^I have set up monitoring for multiple repositories$`, iHaveSetUpMonitoringForMultipleRepositories)
	ctx.Step(`^I open the report in a browser$`, iOpenTheReportInABrowser)
	ctx.Step(`^I parse JSON output with jq$`, iParseJSONOutputWithJq)
	ctx.Step(`^I select a date range \(last (\d+) days\)$`, iSelectADateRangeLastDays)
	ctx.Step(`^I select "([^"]*)" view$`, iSelectView)
	ctx.Step(`^I should be able to export results as JSON$`, iShouldBeAbleToExportResultsAsJSON)
	ctx.Step(`^I should be able to parse with jq$`, iShouldBeAbleToParseWithJq)
	ctx.Step(`^I should see a beautiful visual report$`, iShouldSeeABeautifulVisualReport)
	ctx.Step(`^I should see current score and change from last check$`, iShouldSeeCurrentScoreAndChangeFromLastCheck)
	ctx.Step(`^I should see pass/fail/warning counts$`, iShouldSeePassfailwarningCounts)
	ctx.Step(`^I should see quality score progression$`, iShouldSeeQualityScoreProgression)
	ctx.Step(`^I should see the overall average maturity score$`, iShouldSeeTheOverallAverageMaturityScore)
	ctx.Step(`^I should see tabular results array$`, iShouldSeeTabularResultsArray)
	ctx.Step(`^I track the results over time$`, iTrackTheResultsOverTime)
	ctx.Step(`^I view the trend visualization$`, iViewTheTrendVisualization)
	ctx.Step(`^I want to benchmark team performance$`, iWantToBenchmarkTeamPerformance)
	ctx.Step(`^I want to generate HTML reports$`, iWantToGenerateHTMLReports)
	ctx.Step(`^I want to track quality improvement over time$`, iWantToTrackQualityImprovementOverTime)
	ctx.Step(`^ending scores for all projects$`, endingScoresForAllProjects)
	ctx.Step(`^last check timestamps$`, lastCheckTimestamps)
	ctx.Step(`^new failures since yesterday$`, newFailuresSinceYesterday)
	ctx.Step(`^order-service shows a score drop$`, orderserviceShowsAScoreDrop)
	ctx.Step(`^[Pp]ercentage changes$`, percentageChanges)
	ctx.Step(`^[Tt]rend indicators$`, trendIndicators)
	ctx.Step(`^[Tt]rend indicators \(up/down/same\)$`, trendIndicatorsUpdownsame)

	// Server Mode and Web Interface Steps
	ctx.Step(`^I do not have A(\d+) installed locally$`, iDoNotHaveAInstalledLocally)
	ctx.Step(`^I do not want to install tools$`, iDoNotWantToInstallTools)
	ctx.Step(`^I submit a GitHub URL$`, iSubmitAGitHubURL)
	ctx.Step(`^I submit the project's GitHub URL$`, iSubmitTheProjectsGitHubURL)
	ctx.Step(`^I use the A(\d+) web interface$`, iUseTheAWebInterface)
	ctx.Step(`^the analysis should be comprehensive$`, theAnalysisShouldBeComprehensive)
	ctx.Step(`^the analysis should complete within (\d+) minutes$`, theAnalysisShouldCompleteWithinMinutes)
	ctx.Step(`^A(\d+) runs the quality check$`, aRunsTheQualityCheck)
	ctx.Step(`^I should get complete results$`, iShouldGetCompleteResults)
	ctx.Step(`^I should get immediate quality insights$`, iShouldGetImmediateQualityInsights)

	// Compliance and Healthcare Steps
	ctx.Step(`^I add encryption-at-rest check$`, iAddEncryptionAtRestCheck)
	ctx.Step(`^I add HIPAA audit logging check$`, iAddHIPAAAuditLoggingCheck)
	ctx.Step(`^I add PHI detection check$`, iAddPHIDetectionCheck)
	ctx.Step(`^I add RBAC implementation check$`, iAddRBACImplementationCheck)
	ctx.Step(`^I need HIPAA compliance$`, iNeedHIPAACompliance)
	ctx.Step(`^A(\d+) should enforce HIPAA requirements$`, aShouldEnforceHIPAARequirements)
	ctx.Step(`^the compliance checks should run with built-in checks$`, theComplianceChecksShouldRunWithBuiltinChecks)

	// Gradual Rollout and Adoption Steps
	ctx.Step(`^A(\d+) should enforce appropriate standards for legacy code$`, aShouldEnforceAppropriateStandardsForLegacyCode)
	ctx.Step(`^A(\d+) should fail on warnings in production$`, aShouldFailOnWarningsInProduction)
	ctx.Step(`^A(\d+) should fail only on security issues$`, aShouldFailOnlyOnSecurityIssues)
	ctx.Step(`^A(\d+) should focus on basic functionality$`, aShouldFocusOnBasicFunctionality)
	ctx.Step(`^A(\d+) should focus on CLI-relevant quality$`, aShouldFocusOnCLIrelevantQuality)
	ctx.Step(`^A(\d+) should not fail prematurely$`, aShouldNotFailPrematurely)
	ctx.Step(`^A(\d+) should re-run in CI$`, aShouldRerunInCI)
	ctx.Step(`^A(\d+) should warn on non-critical issues in dev$`, aShouldWarnOnNoncriticalIssuesInDev)
	ctx.Step(`^I am introducing A(\d+) to my team$`, iAmIntroducingAToMyTeam)
	ctx.Step(`^I change to fail only on critical failures for weeks (\d+)-(\d+)$`, iChangeToFailOnlyOnCriticalFailuresForWeeks)
	ctx.Step(`^I don't want to block all PRs immediately$`, iDontWantToBlockAllPRsImmediately)
	ctx.Step(`^I enable full enforcement after week (\d+)$`, iEnableFullEnforcementAfterWeek)
	ctx.Step(`^I have a basic A(\d+) workflow$`, iHaveABasicAWorkflow)
	ctx.Step(`^I have A(\d+) configured in CI$`, iHaveAConfiguredInCI)
	ctx.Step(`^I have A(\d+) configured in the project$`, iHaveAConfiguredInTheProject)
	ctx.Step(`^I have A(\d+) running in CI$`, iHaveARunningInCI)
	ctx.Step(`^I have run "([^"]*)" on my project$`, iHaveRunOnMyProject)
	ctx.Step(`^I include coverage as critical$`, iIncludeCoverageAsCritical)
	ctx.Step(`^I mark all tests as critical$`, iMarkAllTestsAsCritical)
	ctx.Step(`^I mark security checks as critical$`, iMarkSecurityChecksAsCritical)
	ctx.Step(`^I mark tests as warn in dev "([^"]*)"$`, iMarkTestsAsWarnInDev)
	ctx.Step(`^I relax testing requirements \(coverage, cyclomatic\)$`, iRelaxTestingRequirementsCoverageCyclomatic)
	ctx.Step(`^I set "([^"]*)" for weeks (\d+)-(\d+)$`, iSetForWeeks)
	ctx.Step(`^common failures should be identified$`, commonFailuresShouldBeIdentified)
	ctx.Step(`^contributors learn quality standards$`, contributorsLearnQualityStandards)
	ctx.Step(`^the team should adapt gradually$`, theTeamShouldAdaptGradually)
	ctx.Step(`^the team should be notified$`, theTeamShouldBeNotified)
	ctx.Step(`^the transition should not block productivity$`, theTransitionShouldNotBlockProductivity)

	// GitHub Issue Creation Steps
	ctx.Step(`^an issue should be automatically created$`, anIssueShouldBeAutomaticallyCreated)
	ctx.Step(`^the issue should be labeled "([^"]*)" and "([^"]*)"$`, theIssueShouldBeLabeledAnd)
	ctx.Step(`^the issue should include the current score$`, theIssueShouldIncludeTheCurrentScore)
	ctx.Step(`^I want to provide feedback to contributors$`, iWantToProvideFeedbackToContributors)
	ctx.Step(`^the contributor knows exactly what to fix$`, theContributorKnowsExactlyWhatToFix)
	ctx.Step(`^the contributor wants to fix them$`, theContributorWantsToFixThem)

	// Slack Integration Steps
	ctx.Step(`^I post to a Slack webhook$`, iPostToASlackWebhook)
	ctx.Step(`^my team should receive quality reports in Slack$`, myTeamShouldReceiveQualityReportsInSlack)
	ctx.Step(`^I want to post results to Slack$`, iWantToPostResultsToSlack)
	ctx.Step(`^the message should be formatted nicely$`, theMessageShouldBeFormattedNicely)

	// Quality Monitoring Steps
	ctx.Step(`^I check the CI results for all PRs$`, iCheckTheCIResultsForAllPRs)
	ctx.Step(`^I should identify any new failures$`, iShouldIdentifyAnyNewFailures)
	ctx.Step(`^I should identify common failure patterns$`, iShouldIdentifyCommonFailurePatterns)
	ctx.Step(`^I should measure team adoption$`, iShouldMeasureTeamAdoption)
	ctx.Step(`^I should track time to resolution$`, iShouldTrackTimeToResolution)
	ctx.Step(`^I set up weekly automated quality checks$`, iSetUpWeeklyAutomatedQualityChecks)
	ctx.Step(`^I want to reach production-ready over (\d+) months$`, iWantToReachProductionreadyOverMonths)
	ctx.Step(`^the quality score drops below target$`, theQualityScoreDropsBelowTarget)
	ctx.Step(`^the maturity score is below (\d+)%$`, theMaturityScoreIsBelow)
	ctx.Step(`^coverage regressions$`, coverageRegressions)

	// Due Diligence Steps
	ctx.Step(`^my company acquired a startup$`, myCompanyAcquiredAStartup)
	ctx.Step(`^I need to assess their code quality$`, iNeedToAssessTheirCodeQuality)
	ctx.Step(`^I run "([^"]*)" on their codebase$`, iRunOnTheirCodebase)
	ctx.Step(`^I should see build status$`, iShouldSeeBuildStatus)
	ctx.Step(`^I should see coverage levels$`, iShouldSeeCoverageLevels)
	ctx.Step(`^I should see detected languages$`, iShouldSeeDetectedLanguages)
	ctx.Step(`^I should see security vulnerabilities$`, iShouldSeeSecurityVulnerabilities)
	ctx.Step(`^I should see test results$`, iShouldSeeTestResults)

	// Explanation and Help Steps
	ctx.Step(`^I don't understand what a check does$`, iDontUnderstandWhatACheckDoes)
	ctx.Step(`^I should see a description$`, iShouldSeeADescription)
	ctx.Step(`^I should see suggestions for improvement$`, iShouldSeeSuggestionsForImprovement)
	ctx.Step(`^I should see the check name$`, iShouldSeeTheCheckName)
	ctx.Step(`^I should see the requirements to pass$`, iShouldSeeTheRequirementsToPass)
	ctx.Step(`^I should see what tool is used$`, iShouldSeeWhatToolIsUsed)
	ctx.Step(`^I want to investigate a specific check failure$`, iWantToInvestigateASpecificCheckFailure)

	// Output Format Steps
	ctx.Step(`^I should see clearer output$`, iShouldSeeClearerOutput)
	ctx.Step(`^I should see compact encoding$`, iShouldSeeCompactEncoding)
	ctx.Step(`^the output should be in minimal token format$`, theOutputShouldBeInMinimalTokenFormat)
	ctx.Step(`^the output should be valid JSON$`, theOutputShouldBeValidJSON)
	ctx.Step(`^the format should be optimized for parsing$`, theFormatShouldBeOptimizedForParsing)
	ctx.Step(`^I want to process A(\d+) results programmatically$`, iWantToProcessAResultsProgrammatically)

	// Progress and Timing Steps
	ctx.Step(`^each check should have (\d+) minutes to complete$`, eachCheckShouldHaveMinutesToComplete)
	ctx.Step(`^I can track improvement trends over months$`, iCanTrackImprovementTrendsOverMonths)
	ctx.Step(`^I allocate a specific time for PR reviews$`, iAllocateASpecificTimeForPRReviews)
	ctx.Step(`^I am running A(\d+) for the first time$`, iAmRunningAForTheFirstTime)
	ctx.Step(`^[Rr]eal-time progress updates$`, iShouldSeeRealtimeProgressUpdates)

	// Detailed Check Results Steps
	ctx.Step(`^I should see all check results$`, iShouldSeeAllCheckResults)
	ctx.Step(`^I should see detailed race condition information$`, iShouldSeeDetailedRaceConditionInformation)
	ctx.Step(`^I should see new test failures$`, iShouldSeeNewTestFailures)
	ctx.Step(`^I should see the expected vs actual output$`, iShouldSeeTheExpectedVsActualOutput)
	ctx.Step(`^I should see the maturity level$`, iShouldSeeTheMaturityLevel)
	ctx.Step(`^I should see the specific test that failed$`, iShouldSeeTheSpecificTestThatFailed)
	ctx.Step(`^I should see (\d+) tests failing$`, iShouldSeeTestsFailing)
	ctx.Step(`^I should see which aspects need improvement$`, iShouldSeeWhichAspectsNeedImprovement)
	ctx.Step(`^the output should help me fix the issue$`, theOutputShouldHelpMeFixTheIssue)
	ctx.Step(`^the response should explain all issues$`, theResponseShouldExplainAllIssues)

	// Custom Check Steps
	ctx.Step(`^I create custom check scripts$`, iCreateCustomCheckScripts)

	// File Creation and Editing Steps
	ctx.Step(`^I create "([^"]*)"$`, iCreate)
	ctx.Step(`^I create "([^"]*)" for development$`, iCreateForDevelopment)
	ctx.Step(`^I create "([^"]*)" for main branch$`, iCreateForMainBranch)
	ctx.Step(`^I edit "([^"]*)"$`, iEdit)
	ctx.Step(`^I navigate to "([^"]*)"$`, iNavigateTo)
	ctx.Step(`^I create a script "([^"]*)"$`, iCreateAScript)

	// Re-running Steps
	ctx.Step(`^I re-run "([^"]*)"$`, iRerun)
	ctx.Step(`^I run "([^"]*)" again$`, iRunAgain)

	// Dashboard Steps
	ctx.Step(`^I can create GitHub issues directly from the dashboard$`, iCanCreateGitHubIssuesDirectlyFromTheDashboard)
	ctx.Step(`^I can demonstrate ROI of quality efforts$`, iCanDemonstrateROIOfQualityEfforts)
	ctx.Step(`^I can include it in presentations$`, iCanIncludeItInPresentations)
	ctx.Step(`^I can measure the impact of quality initiatives$`, iCanMeasureTheImpactOfQualityInitiatives)
	ctx.Step(`^I can see if all teams are improving$`, iCanSeeIfAllTeamsAreImproving)
	ctx.Step(`^I can spread knowledge across teams$`, iCanSpreadKnowledgeAcrossTeams)
	ctx.Step(`^I can spread successful practices$`, iCanSpreadSuccessfulPractices)
	ctx.Step(`^I should detect any regressions$`, iShouldDetectAnyRegressions)
	ctx.Step(`^I should immediately see the quality issues$`, iShouldImmediatelySeeTheQualityIssues)
	ctx.Step(`^I should see pipeline success rates$`, iShouldSeePipelineSuccessRates)
	ctx.Step(`^I should see relative performance$`, iShouldSeeRelativePerformance)
	ctx.Step(`^I should see updated results$`, iShouldSeeUpdatedResults)
	ctx.Step(`^the new report should show improvement$`, theNewReportShouldShowImprovement)
	ctx.Step(`^the results should show only Go-related items$`, theResultsShouldShowOnlyGorelatedItems)
	ctx.Step(`^all team members should follow same standards$`, allTeamMembersShouldFollowSameStandards)

	// Development Workflow Steps
	ctx.Step(`^I can focus on logic and architecture$`, iCanFocusOnLogicAndArchitecture)
	ctx.Step(`^I can correlate changes with initiatives$`, iCanCorrelateChangesWithInitiatives)
	ctx.Step(`^I should see (\d+)% coverage threshold \((\d+)-(\d+)\)$`, iShouldSeeCoverageThreshold)
	ctx.Step(`^the checks should not block on modern practices$`, theChecksShouldNotBlockOnModernPractices)
	ctx.Step(`^A(\d+) should run checks sequentially$`, aShouldRunChecksSequentially)

	// Test Result Steps
	ctx.Step(`^the test should pass$`, theTestShouldPass)
	ctx.Step(`^the maturity score should improve$`, theMaturityScoreShouldImprove)

	// Team and Culture Steps
	ctx.Step(`^we can celebrate high scores$`, weCanCelebrateHighScores)

	// Open Source Maintenance Steps
	ctx.Step(`^I maintain a popular Go library$`, iMaintainAPopularGoLibrary)

	// Filtering Steps
	ctx.Step(`^I only want to check Go code$`, iOnlyWantToCheckGoCode)

	// Click/UI Interaction Steps
	ctx.Step(`^I click "([^"]*)"$`, iClick)

	// Report Steps
	ctx.Step(`^the report should list critical issues$`, theReportShouldListCriticalIssues)
	ctx.Step(`^the report should provide fix suggestions$`, theReportShouldProvideFixSuggestions)
	ctx.Step(`^the report should show maturity score$`, theReportShouldShowMaturityScore)
	ctx.Step(`^PRs should receive automated quality reports$`, pRsShouldReceiveAutomatedQualityReports)
	ctx.Step(`^the comments should list failed checks$`, theCommentsShouldListFailedChecks)
	ctx.Step(`^the comments should provide fix suggestions$`, theCommentsShouldProvideFixSuggestions)
	ctx.Step(`^the comments should show maturity score$`, theCommentsShouldShowMaturityScore)

	// Enforcement Steps
	ctx.Step(`^the quality gate should be enforced$`, theQualityGateShouldBeEnforced)
	ctx.Step(`^PRs with failures should be blocked$`, pRsWithFailuresShouldBeBlocked)
	ctx.Step(`^PRs with warnings should be allowed$`, pRsWithWarningsShouldBeAllowed)
	ctx.Step(`^I should understand that AI lacks project context$`, iShouldUnderstandThatAILacksProjectContext)

	// Prompt Steps
	ctx.Step(`^the prompt asks for a function with tests and error handling$`, thePromptAsksForAFunctionWithTestsAndErrorHandling)

	// Results Steps
	ctx.Step(`^the results should be available for download$`, theResultsShouldBeAvailableForDownload)

	// Quality Improvement Steps
	ctx.Step(`^the overall codebase quality should improve$`, theOverallCodebaseQualityShouldImprove)

	// Phase-based Enforcement Steps
	ctx.Step(`^Phase (\d+) \(current\) should enforce that build passes$`, phaseCurrentShouldEnforceThatBuildPasses)
	ctx.Step(`^Phase (\d+) \((\d+) months\) should enforce (\d+)% coverage threshold$`, phaseMonthsShouldEnforceCoverageThreshold)
	ctx.Step(`^Phase (\d+) \((\d+) months\) should enforce cyclomatic complexity check$`, phaseMonthsShouldEnforceCyclomaticComplexityCheck)
	ctx.Step(`^Phase (\d+) should enforce (\d+)% coverage threshold$`, phaseShouldEnforceCoverageThreshold)
	ctx.Step(`^Phase (\d+) should enforce health check endpoints$`, phaseShouldEnforceHealthCheckEndpoints)
	ctx.Step(`^Phase (\d+) should enforce metrics instrumentation$`, phaseShouldEnforceMetricsInstrumentation)
	ctx.Step(`^Phase (\d+) should enforce security scan$`, phaseShouldEnforceSecurityScan)
	ctx.Step(`^Phase (\d+) should enforce that critical tests pass$`, phaseShouldEnforceThatCriticalTestsPass)
	ctx.Step(`^the team can track progress through phases$`, theTeamCanTrackProgressThroughPhases)

	// Dashboard Metrics Steps
	ctx.Step(`^individual repository scores$`, individualRepositoryScores)
	ctx.Step(`^key events and milestones$`, keyEventsAndMilestones)
	ctx.Step(`^impact of improvements$`, impactOfImprovements)
	ctx.Step(`^specific recommendations$`, specificRecommendations)

	// Open Source Evaluation Steps
	ctx.Step(`^I can provide a template response$`, iCanProvideATemplateResponse)

	// Additional Steps
	ctx.Step(`^I can break down the work for Week (\d+)-(\d+) to fix failing tests$`, iCanBreakDownTheWorkForWeekToFixFailingTests)
	ctx.Step(`^I can break down the work for Week (\d+)-(\d+) to improve coverage to (\d+)%$`, iCanBreakDownTheWorkForWeekToImproveCoverageTo)
	ctx.Step(`^I can break down the work for Week (\d+)-(\d+) to reach (\d+)% coverage$`, iCanBreakDownTheWorkForWeekToReachCoverage)
	ctx.Step(`^I can commit with confidence in the refactoring$`, iCanCommitWithConfidenceInTheRefactoring)
	ctx.Step(`^I need to understand what's wrong$`, iNeedToUnderstandWhatsWrong)
	ctx.Step(`^I run "([^"]*)" again$`, iRunAgain)
	ctx.Step(`^I should see [Cc]urrent score and change from last check$`, iShouldSeeCurrentScoreAndChangeFromLastCheck)
	ctx.Step(`^I should see [Qq]uality score progression$`, iShouldSeeQualityScoreProgression)
	ctx.Step(`^I should see real-time progress updates$`, iShouldSeeRealtimeProgressUpdates)
	ctx.Step(`^I should see [Ss]tarting scores for all projects$`, iShouldSeeStartingScoresForAllProjects)
	ctx.Step(`^I should see the maturity score$`, iShouldSeeTheMaturityScore)
	ctx.Step(`^I want to enforce production standards$`, iWantToEnforceProductionStandards)
	ctx.Step(`^I want to ensure quality improved$`, iWantToEnsureQualityImproved)
	ctx.Step(`^they address the failing tests$`, theyAddressTheFailingTests)
	ctx.Step(`^they fix the formatting issues$`, theyFixTheFormattingIssues)
	ctx.Step(`^they push their changes$`, theyPushTheirChanges)
	ctx.Step(`^they run "([^"]*)" locally$`, theyRunLocally)
	ctx.Step(`^[Cc]overage regressions$`, coverageRegressions)
	ctx.Step(`^[Ee]nding scores for all projects$`, endingScoresForAllProjects)
	ctx.Step(`^[Ii]mpact of improvements$`, impactOfImprovements)
	ctx.Step(`^[Ii]ndividual repository scores$`, individualRepositoryScores)
	ctx.Step(`^[Kk]ey events and milestones$`, keyEventsAndMilestones)
	ctx.Step(`^[Ll]ast check timestamps$`, lastCheckTimestamps)
	ctx.Step(`^[Nn]ew failures since yesterday$`, newFailuresSinceYesterday)
	ctx.Step(`^[Ss]pecific recommendations$`, specificRecommendations)
	ctx.Step(`^I add encryption at rest check$`, iAddEncryptionAtRestCheck)
	ctx.Step(`^I mark tests as "([^"]*)" \(warn in dev\)$`, iMarkTestsAsWarnInDev)

	// Additional missing step registrations
	ctx.Step(`^I disable container checks \(dockerfile, k8s\)$`, iDisableContainerChecksDockerfileKS)
	ctx.Step(`^I filter PRs by A(\d+) quality score$`, iFilterPRsByAQualityScore)
	ctx.Step(`^I should see (\d+)% coverage \(threshold: (\d+)%\)$`, iShouldSeeCoverageThreshold)
	ctx.Step(`^When I filter PRs by A(\d+) quality score$`, iFilterPRsByAQualityScore)

	// Before/After hooks: create temp dir for scenario, then cleanup
	ctx.Before(beforeScenarioHook)
	ctx.After(afterScenarioHook)
}

func beforeScenarioHook(c context.Context, sc *godog.Scenario) (context.Context, error) {
	ResetState()
	wd, err := os.Getwd()
	if err != nil {
		return c, fmt.Errorf("getwd: %w", err)
	}
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("a2-godogs-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return c, fmt.Errorf("create temp dir: %w", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		return c, fmt.Errorf("chdir to temp: %w", err)
	}
	s := GetState()
	s.SetTempDir(tempDir)
	s.SetOriginalDir(wd)
	return c, nil
}

func afterScenarioHook(c context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	cleanup()
	return c, nil
}

func cleanup() {
	s := GetState()
	originalDir := s.GetOriginalDir()
	tempDir := s.GetTempDir()
	if originalDir != "" {
		_ = os.Chdir(originalDir)
	}
	if tempDir != "" {
		_ = os.RemoveAll(tempDir)
	}
}
