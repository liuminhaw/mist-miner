/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmlog

import (
	"fmt"
	"os"

	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload <group>",
	Short: "Reload history log of a group and recreate reference logger files",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]

		if err := shelf.GenerateHistoryRecords(group, shelf.SHELF_HISTORY_LOGS_PER_PAGE); err != nil {
			fmt.Println("Error generating history records:", err)
			os.Exit(1)
		}
	},
}

func init() {
	LogCmd.AddCommand(reloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
