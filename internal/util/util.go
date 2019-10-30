package util // import "github.com/inabagumi/ytc/v2/internal/util"

import (
	"html"
)

// NormalizeTitle returns a slice of the s with normalized
func NormalizeTitle(s string) string {
	return html.UnescapeString(s)
}
