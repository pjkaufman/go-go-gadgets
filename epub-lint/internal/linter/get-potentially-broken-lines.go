package linter

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

func GetPotentiallyBrokenLines(fileContent string) map[string]string {
	var originalToSuggested = make(map[string]string)
	var parsedLines = map[string]struct{}{}

	parseUnendedParagraphs(fileContent, parsedLines, originalToSuggested)
	parseUnendedDoubleQuotes(fileContent, parsedLines, originalToSuggested)
	parseParagraphsStartingWithLowercaseLetters(fileContent, parsedLines, originalToSuggested)

	return originalToSuggested
}

func parseUnendedParagraphs(fileContent string, parsedLines map[string]struct{}, originalToSuggested map[string]string) {
	var subMatches = unendedParagraphRegex.FindAllStringSubmatch(fileContent, -1)
	if len(subMatches) == 0 {
		return
	}

	for _, groups := range subMatches {
		var currentLine = groups[0]
		if hasParsedLine(parsedLines, currentLine) {
			continue
		}

		addToParsedLines(parsedLines, currentLine)

		var originalString = currentLine
		var suggestedString = groups[1] + groups[3] + " "
		var nextLine = currentLine
		for lineIsPotentiallyBroken := true; lineIsPotentiallyBroken; {
			nextLine = getNextLine(fileContent, nextLine)
			addToParsedLines(parsedLines, nextLine)
			originalString += nextLine

			var nextLineGroups = unendedParagraphRegex.FindStringSubmatch(nextLine)
			lineIsPotentiallyBroken = len(nextLineGroups) > 0
			if lineIsPotentiallyBroken {
				suggestedString += nextLineGroups[3] + " "
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
	var subMatches = paragraphsWithDoubleQuotes.FindAllStringSubmatch(fileContent, -1)
	if len(subMatches) == 0 {
		return
	}

	for _, groups := range subMatches {
		var currentLine = groups[0] + "\n"
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

		var originalString = currentLine
		var suggestedString = groups[1] + groups[3] + groups[4] + groups[5]
		if !strings.HasSuffix(suggestedString, " ") {
			suggestedString += " "
		}

		var i = 1
		var nextLine = currentLine
		for lineIsPotentiallyBroken := true; lineIsPotentiallyBroken; {
			i += 1
			nextLine = getNextLine(fileContent, nextLine)
			addToParsedLines(parsedLines, nextLine)
			originalString += nextLine
			doubleQuoteCount += strings.Count(nextLine, "\"")

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
	var subMatches = paragraphsStartingWithLowercaseLetter.FindAllStringSubmatch(fileContent, -1)
	if len(subMatches) == 0 {
		return
	}

	for _, groups := range subMatches {
		var currentLine = groups[0] + "\n"

		// May need to handle parsed lines to make it so that it does not conflict between the two options that get parsed
		// but for now this should work just fine
		if hasParsedLine(parsedLines, currentLine) {
			continue
		}

		addToParsedLines(parsedLines, currentLine)

		var suggestedString = groups[3]
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

func getNextLine(fileContent, currentLine string) string {
	var endOfLineIndex = strings.Index(fileContent, currentLine)
	if endOfLineIndex == -1 {
		return ""
	}

	endOfLineIndex += len(currentLine)

	var substring = fileContent[endOfLineIndex:]
	var indexOfEndOfLine = strings.Index(substring, "\n")

	if indexOfEndOfLine == -1 {
		return substring
	}

	return substring[0 : indexOfEndOfLine+1]
}

func getPreviousLine(fileContent, currentLine string) string {
	var startOfCurrentLine = strings.Index(fileContent, currentLine)
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
