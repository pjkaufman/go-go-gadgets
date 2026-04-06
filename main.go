package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

// dimension (width, height): (96, 20)
const (
	originalIssueWidth  = 94
	originalIssueHeight = 20 // was initially 45, but I think 20 is a safe minimum
)

func main() {
	var model = simpleModel{
		texts: []string{
			`╭──────────────────────────────────────╮╭────────────────────────────────────────────────────╮
│🗎 Text/Short_Story_Vol_5.xhtml (7/13) ││"<p>That place—was the bedroom of the second        │
│§ Potential Thought Instances (9/9)   ││princess, Hertrauuda(<i>ヘルトラウダ)</i>.</p>"     │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
╰──────────────────────────────────────╯╰───────────────────────────────────┤👁 View├─┤40/45├─╯`,
			`╭──────────────────────────────────────╮╭────────────────────────────────────────────────────╮
│🗎 Text/Short_Story_Vol_5.xhtml (7/13) ││"<p>Sapling-chan｡ﾟ(<i>ﾟ´Д｀ﾟ)</i>ﾟ｡"Horrible!         │
│§ Potential Thought Instances (9/9)   ││This is too much! Everyone also want to hear my     │
│                                      ││voice aren't you? Aren't you!?"</p>"                │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
╰──────────────────────────────────────╯╰───────────────────────────────────┤👁 View├─┤39/45├─╯`,
			`╭──────────────────────────────────────╮╭────────────────────────────────────────────────────╮
│🗎 Text/Short_Story_Vol_5.xhtml (7/13) ││"<p>Sapling-chanヽ(<i>｀Д´#)</i>ﾉ"Actually the      │
│§ Potential Thought Instances (9/9)   ││extra story from the questionnaire of the           │
│                                      ││fourth volume should have me as the leading         │
│                                      ││actress! And yet everyone was asking, give us       │
│                                      ││Marie route~, like that!"</p>"                      │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
│                                      ││                                                    │
╰──────────────────────────────────────╯╰───────────────────────────────────┤👁 View├─┤38/45├─╯`,
		},
	}

	p := tea.NewProgram(&model)
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type simpleModel struct {
	ready         bool
	textIndex     int
	texts         []string
	width, height int
}

func (m simpleModel) Init() tea.Cmd {
	return nil
}

func (m simpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// general logic for handling keys here
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":

			return m, tea.Quit
		case "right":
			if m.textIndex < len(m.texts)-1 {
				m.textIndex++
			}

			return m, nil
		case "left":
			if m.textIndex > 0 {
				m.textIndex--
			}

			return m, nil
		}
	case tea.WindowSizeMsg:
		m.ready = true
		m.height = msg.Height
		m.width = msg.Width
	case error:
		return m, tea.Quit
	}

	return m, nil
}

func (m simpleModel) View() string {
	if m.ready {
		if m.width < originalIssueWidth || m.height < originalIssueHeight {
			return wordwrap.String(fmt.Sprintf("The width and height must be at least %d and %d respectively current width and height are %d and %d", originalIssueWidth, originalIssueHeight, m.width, m.height), m.width)
		}

		lines := strings.Split(m.texts[m.textIndex], "\n")
		// account for possible terminal length and width difference:
		widthPadding := m.width - originalIssueWidth
		for i := range lines {
			lines[i] = strings.Repeat(" ", widthPadding) + lines[i]
		}

		return strings.Join(lines, "\n")
	}

	return ""
}
