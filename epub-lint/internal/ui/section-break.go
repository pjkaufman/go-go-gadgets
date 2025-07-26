package ui

func (m FixableIssuesModel) sectionBreakView() string {
	return m.sectionBreakInfo.input.View()
}

// func (m *FixableIssuesModel) HandleSectionBreakKeys(msg tea.Msg) tea.Cmd {
// 	var (
// 		cmd tea.Cmd
// 	)
// 	m.sectionBreakInput, cmd = m.sectionBreakInput.Update(msg)

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "enter":
// 			m.ContextBreak = strings.TrimSpace(m.sectionBreakInput.Value())
// 			if m.ContextBreak != "" {
// 				m.CurrentStage++

// 				err := m.setupForNextSuggestions()
// 				// TODO: handle better...
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			}
// 		}
// 	}

// 	return cmd
// }
