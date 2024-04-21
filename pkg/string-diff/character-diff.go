package stringdiff

import (
	"bytes"
	"fmt"

	"github.com/andreyvit/diff"
	"github.com/fatih/color"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

var (
	red   = color.New(color.BgRed, color.FgBlack).SprintFunc()
	green = color.New(color.BgGreen, color.FgBlack).SprintFunc()
)

// GetPrettyDiffString gets the diff string of the 2 passed in values where removals have a red background and additions have a green background
func GetPrettyDiffString(original, new string) string {
	diffString := diff.CharacterDiff(original, new)

	var buff bytes.Buffer
	var diffsLen = len(diffString)
	var char, nextChar, nextNextChar, section string
	var inSection bool
	for i := 0; i < len(diffString); {
		char = string(diffString[i])
		if char == "(" && i+2 < diffsLen && !inSection {
			nextChar = string(diffString[i+1])
			nextNextChar = string(diffString[i+2])
			if nextChar == "+" && nextNextChar == "+" {
				inSection = true

				i += 3
				continue
			} else if nextChar == "~" && nextNextChar == "~" {
				inSection = true

				i += 3
				continue
			}
		} else if char == "~" && i+2 < diffsLen && string(diffString[i+1]) == "~" && string(diffString[i+2]) == ")" {
			inSection = false
			buff.WriteString(red(section))
			section = ""

			i += 3
			continue
		} else if char == "+" && i+2 < diffsLen && string(diffString[i+1]) == "+" && string(diffString[i+2]) == ")" {
			inSection = false
			buff.WriteString(green(section))
			section = ""

			i += 3
			continue
		}

		if inSection {
			section += char
		} else {
			buff.WriteString(char)
		}

		i++
	}

	displayString, err := repairLatin1(buff.String())
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to correct any latin1 bad characters: %s`, err))
	}

	return displayString
}

// From https://go.dev/play/p/dBrx_ZmrsMN and https://stackoverflow.com/questions/13510458/golang-convert-iso8859-1-to-utf8
// It seems that the latin1 text was not encoded correctly so we needed to look for any characters that were garbled and go
// ahead and fix those specific characters.
func repairLatin1(s string) (string, error) {
	buf := make([]byte, 0, len(s))
	for i, r := range s {
		if r > 255 {
			return "", fmt.Errorf("character %s at index %d is not part of latin1", string(r), i)
		}
		buf = append(buf, byte(r))
	}
	return string(buf), nil
}
