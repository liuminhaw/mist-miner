package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    tuiWindowSize tea.WindowSizeMsg
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
)
