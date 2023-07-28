package svg

import (
	"regexp"
	"unicode/utf8"
)

// Returns true if the given buffer is a valid SVG.
// Algo based on https://github.com/h2non/go-is-svg/blob/master/svg.go
func IsValidSVG(buf []byte) bool {
	var (
		htmlCommentRegex = regexp.MustCompile(`(?i)\<\!\-\-(?:.|\n|\r)*?-->`)
		svgRegex         = regexp.MustCompile(`(?i)^\s*(?:<\?xml[^>]*>\s*)?(?:<!doctype svg[^>]*>\s*)?<svg[^>]*>[^*]*<\/svg>\s*$`)
	)
	return !isBinary(buf) && svgRegex.Match(htmlCommentRegex.ReplaceAll(buf, []byte{}))
}

// isBinary checks if the given buffer is a binary file.
func isBinary(buf []byte) bool {
	if len(buf) < 24 {
		return false
	}
	for i := 0; i < 24; i++ {
		charCode, _ := utf8.DecodeRuneInString(string(buf[i]))
		if charCode == 65533 || charCode <= 8 {
			return true
		}
	}
	return false
}
