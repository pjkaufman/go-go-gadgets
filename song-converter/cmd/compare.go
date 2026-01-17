package cmd

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var (
	pdfFile, htmlFile string
	numJoinLines      int
	wsCollapse        = regexp.MustCompile(`\s+`)
	tocCollapse       = regexp.MustCompile(`(.+?)  +(\d+)$`) // finds toc page numbers
)

const (
	PdfPathArgEmpty     = "file must have a non-whitespace value"
	HtmlPathArgEmpty    = "source must have a non-whitespace value"
	PdfPathNotPdfFile   = "file must be a pdf file"
	HtmlPathNotHtmlFile = "source must be an html file"
)

// CompareCmd represents the Compare command
var CompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compares the provided HTML and PDF file to see if there are any potentially meaningful difference like linebreaks and whitespace differences",
	// Example: heredoc.Doc(`To write the output of converting the files in the specified folder to html to a file:
	// song-converter  -d working-dir -c cover.md -o songs.html

	// To write the output of converting the files in the specified folder to html to std out:
	// song-converter create-html -d working-dir -s cover.md
	// `),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCompareHtmlFlags(htmlFile, pdfFile)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(htmlFile, "source")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(pdfFile, "file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		pdfLines, err := pdfToTextCleaned(pdfFile, numJoinLines)
		if err != nil {
			log.Fatalf("PDF extraction error: %v", err)
		}

		htmlContent, err := filehandler.ReadInFileContents(htmlFile)
		if err != nil {
			logger.WriteError(err.Error())
		}

		htmlLines := extractTextLinesFromHTML(htmlContent)

		logger.WriteInfo("-- Alignment of PDF vs HTML lines --")
		detectMeaningfulLineDifferences(pdfLines, htmlLines)
	},
}

func init() {
	rootCmd.AddCommand(CompareCmd)

	CompareCmd.Flags().StringVarP(&htmlFile, "source", "s", "", "the html file that was used to generate the pdf file")
	err := CompareCmd.MarkFlagRequired("source")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"source\" as required on compare command: %v\n", err)
	}

	err = CompareCmd.MarkFlagFilename("source", "html")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"source\" as a looking for specific file types on compare command: %v\n", err)
	}

	CompareCmd.Flags().StringVarP(&pdfFile, "file", "f", "", "the pdf file to compare with the html file")
	err = CompareCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as required on compare command: %v\n", err)
	}

	err = CompareCmd.MarkFlagFilename("file", "pdf")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as a looking for specific file types on compare command: %v\n", err)
	}

	CompareCmd.Flags().IntVarP(&numJoinLines, "join-lines", "", 0, "the number of lines at the start of the pdf to join together to help make the html and pdf content as similar as possible")
}

func ValidateCompareHtmlFlags(htmlFilePath, pdfFilePath string) error {
	if strings.TrimSpace(htmlFilePath) == "" {
		return errors.New(HtmlPathArgEmpty)
	}

	if !strings.HasSuffix(htmlFilePath, ".html") {
		return errors.New(HtmlPathNotHtmlFile)
	}

	if strings.TrimSpace(pdfFilePath) == "" {
		return errors.New(PdfPathArgEmpty)
	}

	if !strings.HasSuffix(pdfFilePath, ".pdf") {
		return errors.New(PdfPathNotPdfFile)
	}

	return nil
}

