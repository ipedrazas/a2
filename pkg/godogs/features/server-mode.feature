Feature: Server Mode Journey
  As a user or stakeholder
  I want to analyze repositories through a web interface
  So that I can get quality insights without local installation

  Scenario: Quick repository analysis via web
    Given I navigate to "https://a2.example.com"
    And I do not have A2 installed locally
    When I submit a GitHub URL
    And I select "API" profile
    And I select "Production" target
    And I click "Analyze"
    Then I should see real-time progress updates
    And the analysis should complete within 5 minutes
    And I should see a maturity score
    And I should see pass/fail/warning counts
    And I should be able to export results as JSON

  Scenario: Team quality dashboard
    Given I am a DevOps engineer
    And I have set up monitoring for multiple repositories
    When I access the team dashboard
    Then I should see:
      - Overall average maturity score
      - Individual repository scores
      - Trend indicators (up/down/same)
      - Last check timestamps
    And I should see alerts for score drops
    And I should see new test failures
    And I can identify projects needing attention

  Scenario: Drill down into project details
    Given I am viewing the team dashboard
    And order-service shows a score drop
    When I click on the order-service project
    Then I should see:
      - Current score and change from last check
      - New failures since yesterday
      - Coverage regressions
      - Specific recommendations
    And I can create GitHub issues directly from the dashboard
    And I can assign work to team members

  Scenario: Historical trend analysis
    Given I want to track quality improvement over time
    When I select a date range (last 90 days)
    And I view the trend visualization
    Then I should see:
      - Quality score progression
      - Key events and milestones
      - Impact of improvements
    And I can correlate changes with initiatives
    And I can demonstrate ROI of quality efforts

  Scenario: Compare repositories
    Given I want to benchmark team performance
    When I select "Team Comparison" view
    And I choose the last 30 days
    Then I should see:
      - Starting scores for all projects
      - Ending scores for all projects
      - Percentage changes
      - Trend indicators
    And I can identify best performers
    And I can spread successful practices
    And I can see if all teams are improving

  Scenario: No-installation quick checks
    Given I am evaluating an open source project
    And I do not want to install tools
    When I use the A2 web interface
    And I submit the project's GitHub URL
    Then I should get immediate quality insights
    And I can make quick adoption decisions
    And I can share results with my team
    And the analysis should be comprehensive
