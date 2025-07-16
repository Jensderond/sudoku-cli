package ui

import "github.com/charmbracelet/lipgloss"

// UI styles
var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	LivesStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	TimerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	// Cell styles
	CursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")). // Bright cyan
		Bold(true)

	InitialCellStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	CorrectCellStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	IncorrectCellStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	HighlightedCellStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("45")). // Lighter cyan for highlighted
		Bold(true)
)
