// based on https://github.com/yorukot/superfile/blob/1b8a08d11ceca9b97af2fd531f806a29abc5560f/src/internal/ui/rendering/border.go

package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
)

type BorderConfig struct {
	// Optional info items at the bottom of the border
	infoItems []string
	// Including corners. Both should be >= 2
	width           int
	height          int
	titleLeftMargin int
}

func (b *BorderConfig) SetInfoItems(infoItems ...string) {
	for i := range infoItems {
		infoItems[i] = ansi.Strip(infoItems[i])
	}
	b.infoItems = infoItems
}

func (b *BorderConfig) AreInfoItemsTruncated() bool {
	cnt := len(b.infoItems)
	if cnt == 0 {
		return false
	}
	actualWidth := b.width - 2
	// border.MiddleLeft <content> border.MiddleRight border.Bottom
	availWidth := actualWidth/cnt - 3
	for i := range b.infoItems {
		if ansi.StringWidth(b.infoItems[i]) > availWidth {
			return true
		}
	}
	return false
}

// border.Top with something that takes up more than 1 runewidth will not work, so
// we only allow 1 runewidth for now, in the config. multiple things like
// border corner characters must be single rune, or else it would break rendering.
// This is all filled in one function to prevent passing around too many values
// in helper functions
func (b *BorderConfig) GetBorder(borderStrings lipgloss.Border) lipgloss.Border {
	res := borderStrings
	// excluding corners. Maybe we can move this to a utility function
	actualWidth := b.width - 2

	cnt := len(b.infoItems)
	// Minimum 4 character for each info item so that at least first character is rendered
	if cnt > 0 && actualWidth >= cnt*4 {
		// Max available width for each item's actual content
		// border.MiddleLeft <content> border.MiddleRight border.Bottom
		availWidth := actualWidth/cnt - 3
		infoText := ""
		for _, item := range b.infoItems {
			item = ansi.Truncate(item, availWidth, "")
			infoText += borderStrings.MiddleRight + item + borderStrings.MiddleLeft + borderStrings.Bottom
		}
		// Fill the rest with border char.
		remainingWidth := actualWidth - ansi.StringWidth(infoText)
		res.Bottom = strings.Repeat(borderStrings.Bottom, remainingWidth) + infoText
	}

	return res
}

func NewBorderConfig(height int, width int) BorderConfig {
	return BorderConfig{
		height:          height,
		width:           width,
		titleLeftMargin: 1,
	}
}
