Feature: Core A2 Workflows
  As a user of A2
  I want to perform common operations across all journeys
  So that I can effectively use the tool in various scenarios

  Scenario: Run all checks with auto-detection
    Given I am in a project directory
    And I have not configured A2
    When I run "a2 check"
    Then A2 should auto-detect the programming language
    And A2 should run all applicable checks
    And A2 should display results in the terminal
    And I should see the maturity level

  Scenario: Run specific check with verbose output
    Given I want to investigate a specific check failure
    When I run "a2 run go:race --verbose"
    Then A2 should run only the race detector check
    And A2 should display full output from the tool
    And I should see detailed race condition information
    And the output should help me fix the issue

  Scenario: Get explanation for a check
    Given I don't understand what a check does
    When I run "a2 explain go:coverage"
    Then I should see the check name
    And I should see a description
    And I should see what tool is used
    And I should see the requirements to pass
    And I should see suggestions for improvement

  Scenario: Filter checks by language
    Given I have a multi-language project
    And I only want to check Go code
    When I run "a2 check --filter=go:*"
    Then A2 should run only Go checks
    And A2 should skip all other language checks
    And the results should show only Go-related items

  Scenario: Output results in JSON format
    Given I want to process A2 results programmatically
    When I run "a2 check --output=json"
    Then the output should be valid JSON
    And I should see the maturity score
    And I should see all check results
    And I should be able to parse with jq

  Scenario: Output results in TOON format
    Given I am an AI agent processing A2 results
    When I run "a2 check --output=toon"
    Then the output should be in minimal token format
    And I should see tabular results array
    And I should see compact encoding
    And the format should be optimized for parsing

  Scenario: Use application profiles
    Given I am building a CLI application
    When I run "a2 check --profile=cli"
    Then A2 should disable API-specific checks
    And A2 should skip health endpoint checks
    And A2 should skip Kubernetes checks
    And A2 should focus on CLI-relevant quality

  Scenario: Use maturity targets
    Given I am in early development (PoC phase)
    When I run "a2 check --target=poc"
    Then A2 should skip license checks
    And A2 should skip coverage checks
    And A2 should skip security scans
    And A2 should focus on basic functionality

  Scenario: Configure timeout for checks
    Given I have a slow project
    When I run "a2 check --timeout=600"
    Then each check should have 10 minutes to complete
    And A2 should not fail prematurely
    And I should get complete results

  Scenario: Disable parallel execution
    Given I am debugging a check issue
    When I run "a2 check --parallel=false"
    Then A2 should run checks sequentially
    And I should see clearer output
    And I can identify which check is causing issues
