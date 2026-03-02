package converter

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	wsCollapse  = regexp.MustCompile(`\s+`)
	tocCollapse = regexp.MustCompile(`(.+?)  +(\d+)$`) // finds toc page numbers
)

// PdfTextCleanup Takes in the pdf text and then converts it to cleaned up text lines
// - combineN: if >0, combines the first N lines into a single line at the beginning of the result slice.
func PdfTextCleanup(pdfText string, combineN int, stripTocLineNums bool) []string {
	var (
		lines   = strings.Split(strings.ReplaceAll(pdfText, "\f", ""), "\n")
		cleaned []string
	)

	// Clean and filter lines
	for _, origLine := range lines {
		line := origLine
		if strings.TrimSpace(line) == "" {
			continue // skip blank
		}
		if _, err := strconv.Atoi(strings.TrimSpace(line)); err == nil {
			continue // skip page numbers
		}

		if len(line) > 3 && strings.HasPrefix(line, "    ") { // 4+ spaces
			line = wsCollapse.ReplaceAllString(line, " ")
		}

		line = strings.TrimLeft(line, " \t")

		// Remove any spaces between text and a trailing number (if two or more spaces)
		if m := tocCollapse.FindStringSubmatch(line); m != nil {
			if stripTocLineNums {
				line = m[1]
			} else {
				line = m[1] + m[2]
			}
		}

		cleaned = append(cleaned, line)
	}

	// Optionally combine first N lines into the first result line
	if combineN > 1 && len(cleaned) >= combineN {
		combined := strings.Join(cleaned[:combineN], " ")
		// Optionally collapse spaces in the combined line
		combined = wsCollapse.ReplaceAllString(combined, " ")
		cleaned = append([]string{combined}, cleaned[combineN:]...)
	}
	return cleaned
}
