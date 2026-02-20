Feature: Project notification workflow
  The system sends welcome and reminder messages to upcoming NYC Cares
  project attendees based on how far away the project date is.

  Background:
    Given the notification system is ready

  Scenario: Projects too far away are skipped
    Given a project "Test Project" scheduled 30 days from now
    When the workflow runs
    Then the execution should succeed
    And the project should be skipped
    And no notification should be recorded

  Scenario: Welcome message is denied
    Given a project "Test Project" scheduled 6 days from now
    When the workflow runs
    And it requests approval
    And the message is denied
    Then the execution should succeed
    And the workflow should route to error handling
    And no notification should be recorded

  Scenario: Welcome message is approved
    Given a project "Test Project" scheduled 5 days from now
    When the workflow runs
    And it requests approval
    And the message is approved
    Then the execution should succeed
    And a welcome notification should be recorded
    And no reminder notification should be recorded

  Scenario: Welcome already sent is skipped
    Given a project "Test Project" scheduled 4 days from now
    And a welcome has already been sent
    When the workflow runs
    Then the execution should succeed

  Scenario: Reminder message is approved
    Given a project "Test Project" scheduled 1 day from now
    And a welcome has already been sent
    When the workflow runs
    And it requests approval
    And the message is approved
    Then the execution should succeed
    And a welcome notification should be recorded
    And a reminder notification should be recorded

  Scenario: Reminder already sent is skipped
    Given a project "Test Project" scheduled 2 days from now
    And both welcome and reminder have already been sent
    When the workflow runs
    Then the execution should succeed

  Scenario: Past projects are skipped
    Given a project "Test Project" scheduled -1 days from now
    When the workflow runs
    Then the execution should succeed
    And the project should be skipped
    And no notification should be recorded
