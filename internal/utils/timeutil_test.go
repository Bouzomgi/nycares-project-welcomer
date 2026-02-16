package utils

import (
	"testing"
	"time"
)

func TestStringToDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		wantErr  bool
	}{
		{"2025-03-15", time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC), false},
		{"2024-01-01", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"invalid", time.Time{}, true},
		{"2025/03/15", time.Time{}, true},
		{"", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := StringToDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("StringToDate(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDateToString(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC), "2025-03-15"},
		{time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC), "2024-01-01"},
		{time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), "2025-12-31"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := DateToString(tt.input)
			if got != tt.expected {
				t.Errorf("DateToString(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	original := "2025-06-15"
	parsed, err := StringToDate(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := DateToString(parsed)
	if result != original {
		t.Errorf("round trip failed: %q -> %v -> %q", original, parsed, result)
	}
}
