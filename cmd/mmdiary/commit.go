/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmdiary

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/cmd/mmerr"
	"github.com/liuminhaw/mist-miner/tui"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var CommitCmd = &cobra.Command{
	Use:          "commit <group>",
	Short:        "Commit udpated diary records to history",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return mmerr.NewArgsError(
				mmerr.DiaryCmdType,
				fmt.Sprintf("accepts 1 args, received %d", len(args)),
			)
		}
		group := args[0]

		model, err := tui.InitCommitDiaryModel(group)
		if err != nil {
			return fmt.Errorf("diary commit sub-command failed: %w", err)
		}

		if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("diary commit sub-command failed: %w", err)
		}

		return nil
	},
}

func init() {
	DiaryCmd.AddCommand(CommitCmd)
	// diaryCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
