Feature: CI/CD Integration Journey
  As a DevOps engineer
  I want to integrate A2 quality gates into our CI/CD pipeline
  So that low-quality code cannot reach production

  Scenario: Setup GitHub Actions workflow
    Given I have a GitHub repository
    And I have A2 configured in the project
    When I create ".github/workflows/a2-check.yml"
    And I configure it to run on pull requests and pushes
    And I install A2 in the workflow
    And I run "a2 check --output=json > results.json"
    And I upload the results as artifacts
    Then the workflow should run on every PR
    And the results should be available for download
    And the check status should appear on GitHub

  Scenario: Set strict quality enforcement
    Given I have a basic A2 workflow
    And I want to enforce production standards
    When I update the workflow to check exit codes
    And I configure it to fail on exit code 2 (failures)
    And I enable branch protection rules
    Then PRs with failures should be blocked
    And PRs with warnings should be allowed
    And the quality gate should be enforced

  Scenario: Add PR comment bot
    Given I have A2 running in CI
    And I want to provide feedback to contributors
    When I add a GitHub Action script
    And I configure it to parse results.json
    And I set it to comment on pull requests
    Then PRs should receive automated quality reports
    And the comments should show maturity score
    And the comments should list failed checks
    And the comments should provide fix suggestions

  Scenario: Gradual rollout strategy
    Given I am introducing A2 to my team
    And I don't want to block all PRs immediately
    When I set "continue-on-error: true" for weeks 1-2
    And I change to fail only on critical failures for weeks 3-4
    And I enable full enforcement after week 5
    Then the team should adapt gradually
    And common failures should be identified
    And the transition should not block productivity

  Scenario: Monitor quality trends
    Given A2 is running in CI for all projects
    When I track the results over time
    Then I should see pipeline success rates
    And I should identify common failure patterns
    And I should measure team adoption
    And I should track time to resolution
