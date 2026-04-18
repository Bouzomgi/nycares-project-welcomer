package snsservice

import (
	"context"
	"testing"
)

func TestPublishHTMLEmailNotification_ValidationErrors(t *testing.T) {
	// SNSService with a nil client — validation errors fire before the client is
	// ever called, so these tests are safe without a real AWS connection.
	tests := []struct {
		name      string
		topicARN  string
		plainText string
		htmlBody  string
		subject   string
		wantErr   string
	}{
		{
			name:      "empty topicARN returns error",
			topicARN:  "",
			plainText: "hello",
			htmlBody:  "<p>hello</p>",
			subject:   "Subject",
			wantErr:   "topicARN cannot be empty",
		},
		{
			name:      "empty subject returns error",
			topicARN:  "arn:aws:sns:us-east-1:123456789012:test",
			plainText: "hello",
			htmlBody:  "<p>hello</p>",
			subject:   "",
			wantErr:   "subject cannot be empty",
		},
		{
			name:      "empty plainText returns error",
			topicARN:  "arn:aws:sns:us-east-1:123456789012:test",
			plainText: "",
			htmlBody:  "<p>hello</p>",
			subject:   "Subject",
			wantErr:   "plainText cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSNSService(nil, tt.topicARN)
			_, err := svc.PublishHTMLEmailNotification(context.Background(), tt.plainText, tt.htmlBody, tt.subject)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}
