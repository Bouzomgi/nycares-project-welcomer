package computemessage

import (
	"testing"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

func TestShouldSendWelcome(t *testing.T) {
	projectDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{"8 days before - too early", projectDate.AddDate(0, 0, -8), false},
		{"7 days before - boundary", projectDate.AddDate(0, 0, -7).Add(time.Hour), true},
		{"5 days before - in window", projectDate.AddDate(0, 0, -5), true},
		{"1 day before - in window", projectDate.AddDate(0, 0, -1), true},
		{"project day - too late", projectDate, false},
		{"after project - too late", projectDate.AddDate(0, 0, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSendWelcome(tt.now, projectDate)
			if got != tt.expected {
				t.Errorf("shouldSendWelcome(%v, %v) = %v, want %v", tt.now, projectDate, got, tt.expected)
			}
		})
	}
}

func TestShouldSendReminder(t *testing.T) {
	projectDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{"3 days before - too early", projectDate.AddDate(0, 0, -3), false},
		{"2 days before - boundary", projectDate.AddDate(0, 0, -2).Add(time.Hour), true},
		{"1 day before - in window", projectDate.AddDate(0, 0, -1), true},
		{"project day - too late", projectDate, false},
		{"after project - too late", projectDate.AddDate(0, 0, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSendReminder(tt.now, projectDate)
			if got != tt.expected {
				t.Errorf("shouldSendReminder(%v, %v) = %v, want %v", tt.now, projectDate, got, tt.expected)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Park Cleanup", "park-cleanup"},
		{"Food Bank Volunteer", "food-bank-volunteer"},
		{"UPPER CASE", "upper-case"},
		{"single", "single"},
		{"", ""},
		{"  extra   spaces  ", "extra-spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toKebabCase(tt.input)
			if got != tt.expected {
				t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidS3URI(t *testing.T) {
	tests := []struct {
		uri      string
		expected bool
	}{
		{"s3://my-bucket/path/to/file.md", true},
		{"s3://bucket123/file", true},
		{"s3://ab/file", false},   // bucket name too short
		{"http://bucket/file", false},
		{"s3:///file", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			got := isValidS3URI(tt.uri)
			if got != tt.expected {
				t.Errorf("isValidS3URI(%q) = %v, want %v", tt.uri, got, tt.expected)
			}
		})
	}
}

func TestComputeS3MessageRefPath(t *testing.T) {
	ref, err := computeS3MessageRefPath("my-bucket", "Park Cleanup", domain.Welcome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "s3://my-bucket/park-cleanup/welcome.md"
	if ref != expected {
		t.Errorf("got %q, want %q", ref, expected)
	}

	ref, err = computeS3MessageRefPath("my-bucket", "Food Bank", domain.Reminder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "s3://my-bucket/food-bank/reminder.md"
	if ref != expected {
		t.Errorf("got %q, want %q", ref, expected)
	}
}
