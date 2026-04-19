package logger

import "charm.land/lipgloss/v2"

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Red)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Yellow)
)
