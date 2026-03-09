package compare

import (
	"fmt"
	"strings"
)

type Difference struct {
	Message  string
	DiffType DifferenceType
}

func (d *Difference) ToDisplayText() string {
	var differenceType = "unknown"
	switch d.DiffType {
	case LikelyMismatch:
		differenceType = "Likely Mismatch"
	case DefiniteMismatch:
		differenceType = "Definite Mismatch"
	case WrappedLine:
		differenceType = "Wrapped"
	case PartiallyWrappedLine:
		differenceType = "Partially Wrapped"
	case Whitespace:
		differenceType = "Whitespace"
	case Line:
		differenceType = "Line Mismatch"
	}

	return fmt.Sprintf("[%s]: %s", differenceType, d.Message)
}

// Align PDF lines with HTML lines and detect explicit linebreaks vs wraps
func CompareLines(pdfLines, htmlLines []string) (differences []Difference) {
	if len(pdfLines) != len(htmlLines) {
		differences = append(differences, Difference{
			Message:  fmt.Sprintf("Line count mismatch for HTML and PDF file: expected %d but was %d", len(htmlLines), len(pdfLines)),
			DiffType: LikelyMismatch,
		})
	}

	var pdfIdx int
	for i, htmlLine := range htmlLines {
		if pdfIdx >= len(pdfLines) {
			remainingCount := len(htmlLines) - i
			lineText := "line"
			if remainingCount != 1 {
				lineText += "s"
			}

			differences = append(differences, Difference{
				Message:  fmt.Sprintf("Ran out of lines in the PDF to compare to the HTML: had %d %s to go", remainingCount, lineText),
				DiffType: DefiniteMismatch,
			})
			break
		}

		pdfLine := pdfLines[pdfIdx]
		if htmlLine == pdfLine { // the lines match, so we can continue to the next line
			pdfIdx++
			continue
		}

		// Check if the lines have wrapped between PDF and HTML.
		if strings.HasPrefix(htmlLine, pdfLine) {
			// Try to concatenate additional PDF lines to see if together they match the HTML line
			var (
				combined    = pdfLine
				nextIdx     = pdfIdx + 1
				wrapped     = false
				partialWrap = false
			)
			for nextIdx < len(pdfLines) {
				if strings.HasSuffix(combined, "-") {
					combined += pdfLines[nextIdx]
				} else {
					combined += " " + pdfLines[nextIdx]
				}

				if combined == htmlLine {
					differences = append(differences, Difference{
						Message:  fmt.Sprintf("HTML line %d matches across %d PDF lines: %q", i+1, nextIdx-pdfIdx+1, htmlLine),
						DiffType: WrappedLine,
					})
					pdfIdx = nextIdx + 1
					wrapped = true
					break
				}

				// If still a prefix, keep going; otherwise stop
				if !strings.HasPrefix(htmlLine, combined) {
					break
				}

				partialWrap = true

				nextIdx++
			}

			if wrapped {
				continue
			}

			if partialWrap {
				differences = append(differences, Difference{
					Message:  fmt.Sprintf("HTML line %d partially across %d PDF lines: %q", i+1, nextIdx-pdfIdx, htmlLine),
					DiffType: PartiallyWrappedLine,
				})
				pdfIdx = nextIdx
				continue
			}

			// No real further match other than start of line, so check the remaining output
		}

		// Check for single whitespace difference
		htmlNorm := strings.ReplaceAll(htmlLine, " ", "")
		pdfNorm := strings.ReplaceAll(pdfLine, " ", "")
		if htmlNorm == pdfNorm {
			differences = append(differences, Difference{
				Message:  fmt.Sprintf("Line %d vs. %d differs only by whitespace (HTML: %q | PDF: %q)", i+1, pdfIdx+1, htmlLine, pdfLine),
				DiffType: Whitespace,
			})
			pdfIdx++
			continue
		}

		// For now we will handle the following scenario as is, but in the future I may want to change this...
		// If none of the above, log as a mismatch
		differences = append(differences, Difference{
			Message:  fmt.Sprintf("Line %d does not match:\n  HTML: %q\n  PDF:  %q", i+1, htmlLine, pdfLine),
			DiffType: Line,
		})
		pdfIdx++
	}

	if pdfIdx < len(pdfLines) {
		remainingCount := len(pdfLines) - pdfIdx
		lineText := "line"
		if remainingCount != 1 {
			lineText += "s"
		}

		differences = append(differences, Difference{
			Message:  fmt.Sprintf("Ran out of lines in the HTML to compare to the PDF: had %d %s to go", remainingCount, lineText),
			DiffType: DefiniteMismatch,
		})
	}

	return
}
