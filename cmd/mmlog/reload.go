/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmlog

import (
	"fmt"

	"github.com/liuminhaw/mist-miner/cmd/mmerr"
	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

// reloadCmd represents the reload command
var ReloadCmd = &cobra.Command{
	Use:          "reload <group>",
	Short:        "Reload history log of a group to recreate reference logger files and pointer files",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return mmerr.NewArgsError(
				mmerr.LogReloadCmdType,
				fmt.Sprintf("accepts 1 args, received %d", len(args)),
			)
		}
		group := args[0]

		if err := shelf.GenerateHistoryRecords(group, shelf.SHELF_HISTORY_LOGS_PER_PAGE); err != nil {
			return fmt.Errorf("log reload sub-command failed: %w", err)
		}

		if err := shelf.GenerateHistoryPointers(group); err != nil {
			return fmt.Errorf("log reload sub-command failed: %w", err)
		}

		return nil
	},
}

func init() {
	LogCmd.AddCommand(ReloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
