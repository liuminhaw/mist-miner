package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorTitleBackground = "62"
	colorTitleForeground = "230"
	colorMainText        = "#dddddd"
)

var (
	tuiWindowSize tea.WindowSizeMsg

	// Style for the list view
	listStyle = lipgloss.NewStyle().Margin(1, 2)

	// Style for resource detail view
	detailTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(colorTitleBackground)).
				Foreground(lipgloss.Color(colorTitleForeground)).
				Padding(0, 1).
				Margin(1, 0)

	detailFooterStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(colorTitleBackground)).
				Foreground(lipgloss.Color(colorTitleForeground)).
				Padding(0, 1).
				Margin(0, 0, 1, 0)
	detailDocStyle = lipgloss.NewStyle().Margin(0, 3).Foreground(lipgloss.Color(colorMainText))
)
