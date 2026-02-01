Feature: Quick Start Journey
  As a new user
  I want to run A2 checks on my project within 5 minutes
  So that I can get immediate value without complex setup

  Scenario: Successfully install and run first check
    Given I have Go installed on my system
    And I have an existing Go project directory
    When I install A2 using "go install github.com/ipedrazas/a2@latest"
    And I verify the installation with "a2 --version"
    And I navigate to my project directory
    And I run "a2 check"
    Then A2 should auto-detect the project language
    And A2 should run appropriate checks in parallel
    And A2 should display results with color coding
    And I should see a maturity score
    And I should receive clear suggestions for improvement

  Scenario: Interpret check results correctly
    Given I have run "a2 check" on my project
    When I view the output
    Then I should see green checkmarks for passed checks
    And I should see red X marks for critical failures
    And I should see yellow warnings for recommendations
    And I should see blue info for optional tools not installed

  Scenario: Fix and re-run checks
    Given A2 detected issues in my code
    And I received actionable suggestions
    When I fix the identified issues
    And I run "a2 check" again
    Then I should see updated results
    And my maturity score should improve
    And I should be able to commit with confidence

  Scenario: Handle missing required tools
    Given I am running A2 for the first time
    When some required tools are not installed
    Then A2 should display Info status for missing tools
    And A2 should continue running available checks
    And A2 should suggest how to install missing tools
    And the check should complete successfully
