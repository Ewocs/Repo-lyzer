package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00E5FF"))

	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 4)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF87")).
		Bold(true)

	NormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	InputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)

	SubtleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF87")).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB6C1")).
		Italic(true)

	DimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	// Theme styles
	DarkBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1, 4)

	LightBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CCCCCC")).
		Padding(1, 4).
		Foreground(lipgloss.Color("#000000"))

	HighlightStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Bold(true)

	MetricStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D7FF")).
		Bold(true)
)
