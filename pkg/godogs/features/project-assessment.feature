Feature: Project Assessment Journey
  As an engineering manager
  I want to assess the quality of multiple projects
  So that I can make data-driven decisions about improvements

  Scenario: Team-wide quality audit
    Given my team maintains 8 microservices
    And I want to assess the current state
    When I create a script to run A2 on all projects
    And I collect results in JSON format
    And I analyze the maturity scores
    Then I should see a score for each project
    And I can identify top performers (90%+)
    And I can identify projects needing attention (<70%)
    And I can calculate the team average
    And I should have clear visibility into quality

  Scenario: Deep dive on low-performing projects
    Given I identified analytics-platform at 58% maturity
    And I need to understand the specific issues
    When I examine the detailed results
    And I filter for failed checks
    Then I should see 12 tests failing
    And I should see 35% coverage (threshold: 80%)
    And I can create an actionable improvement plan
    And I can assign resources to the most critical issues

  Scenario: Create data-driven improvement plan
    Given a project scores 58% with multiple failures
    When I create a 6-week improvement plan
    Then I can break down the work for Week 1-2 to fix failing tests
    And I can break down the work for Week 3-4 to improve coverage to 60%
    And I can break down the work for Week 5-6 to reach 80% coverage
    And I can set clear milestones
    And I can assign senior developers to critical tasks
    And I can track progress weekly

  Scenario: Due diligence for acquired codebase
    Given my company acquired a startup
    And I need to assess their code quality
    When I run "a2 check" on their codebase
    Then I should see detected languages
    And I should see build status
    And I should see test results
    And I should see coverage levels
    And I should see security vulnerabilities
    And I can create an executive summary
    And I can estimate remediation investment
    And I can assess technical debt risk
    And I can make informed integration decisions

  Scenario: Track quality progress over time
    Given I set up weekly automated quality checks
    And I create GitHub issues for low scores
    When the quality score drops below target
    Then an issue should be automatically created
    And the issue should include the current score
    And the issue should be labeled "quality" and "automated"
    And the team should be notified
    And I can track improvement trends over months

  Scenario: Compare team performance
    Given I have maturity scores for all projects
    When I create a visualization of the data
    Then I should see relative performance
    And I can identify best practices
    And I can spread knowledge across teams
    And I can set realistic quality targets
    And I can measure the impact of quality initiatives
