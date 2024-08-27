/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmlog

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liuminhaw/mist-miner/tui"

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

		model, err := tui.InitLogModel(group, 0)
		if err != nil {
			fmt.Println("Error initializing log model:", err)
			os.Exit(1)
		}

		if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
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
