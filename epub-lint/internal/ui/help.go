package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
)

type footerKeyMap struct {
	PrevNextSuggestion key.Binding
	PrevNextIssueType  key.Binding
	PrevNextFile       key.Binding
	Edit               key.Binding
	Copy               key.Binding
	Accept             key.Binding
	Quit               key.Binding
	ExitWithoutSaving  key.Binding
	Reset              key.Binding
	Original           key.Binding
	CancelEdit         key.Binding
}

type footerKeys struct {
	Long  footerKeyMap
	Short footerKeyMap
}

var (
	masterKeyList = footerKeys{
		Long: footerKeyMap{
			PrevNextSuggestion: key.NewBinding(
				key.WithKeys("left/right"),
				key.WithHelp("←/→", "Previous/Next Suggestion"),
			),
			PrevNextIssueType: key.NewBinding(
				key.WithKeys("ctrl+u/ctrl+d"),
				key.WithHelp("^U/^D", "Previous/Next Issue Type"),
			),
			PrevNextFile: key.NewBinding(
				key.WithKeys("pgup/pgdn"),
				key.WithHelp("PgUp/PgDn", "Previous/Next File"),
			),
			Edit: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("E", "Edit"),
			),
			Copy: key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("C", "Copy"),
			),
			Accept: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("Enter", "Accept"),
			),
			Quit: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("Esc", "Quit"),
			),
			ExitWithoutSaving: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("Ctrl+C", "Exit without saving"),
			),
			Reset: key.NewBinding(
				key.WithKeys("ctrl+r"),
				key.WithHelp("Ctrl+R", "Reset"),
			),
			Original: key.NewBinding(
				key.WithKeys("ctrl+o"),
				key.WithHelp("Ctrl+O", "Original content"),
			),
			CancelEdit: key.NewBinding(
				key.WithKeys("ctrl+e"),
				key.WithHelp("Ctrl+E", "Cancel edit"),
			),
		},
		Short: footerKeyMap{
			PrevNextSuggestion: key.NewBinding(
				key.WithKeys("left/right"),
				key.WithHelp("←/→", "Suggestion"),
			),
			PrevNextIssueType: key.NewBinding(
				key.WithKeys("ctrl+u/ctrl+d"),
				key.WithHelp("^U/^D", "Issue Type"),
			),
			PrevNextFile: key.NewBinding(
				key.WithKeys("ctrl+pgup/ctrl+pgdn"),
				key.WithHelp("^PgUp/^PgDn", "File"),
			),
			Edit: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("E", "Edit"),
			),
			Copy: key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("C", "Copy"),
			),
			Accept: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("Enter", "Accept"),
			),
			Quit: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("Esc", "Quit"),
			),
			ExitWithoutSaving: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("Ctrl+C", "Exit"),
			),
			Reset: key.NewBinding(
				key.WithKeys("ctrl+r"),
				key.WithHelp("Ctrl+R", "Reset"),
			),
			Original: key.NewBinding(
				key.WithKeys("ctrl+o"),
				key.WithHelp("Ctrl+O", "Original"),
			),
			CancelEdit: key.NewBinding(
				key.WithKeys("ctrl+e"),
				key.WithHelp("Ctrl+E", "Cancel"),
			),
		},
	}
)

func (m FixableIssuesModel) footerView() string {
	var (
		keys     = masterKeyList.Long
		maxWidth = m.width - footerBorderStyle.GetHorizontalBorderSize()
	)
	if maxWidth < minLargeLayoutThreshold {
		keys = masterKeyList.Short
	}

	var (
		s    strings.Builder
		line strings.Builder
		help string
	)
	s.WriteString(fillLine(controlsStyle.Render("Controls:"), maxWidth) + "\n")
	for _, key := range m.currentFooterBindings(keys) {
		help = fmt.Sprintf("%s: %s", key.Help().Key, key.Help().Desc)
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

func (m FixableIssuesModel) currentFooterBindings(keys footerKeyMap) []key.Binding {
	switch m.currentStage {
	case sectionBreak:
		return []key.Binding{
			keys.Accept,
			keys.Quit,
			keys.ExitWithoutSaving,
		}

	case suggestionsProcessing:
		if m.PotentiallyFixableIssuesInfo.isEditing {
			return []key.Binding{
				keys.Reset,
				keys.Original,
				keys.CancelEdit,
				keys.Accept,
				keys.Quit,
				keys.ExitWithoutSaving,
			}
		}

		if m.PotentiallyFixableIssuesInfo.currentSuggestionState != nil &&
			m.PotentiallyFixableIssuesInfo.currentSuggestionState.isAccepted {

			return []key.Binding{
				keys.PrevNextSuggestion,
				keys.PrevNextIssueType,
				keys.PrevNextFile,
				keys.Copy,
				keys.Quit,
				keys.ExitWithoutSaving,
			}
		}

		return []key.Binding{
			keys.PrevNextSuggestion,
			keys.PrevNextIssueType,
			keys.PrevNextFile,
			keys.Edit,
			keys.Copy,
			keys.Accept,
			keys.Quit,
			keys.ExitWithoutSaving,
		}

	case stageCssSelection:
		return []key.Binding{
			keys.PrevNextSuggestion,
			keys.Accept,
			keys.ExitWithoutSaving,
		}
	}

	return nil
}
