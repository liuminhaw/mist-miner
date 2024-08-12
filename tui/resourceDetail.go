package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/liuminhaw/mist-miner/shared"
	"github.com/liuminhaw/mist-miner/shelf"
)

// Custom struct for tea messages
type prevPageMsg struct{}
// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
// var useHighPerformanceRenderer = true

type resourceDetailModel struct {
	group           string
	hash            string
	content         string
	ready           bool
	viewport        viewport.Model
	highPerformance bool

	prevModel tea.Model
}

func InitResourceDetailModel(
	group, hash string,
	prev tea.Model,
) (tea.Model, error) {
	content, err := shelf.ObjectRead(group, hash)
	if err != nil {
		return nil, fmt.Errorf("InitResourceDetailModel(%s, %s): %w", group, hash, err)
	}

	resource := shared.MinerResource{}
	if err := json.Unmarshal([]byte(content), &resource); err != nil {
		return nil, fmt.Errorf("InitResourceDetailModel(%s, %s): %w", group, hash, err)
	}
	resourceMd, err := resource.RenderMarkdown()
	if err != nil {
		return nil, fmt.Errorf("InitResourceDetailModel(%s, %s): %w", group, hash, err)
	}

	model := resourceDetailModel{
		group: group,
		hash:  hash,
		// content:   content,
		content:         resourceMd,
		ready:           false,
		prevModel:       prev,
		highPerformance: true,
	}

	return model, nil
}

func (m resourceDetailModel) Init() tea.Cmd {
	return nil
}

func (m resourceDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		h, v := detailDocStyle.GetFrameSize()

		if !m.ready {
			m.viewport = viewport.New(msg.Width-h, msg.Height-v-verticalMarginHeight)
			m.viewport.YPosition = headerHeight

			// Set glamour viewport content
			glamourContent := ""
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(m.viewport.Width-2),
			)
			if err != nil {
				m.highPerformance = false
				glamourContent = fmt.Sprintf("glamour renderer error: %s", err)
			} else {
				glamourContent, err = renderer.Render(m.content)
				if err != nil {
					m.highPerformance = false
					glamourContent = fmt.Sprintf("glamour render error: %s", err)
				}
			}

			// m.viewport.SetContent(m.content)
			m.viewport.HighPerformanceRendering = m.highPerformance
			m.viewport.SetContent(glamourContent)
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width - h
			m.viewport.Height = msg.Height - verticalMarginHeight - v
		}

		if m.highPerformance {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			cmds = append(cmds, tea.ClearScrollArea, func() tea.Msg { return prevPageMsg{} })
			return m, tea.Batch(cmds...)
		}
	case prevPageMsg:
		return m.prevModel.Update(tuiWindowSize)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m resourceDetailModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return detailDocStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.viewport.View(), m.footerView()),
	)
}

func (m resourceDetailModel) headerView() string {
	title := detailTitleStyle.Render(fmt.Sprintf("Resource: %s", m.hash[:12]))
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colorTitleBackground)).Margin(1, 0, 1, 0)
	line := style.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m resourceDetailModel) footerView() string {
	info := detailFooterStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colorTitleBackground)).Margin(0, 0, 1, 0)
	line := style.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

