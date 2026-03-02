package converter

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// HtmlToText Extracts visible lines from HTML (treat <p>, <br>, <div>, <li> as new lines)
func HtmlToText(source string) []string {
	// the parser only recognizes <br>, so make sure that is the only line break present
	source = strings.ReplaceAll(source, "<br/>", "<br>")

	var (
		lines            []string
		inMetadataDiv    bool
		metadataDivDepth int
		lastFoundIdIndex = -1
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
			} else if tag == "a" && hasAttr {
				for {
					key, val, more := tokenizer.TagAttr()
					if string(key) == "href" {
						var (
							h1Text string
							id, _  = strings.CutPrefix(string(val), "#")
						)
						h1Text, lastFoundIdIndex = extractIdText(source, id, lastFoundIdIndex)
						if h1Text != "" {
							lineBuf.WriteString(h1Text)
						}
					}
					if !more {
						break
					}
				}
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
			} else if !inMetadataDiv && (tag == "div" || tag == "li" || tag == "p" || tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6") {
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

func extractIdText(source, id string, lastFoundIdIndex int) (string, int) {
	if lastFoundIdIndex > 0 {
		source = source[lastFoundIdIndex:]
	}
	index := strings.Index(source, fmt.Sprintf(`id="%s"`, id))
	if index == -1 {
		return "", lastFoundIdIndex
	}

	endOfOpeningEl := strings.Index(source[index:], ">")
	if endOfOpeningEl == -1 { // this should not happen, but just in case...
		return "", index + lastFoundIdIndex
	}

	endIndex := strings.Index(source[index+endOfOpeningEl:], "</h1>")
	if endIndex == -1 { // this should not happen, but just in case...
		return "", index + lastFoundIdIndex
	}

	var (
		startOfText = index + endOfOpeningEl + 1
		endOfText   = index + endOfOpeningEl + endIndex
	)

	var h1Text = source[startOfText:endOfText]
	// remove any wrapper around the title and subtitle of the heading
	h1Text = strings.Replace(h1Text, "</span>", "", 1)
	h1Text = strings.Replace(h1Text, `<span class="other-title">`, "", 1)

	return html.UnescapeString(h1Text), lastFoundIdIndex + endOfText + 5 // the end of the ending h1 el
}
