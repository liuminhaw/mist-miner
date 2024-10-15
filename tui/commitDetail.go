package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/liuminhaw/mist-miner/shelf"
)

type commitDetailModel struct {
	filepath        string
	identifier      string
	content         string
	ready           bool
	viewport        viewport.Model
	highPerformance bool

	prevModel tea.Model
}

func InitCommitDetailModel(staticTempPath string, prev tea.Model) (tea.Model, error) {
	staticTmp, err := shelf.NewDiaryStaticTempFile(staticTempPath)
	if err != nil {
		return nil, fmt.Errorf("InitCommitDetailModel: %w", err)
	}

	content, err := staticTmp.Read()
	if err != nil {
		return nil, fmt.Errorf("InitCommitDetailModel: %w", err)
	}

	return commitDetailModel{
		filepath:        staticTempPath,
		identifier:      staticTmp.Meta.Identifier,
		content:         content,
		ready:           false,
		prevModel:       prev,
		highPerformance: true,
	}, nil
}

func (m commitDetailModel) Init() tea.Cmd {
	return nil
}

func (m commitDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tuiWindowSize = msg

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

			glamourContent = detailDocStyle.Render(glamourContent)

			m.viewport.HighPerformanceRendering = m.highPerformance
			m.viewport.SetContent(glamourContent)
			m.ready = true
			// m.viewport.YPosition = headerHeight + 1
		} else {
			cmds = append(cmds, tea.ClearScrollArea, func() tea.Msg { return reloadDetailMsg{} })
			return m, tea.Sequence(cmds...)
		}

		if m.highPerformance {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "b":
			cmds = append(cmds, tea.ClearScrollArea, func() tea.Msg { return prevPageMsg{} })
			return m, tea.Batch(cmds...)
		}
	case prevPageMsg:
		return m.prevModel.Update(tuiWindowSize)
	case reloadDetailMsg:
		detail, _ := InitCommitDetailModel(m.filepath, m.prevModel)
		return detail.Update(tuiWindowSize)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m commitDetailModel) View() string {
	if !m.ready {
		return "\n Initializing..."
	}
	return detailDocStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.viewport.View(), m.footerView()),
	)
}

func (m commitDetailModel) headerView() string {
	title := detailTitleStyle.Render(fmt.Sprintf("Resource: %s", m.identifier))
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colorTitleBackground)).Margin(1, 0, 1, 0)
	line := style.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m commitDetailModel) footerView() string {
	info := detailFooterStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colorTitleBackground)).Margin(0, 0, 1, 0)
	line := style.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
