package potentiallyfixableissue

import (
	"regexp"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

// this loop limit is meant to make sure that bad loops are ignored.
// If it needs to be more than 10, we can increase it. But for now, 10 works.
const maxQuoteLoops = 10

var unendedParagraphRegex = regexp.MustCompile(`((^|\n)[ \t]*<p[^>]*>)([^\n]*(Dr\.|Esq\.|Hon\.|Jr\.|Mr\.|Mrs\.|Ms\.|Messrs\.|Mmes\.|Msgr\.|Prof\.|Rev\.|Sr\.|St\.|Capt\.|Lt\.|Mt\.|Mtn\.|Gen\.|Sen\.|[a-zA-z,\d%–-—])["']?)( ?)(</p>\n)`)
var paragraphsWithDoubleQuotes = regexp.MustCompile(`((^|\n)[ \t]*<p[^>]*>)([^\n]*)(")([^\n]*)(</p>)`)
var paragraphsStartingWithLowercaseLetter = regexp.MustCompile(`((^|\n)[ \t]*<p[^>]*>)(\s*[a-z][^\n]*</p>)`)

func GetPotentiallyBrokenLines(fileContent string) (map[string]string, error) {
	var originalToSuggested = make(map[string]string)
	var parsedLines = map[string]struct{}{}

	parseUnendedParagraphs(fileContent, parsedLines, originalToSuggested)
	parseUnendedDoubleQuotes(fileContent, parsedLines, originalToSuggested)
	parseParagraphsStartingWithLowercaseLetters(fileContent, parsedLines, originalToSuggested)

	return originalToSuggested, nil
}

func parseUnendedParagraphs(fileContent string, parsedLines map[string]struct{}, originalToSuggested map[string]string) {
	var subMatches = unendedParagraphRegex.FindAllStringSubmatchIndex(fileContent, -1)
	if len(subMatches) == 0 {
		return
	}

	for _, groups := range subMatches {
		var currentLine = fileContent[groups[0]:groups[1]]
		if hasParsedLine(parsedLines, currentLine) {
			continue
		}

		addToParsedLines(parsedLines, currentLine)

		var (
			originalString  = currentLine
			suggestedString = fileContent[groups[2]:groups[3]] + fileContent[groups[6]:groups[7]] + " "
			nextLine        = currentLine
			endOfLineIndex  = groups[1]
		)
		for lineIsPotentiallyBroken := true; lineIsPotentiallyBroken; {
			nextLine = getNextLine(fileContent, endOfLineIndex)

			addToParsedLines(parsedLines, nextLine)
			originalString += nextLine

			var nextLineGroups = unendedParagraphRegex.FindAllStringSubmatchIndex(nextLine, 1)
			lineIsPotentiallyBroken = len(nextLineGroups) > 0
			if lineIsPotentiallyBroken {
				suggestedString += nextLine[nextLineGroups[0][6]:nextLineGroups[0][7]] + " "
				endOfLineIndex += len(nextLine)
			} else {
				var endOfOpeningTag = strings.Index(nextLine, ">")

				if endOfOpeningTag == -1 {
					suggestedString += nextLine
				} else {
					suggestedString += nextLine[endOfOpeningTag+1:]
				}
			}
		}

		// we included an ending newline character for the next lines that we pulled back
		// we do not need them when it comes to the ending of the original and suggested strings
		originalString = strings.TrimRight(originalString, "\n")
		suggestedString = strings.TrimRight(suggestedString, "\n")

		originalToSuggested[originalString] = suggestedString
	}
}

func parseUnendedDoubleQuotes(fileContent string, parsedLines map[string]struct{}, originalToSuggested map[string]string) {
	var subMatches = paragraphsWithDoubleQuotes.FindAllStringSubmatchIndex(fileContent, -1)
	if len(subMatches) == 0 {
		return
	}

	for _, groups := range subMatches {
		var currentLine = fileContent[groups[0]:groups[1]] + "\n"
		var doubleQuoteCount = strings.Count(currentLine, "\"")
		if doubleQuoteCount%2 == 0 {
			continue
		}

		// May need to handle parsed lines to make it so that it does not conflict between the two options that get parsed
		// but for now this should work just fine
		if hasParsedLine(parsedLines, currentLine) {
			continue
		}

		addToParsedLines(parsedLines, currentLine)

		var (
			originalString  = currentLine
			suggestedString = fileContent[groups[2]:groups[3]] + fileContent[groups[6]:groups[7]] + fileContent[groups[8]:groups[9]] + fileContent[groups[10]:groups[11]]
		)
		if !strings.HasSuffix(suggestedString, " ") {
			suggestedString += " "
		}

		var (
			i              = 1
			nextLine       = currentLine
			endOfLineIndex = groups[1] + 1
		)
		for lineIsPotentiallyBroken := true; lineIsPotentiallyBroken; {
			i += 1
			nextLine = getNextLine(fileContent, endOfLineIndex)
			addToParsedLines(parsedLines, nextLine)
			originalString += nextLine
			doubleQuoteCount += strings.Count(nextLine, "\"")
			endOfLineIndex += len(nextLine)

			lineIsPotentiallyBroken = doubleQuoteCount%2 != 0 && nextLine != "" && i < maxQuoteLoops

			var endOfOpeningTag = strings.Index(nextLine, ">")
			var lineContent = nextLine
			if endOfOpeningTag != -1 {
				lineContent = nextLine[endOfOpeningTag+1:]
			}

			if lineIsPotentiallyBroken {
				var startOfEndingTag = strings.LastIndex(lineContent, "<")

				if startOfEndingTag != -1 {
					lineContent = lineContent[0:startOfEndingTag]
				}
			}

			suggestedString += lineContent
		}

		// we included an ending newline character for the next lines that we pulled back
		// we do not need them when it comes to the ending of the original and suggested strings
		originalString = strings.TrimRight(originalString, "\n")
		suggestedString = strings.TrimRight(suggestedString, "\n")
		suggestedString = strings.ReplaceAll(suggestedString, "  ", " ")

		originalToSuggested[originalString] = suggestedString
	}
}

func parseParagraphsStartingWithLowercaseLetters(fileContent string, parsedLines map[string]struct{}, originalToSuggested map[string]string) {
	var subMatches = paragraphsStartingWithLowercaseLetter.FindAllStringSubmatchIndex(fileContent, -1)
	if len(subMatches) == 0 {
		return
	}

	for _, groups := range subMatches {
		var currentLine = fileContent[groups[0]:groups[1]] + "\n"

		// May need to handle parsed lines to make it so that it does not conflict between the two options that get parsed
		// but for now this should work just fine
		if hasParsedLine(parsedLines, currentLine) {
			continue
		}

		addToParsedLines(parsedLines, currentLine)
		var suggestedString = fileContent[groups[6]:groups[7]]
		if !strings.HasPrefix(suggestedString, " ") {
			suggestedString = " " + suggestedString
		}

		var previousLine = getPreviousLine(fileContent, groups[0])
		addToParsedLines(parsedLines, previousLine)

		var endingTagStart = strings.Index(previousLine, "</p>")
		// we cannot continue since the ending tag is missing, but we don't need to error out here
		if endingTagStart == -1 {
			logger.WriteWarnf(`failed to find ending paragraph tag for line %q\n`, previousLine)
			continue
		}

		var originalString = previousLine + currentLine
		suggestedString = previousLine[0:endingTagStart] + suggestedString

		// we included an ending newline character for the next lines that we pulled back
		// we do not need them when it comes to the ending of the original and suggested strings
		originalString = strings.TrimRight(originalString, "\n")
		suggestedString = strings.TrimRight(suggestedString, "\n")
		suggestedString = strings.ReplaceAll(suggestedString, "  ", " ")

		originalToSuggested[originalString] = suggestedString
	}
}

func hasParsedLine(parsedLines map[string]struct{}, line string) bool {
	var trimmedLine = strings.TrimSpace(line)
	_, alreadyParsed := parsedLines[trimmedLine]

	return alreadyParsed
}

func addToParsedLines(parsedLines map[string]struct{}, line string) {
	parsedLines[strings.TrimSpace(line)] = struct{}{}
}

func getNextLine(fileContent string, endOfLineIndex int) string {
	if endOfLineIndex == -1 {
		return ""
	}

	var substring = fileContent[endOfLineIndex:]
	var indexOfEndOfLine = strings.Index(substring, "\n")

	if indexOfEndOfLine == -1 {
		return substring
	}

	return substring[0 : indexOfEndOfLine+1]
}

func getPreviousLine(fileContent string, startOfCurrentLine int) string {
	if startOfCurrentLine == -1 {
		return ""
	}

	var startOfPreviousLine int
	for i := startOfCurrentLine - 1; i > 0; i-- {
		if fileContent[i:i+1] == "\n" {
			startOfPreviousLine = i
			break
		}
	}

	return fileContent[startOfPreviousLine:startOfCurrentLine]
}
