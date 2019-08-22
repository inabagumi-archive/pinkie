package util

import (
	"html"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\s*【[^】]+\s*\/\s*あにまーれ】?\s*`)

// NormalizeTitle returns a slice of the s with normalized
func NormalizeTitle(s string) string {
	s = re.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	s = html.UnescapeString(s)

	return s
}
