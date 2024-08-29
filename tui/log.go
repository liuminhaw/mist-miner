package tui

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/locks"
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
	group    string
	logIndex int
	list     list.Model
	width    int
	height   int
}

func InitLogModel(group string, logIdx int) (tea.Model, error) {
	list, err := readLogItems(group, logIdx)
	if err != nil {
		return nil, fmt.Errorf("InitLogModel(%s): %w", group, err)
	}

	model := logModel{
		group:    group,
		list:     list,
		logIndex: logIdx,
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
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem().(logItem)
			if strings.Contains(selectedItem.hash, shelf.SHELF_HISTORY_LOGS_PREV) {
				model, _ := InitLogModel(m.group, m.logIndex-1)
				logModel, _ := model.(logModel)
				logModel.list.Select(len(logModel.list.Items()) - 1)
				return logModel.Update(tuiWindowSize)
			}
			if strings.Contains(selectedItem.hash, shelf.SHELF_HISTORY_LOGS_NEXT) {
				logModel, _ := InitLogModel(m.group, m.logIndex+1)
				return logModel.Update(tuiWindowSize)
			}
			mark, _ := InitMarkModel(m.group, selectedItem.hash, m)
			return mark.Update(tuiWindowSize)
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m logModel) View() string {
	return listStyle.Render(m.list.View())
}

func readLogItems(group string, logIdx int) (list.Model, error) {
	items := []list.Item{}

	record, err := shelf.NewHistoryRecord(group, logIdx)
	if err != nil {
		return list.Model{}, fmt.Errorf("readLogItems(%s, %d): %w", group, logIdx, err)
	}

	fileLock, err := locks.NewLock(locks.HISTORY_LOCKFILE)
	if err != nil {
		return list.Model{}, fmt.Errorf("readLogItems(%s, %d): %w", group, logIdx, err)
	}
	locked, err := fileLock.TryRLock()
	if err != nil {
		return list.Model{}, fmt.Errorf("readLogItems(%s, %d): %w", group, logIdx, err)
	}
	defer fileLock.Unlock()

	if !locked {
		return list.Model{}, locks.ErrIsLocked
	}

	recordReader, err := record.Read()
	if err != nil {
		return list.Model{}, fmt.Errorf("readLogItems(%s, %d): %w", group, logIdx, err)
	}
	defer record.CloseFile()
	defer recordReader.Close()

	scanner := bufio.NewScanner(recordReader)
	for scanner.Scan() {
		recordFields := strings.Split(scanner.Text(), " ")
		recordHash := recordFields[0]
		recordTimestamp, err := time.Parse(time.RFC3339, recordFields[1])
		if err != nil {
			return list.Model{}, fmt.Errorf("readLogItems(%s, %d): %w", group, logIdx, err)
		}
		items = append(items, logItem{hash: recordHash, timestamp: recordTimestamp})
	}

	if err := scanner.Err(); err != nil {
		return list.Model{}, fmt.Errorf("readLogItems(%s, %d): %w", group, logIdx, err)
	}

	return list.New(items, list.NewDefaultDelegate(), 0, 0), nil
}
