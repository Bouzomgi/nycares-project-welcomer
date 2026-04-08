package computemessage

import (
	"testing"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

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
		{"s3://ab/file", false}, // bucket name too short
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

	ref, err = computeS3MessageRefPath("my-bucket", "Any Project", domain.ThankYou)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "s3://my-bucket/thankyou.md"
	if ref != expected {
		t.Errorf("got %q, want %q", ref, expected)
	}
}
