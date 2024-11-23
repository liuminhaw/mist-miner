package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shared"
	"github.com/liuminhaw/mist-miner/shelf"
)

type editorFinishedMsg struct {
	diary shelf.DiaryTempFile
	err   error
}

func openEditor(tempDiary shelf.DiaryTempFile) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, tempDiary.Path) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{
			diary: tempDiary,
			err:   err,
		}
	})
}

type diaryItem struct {
	identifier     string
	identifierHash string
	diaryHash      string
	alias          string
	hasNote        bool
}

func (i diaryItem) Title() string { return i.identifier }
func (i diaryItem) Description() string {
	var msg string
	if i.alias == "" {
		msg = fmt.Sprintf("hash: %s", i.identifierHash[:12])
	} else {
		msg = fmt.Sprintf("alias: %s, hash: %s", i.alias, i.diaryHash[:12])
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

	// editor diaryEditor

	err error
}

func InitDiaryModel(group, plugin string) (tea.Model, error) {
	list, err := readDiaryItems(group, plugin)
	if err != nil {
		return nil, fmt.Errorf("InitDiaryModel(%s): %w", group, err)
	}

	model := diaryModel{
		group:  group,
		plugin: plugin,
		list:   list,
	}
	model.list.Title = fmt.Sprintf("Diary for group %s", group)

	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.quit,
			customKeys.editor,
		}
	}
	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.quit,
			customKeys.editor,
		}
	}

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
		case "e":
			selectedItem := m.list.SelectedItem().(diaryItem)
			mDiary, err := readMinerDiary(m.group, selectedItem.identifierHash)
			if err != nil {
				m.err = err
			}

			// diary := shelf.NewDiary(m.group, m.plugin, mDiary.Hash)
			diary := shelf.NewDiary(
				m.group,
				m.plugin,
				selectedItem.identifier,
				selectedItem.alias,
				mDiary.Hash,
			)
			if diary.Exist() {
				// Get the diary record
				// Copy the diary record to a temporary file
				// Open the editor with the temporary file
			} else {
				// Create a new diary record in temporary file
				tempDiary, err := diary.NewTempFile()
				if err != nil {
					m.err = err
				}

				if !tempDiary.StaticExist() {
					if err := tempDiary.Init(); err != nil {
						m.err = err
					}
					defer tempDiary.Close()
				} else {
					if err := tempDiary.CopyFromStatic(); err != nil {
						m.err = err
					}
				}

				// Open the editor with the new diary record
				return m, openEditor(tempDiary)
			}

			// return m, openEditor()

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
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		if err := msg.diary.ToStaticTemp(); err != nil {
			m.err = err
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m diaryModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s\n", m.err)
	}
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
				identifier: idHashMap.Identifier,
				// identifierHash: outline.ResourceHash,
				identifierHash: outline.ResourceHash,
				diaryHash:      outline.DiaryHash,
				alias:          idHashMap.Alias,
			})

			// TODO: See if has note
		}
	}

	// TODO: Handle and show message if no matching plugin found

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}

// readMinerDiary reads the miner diary record from the shelf of given group and hash
func readMinerDiary(group, hash string) (shared.MinerDiary, error) {
	content, err := shelf.NewObjectRecord(group, hash).RecordRead()
	if err != nil {
		return shared.MinerDiary{}, fmt.Errorf("readMinerDiary(%s, %s): %w", group, hash, err)
	}

	diary := shared.MinerDiary{}
	if err := json.Unmarshal([]byte(content), &diary); err != nil {
		return shared.MinerDiary{}, fmt.Errorf("readMinerDiary(%s, %s): %w", group, hash, err)
	}

	return diary, nil
}
