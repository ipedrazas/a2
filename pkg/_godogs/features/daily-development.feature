Feature: Daily Development Journey
  As an AI-assisted developer
  I want to validate my code quality before committing
  So that I can catch issues early and maintain high quality

  Scenario: Validate AI-generated code
    Given I use an AI assistant to generate code
    And I have not reviewed the code yet
    When I run "a2 check" immediately after generation
    Then A2 should detect any build failures
    And A2 should identify missing tests
    And A2 should flag formatting issues
    And A2 should check for security vulnerabilities
    And I should receive actionable feedback

  Scenario: Iterative improvement cycle
    Given A2 detected issues with my code
    And I received specific failure messages
    When I fix the build issues
    And I run "a2 check" again
    Then I should see progress indicators
    And the number of failures should decrease
    And I should continue fixing until all critical issues pass

  Scenario: Quick pre-push validation
    Given I have been working on a feature for 2 hours
    And I want to push my changes
    When I run "a2 check --output=toon"
    Then I should receive a minimal token format output
    And I should see the status of all checks
    And I should identify any remaining issues
    When I address the warnings
    And I run "a2 check" one final time
    Then all checks should pass
    And my maturity score should be 100%
    And I can push with confidence

  Scenario: Get detailed check information
    Given A2 reported a failing test check
    And I don't know what's required
    When I run "a2 explain go:tests"
    Then I should see the check description
    And I should see the required tool command
    And I should see the requirements for passing
    And I should receive specific suggestions for fixing

  Scenario: Incremental development workflow
    Given I am implementing a new feature step by step
    When I run "a2 check" after each small change
    Then I should receive fast feedback
    And I should catch issues early
    And the feedback loop should be less than 2 minutes
    And my development velocity should be maintained
