package tui

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shelf"
)

type commitDiaryItem struct {
	filename   string
	filepath   string
	group      string
	plugin     string
	identifier string
}

func (i commitDiaryItem) Title() string { return i.identifier }
func (i commitDiaryItem) Description() string {
	return fmt.Sprintf("Group: %s, plugin: %s", i.group, i.plugin)
}

func (i commitDiaryItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.group, i.plugin, i.identifier)
}

// remove removes the commit diary item from the disk
func (i commitDiaryItem) remove() error {
	if err := os.Remove(i.filepath); err != nil {
		return fmt.Errorf("failed to remove commit diary item: %w", err)
	}
	return nil
}

type commitDiaryModel struct {
	list   list.Model
	width  int
	height int
}

func InitCommitDiaryModel(group string) (tea.Model, error) {
	list, err := readCommitDiaryItems(group)
	if err != nil {
		return nil, fmt.Errorf("init commit diary model: %w", err)
	}

	model := commitDiaryModel{
		list: list,
	}
	model.list.Title = "Diaries to commit"

	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.enter,
			customKeys.quit,
		}
	}
	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.enter,
			customKeys.quit,
		}
	}

	model.list.SetStatusBarItemName("diary", "diaries")
	model.list.SetFilteringEnabled(true)
	model.list.DisableQuitKeybindings()

	return model, nil
}

func (m commitDiaryModel) Init() tea.Cmd {
	return nil
}

func (m commitDiaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
		case "enter":
			selectedItem := m.list.SelectedItem().(commitDiaryItem)
			detail, _ := InitCommitDetailModel(selectedItem.filepath, m)
			return detail.Update(tuiWindowSize)
		case "c":
			submitModel, _ := InitCommitSubmitModel(m.list.Items())
			cmds = append(cmds, tea.ExitAltScreen, submitModel.Init())
			return submitModel, tea.Batch(cmds...)

			// Commit current updated diaries
			// On each selected identifier object
			// - Get current diary record
			// - Create a new diary record
			// - Update stuff outline content
			// - Update identifier hash maps content
			// Update label mark content after all identifier hash maps are updated
			// Create new log entry ()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m commitDiaryModel) View() string {
	return listStyle.Render(m.list.View())
}

func readCommitDiaryItems(group string) (list.Model, error) {
	items := []list.Item{}

	// groups, err := os.ReadDir(shelf.ShelfTempDiary())
	// if err != nil {
	// 	return list.Model{}, fmt.Errorf("read commit diary items: %v", err)
	// }

	// for _, group := range groups {
	plugins, err := os.ReadDir(filepath.Join(shelf.ShelfTempDiary(), group))
	if err != nil {
		return list.Model{}, fmt.Errorf("read commit diary items: %v", err)
	}

	for _, plugin := range plugins {
		tempDir := filepath.Join(shelf.ShelfTempDiary(), group, plugin.Name(), "static")
		files, err := os.ReadDir(tempDir)
		if err != nil {
			return list.Model{}, fmt.Errorf("read commit diary items: %v", err)
		}

		for _, file := range files {
			filenameWithoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			decodedBytes, err := base64.RawURLEncoding.DecodeString(filenameWithoutExt)
			if err != nil {
				return list.Model{}, fmt.Errorf(
					"read commit diary items: decode filename: %w",
					err,
				)
			}

			items = append(items, commitDiaryItem{
				filename:   file.Name(),
				filepath:   filepath.Join(tempDir, file.Name()),
				group:      group,
				plugin:     plugin.Name(),
				identifier: string(decodedBytes),
			})
		}
	}
	// }

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}
