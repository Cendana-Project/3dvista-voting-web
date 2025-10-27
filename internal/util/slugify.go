package util

import (
	"regexp"
	"strings"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)
)

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with hyphens
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")

	// Trim leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}


