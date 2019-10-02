package util

import (
	"html"
)

// NormalizeTitle returns a slice of the s with normalized
func NormalizeTitle(s string) string {
	return html.UnescapeString(s)
}
