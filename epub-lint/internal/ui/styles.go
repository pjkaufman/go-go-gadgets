package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle            = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CBACF7"))
	inactiveStyle         = lipgloss.NewStyle().Faint(true)
	activeStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#fab387"))
	headerBorderStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("12"))
	controlsStyle         = lipgloss.NewStyle().Faint(true).Bold(true)
	footerBorderStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("12"))
	fileNameStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	suggestionNameStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f5e0dc"))
	leftStatusBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("12"))
	suggestionBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true) //.BorderForeground(lipgloss.Color("#b4befe"))

	acceptedChangeTitleStyle = lipgloss.NewStyle().Bold(true)
	displayStyle             = lipgloss.NewStyle()
)
