//go:build integration

package integration

import (
	"fmt"

	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/cucumber/godog"
)

// scenarioContext holds per-scenario state, reset between scenarios.
type scenarioContext struct {
	tc          *testClients
	today       string
	projectDate string
	execArn     string
	taskToken   string
}

// --- Background ---

func (sc *scenarioContext) theNotificationSystemIsReady() error {
	tc, err := initTestClients()
	if err != nil {
		return err
	}
	sc.tc = tc
	sc.today = currentDateStr()
	return nil
}

// --- Given ---

func (sc *scenarioContext) aProjectScheduledDaysFromNow(name string, days int) error {
	sc.projectDate = dateOffset(sc.today, days)

	if err := sc.tc.deleteNotification(name, sc.projectDate); err != nil {
		return err
	}

	return sc.tc.setMockProjects([]projectInput{
		{Name: name, Date: sc.projectDate, Id: testProjectId, CampaignId: testCampaignId},
	})
}

func (sc *scenarioContext) aWelcomeHasAlreadyBeenSent() error {
	return sc.tc.seedNotification(testProjectName, sc.projectDate, testProjectId, true, false)
}

func (sc *scenarioContext) bothWelcomeAndReminderHaveAlreadyBeenSent() error {
	return sc.tc.seedNotification(testProjectName, sc.projectDate, testProjectId, true, true)
}

// --- When ---

func (sc *scenarioContext) theWorkflowRuns() error {
	arn, err := sc.tc.startExecution()
	if err != nil {
		return err
	}
	sc.execArn = arn
	return nil
}

func (sc *scenarioContext) itRequestsApproval() error {
	token, err := sc.tc.pollForTaskToken(sc.execArn)
	if err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("expected non-empty task token")
	}
	sc.taskToken = token
	return nil
}

func (sc *scenarioContext) theMessageIsDenied() error {
	return sc.tc.rejectTask(sc.taskToken)
}

func (sc *scenarioContext) theMessageIsApproved() error {
	return sc.tc.approveTask(sc.taskToken)
}

// --- Then ---

func (sc *scenarioContext) theExecutionShouldSucceed() error {
	status, err := sc.tc.waitForExecutionComplete(sc.execArn)
	if err != nil {
		return err
	}
	if status != sfntypes.ExecutionStatusSucceeded {
		return fmt.Errorf("expected SUCCEEDED, got %s", status)
	}
	return nil
}

func (sc *scenarioContext) theProjectShouldBeSkipped() error {
	skipped, err := sc.tc.executionEndedWithSkip(sc.execArn)
	if err != nil {
		return err
	}
	if !skipped {
		return fmt.Errorf("expected execution to skip via EndProjectIteration")
	}
	return nil
}

func (sc *scenarioContext) theWorkflowShouldRouteToErrorHandling() error {
	entered, err := sc.tc.executionEnteredState(sc.execArn, "DLQNotifier")
	if err != nil {
		return err
	}
	if !entered {
		return fmt.Errorf("expected execution to route to error handling")
	}
	return nil
}

func (sc *scenarioContext) noNotificationShouldBeRecorded() error {
	item, err := sc.tc.getNotification(testProjectName, sc.projectDate)
	if err != nil {
		return err
	}
	if item != nil {
		return fmt.Errorf("expected no notification for %q on %s, but found one", testProjectName, sc.projectDate)
	}
	return nil
}

func (sc *scenarioContext) aWelcomeNotificationShouldBeRecorded() error {
	return sc.assertNotificationBoolAttr("HasSentWelcome", true)
}

func (sc *scenarioContext) noReminderNotificationShouldBeRecorded() error {
	return sc.assertNotificationBoolAttr("HasSentReminder", false)
}

func (sc *scenarioContext) aReminderNotificationShouldBeRecorded() error {
	return sc.assertNotificationBoolAttr("HasSentReminder", true)
}

func (sc *scenarioContext) assertNotificationBoolAttr(attr string, expected bool) error {
	item, err := sc.tc.getNotification(testProjectName, sc.projectDate)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("expected notification for %q on %s, but found none", testProjectName, sc.projectDate)
	}

	val, err := getBoolAttr(item, attr)
	if err != nil {
		return err
	}
	if val != expected {
		return fmt.Errorf("attribute %q = %v, want %v", attr, val, expected)
	}
	return nil
}

// InitializeScenario wires Gherkin steps to their implementations.
func InitializeScenario(ctx *godog.ScenarioContext) {
	sc := &scenarioContext{}

	// Background
	ctx.Given(`^the notification system is ready$`, sc.theNotificationSystemIsReady)

	// Given
	ctx.Given(`^a project "([^"]*)" scheduled (\d+) days? from now$`, sc.aProjectScheduledDaysFromNow)
	ctx.Given(`^a project "([^"]*)" scheduled (-\d+) days? from now$`, sc.aProjectScheduledDaysFromNow)
	ctx.Given(`^a welcome has already been sent$`, sc.aWelcomeHasAlreadyBeenSent)
	ctx.Given(`^both welcome and reminder have already been sent$`, sc.bothWelcomeAndReminderHaveAlreadyBeenSent)

	// When
	ctx.When(`^the workflow runs$`, sc.theWorkflowRuns)
	ctx.When(`^it requests approval$`, sc.itRequestsApproval)
	ctx.When(`^the message is denied$`, sc.theMessageIsDenied)
	ctx.When(`^the message is approved$`, sc.theMessageIsApproved)

	// Then
	ctx.Then(`^the execution should succeed$`, sc.theExecutionShouldSucceed)
	ctx.Then(`^the project should be skipped$`, sc.theProjectShouldBeSkipped)
	ctx.Then(`^the workflow should route to error handling$`, sc.theWorkflowShouldRouteToErrorHandling)
	ctx.Then(`^no notification should be recorded$`, sc.noNotificationShouldBeRecorded)
	ctx.Then(`^a welcome notification should be recorded$`, sc.aWelcomeNotificationShouldBeRecorded)
	ctx.Then(`^no reminder notification should be recorded$`, sc.noReminderNotificationShouldBeRecorded)
	ctx.Then(`^a reminder notification should be recorded$`, sc.aReminderNotificationShouldBeRecorded)
}
