package tui

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/liuminhaw/mist-miner/shared"
	"github.com/liuminhaw/mist-miner/shelf"
)

type commitCache struct {
	head            shelf.RefMark
	labelMark       shelf.LabelMark
	groupIdHashMaps map[string]shelf.IdentifierHashMaps
	stuffOutline    shelf.StuffOutline
	isCached        bool
}

type commitSubmitModel struct {
	diaries  []commitDiaryItem
	cache    commitCache
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
}

var (
	currentDiaryNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle             = lipgloss.NewStyle().Margin(1, 2)
	checkMark             = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func InitCommitSubmitModel(items []list.Item) (tea.Model, error) {
	var diaryItems []commitDiaryItem
	for _, item := range items {
		diary, ok := item.(commitDiaryItem)
		if !ok {
			return nil, fmt.Errorf("InitCommitSubmitModel: item is not commitDiaryItem")
		}
		diaryItems = append(diaryItems, diary)
	}

	// Sort the diaries for update comparibility later
	slices.SortStableFunc(diaryItems, func(a, b commitDiaryItem) int {
		if a.group == b.group {
			if a.plugin == b.plugin {
				return strings.Compare(a.identifier, b.identifier)
			}
			return strings.Compare(a.plugin, b.plugin)
		}
		return strings.Compare(a.group, b.group)
	})

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return commitSubmitModel{
		diaries: diaryItems,
		cache: commitCache{
			groupIdHashMaps: make(map[string]shelf.IdentifierHashMaps),
		},
		spinner:  s,
		progress: p,
	}, nil
}

func (m commitSubmitModel) Init() tea.Cmd {
	return tea.Batch(updateDiaryLog(m.diaries[m.index], &m.cache), m.spinner.Tick)
}

func (m commitSubmitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case updateDiaryLogMsg:
		if msg.err != nil {
			return m, tea.Sequence(
				tea.Printf(msg.err.Error()),
				tea.Quit,
			)
		}

		// diary := m.diaries[m.index]
		if m.index >= len(m.diaries)-1 {
			// Everything's been processed
			m.done = true
			return m, tea.Sequence(
				// tea.Printf("%s Group: %s, Plugin: %s, Id: %s", checkMark, diary.group, diary.plugin, diary.Title()),
				tea.Printf(msg.msg),
				tea.Quit,
			)
		}

		// Update progress bar
		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.diaries)))
		return m, tea.Batch(
			progressCmd,
			// tea.Printf("%s %s", checkMark, diary.Title()),
			// tea.Printf("%s Group: %s, Plugin: %s, Id: %s", checkMark, diary.group, diary.plugin, diary.Title()),
			tea.Printf(msg.msg),
			updateDiaryLog(m.diaries[m.index], &m.cache),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}

	return m, nil
}

func (m commitSubmitModel) View() string {
	n := len(m.diaries)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! %d resource diaries committed.\n", n))
	}

	diaryCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)
	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+diaryCount))

	diaryTitle := currentDiaryNameStyle.Render(m.diaries[m.index].Title())
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Committing " + diaryTitle)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+prog+diaryCount+info))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + diaryCount
}

type updateDiaryLogMsg struct {
	msg string
	err error
}

