package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	groupStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	fileStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
	// inactiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// diffAddStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	// diffRemoveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func clearScreen(s *strings.Builder) {
	s.WriteString("\r\033[K") // Clear current line
	s.WriteString("\033[2J")  // Clear entire screen
	s.WriteString("\033[H")   // Move cursor to top-left
}
