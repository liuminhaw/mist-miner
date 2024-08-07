/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmlog

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/shelf"

	// "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var LogCmd = &cobra.Command{
	Use:   "log <group>",
	Short: "Show mining result log of a group",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]

		listModel, err := initialModel(group)
		if err != nil {
			fmt.Println("Failed to initialize model:", err)
			os.Exit(1)
		}
		listModel.Title = fmt.Sprintf("mining log for group %s", group)
		listModel.SetStatusBarItemName("entry", "entries")
		listModel.SetFilteringEnabled(true)

		m := model{
			logList: listModel,
			state:   logView,
		}
		m.logList.DisableQuitKeybindings()
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Println("Error running log command:", err)
			os.Exit(1)
		}
	},
}

func init() {
	// rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initialModel(group string) (list.Model, error) {
	items := []list.Item{}

	var err error
	head := shelf.RefMark{
		Name:  "HEAD",
		Group: group,
	}
	head.Reference, err = head.CurrentRef()
	if err != nil {
		return list.Model{}, fmt.Errorf("initial model: %w", err)
	}

	reference := string(head.Reference)
	for {
		mark, err := shelf.ReadMark(group, reference)
		if err != nil {
			return list.Model{}, fmt.Errorf("initial model: read mark %w", err)
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
