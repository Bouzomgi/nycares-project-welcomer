package email

import (
	"bytes"
	"fmt"
	"html/template"
)

type workflowFailedData struct {
	FailedStep   string
	ErrorMessage string
}

type completionData struct {
	MessageType string
	ProjectName string
	ProjectDate string
	Destination string
}

type approvalRequestData struct {
	ProjectName    string
	ProjectDate    string
	MessageType    string
	MessageContent string
	Destination    string
	ApproveLink    string
	RejectLink     string
	RegenerateLink string
	RefineFormBase string
	IsThankYou     bool
}

var workflowFailedTmpl = template.Must(template.New("workflowFailed").Parse(
	`<h2>Workflow Step Failed</h2>
<table>
  <tr><td><strong>Step</strong></td><td>{{.FailedStep}}</td></tr>
  <tr><td><strong>Error</strong></td><td>{{.ErrorMessage}}</td></tr>
</table>`))

var completionTmpl = template.Must(template.New("completion").Parse(
	`<p>Successfully sent <strong>{{.MessageType}}</strong> message to <strong>{{.ProjectName}}</strong> on {{.ProjectDate}}!</p>` +
		`<p><em>Sending to: {{.Destination}}</em></p>`))

var approvalRequestTmpl = template.Must(template.New("approvalRequest").Parse(
	`<p><strong>Project:</strong> {{.ProjectName}}<br><strong>Date:</strong> {{.ProjectDate}}<br><strong>Message Type:</strong> {{.MessageType}}<br><strong>Destination:</strong> {{.Destination}}</p>` +
		`<p><strong>Message Content:</strong></p>` +
		`<pre>{{.MessageContent}}</pre>` +
		`{{if .IsThankYou}}` +
		`<p><a href="{{.ApproveLink}}">Approve</a> &nbsp; <a href="{{.RejectLink}}">Reject</a> &nbsp; <a href="{{.RegenerateLink}}">Regenerate</a></p>` +
		`<p><strong>Refine &amp; Regenerate:</strong></p>` +
		`<form method="get" action="{{.RefineFormBase}}">` +
		`<input type="hidden" name="action" value="refine">` +
		`<textarea name="context" rows="3" cols="60" placeholder="e.g. it was raining today"></textarea><br>` +
		`<input type="submit" value="Refine &amp; Regenerate">` +
		`</form>` +
		`{{else}}` +
		`<p><a href="{{.ApproveLink}}">Approve</a> &nbsp; <a href="{{.RejectLink}}">Reject</a></p>` +
		`{{end}}`))

// WorkflowFailed returns the subject, plain text, and HTML body for a workflow step failure email.
// errorMessage should be the human-readable error (already extracted from any JSON Cause blob).
func WorkflowFailed(failedStep, errorMessage string) (subject, plainText, htmlBody string) {
	subject = "NYC Cares Project Welcomer \u2014 Workflow Step Failed"

	plainText = fmt.Sprintf("Workflow step failed.\nStep: %s\nError: %s", failedStep, errorMessage)

	var buf bytes.Buffer
	_ = workflowFailedTmpl.Execute(&buf, workflowFailedData{
		FailedStep:   failedStep,
		ErrorMessage: errorMessage,
	})
	htmlBody = buf.String()

	return
}

// ApprovalRequest returns the subject, plain text, and HTML body for a message approval email.
// mockMode indicates whether send/pin requests will go to the mock server or the real NYC Cares platform.
// regenerateLink triggers a fresh generation with no additional context (thankYou only).
// refineFormBase is the callback URL (with token and secret) used as the HTML form action for refinement (thankYou only).
func ApprovalRequest(projectName, projectDate, messageType, messageContent, approveLink, rejectLink, regenerateLink, refineFormBase string, mockMode bool) (subject, plainText, htmlBody string) {
	subject = "Project Message Approval"

	destination := "real NYC Cares platform"
	if mockMode {
		destination = "mock server"
	}

	isThankYou := messageType == "thankYou"

	if isThankYou {
		plainText = fmt.Sprintf(
			"Project: %s\nDate: %s\nMessage Type: %s\nDestination: %s\n\nMessage Content:\n%s\n\nApprove: %s\n\nReject: %s\n\nRegenerate: %s\n\nRefine: %s&action=refine&context=YOUR+CONTEXT+HERE",
			projectName, projectDate, messageType, destination, messageContent, approveLink, rejectLink, regenerateLink, refineFormBase,
		)
	} else {
		plainText = fmt.Sprintf(
			"Project: %s\nDate: %s\nMessage Type: %s\nDestination: %s\n\nMessage Content:\n%s\n\nApprove: %s\n\nReject: %s",
			projectName, projectDate, messageType, destination, messageContent, approveLink, rejectLink,
		)
	}

	var buf bytes.Buffer
	_ = approvalRequestTmpl.Execute(&buf, approvalRequestData{
		ProjectName:    projectName,
		ProjectDate:    projectDate,
		MessageType:    messageType,
		MessageContent: messageContent,
		Destination:    destination,
		ApproveLink:    approveLink,
		RejectLink:     rejectLink,
		RegenerateLink: regenerateLink,
		RefineFormBase: refineFormBase,
		IsThankYou:     isThankYou,
	})
	htmlBody = buf.String()

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

	var buf bytes.Buffer
	_ = completionTmpl.Execute(&buf, completionData{
		MessageType: messageType,
		ProjectName: projectName,
		ProjectDate: projectDate,
		Destination: destination,
	})
	htmlBody = buf.String()

	return
}
