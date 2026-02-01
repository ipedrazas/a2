Feature: AI-Assisted Development Journey
  As an AI-assisted developer
  I want to validate AI-generated code with A2
  So that I can maintain quality while using AI tools

  Scenario: Build feature with AI assistance
    Given I use an AI assistant to generate code
    And the prompt asks for a function with tests and error handling
    When I run "a2 check --filter=go:*" before reviewing the code
    Then A2 should validate the build status
    And A2 should check test coverage
    And A2 should analyze cyclomatic complexity
    And I should see which aspects need improvement

  Scenario: Diagnose and fix AI-generated issues
    Given A2 detected a failing test
    And I need to understand what's wrong
    When I run "a2 run go:tests --verbose"
    Then I should see the specific test that failed
    And I should see the expected vs actual output
    When I ask the AI to fix the specific issue
    And I re-run "a2 check"
    Then the test should pass
    And the maturity score should improve

  Scenario: Bulk AI refactoring validation
    Given I asked AI to refactor a large module
    And I want to ensure quality improved
    When I run "a2 check --output=json > refactor-results.json"
    And I compare the score with the baseline
    Then I should see if the maturity score increased
    And I should identify any new failures
    And I should detect any regressions
    And I can commit with confidence in the refactoring

  Scenario: Learn from A2 feedback
    Given I relied too heavily on AI suggestions
    And I didn't review the code carefully
    When I run "a2 check"
    And A2 detects multiple critical issues
    Then I should understand that AI lacks project context
    And I should establish a new workflow to generate code with AI
    And I should run "a2 check" immediately after generation
    And I should review and fix issues
    And I should iterate with AI if needed
    And I should re-check before commit
    And the quality of my AI-assisted code should improve

  Scenario: Rapid iteration cycle
    Given I am using AI to develop features quickly
    When I follow the validate-fix-recheck pattern
    And I run "a2 check" after each AI iteration
    Then my development cycle should be fast
    And I should catch issues before commit
    And I should maintain high quality at high velocity
    And A2 should prevent "AI slop" from entering the codebase
