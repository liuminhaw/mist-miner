package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shelf"
)

type diaryItem struct {
	identifier     string
	identifierHash string
	alias          string
	hasNote        bool
}

func (i diaryItem) Title() string { return i.identifier }
func (i diaryItem) Description() string {
	var msg string
	if i.alias == "" {
		msg = fmt.Sprintf("hash: %s", i.identifierHash[:12])
	} else {
		msg = fmt.Sprintf("alias: %s, hash: %s", i.alias, i.identifierHash[:12])
	}

	return msg
}

func (i diaryItem) FilterValue() string {
	return fmt.Sprintf("%s %s %t", i.identifier, i.alias, i.hasNote)
}

type diaryModel struct {
	group  string
	plugin string
	list   list.Model
	width  int
	height int
}

func InitDiaryModel(group, plugin string) (tea.Model, error) {
	list, err := readDiaryItems(group, plugin)
	if err != nil {
		return nil, fmt.Errorf("InitDiaryModel(%s): %w", group, err)
	}

	model := diaryModel{
		group: group,
		list:  list,
	}
	model.list.Title = fmt.Sprintf("Diary for group %s", group)
	model.list.SetStatusBarItemName("resource", "resources")
	model.list.SetFilteringEnabled(true)
	model.list.DisableQuitKeybindings()

	return model, nil
}

func (m diaryModel) Init() tea.Cmd {
	return nil
}

func (m diaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// case "enter":
			// 	selectedItem := m.list.SelectedItem().(logItem)
			// 	if strings.Contains(selectedItem.hash, shelf.SHELF_HISTORY_LOGS_PREV) {
			// 		model, _ := InitLogModel(m.group, m.logIndex-1)
			// 		logModel, _ := model.(logModel)
			// 		logModel.list.Select(len(logModel.list.Items()) - 1)
			// 		return logModel.Update(tuiWindowSize)
			// 	}
			// 	if strings.Contains(selectedItem.hash, shelf.SHELF_HISTORY_LOGS_NEXT) {
			// 		logModel, _ := InitLogModel(m.group, m.logIndex+1)
			// 		return logModel.Update(tuiWindowSize)
			// 	}
			// 	mark, _ := InitMarkModel(m.group, selectedItem.hash, m)
			// 	return mark.Update(tuiWindowSize)
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m diaryModel) View() string {
	return listStyle.Render(m.list.View())
}

// readDiaryItems reads resources information from plugin in the group
// and returns a list model for tui view with the items.
func readDiaryItems(group, plugin string) (list.Model, error) {
	items := []list.Item{}

	head, err := shelf.NewRefMark(shelf.SHELF_MARK_FILE, group)
	if err != nil {
		return list.Model{}, fmt.Errorf("readDiaryItems(%s, %s): %w", group, plugin, err)
	}

	mark, err := shelf.ReadMark(group, string(head.Reference))
	if err != nil {
		return list.Model{}, fmt.Errorf("readDiaryItems(%s, %s): %w", group, plugin, err)
	}

	for _, mapping := range mark.Mappings {
		if mapping.Module != plugin {
			continue
		}
		idHashMaps, err := shelf.ReadIdentifierHashMaps(group, mapping.Hash)
		if err != nil {
			return list.Model{}, fmt.Errorf("readDiaryItems(%s, %s): %w", group, plugin, err)
		}

		for _, idHashMap := range idHashMaps.Maps {
			outline, err := shelf.ReadStuffOutline(group, idHashMap.Hash)
			if err != nil {
				return list.Model{}, fmt.Errorf("readDiaryItems(%s, %s): %w", group, plugin, err)
			}

			items = append(items, diaryItem{
				identifier:     idHashMap.Identifier,
				identifierHash: outline.ResourceHash,
				alias:          idHashMap.Alias,
			})

			// TODO: See if has note
		}
	}

	// TODO: Handle and show message if no matching plugin found

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}
