package email

import (
	"strings"
	"testing"
)

func TestWorkflowFailed_Subject(t *testing.T) {
	subject, _, _ := WorkflowFailed("Login", "connection refused")
	want := "NYC Cares Project Welcomer \u2014 Workflow Step Failed"
	if subject != want {
		t.Errorf("subject = %q, want %q", subject, want)
	}
}

func TestWorkflowFailed_ContainsFields(t *testing.T) {
	_, plainText, htmlBody := WorkflowFailed("FetchProjects", "timeout after 30s")

	for _, s := range []string{"FetchProjects", "timeout after 30s"} {
		if !strings.Contains(plainText, s) {
			t.Errorf("plainText missing %q", s)
		}
		if !strings.Contains(htmlBody, s) {
			t.Errorf("htmlBody missing %q", s)
		}
	}
}

func TestWorkflowFailed_HTMLEscaping(t *testing.T) {
	_, _, htmlBody := WorkflowFailed("<b>Step</b>", `<script>alert(1)</script>`)

	if strings.Contains(htmlBody, "<script>") {
		t.Error("htmlBody should escape <script> in errorMessage")
	}
	if strings.Contains(htmlBody, "<b>Step</b>") {
		t.Error("htmlBody should escape HTML in failedStep")
	}
	if !strings.Contains(htmlBody, "&lt;script&gt;") {
		t.Error("htmlBody should contain escaped script tag")
	}
}

func TestWorkflowFailed_NoErrorTypeNoise(t *testing.T) {
	_, plainText, htmlBody := WorkflowFailed("Login", "connection refused")

	for _, noise := range []string{"errorType", "errorString", "stackTrace"} {
		if strings.Contains(plainText, noise) {
			t.Errorf("plainText should not contain %q", noise)
		}
		if strings.Contains(htmlBody, noise) {
			t.Errorf("htmlBody should not contain %q", noise)
		}
	}
}

func TestApprovalRequest_Subject(t *testing.T) {
	subject, _, _ := ApprovalRequest("Park Cleanup", "2026-04-10", "welcome", "content", "http://approve", "http://reject")
	if subject != "Project Message Approval" {
		t.Errorf("subject = %q, want %q", subject, "Project Message Approval")
	}
}

func TestApprovalRequest_ContainsFields(t *testing.T) {
	_, plainText, htmlBody := ApprovalRequest("Park Cleanup", "2026-04-10", "welcome", "Hello volunteers!", "http://approve", "http://reject")

	checks := []string{"Park Cleanup", "2026-04-10", "welcome", "Hello volunteers!", "http://approve", "http://reject"}
	for _, s := range checks {
		if !strings.Contains(plainText, s) {
			t.Errorf("plainText missing %q", s)
		}
		if !strings.Contains(htmlBody, s) {
			t.Errorf("htmlBody missing %q", s)
		}
	}
}

func TestApprovalRequest_HTMLEscaping(t *testing.T) {
	_, _, htmlBody := ApprovalRequest(
		"<Project>", "<date>", "<type>", "<script>xss</script>",
		"http://approve", "http://reject",
	)

	for _, raw := range []string{"<Project>", "<date>", "<type>", "<script>"} {
		if strings.Contains(htmlBody, raw) {
			t.Errorf("htmlBody should escape %q", raw)
		}
	}
}

func TestCompletion_Subject(t *testing.T) {
	subject, _, _ := Completion("welcome", "Park Cleanup", "2026-04-10")
	if subject != "Message Sent!" {
		t.Errorf("subject = %q, want %q", subject, "Message Sent!")
	}
}

func TestCompletion_ContainsFields(t *testing.T) {
	_, plainText, htmlBody := Completion("reminder", "Food Bank", "2026-04-15")

	for _, s := range []string{"reminder", "Food Bank", "2026-04-15"} {
		if !strings.Contains(plainText, s) {
			t.Errorf("plainText missing %q", s)
		}
		if !strings.Contains(htmlBody, s) {
			t.Errorf("htmlBody missing %q", s)
		}
	}
}

func TestCompletion_HTMLEscaping(t *testing.T) {
	_, _, htmlBody := Completion("<b>welcome</b>", "<Project & Name>", "2026-04-15")

	if strings.Contains(htmlBody, "<b>welcome</b>") {
		t.Error("htmlBody should escape HTML in messageType")
	}
	if strings.Contains(htmlBody, "<Project & Name>") {
		t.Error("htmlBody should escape HTML in projectName")
	}
}