// Extract visible lines from HTML (treat <p>, <br>, <div>, <li> as new lines)
func extractTextLinesFromHTML(source string) []string {
	// the parser only recognizes <br>, so make sure that is the only line break present
	source = strings.ReplaceAll(source, "<br/>", "<br>")

	var (
		lines            []string
		inMetadataDiv    bool
		metadataDivDepth int
		lineBuf          bytes.Buffer
		metadataBuf      bytes.Buffer
		flushMetadataBuf = func() {
			if metadataBuf.Len() > 0 {
				s := collapseWhitespace(metadataBuf.String())
				if s != "" {
					lines = append(lines, s)
				}

				metadataBuf.Reset()
			}
		}
	)

	tokenizer := html.NewTokenizer(strings.NewReader(source))
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			// On EOF, flush any remaining line or metadata
			flushMetadataBuf()
			if lineBuf.Len() > 0 {
				s := collapseWhitespace(lineBuf.String())
				if s != "" {
					lines = append(lines, s)
				}
			}
			return lines

		case html.TextToken:
			txt := html.UnescapeString(string(tokenizer.Text()))
			txt = strings.ReplaceAll(txt, "\u00A0", " ") // Convert &nbsp; (U+00A0) to space
			if inMetadataDiv {
				// adding the space makes sure that when we add the different metadata, it will have spaces between content even if there are no empty pieces of metadata between them
				metadataBuf.WriteString(" " + txt)
			} else {
				lineBuf.WriteString(txt)
			}

		case html.StartTagToken:
			tagName, hasAttr := tokenizer.TagName()
			tag := string(tagName)
			isMeta := false
			if tag == "div" && hasAttr {
				for {
					key, val, more := tokenizer.TagAttr()
					if string(key) == "class" && strings.Contains(string(val), "metadata") {
						isMeta = true
					}
					if !more {
						break
					}
				}
			} else if tag == "p" && inMetadataDiv {
				// if we hit a paragraph tag and we still think we are in the metadata, we are not, so we need to flush that data out...
				flushMetadataBuf()

				inMetadataDiv = false
				metadataDivDepth = 0
				isMeta = false
			}

			if isMeta {
				// Starting .metadata div
				if inMetadataDiv {
					flushMetadataBuf()
				}

				inMetadataDiv = true
				metadataDivDepth = 1
			} else if inMetadataDiv && tag == "div" {
				metadataDivDepth++
			} else if !inMetadataDiv && (tag == "br" || tag == "div" || tag == "li" || tag == "p") {
				s := collapseWhitespace(lineBuf.String())
				if s != "" {
					lines = append(lines, s)
				}
				lineBuf.Reset()
			}

		case html.EndTagToken:
			tagName, _ := tokenizer.TagName()
			tag := string(tagName)
			if tag == "div" && inMetadataDiv {
				metadataDivDepth--
				if metadataDivDepth == 0 {
					// Leaving .metadata block
					flushMetadataBuf()
					inMetadataDiv = false
				}
			} else if !inMetadataDiv && (tag == "div" || tag == "li" || tag == "p") {
				s := collapseWhitespace(lineBuf.String())
				if s != "" {
					lines = append(lines, s)
				}
				lineBuf.Reset()
			}
		}
	}
}

// Utility: Collapse whitespace, remove leading/trailing, squish internal whitespace
func collapseWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// Align PDF lines with HTML lines and detect explicit linebreaks vs wraps
// TODO: convert this into a function that will return a list of lines responses that then get transformed into strings for output
func detectMeaningfulLineDifferences(pdfLines, htmlLines []string) {
	if len(pdfLines) != len(htmlLines) {
		logger.WriteInfof("[Likely Mismatch]: Line count mismatch for HTML and PDF file: expected %d but was %d\n", len(htmlLines), len(pdfLines))
	}

	var pdfIdx int
	for i, htmlLine := range htmlLines {
		if pdfIdx >= len(pdfLines) {
			remainingCount := len(htmlLines) - i
			lineText := "line"
			if remainingCount != 1 {
				lineText += "s"
			}
			logger.WriteInfof("[Definite Mismatch]: Ran out of lines in the PDF to compare to the HTML: had %d %s to go\n", remainingCount, lineText)
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
					logger.WriteInfof("[Wrapped]: HTML line %d matches across %d PDF lines: %q\n", i+1, nextIdx-pdfIdx+1, htmlLine)
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
				logger.WriteInfof("[Partially Wrapped]: HTML line %d partially across %d PDF lines: %q\n", i+1, nextIdx-pdfIdx+1, htmlLine)
				pdfIdx = nextIdx
				continue
			}

			// No real further match other than start of line, so check the remaining output
		}

		// Check for single whitespace difference
		htmlNorm := strings.ReplaceAll(htmlLine, " ", "")
		pdfNorm := strings.ReplaceAll(pdfLine, " ", "")
		if htmlNorm == pdfNorm {
			logger.WriteInfof("[Whitespace]: Line %d vs. %d differs only by whitespace (HTML: %q | PDF: %q)\n", i+1, pdfIdx+1, htmlLine, pdfLine)
			pdfIdx++
			continue
		}

		// TODO: decide if the below is how I want the mismatch handled, but for now it should do

		// If none of the above, log as a mismatch
		logger.WriteInfof("[Line Mismatch]: Line %d does not match:\n  HTML: %q\n  PDF:  %q\n", i+1, htmlLine, pdfLine)
		pdfIdx++
	}

	if pdfIdx < len(pdfLines) {
		remainingCount := len(pdfLines) - pdfIdx
		lineText := "line"
		if remainingCount != 1 {
			lineText += "s"
		}
		logger.WriteInfof("[Definite Mismatch]: Ran out of lines in the HTML to compare to the PDF: had %d %s to go\n", remainingCount, lineText)
	}
}

// Extract and clean lines from PDF using pdftotext.
// - combineN: if >0, combines the first N lines into a single line at the beginning of the result slice.
func pdfToTextCleaned(pdfPath string, combineN int) ([]string, error) {
	// TODO: update this to actually be run the way I run other cli tools in this repo...
	out, err := exec.Command("pdftotext", "-layout", pdfPath, "-").Output()
	if err != nil {
		return nil, err
	}
	var (
		lines   = strings.Split(strings.ReplaceAll(string(out), "\f", ""), "\n")
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
			line = m[1] + m[2]
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
	return cleaned, nil
}
