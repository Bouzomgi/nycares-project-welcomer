package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/*.html
var templateFiles embed.FS

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
	Secret         string
}

var workflowFailedTmpl = template.Must(template.New("workflow_failed.html").ParseFS(templateFiles, "templates/workflow_failed.html"))
var completionTmpl = template.Must(template.New("completion.html").ParseFS(templateFiles, "templates/completion.html"))
var approvalRequestTmpl = template.Must(template.New("approval_request.html").ParseFS(templateFiles, "templates/approval_request.html"))

// WorkflowFailed returns the subject, plain text, and HTML body for a workflow step failure email.
// errorMessage should be the human-readable error (already extracted from any JSON Cause blob).
func WorkflowFailed(failedStep, errorMessage string) (subject, plainText, htmlBody string, err error) {
	subject = "NYCares Project Welcomer \u2014 Workflow Step Failed"

	plainText = fmt.Sprintf("Workflow step failed.\nStep: %s\nError: %s", failedStep, errorMessage)

	var buf bytes.Buffer
	if err = workflowFailedTmpl.Execute(&buf, workflowFailedData{
		FailedStep:   failedStep,
		ErrorMessage: errorMessage,
	}); err != nil {
		return
	}
	htmlBody = strings.TrimSpace(buf.String())

	return
}

// ApprovalRequest returns the subject, plain text, and HTML body for a message approval email.
// mockMode indicates whether send/pin requests will go to the mock server or the real NYCares platform.
// regenerateLink triggers a fresh generation with no additional context (thankYou only).
// refineFormBase is the callback URL (with token and secret) used as the HTML form action for refinement (thankYou only).
func ApprovalRequest(projectName, projectDate, messageType, messageContent, approveLink, rejectLink, regenerateLink, refineFormBase, secret string, mockMode bool) (subject, plainText, htmlBody string, err error) {
	subject = "Project Message Approval"

	destination := "real NYCares platform"
	if mockMode {
		destination = "mock server"
	}

	isThankYou := messageType == "thankYou"

	if isThankYou {
		plainText = fmt.Sprintf(
			"Project: %s\nDate: %s\nMessage Type: %s\nDestination: %s\n\nMessage Content:\n%s\n\nPlease use the HTML version of this email to approve, reject, regenerate, or refine this message.",
			projectName, projectDate, messageType, destination, messageContent,
		)
	} else {
		plainText = fmt.Sprintf(
			"Project: %s\nDate: %s\nMessage Type: %s\nDestination: %s\n\nMessage Content:\n%s\n\nPlease use the HTML version of this email to approve or reject this message.",
			projectName, projectDate, messageType, destination, messageContent,
		)
	}

	var buf bytes.Buffer
	if err = approvalRequestTmpl.Execute(&buf, approvalRequestData{
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
		Secret:         secret,
	}); err != nil {
		return
	}
	htmlBody = strings.TrimSpace(buf.String())

	return
}

// Completion returns the subject, plain text, and HTML body for a successful message send notification.
// mockMode indicates whether send/pin requests will go to the mock server or the real nycares platform.
func Completion(messageType, projectName, projectDate string, mockMode bool) (subject, plainText, htmlBody string, err error) {
	subject = "Message Sent!"

	destination := "real NYCares platform"
	if mockMode {
		destination = "mock server"
	}

	plainText = fmt.Sprintf(
		"Successfully sent %s message to %s on %s!\n\nSent to: %s",
		messageType, projectName, projectDate, destination,
	)

	var buf bytes.Buffer
	if err = completionTmpl.Execute(&buf, completionData{
		MessageType: messageType,
		ProjectName: projectName,
		ProjectDate: projectDate,
		Destination: destination,
	}); err != nil {
		return
	}
	htmlBody = strings.TrimSpace(buf.String())

	return
}
