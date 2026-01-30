Feature: Customization Journey
  As an advanced user
  I want to tailor A2 to my team's specific needs
  So that quality standards match our requirements

  Scenario: Create custom profile for legacy project
    Given I have a legacy monolith application
    And I cannot easily add tests or containerization
    When I create a custom profile in ".a2.yaml"
    And I disable container checks (dockerfile, k8s)
    And I disable cloud-native checks (health, metrics, tracing)
    And I relax testing requirements (coverage, cyclomatic)
    And I set target to "development"
    Then A2 should enforce appropriate standards for legacy code
    And the checks should not block on modern practices
    And I can create a gradual improvement roadmap

  Scenario: Industry-specific compliance checks
    Given I am building a healthcare application
    And I need HIPAA compliance
    When I define external checks in ".a2.yaml"
    And I add HIPAA audit logging check
    And I add encryption at rest check
    And I add PHI detection check
    And I add RBAC implementation check
    And I create custom check scripts
    Then A2 should enforce HIPAA requirements
    And all team members should follow same standards
    And the compliance checks should run with built-in checks

  Scenario: Environment-specific severity levels
    Given I have different requirements per branch
    When I create ".a2.yaml" for development
    And I set severity_mode to "relaxed"
    And I mark security checks as critical
    And I mark tests as "important" (warn in dev)
    Then A2 should warn on non-critical issues in dev
    And A2 should fail only on security issues

    Given I create ".a2.production.yaml" for main branch
    When I set severity_mode to "strict"
    And I mark all tests as critical
    And I include coverage as critical
    Then A2 should fail on warnings in production
    And CI should use different configs per branch

  Scenario: Custom output integrations
    Given I want to post results to Slack
    When I create a script "a2-to-slack.sh"
    And I parse JSON output with jq
    And I extract score, level, failures, warnings
    And I post to a Slack webhook
    Then my team should receive quality reports in Slack
    And the message should be formatted nicely
    And we can celebrate high scores

    Given I want to generate HTML reports
    When I create "a2-report.sh"
    And I convert JSON to HTML with charts
    And I open the report in a browser
    Then I should see a beautiful visual report
    And I can share it with stakeholders
    And I can include it in presentations

  Scenario: Gradual improvement with custom phases
    Given I have a legacy project at 40% maturity
    And I want to reach production-ready over 12 months
    When I create improvement phases in ".a2.yaml"
    Then Phase 1 (current) should enforce:
      - Build passes
      - Critical tests pass
      - Security scan
    And Phase 2 (3 months) should enforce:
      - 40% coverage threshold
    And Phase 3 (6 months) should enforce:
      - Cyclomatic complexity check
      - 60% coverage threshold
    And Phase 4 (12 months) should enforce:
      - 80% coverage threshold
      - Metrics instrumentation
      - Health check endpoints
    And the team can track progress through phases
