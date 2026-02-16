package domain

import (
	"testing"
)

func TestNotificationTypeString(t *testing.T) {
	tests := []struct {
		input    NotificationType
		expected string
	}{
		{Welcome, "welcome"},
		{Reminder, "reminder"},
		{NotificationType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.input.String()
			if got != tt.expected {
				t.Errorf("NotificationType(%d).String() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseNotificationType(t *testing.T) {
	tests := []struct {
		input    string
		expected NotificationType
		wantErr  bool
	}{
		{"welcome", Welcome, false},
		{"reminder", Reminder, false},
		{"invalid", 0, true},
		{"", 0, true},
		{"Welcome", 0, true}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseNotificationType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNotificationType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("ParseNotificationType(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
