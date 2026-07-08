package ui

import (
	"fmt"
	"strings"
)

type helpKeys struct {
	prevNextSuggestion helpKey
	prevNextIssueType  helpKey
	prevNextFile       helpKey
	edit               helpKey
	copy               helpKey
	accept             helpKey
	acceptEdit         helpKey
	quit               helpKey
	exitWithoutSaving  helpKey
	reset              helpKey
	original           helpKey
	cancelEdit         helpKey
}

type helpKey struct {
	keys  string
	short string
	long  string
}

var (
	masterKeyList = helpKeys{
		prevNextSuggestion: helpKey{
			keys:  "←/→",
			long:  "Previous/Next Suggestion",
			short: "Suggestion",
		},
		prevNextIssueType: helpKey{
			keys:  "^U/^D",
			long:  "Previous/Next Issue Type",
			short: "Issue Type",
		},
		prevNextFile: helpKey{
			keys:  "PgUp/PgDn",
			long:  "Previous/Next File",
			short: "File",
		},
		edit: helpKey{
			keys:  "E",
			long:  "Edit",
			short: "Edit",
		},
		copy: helpKey{
			keys:  "C",
			long:  "Copy",
			short: "Copy",
		},
		accept: helpKey{
			keys:  "Enter",
			long:  "Accept",
			short: "Accept",
		},
		acceptEdit: helpKey{
			keys:  "Ctrl+S",
			long:  "Accept",
			short: "Accept",
		},
		quit: helpKey{
			keys:  "Esc",
			long:  "Quit",
			short: "Quit",
		},
		exitWithoutSaving: helpKey{
			keys:  "Ctrl+C",
			long:  "Exit without saving",
			short: "Exit",
		},
		reset: helpKey{
			keys:  "Ctrl+R",
			long:  "Reset",
			short: "Reset",
		},
		original: helpKey{
			keys:  "Ctrl+O",
			long:  "Original content",
			short: "Original",
		},
		cancelEdit: helpKey{
			keys:  "Ctrl+E",
			long:  "Cancel edit",
			short: "Cancel",
		},
	}
)

func (m FixableIssuesModel) footerView() string {
	var (
		maxWidth = m.width - footerBorderStyle.GetHorizontalBorderSize()
		useShort = maxWidth < minLargeLayoutThreshold
	)

	var (
		s          strings.Builder
		line       strings.Builder
		help, desc string
	)
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), maxWidth))
	s.WriteString("\n")
	for _, key := range m.currentFooterBindings(masterKeyList) {
		desc = key.long
		if useShort {
			desc = key.short
		}

		help = fmt.Sprintf("%s: %s", key.keys, desc)
		if line.Len() == 0 {
			line.WriteString(help)
			s.WriteString(help)
		} else if line.Len()+len(help)+3 <= maxWidth {
			s.WriteString(controlsStyle.Render(" • "))
			line.WriteString(controlsStyle.Render(" • "))
			s.WriteString(help)
			line.WriteString(help)
		} else {
			s.WriteString("\n")
			line.Reset()

			line.WriteString(help)
			s.WriteString(help)
		}
	}

	return footerBorderStyle.Render(s.String())
}

func (m FixableIssuesModel) currentFooterBindings(keys helpKeys) []helpKey {
	switch m.currentStage {
	case sectionBreak:
		return []helpKey{
			keys.accept,
			keys.quit,
			keys.exitWithoutSaving,
		}

	case suggestionsProcessing:
		if m.PotentiallyFixableIssuesInfo.isEditing {
			return []helpKey{
				keys.reset,
				keys.original,
				keys.cancelEdit,
				keys.acceptEdit,
				keys.quit,
				keys.exitWithoutSaving,
			}
		}

		var currentSuggestion = m.PotentiallyFixableIssuesInfo.SuggestionManager.CurrentSuggestionState
		if currentSuggestion != nil && currentSuggestion.IsAccepted {
			return []helpKey{
				keys.prevNextSuggestion,
				keys.prevNextIssueType,
				keys.prevNextFile,
				keys.copy,
				keys.quit,
				keys.exitWithoutSaving,
			}
		}

		return []helpKey{
			keys.prevNextSuggestion,
			keys.prevNextIssueType,
			keys.prevNextFile,
			keys.edit,
			keys.copy,
			keys.accept,
			keys.quit,
			keys.exitWithoutSaving,
		}

	case stageCssSelection:
		return []helpKey{
			keys.prevNextSuggestion,
			keys.accept,
			keys.exitWithoutSaving,
		}
	}

	return nil
}
