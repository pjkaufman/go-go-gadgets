package stringdiff

import (
	"bytes"
	"fmt"

	"github.com/andreyvit/diff"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	cliRed   = color.New(color.BgRed, color.FgBlack).SprintFunc()
	cliGreen = color.New(color.BgGreen, color.FgBlack).SprintFunc()
	tuiRed   = lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprint(color.FgBlack))).Background(lipgloss.Color(fmt.Sprint(color.BgRed)))
	tuiGreen = lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprint(color.FgBlack))).Background(lipgloss.Color(fmt.Sprint(color.BgGreen)))
)

// GetPrettyDiffString gets the diff string of the 2 passed in values where removals have a red background and additions have a green background
func GetPrettyDiffString(original, new string, isCli bool) (string, error) {
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

			if isCli {
				buff.WriteString(cliRed(section))
			} else {
				buff.WriteString(tuiRed.Render(section))
			}

			section = ""

			i += 3
			continue
		} else if char == "+" && i+2 < diffsLen && string(diffString[i+1]) == "+" && string(diffString[i+2]) == ")" {
			inSection = false
			if isCli {
				buff.WriteString(cliGreen(section))
			} else {
				buff.WriteString(tuiGreen.Render(section))
			}

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
		return "", fmt.Errorf(`failed to correct any latin1 bad characters: %w`, err)
	}

	return displayString, nil
}

// From https://go.dev/play/p/dBrx_ZmrsMN and https://stackoverflow.com/questions/13510458/golang-convert-iso8859-1-to-utf8
// It seems that the latin1 text was not encoded correctly so we needed to look for any characters that were garbled and go
// ahead and fix those specific characters.
func repairLatin1(s string) (string, error) {
	buf := make([]byte, 0, len(s))
	for i, r := range s {
		if r > 255 {
			return "", fmt.Errorf("character %q at index %d is not part of latin1", string(r), i)
		}
		buf = append(buf, byte(r))
	}
	return string(buf), nil
}
