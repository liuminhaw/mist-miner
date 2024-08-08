package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shelf"
)

type markItem struct {
	hash   string
	plugin string
}

func (i markItem) Title() string { return i.plugin }
func (i markItem) Description() string {
	return fmt.Sprintf("hash: %s", i.hash)
}

func (i markItem) FilterValue() string {
	return fmt.Sprintf("%s %s", i.hash, i.plugin)
}

type markModel struct {
	group  string
	hash   string
	list   list.Model
	width  int
	height int
}

func InitMarkModel(group, markHash string) (tea.Model, error) {
	list, err := readMarkItems(group, markHash)
	if err != nil {
		return nil, fmt.Errorf("InitMarkModel(%s, %s): %w", group, markHash, err)
	}

	model := markModel{
		hash:  markHash,
		group: group,
		list:  list,
	}
	model.list.Title = fmt.Sprintf("Mark: %s in group %s", markHash[:8], group)
	model.list.SetStatusBarItemName("entry", "entries")
	model.list.SetFilteringEnabled(true)
	model.list.DisableQuitKeybindings()

	return model, nil
}

func (m markModel) Init() tea.Cmd {
	return nil
}

func (m markModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "ctrl+z":
			logModel, _ := InitLogModel(m.group)
			return logModel.Update(tuiWindowSize)
		case "enter":
			selectedItem := m.list.SelectedItem().(markItem)
			resource, _ := InitResourceModel(m.group, selectedItem.hash, m.hash)
			return resource.Update(tuiWindowSize)
		}
	case markReadMsg:
		m.list.SetItems(msg.items)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m markModel) View() string {
	return docStyle.Render(m.list.View())
}

func readMarkItems(group, hash string) (list.Model, error) {
	mark, err := shelf.ReadMark(group, hash)
	if err != nil {
		return list.Model{}, fmt.Errorf("readMarkItems(%s, %s): %w", group, hash, err)
	}

	items := []list.Item{}
	for _, m := range mark.Mappings {
		items = append(items, markItem{hash: m.Hash, plugin: m.Module})
	}

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}

type markReadMsg struct {
	items []list.Item
	err   error
}
