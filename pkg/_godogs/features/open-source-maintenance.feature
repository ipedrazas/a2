Feature: Open Source Maintenance Journey
  As an open source maintainer
  I want to automatically validate pull requests with A2
  So that I can review contributions efficiently

  Scenario: Automated PR quality check
    Given I maintain a popular Go library
    And I receive 20+ PRs per week
    And I have A2 configured in CI
    When a contributor submits a new PR
    Then GitHub Actions should automatically run "a2 check"
    And A2 should comment on the PR with a quality report
    And the report should show maturity score
    And the report should list critical issues
    And the report should provide fix suggestions
    And I should save 30+ minutes of manual review time

  Scenario: Identify and reject low-quality PRs
    Given a first-time contributor submits a PR
    And the code appears to be AI-generated
    When A2 runs the quality check
    And the maturity score is below 50%
    And A2 detects security issues (e.g., leaked API key)
    Then I should immediately see the quality issues
    And I can provide a template response
    And the response should explain all issues
    And the contributor knows exactly what to fix

  Scenario: Prioritize review queue by quality
    Given I have 15 PRs to review
    And I have limited time
    When I check the CI results for all PRs
    Then I can categorize green PRs (90%+) to review first
    And I can categorize yellow PRs (70-89%) to review second
    And I can categorize red PRs (<70%) to request fixes
    And I can focus my time on high-quality PRs
    And I should save 4.5 hours per week

  Scenario: Guide contributors to self-service
    Given A2 detected issues in a PR
    And the contributor wants to fix them
    When they run "a2 check" locally
    And they address the failing tests
    And they fix the formatting issues
    And they push their changes
    Then A2 should re-run in CI
    And the new report should show improvement
    And the PR should be ready for review
    And I can focus on logic and architecture

  Scenario: Batch review efficiency
    Given I allocate a specific time for PR reviews
    When I filter PRs by A2 quality score
    And I review high-quality PRs first
    Then my review time should be significantly reduced
    And low-quality PRs should be caught automatically
    And contributors learn quality standards
    And the overall codebase quality should improve
