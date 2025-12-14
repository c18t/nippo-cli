package tui

import "github.com/charmbracelet/lipgloss"

var (
	// SuccessStyle is used for success messages
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	// WarningStyle is used for warning messages
	WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	// InfoStyle is used for informational messages
	InfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	// SpinnerStyle is used for spinner components
	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	// PromptStyle is used for input prompts
	PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	// DimStyle is used for less prominent text
	DimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
