package mmlog

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	hash      string
	timestamp time.Time
}

func (i item) Title() string { return i.hash }
func (i item) Description() string {
	return fmt.Sprintf("timestamp: %s", i.timestamp.Format("2006-01-02 15:04:05 -0700"))
}

func (i item) FilterValue() string {
	return fmt.Sprintf("%s %s", i.hash, i.timestamp.Format("2006-01-02 15:04:05 -0700"))
}

type model struct {
	list   list.Model
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
