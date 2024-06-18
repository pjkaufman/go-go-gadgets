package linter

import (
	"fmt"
	"strings"
)

func ParseTextReplacements(text string) (map[string]string, error) {
	replaceValueToReplacement := make(map[string]string)

	var lines = strings.Split(text, "\n")
	var numLines = len(lines)
	if numLines <= 2 {
		return replaceValueToReplacement, nil
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
			return nil, fmt.Errorf("could not parse %q because it does not have the proper amount of \"|\"s in it", line)
		}

		replaceValueToReplacement[strings.Trim(lineParts[1], " ")] = strings.Trim(lineParts[2], " ")
	}

	return replaceValueToReplacement, nil
}
