package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shelf"
)

type resourceItem struct {
	alias      string
	hash       string
	identifier string
}

func (i resourceItem) Title() string { return i.identifier }
func (i resourceItem) Description() string {
	if i.alias != "" {
		return fmt.Sprintf("hash: %s, alias: %s", i.hash[:12], i.alias)
	} else {
		return fmt.Sprintf("hash: %s", i.hash[:12])
	}
}

func (i resourceItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.identifier, i.alias, i.hash[:12])
}

type resourceModel struct {
	group  string
	hash   string
	list   list.Model
	width  int
	height int

	prevModel tea.Model
}

func InitResourceModel(group, resourceHash string, prev tea.Model) (tea.Model, error) {
	list, err := readResourceItems(group, resourceHash)
	if err != nil {
		return nil, fmt.Errorf("InitResourceModel(%s, %s): %w", group, resourceHash, err)
	}

	model := resourceModel{
		hash:      resourceHash,
		group:     group,
		list:      list,
		prevModel: prev,
	}
	model.list.Title = fmt.Sprintf("Plugin: %s in group %s", resourceHash[:8], group)
	model.list.SetStatusBarItemName("entry", "entries")
	model.list.SetFilteringEnabled(true)
	model.list.DisableQuitKeybindings()

	return model, nil
}

func (m resourceModel) Init() tea.Cmd {
	return nil
}

func (m resourceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tuiWindowSize = msg
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m.prevModel.Update(tuiWindowSize)
		case "enter":
			selectedItem := m.list.SelectedItem().(resourceItem)
			detail, _ := InitResourceDetailModel(m.group, selectedItem.hash, m)
			return detail.Update(tuiWindowSize)
		}
	case markReadMsg:
		m.list.SetItems(msg.items)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m resourceModel) View() string {
	return listStyle.Render(m.list.View())
}

func readResourceItems(group, hash string) (list.Model, error) {
	idHashMaps, err := shelf.ReadIdentifierHashMaps(group, hash)
	if err != nil {
		return list.Model{}, fmt.Errorf("readResourceItems(%s, %s): %w", group, hash, err)
	}

	items := []list.Item{}
	for _, m := range idHashMaps.Maps {
		outline, err := shelf.ReadStuffOutline(group, m.Hash)
		if err != nil {
			return list.Model{}, fmt.Errorf("readResourceItems(%s, %s): %w", group, hash, err)
		}

		items = append(
			items,
			resourceItem{alias: m.Alias, hash: outline.ResourceHash, identifier: m.Identifier},
		)
	}

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}
