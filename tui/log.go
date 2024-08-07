package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shelf"
)

// logItem is a list item for showing log entries in list view
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

type logModel struct {
	group  string
	list   list.Model
	width  int
	height int
}

func InitLogModel(group string) (tea.Model, error) {
	list, err := readLogItems(group)
	if err != nil {
		return nil, fmt.Errorf("InitLogModel(%s): %w", group, err)
	}

	model := logModel{
		group: group,
		list:  list,
	}
	model.list.Title = fmt.Sprintf("Mined logs for group %s", group)
	model.list.SetStatusBarItemName("entry", "entries")
	model.list.SetFilteringEnabled(true)
	model.list.DisableQuitKeybindings()

	return model, nil
}

func (m logModel) Init() tea.Cmd {
	return nil
}

func (m logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
        tuiWindowSize = msg
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
        case "enter":
            selectedItem := m.list.SelectedItem().(logItem)
            mark, _ := InitMarkModel(m.group, selectedItem.hash)
            return mark.Update(tuiWindowSize)
			// case "esc":
			// 	m.state = logView
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m logModel) View() string {
	return docStyle.Render(m.list.View())
}

func readLogItems(group string) (list.Model, error) {
	items := []list.Item{}
	var err error
	head := shelf.RefMark{
		Name:  "HEAD",
		Group: group,
	}
	head.Reference, err = head.CurrentRef()
	if err != nil {
		return list.Model{}, fmt.Errorf("readLogItems: %w", err)
	}

	reference := string(head.Reference)
	for {
		mark, err := shelf.ReadMark(group, reference)
		if err != nil {
			return list.Model{}, fmt.Errorf("readLogItems: %w", err)
		}
		items = append(items, logItem{hash: mark.Hash, timestamp: mark.TimeStamp})

		if mark.Parent == "nil" {
			break
		} else {
			reference = mark.Parent
		}
	}

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}
