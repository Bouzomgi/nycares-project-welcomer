package email

import (
	"strings"
	"testing"
)

func TestWorkflowFailed_Subject(t *testing.T) {
	subject, _, _, err := WorkflowFailed("Login", "connection refused")
	if err != nil {
		t.Fatal(err)
	}
	want := "NYC Cares Project Welcomer \u2014 Workflow Step Failed"
	if subject != want {
		t.Errorf("subject = %q, want %q", subject, want)
	}
}

func TestWorkflowFailed_ContainsFields(t *testing.T) {
	_, plainText, htmlBody, err := WorkflowFailed("FetchProjects", "timeout after 30s")
	if err != nil {
		t.Fatal(err)
	}

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
	_, _, htmlBody, err := WorkflowFailed("<b>Step</b>", `<script>alert(1)</script>`)
	if err != nil {
		t.Fatal(err)
	}

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
	_, plainText, htmlBody, err := WorkflowFailed("Login", "connection refused")
	if err != nil {
		t.Fatal(err)
	}

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
	subject, _, _, err := ApprovalRequest("Park Cleanup", "2026-04-10", "welcome", "content", "http://approve", "http://reject", "http://regenerate", "http://refine-base", "", false)
	if err != nil {
		t.Fatal(err)
	}
	if subject != "Project Message Approval" {
		t.Errorf("subject = %q, want %q", subject, "Project Message Approval")
	}
}

func TestApprovalRequest_ContainsFields(t *testing.T) {
	_, plainText, htmlBody, err := ApprovalRequest("Park Cleanup", "2026-04-10", "welcome", "Hello volunteers!", "http://approve", "http://reject", "http://regenerate", "http://refine-base", "", false)
	if err != nil {
		t.Fatal(err)
	}

	plainChecks := []string{"Park Cleanup", "2026-04-10", "welcome", "Hello volunteers!"}
	for _, s := range plainChecks {
		if !strings.Contains(plainText, s) {
			t.Errorf("plainText missing %q", s)
		}
	}
	htmlChecks := []string{"Park Cleanup", "2026-04-10", "welcome", "Hello volunteers!", "http://approve", "http://reject"}
	for _, s := range htmlChecks {
		if !strings.Contains(htmlBody, s) {
			t.Errorf("htmlBody missing %q", s)
		}
	}

	// Regenerate/Refine must not appear for non-thankYou messages
	for _, s := range []string{"http://regenerate", "http://refine-base", "Regenerate", "refine"} {
		if strings.Contains(plainText, s) {
			t.Errorf("plainText should not contain %q for welcome message", s)
		}
		if strings.Contains(htmlBody, s) {
			t.Errorf("htmlBody should not contain %q for welcome message", s)
		}
	}
}

func TestApprovalRequest_HTMLEscaping(t *testing.T) {
	_, _, htmlBody, err := ApprovalRequest(
		"<Project>", "<date>", "<type>", "<script>xss</script>",
		"http://approve", "http://reject", "http://regenerate", "http://refine-base", "", false,
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, raw := range []string{"<Project>", "<date>", "<type>", "<script>"} {
		if strings.Contains(htmlBody, raw) {
			t.Errorf("htmlBody should escape %q", raw)
		}
	}
}

func TestApprovalRequest_MockMode(t *testing.T) {
	_, plainText, htmlBody, err := ApprovalRequest("Park Cleanup", "2026-04-10", "welcome", "content", "http://approve", "http://reject", "http://regenerate", "http://refine-base", "", true)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(plainText, "mock server") {
		t.Error("plainText should indicate mock server when mockMode=true")
	}
	if !strings.Contains(htmlBody, "mock server") {
		t.Error("htmlBody should indicate mock server when mockMode=true")
	}

	_, plainText2, htmlBody2, err := ApprovalRequest("Park Cleanup", "2026-04-10", "welcome", "content", "http://approve", "http://reject", "http://regenerate", "http://refine-base", "", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(plainText2, "real NYC Cares platform") {
		t.Error("plainText should indicate real platform when mockMode=false")
	}
	if !strings.Contains(htmlBody2, "real NYC Cares platform") {
		t.Error("htmlBody should indicate real platform when mockMode=false")
	}
}

func TestApprovalRequest_ContainsRegenerateAndRefine(t *testing.T) {
	_, plainText, htmlBody, err := ApprovalRequest("Park Cleanup", "2026-04-10", "thankYou", "Thanks!", "http://approve", "http://reject", "http://regenerate", "http://refine-base", "", false)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(plainText, "HTML version") {
		t.Error("plainText should direct user to HTML version")
	}
	if !strings.Contains(htmlBody, "Regenerate") {
		t.Error("htmlBody missing Regenerate option")
	}
	if !strings.Contains(htmlBody, `value="refine"`) {
		t.Error("htmlBody missing refine form hidden action field")
	}
	if !strings.Contains(htmlBody, "http://refine-base") {
		t.Error("htmlBody missing refine form base URL")
	}
}

func TestCompletion_Subject(t *testing.T) {
	subject, _, _, err := Completion("welcome", "Park Cleanup", "2026-04-10", false)
	if err != nil {
		t.Fatal(err)
	}
	if subject != "Message Sent!" {
		t.Errorf("subject = %q, want %q", subject, "Message Sent!")
	}
}

func TestCompletion_ContainsFields(t *testing.T) {
	_, plainText, htmlBody, err := Completion("reminder", "Food Bank", "2026-04-15", false)
	if err != nil {
		t.Fatal(err)
	}

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
	_, _, htmlBody, err := Completion("<b>welcome</b>", "<Project & Name>", "2026-04-15", false)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(htmlBody, "<b>welcome</b>") {
		t.Error("htmlBody should escape HTML in messageType")
	}
	if strings.Contains(htmlBody, "<Project & Name>") {
		t.Error("htmlBody should escape HTML in projectName")
	}
}

func TestCompletion_MockMode(t *testing.T) {
	_, plainText, htmlBody, err := Completion("welcome", "Park Cleanup", "2026-04-10", true)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(plainText, "mock server") {
		t.Error("plainText should indicate mock server when mockMode=true")
	}
	if !strings.Contains(htmlBody, "mock server") {
		t.Error("htmlBody should indicate mock server when mockMode=true")
	}

	_, plainText2, htmlBody2, err := Completion("welcome", "Park Cleanup", "2026-04-10", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(plainText2, "real NYC Cares platform") {
		t.Error("plainText should indicate real platform when mockMode=false")
	}
	if !strings.Contains(htmlBody2, "real NYC Cares platform") {
		t.Error("htmlBody should indicate real platform when mockMode=false")
	}
}
