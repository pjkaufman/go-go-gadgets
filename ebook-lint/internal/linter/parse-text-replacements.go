package linter

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

// TODO: have this return an error
func ParseTextReplacements(text string) map[string]string {
	replaceValueToReplacement := make(map[string]string)

	var lines = strings.Split(text, "\n")
	var numLines = len(lines)
	if numLines <= 2 {
		return replaceValueToReplacement
	}

	// start after the markdown table header and divider lines
	var i = 2
	for i < numLines {
		var line = lines[i]
		i++
		var lineParts = strings.Split(line, "|")
		var numParts = len(lineParts)
		if numParts == 1 {
			continue
		} else if numParts != 4 {
			logger.WriteError(fmt.Sprintf("Could not parse %q because it does not have the proper amount of \"|\"s in it", line))
			continue
		}

		replaceValueToReplacement[strings.Trim(lineParts[1], " ")] = strings.Trim(lineParts[2], " ")
	}

	return replaceValueToReplacement
}
