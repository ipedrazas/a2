Feature: First-Time Setup Journey
  As a senior software engineer
  I want to configure A2 for my new API project
  So that my team can maintain consistent quality standards

  Scenario: Interactive configuration for new project
    Given I have A2 installed
    And I have a new API project
    When I run "a2 add -i"
    And I select "API" as the application type
    And I select "Production" as the maturity level
    And I select "Auto-detect" for language detection
    And I choose not to enable external checks
    Then A2 should create ".a2.yaml" with sensible defaults
    And the configuration should include API profile
    And the configuration should include production target
    And E2E tests should be disabled in the configuration

  Scenario: Customize coverage thresholds
    Given I have created an initial ".a2.yaml" configuration
    And my team maintains high standards
    When I edit the configuration file
    And I set go.coverage_threshold to 85
    And I set go.cyclomatic_threshold to 10
    And I save the file
    Then A2 should use the stricter thresholds
    And checks should enforce the new standards

  Scenario: Add required documentation files
    Given I have a basic ".a2.yaml" configuration
    When I add required documentation files to the configuration
    And I include "README.md"
    And I include "LICENSE"
    And I include "CONTRIBUTING.md"
    And I include "docs/api.md"
    And I include ".env.example"
    Then A2 should verify these files exist
    And A2 should fail if required files are missing

  Scenario: Disable irrelevant checks
    Given I have a Go-only project
    When I edit ".a2.yaml"
    And I disable "*:e2e" checks
    And I disable "common:k8s" checks
    And I disable "python:*" checks
    Then A2 should skip these checks
    And the check execution should be faster
    And the results should only show relevant checks

  Scenario: Team adoption workflow
    Given I have configured A2 for my project
    When I commit the configuration file
    And I push to the repository
    And I communicate the setup to my team
    Then team members should be able to run "a2 check"
    And everyone should see the same quality standards
    And the configuration should be version controlled
