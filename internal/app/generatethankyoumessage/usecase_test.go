package generatethankyoumessage

import (
	"testing"
)

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Central Park Cleanup", "central-park-cleanup"},
		{"after-school tutoring", "after-school-tutoring"},
		{"food_bank_sort", "food-bank-sort"},
		{"Meals on Wheels!", "meals-on-wheels"},
		{"NYC---Cares Event", "nyc-cares-event"},
		{"  leading spaces  ", "leading-spaces"},
		{"", ""},
		{"single", "single"},
		{"Already-Kebab", "already-kebab"},
	}

	for _, tt := range tests {
		got := toKebabCase(tt.input)
		if got != tt.expected {
			t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
