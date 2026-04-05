// emailpreview renders all email templates with sample data and writes them to email-preview/.
// Run with: go run ./cmd/emailpreview
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/email"
)

func main() {
	outDir := "email-preview"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	type entry struct {
		name      string
		subject   string
		plainText string
		htmlBody  string
	}

	var emails []entry

	for _, mockMode := range []bool{false, true} {
		suffix := ""
		if mockMode {
			suffix = "-mock"
		}

		s, p, h := email.WorkflowFailed(
			"SendAndPinMessage",
			"connection to NYC Cares API timed out after 30s",
		)
		emails = append(emails, entry{"workflow-failed" + suffix, s, p, h})

		s, p, h = email.ApprovalRequest(
			"Central Park Cleanup",
			"2026-04-10",
			"welcome",
			"Hi volunteers! We're excited to have you join us for Central Park Cleanup on April 10th. Please arrive at the Visitor Center by 9am.",
			"http://localhost:4566/callback?token=abc123&action=approve&secret=test-secret",
			"http://localhost:4566/callback?token=abc123&action=reject&secret=test-secret",
			mockMode,
		)
		emails = append(emails, entry{"approval-request" + suffix, s, p, h})

		s, p, h = email.Completion("reminder", "Brooklyn Food Bank", "2026-04-15", mockMode)
		emails = append(emails, entry{"completion" + suffix, s, p, h})
	}

	for _, e := range emails {
		htmlFile := filepath.Join(outDir, e.name+".html")
		txtFile := filepath.Join(outDir, e.name+".txt")

		fullHTML := fmt.Sprintf("<!DOCTYPE html>\n<html><head><meta charset=\"utf-8\"><title>%s</title></head><body>\n%s\n</body></html>", e.subject, e.htmlBody)

		if err := os.WriteFile(htmlFile, []byte(fullHTML), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", htmlFile, err)
			os.Exit(1)
		}
		if err := os.WriteFile(txtFile, []byte("Subject: "+e.subject+"\n\n"+e.plainText), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", txtFile, err)
			os.Exit(1)
		}

		fmt.Printf("  %s\n  %s\n", htmlFile, txtFile)
	}

	fmt.Println("\nOpen the .html files in a browser to preview.")
}
