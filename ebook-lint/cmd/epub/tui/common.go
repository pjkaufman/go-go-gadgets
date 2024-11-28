package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	groupStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	fileStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
	// inactiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// diffAddStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	// diffRemoveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)
