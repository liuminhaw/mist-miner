/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmlog

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/cmd/mmerr"
	"github.com/liuminhaw/mist-miner/tui"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var LogCmd = &cobra.Command{
	Use:          "log <group>",
	Short:        "Show mining result log of a group",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return mmerr.NewArgsError(
				mmerr.LogCmdType,
				fmt.Sprintf("accepts 1 args, received %d", len(args)),
			)
		}
		group := args[0]

		model, err := tui.InitLogModel(group, 0)
		if err != nil {
			return fmt.Errorf("log sub-command failed: %w", err)
		}

		if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
			return fmt.Errorf("log sub-command failed: %w", err)
		}

		return nil
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
