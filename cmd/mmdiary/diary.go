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

// diaryCmd represents the diary command
var DiaryCmd = &cobra.Command{
	Use:          "diary <group> <plugin>",
	Short:        "Show notes of resources in a group",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return mmerr.NewArgsError(
				mmerr.DiaryCmdType,
				fmt.Sprintf("accepts 2 args, received %d", len(args)),
			)
		}
		group := args[0]
		plugin := args[1]

		model, err := tui.InitDiaryModel(group, plugin)
		if err != nil {
			return fmt.Errorf("diary sub-command failed: %w", err)
		}

		if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("diary sub-command failed: %w", err)
		}

		return nil
	},
}

func init() {
	// rootCmd.AddCommand(diaryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diaryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diaryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
