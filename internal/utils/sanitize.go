package utils

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	// StrictPolicy removes all HTML tags
	StrictPolicy = bluemonday.StrictPolicy()
	// UGCPolicy allows safe HTML commonly used in user-generated content
	UGCPolicy = bluemonday.UGCPolicy()
)

// SanitizeString removes all HTML tags and escapes special characters
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}
	// First remove HTML tags
	sanitized := StrictPolicy.Sanitize(input)
	// Then escape HTML entities
	sanitized = html.EscapeString(sanitized)
	// Trim whitespace
	return strings.TrimSpace(sanitized)
}

// SanitizeURL sanitizes URL strings, keeping only safe characters
func SanitizeURL(input string) string {
	if input == "" {
		return ""
	}
	// Remove HTML tags
	sanitized := StrictPolicy.Sanitize(input)
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

// SanitizeDescription allows safe HTML in descriptions but removes dangerous content
func SanitizeDescription(input *string) *string {
	if input == nil || *input == "" {
		return input
	}
	sanitized := UGCPolicy.Sanitize(*input)
	sanitized = strings.TrimSpace(sanitized)
	return &sanitized
}
