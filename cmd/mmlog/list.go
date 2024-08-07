package mmlog

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type logItem struct {
	hash      string
	timestamp time.Time
}

func (i logItem) Title() string { return i.hash }
func (i logItem) Description() string {
	return fmt.Sprintf("timestamp: %s", i.timestamp.Format("2006-01-02 15:04:05 -0700"))
}

func (i logItem) FilterValue() string {
	return fmt.Sprintf("%s %s", i.hash, i.timestamp.Format("2006-01-02 15:04:05 -0700"))
}

type viewState int

const (
	logView viewState = iota
	pluginView
)

type prevListInfo struct {
	page          int
	index         int
	pageSize      int
	absoluteIndex int
}

type model struct {
	logList    list.Model
	state      viewState
	listMemory prevListInfo
    listsMemory map[int]list.Model
	detail     string
	width      int
	height     int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.logList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case logView:
				m.detail = "Change to plugin view"
				m.state = pluginView
				m.listMemory = prevListInfo{
					page:     m.logList.Paginator.Page,
					index:    m.logList.Index(),
					pageSize: m.logList.Paginator.PerPage,
				}
				m.listMemory.absoluteIndex = m.listMemory.page*m.listMemory.pageSize + m.listMemory.index
			case pluginView:
				m.state = logView
			}
		case "esc":
			m.state = logView
		}
	}

	var cmd tea.Cmd
	if m.state == logView {
		m.logList, cmd = m.logList.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case logView:
		return docStyle.Render(m.logList.View())
	case pluginView:
		msg := fmt.Sprintf(
            "Hash: %s\nPage: %d\nIndex: %d\nPageSize: %d\nAbsoluteIndex: %d\n",
            m.logList.SelectedItem().(logItem).hash,
			m.listMemory.page,
			m.listMemory.index,
			m.listMemory.pageSize,
			m.listMemory.absoluteIndex,

		)
		return docStyle.Render(msg)
	}

	return ""
}