// updateDiaryLog updates the diary log record of stuff diary and stuff outline
func updateDiaryLog(item commitDiaryItem, cache *commitCache) tea.Cmd {
	// Read and fill cache from log record if cache does not exist
	if !cache.isCached {
		head, err := shelf.NewRefMark(shelf.SHELF_MARK_FILE, item.group)
		if err != nil {
			return func() tea.Msg {
				return updateDiaryLogMsg{
					err: fmt.Errorf("Update diary: failed to get ref head"),
				}
			}
		}
		cache.head = head

		mark, err := shelf.ReadMark(item.group, string(head.Reference))
		if err != nil {
			return func() tea.Msg {
				return updateDiaryLogMsg{
					err: fmt.Errorf("Update diary: failed to read mark cache"),
				}
			}
		}
		cache.labelMark = *mark

		cache.isCached = true
	}

	var diary shelf.Diary
	var diaryLogger struct {
		stuffHash      string
		outlineHash    string
		mappingMatched bool
		idMatched      bool
		hashMap        struct {
			mapId  []string
			itemId []string
		}
	}
mappingsLoop:
	// for _, mapping := range mark.Mappings {
	for _, mapping := range cache.labelMark.Mappings {
		if mapping.Module != item.plugin {
			continue
		}

		diaryLogger.mappingMatched = true

		// var idHashMaps *shelf.IdentifierHashMaps
		key := fmt.Sprintf("%s_%s", item.group, mapping.Hash)
		if _, ok := cache.groupIdHashMaps[key]; !ok {
			// var err error
			idHashMaps, err := shelf.ReadIdentifierHashMaps(item.group, mapping.Hash)
			if err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to read identifier hash maps"),
					}
				}
			}
			cache.groupIdHashMaps[item.group] = *idHashMaps
		}

		for _, idHashMap := range cache.groupIdHashMaps[item.group].Maps {
			diaryLogger.hashMap.mapId = append(diaryLogger.hashMap.mapId, idHashMap.Identifier)
			diaryLogger.hashMap.itemId = append(diaryLogger.hashMap.itemId, item.identifier)

			if idHashMap.Identifier != item.identifier {
				continue
			}

			diaryLogger.idMatched = true

			// TODO: Update stuff diary and outline
			outline, err := shelf.ReadStuffOutline(item.group, idHashMap.Hash)
			if err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to read stuff outline"),
					}
				}
			}
			mDiary, err := readMinerDiary(item.group, outline.DiaryHash)
			if err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to read miner diary"),
					}
				}
			}

			staticTmpDiaryFile, err := shelf.NewDiaryStaticTempFile(item.filepath)
			if err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to create static temp diary"),
					}
				}
			}
			diary, err = staticTmpDiaryFile.WriteDiary()
			if err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to write static temp diary: %w", err),
					}
				}
			}

			// Update stuff diary record
			mDiaryNew := shared.NewMinerDiary(diary.Hash, string(cache.head.Reference), mDiary.Logs.Curr)
			diaryResource, err := shelf.NewStuff(item.group, &mDiaryNew)
			if err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to create diary resource: %w", err),
					}
				}
			}
			if _, err := diaryResource.Write(); err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to write diary resource: %w", err),
					}
				}
			}

			// Update stuff outline record
			newOutline := shelf.NewStuffOutline(item.group, outline.ResourceHash, diaryResource.Hash)
			if err := newOutline.Write(); err != nil {
				return func() tea.Msg {
					return updateDiaryLogMsg{
						err: fmt.Errorf("Update diary: failed to write stuff outline: %w", err),
					}
				}
			}

			diaryLogger.stuffHash = diaryResource.Hash
			diaryLogger.outlineHash = newOutline.Hash

			break mappingsLoop
		}
	}

	d := time.Millisecond * time.Duration(rand.Intn(300))
	return tea.Tick(d, func(time.Time) tea.Msg {
		return updateDiaryLogMsg{
			msg: fmt.Sprintf(
				"%s Diary %s committed, hash: %s\n  Diary stuff hash: %s\n  Stuff outline hash: %s,\n Mapping matched: %v,\n Id matched: %v,\n Map id: %+v,\n Item id: %+v\n",
				checkMark,
				item.Title(),
				diary.Hash,
				diaryLogger.stuffHash,
				diaryLogger.outlineHash,
				diaryLogger.mappingMatched,
				diaryLogger.idMatched,
				diaryLogger.hashMap.mapId,
				diaryLogger.hashMap.itemId,
			),
		}
	})

	// Get newest log record

	// Below code is for testing progress tui effect
	// d := time.Millisecond * time.Duration(rand.Intn(100))
	// return tea.Tick(d, func(time.Time) tea.Msg {
	// 	return updateDiaryLogMsg(diary.Title())
	// })
}
