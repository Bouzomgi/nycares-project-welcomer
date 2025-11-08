package main

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	var words []string
	var current strings.Builder

	// Helper: check if rune is a separator
	isSeparator := func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	}

	for _, r := range s {
		if isSeparator(r) {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	if len(words) == 0 {
		return ""
	}

	// First word is all lowercase
	result := strings.ToLower(words[0])

	// Remaining words: capitalize first rune, lowercase rest
	for _, word := range words[1:] {
		if word == "" {
			continue
		}
		runes := []rune(word)
		result += string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
	}

	return result
}
