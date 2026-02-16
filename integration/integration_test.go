//go:build integration

package integration

import (
	"testing"

	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

// These integration tests run sequentially against a live LocalStack environment.
// They require docker compose to be up with all services running.
//
// Run with: go test ./integration/ -v -count=1
//
// The tests build on each other's DynamoDB state, so they must run in order.

const (
	testProjectName = "Test Project"
	testProjectId   = "a1Bxx0000001XYZ"
	testCampaignId  = "11111111-1111-1111-1111-111111111111"
)

func TestIntegration(t *testing.T) {
	tc := newTestClients(t)
	today := currentDateStr()

	// Clean up any pre-existing DynamoDB state for our test project
	projectDateWelcome := dateOffset(today, 5) // 5 days from "today" -> in welcome range
	tc.deleteNotification(t, testProjectName, projectDateWelcome)

	t.Run("1_all_projects_too_far_away", func(t *testing.T) {
		// Set project date to 30 days from now (outside welcome window of 7 days)
		farDate := dateOffset(today, 30)
		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: farDate, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)
		status := tc.waitForExecutionComplete(t, execArn)

		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED, got %s", status)
		}

		// Verify the execution hit EndProjectIteration (skipped)
		if !tc.executionEndedWithSkip(t, execArn) {
			t.Error("expected execution to skip via EndProjectIteration")
		}

		// Verify no DynamoDB entry was created
		item := tc.getNotification(t, testProjectName, farDate)
		if item != nil {
			t.Error("expected no DynamoDB entry for far-away project")
		}
	})

	t.Run("2_welcome_range_denied", func(t *testing.T) {
		// Set project in welcome range (5 days from now)
		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: projectDateWelcome, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)

		// Wait for the task token (approval step)
		taskToken := tc.pollForTaskToken(t, execArn)
		if taskToken == "" {
			t.Fatal("expected non-empty task token")
		}

		// Reject the approval
		tc.rejectTask(t, taskToken)

		// Execution should complete (DLQNotifier catches the rejection)
		status := tc.waitForExecutionComplete(t, execArn)
		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED (DLQ handled), got %s", status)
		}

		// Verify DLQNotifier was invoked
		if !tc.executionEnteredState(t, execArn, "DLQNotifier") {
			t.Error("expected execution to enter DLQNotifier state")
		}

		// No DynamoDB record should exist since the message was rejected
		item := tc.getNotification(t, testProjectName, projectDateWelcome)
		if item != nil {
			t.Error("expected no DynamoDB entry after rejection")
		}
	})

	t.Run("3_welcome_range_approved", func(t *testing.T) {
		// Same project still in welcome range
		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: projectDateWelcome, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)

		// Wait for approval step
		taskToken := tc.pollForTaskToken(t, execArn)
		tc.approveTask(t, taskToken)

		status := tc.waitForExecutionComplete(t, execArn)
		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED, got %s", status)
		}

		// Verify NotifyCompletion was reached
		if !tc.executionEnteredState(t, execArn, "NotifyCompletion") {
			t.Error("expected execution to reach NotifyCompletion")
		}

		// Verify DynamoDB has the welcome record
		item := tc.getNotification(t, testProjectName, projectDateWelcome)
		if item == nil {
			t.Fatal("expected DynamoDB entry after approved welcome")
		}
		assertBoolAttr(t, *item, "HasSentWelcome", true)
		assertBoolAttr(t, *item, "HasSentReminder", false)
	})

	t.Run("4_welcome_already_sent_skipped", func(t *testing.T) {
		// Same project, still in welcome range, but welcome already sent
		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: projectDateWelcome, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)
		status := tc.waitForExecutionComplete(t, execArn)

		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED, got %s", status)
		}

		// Should skip since welcome was already sent and it's not yet in reminder range
		// The execution should hit EndProjectIteration or DLQNotifier
		// (depending on whether "all notifications already sent" is caught)
	})

	t.Run("5_reminder_range_approved", func(t *testing.T) {
		// Move project to reminder range (1 day from now, within the 2-day reminder window)
		projectDateReminder := dateOffset(today, 1)

		// We need to update the mock projects and also ensure DynamoDB has the right state
		// First, clean up old entry and create one with HasSentWelcome=true for this date
		tc.deleteNotification(t, testProjectName, projectDateReminder)
		tc.seedNotification(t, testProjectName, projectDateReminder, testProjectId, true, false)

		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: projectDateReminder, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)

		// Wait for approval step
		taskToken := tc.pollForTaskToken(t, execArn)
		tc.approveTask(t, taskToken)

		status := tc.waitForExecutionComplete(t, execArn)
		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED, got %s", status)
		}

		// Verify DynamoDB has both welcome and reminder sent
		item := tc.getNotification(t, testProjectName, projectDateReminder)
		if item == nil {
			t.Fatal("expected DynamoDB entry after approved reminder")
		}
		assertBoolAttr(t, *item, "HasSentWelcome", true)
		assertBoolAttr(t, *item, "HasSentReminder", true)
	})

	t.Run("6_reminder_already_sent_skipped", func(t *testing.T) {
		// Same project in reminder range, reminder already sent
		projectDateReminder := dateOffset(today, 1)

		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: projectDateReminder, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)
		status := tc.waitForExecutionComplete(t, execArn)

		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED, got %s", status)
		}

		// Should skip since both messages already sent
	})

	t.Run("7_project_in_past_skipped", func(t *testing.T) {
		// Set project date to yesterday (in the past)
		pastDate := dateOffset(today, -1)
		tc.setMockProjects(t, []projectInput{
			{Name: testProjectName, Date: pastDate, Id: testProjectId, CampaignId: testCampaignId},
		})

		execArn := tc.startExecution(t)
		status := tc.waitForExecutionComplete(t, execArn)

		if status != sfntypes.ExecutionStatusSucceeded {
			t.Fatalf("expected SUCCEEDED, got %s", status)
		}

		// Should skip since project is in the past
		item := tc.getNotification(t, testProjectName, pastDate)
		if item != nil {
			t.Error("expected no DynamoDB entry for past project")
		}
	})
}
