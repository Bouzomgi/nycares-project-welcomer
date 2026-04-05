package email

import (
	"fmt"
	"html"
)

// WorkflowFailed returns the subject, plain text, and HTML body for a workflow step failure email.
// errorMessage should be the human-readable error (already extracted from any JSON Cause blob).
func WorkflowFailed(failedStep, errorMessage string) (subject, plainText, htmlBody string) {
	subject = "NYC Cares Project Welcomer \u2014 Workflow Step Failed"

	plainText = fmt.Sprintf("Workflow step failed.\nStep: %s\nError: %s", failedStep, errorMessage)

	htmlBody = fmt.Sprintf(`<h2>Workflow Step Failed</h2>
<table>
  <tr><td><strong>Step</strong></td><td>%s</td></tr>
  <tr><td><strong>Error</strong></td><td>%s</td></tr>
</table>`,
		html.EscapeString(failedStep),
		html.EscapeString(errorMessage),
	)

	return
}

// ApprovalRequest returns the subject, plain text, and HTML body for a message approval email.
// mockMode indicates whether send/pin requests will go to the mock server or the real NYC Cares platform.
func ApprovalRequest(projectName, projectDate, messageType, messageContent, approveLink, rejectLink string, mockMode bool) (subject, plainText, htmlBody string) {
	subject = "Project Message Approval"

	destination := "real NYC Cares platform"
	if mockMode {
		destination = "mock server"
	}

	plainText = fmt.Sprintf(
		"Project: %s\nDate: %s\nMessage Type: %s\nDestination: %s\n\nMessage Content:\n%s\n\nApprove: %s\n\nReject: %s",
		projectName, projectDate, messageType, destination, messageContent, approveLink, rejectLink,
	)

	htmlBody = fmt.Sprintf(
		`<p><strong>Project:</strong> %s<br><strong>Date:</strong> %s<br><strong>Message Type:</strong> %s<br><strong>Destination:</strong> %s</p>`+
			`<p><strong>Message Content:</strong></p>`+
			`<pre>%s</pre>`+
			`<p><a href="%s">Approve</a> &nbsp; <a href="%s">Reject</a></p>`,
		html.EscapeString(projectName),
		html.EscapeString(projectDate),
		html.EscapeString(messageType),
		html.EscapeString(destination),
		html.EscapeString(messageContent),
		approveLink,
		rejectLink,
	)

	return
}

// Completion returns the subject, plain text, and HTML body for a successful message send notification.
// mockMode indicates whether send/pin requests will go to the mock server or the real NYC Cares platform.
func Completion(messageType, projectName, projectDate string, mockMode bool) (subject, plainText, htmlBody string) {
	subject = "Message Sent!"

	destination := "real NYC Cares platform"
	if mockMode {
		destination = "mock server"
	}

	plainText = fmt.Sprintf(
		"Successfully sent %s message to %s on %s!\n\nSending to: %s",
		messageType, projectName, projectDate, destination,
	)

	htmlBody = fmt.Sprintf(
		`<p>Successfully sent <strong>%s</strong> message to <strong>%s</strong> on %s!</p>`+
			`<p><em>Sending to: %s</em></p>`,
		html.EscapeString(messageType),
		html.EscapeString(projectName),
		html.EscapeString(projectDate),
		html.EscapeString(destination),
	)

	return
}
