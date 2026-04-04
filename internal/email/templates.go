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
func ApprovalRequest(projectName, projectDate, messageType, messageContent, approveLink, rejectLink string) (subject, plainText, htmlBody string) {
	subject = "Project Message Approval"

	plainText = fmt.Sprintf(
		"Project: %s\nDate: %s\nMessage Type: %s\n\nMessage Content:\n%s\n\nApprove: %s\n\nReject: %s",
		projectName, projectDate, messageType, messageContent, approveLink, rejectLink,
	)

	htmlBody = fmt.Sprintf(
		`<p><strong>Project:</strong> %s<br><strong>Date:</strong> %s<br><strong>Message Type:</strong> %s</p>`+
			`<p><strong>Message Content:</strong></p>`+
			`<pre>%s</pre>`+
			`<p><a href="%s">Approve</a> &nbsp; <a href="%s">Reject</a></p>`,
		html.EscapeString(projectName),
		html.EscapeString(projectDate),
		html.EscapeString(messageType),
		html.EscapeString(messageContent),
		approveLink,
		rejectLink,
	)

	return
}

// Completion returns the subject, plain text, and HTML body for a successful message send notification.
func Completion(messageType, projectName, projectDate string) (subject, plainText, htmlBody string) {
	subject = "Message Sent!"

	plainText = fmt.Sprintf(
		"Successfully sent %s message to %s on %s!",
		messageType, projectName, projectDate,
	)

	htmlBody = fmt.Sprintf(
		`<p>Successfully sent <strong>%s</strong> message to <strong>%s</strong> on %s!</p>`,
		html.EscapeString(messageType),
		html.EscapeString(projectName),
		html.EscapeString(projectDate),
	)

	return
}
