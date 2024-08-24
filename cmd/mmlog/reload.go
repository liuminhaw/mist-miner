/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mmlog

import (
	"compress/zlib"
	"fmt"
	"os"

	"github.com/gofrs/flock"
	"github.com/liuminhaw/mist-miner/locks"
	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

const LOGS_PER_PAGE = 1000

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload <group>",
	Short: "Reload history log of a group and recreate reference logger files",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]

		var err error
		head := shelf.RefMark{
			Name:  "HEAD",
			Group: group,
		}
		// TODO: Implement Read flock in CurrentRef function
		head.Reference, err = head.CurrentRef()
		if err != nil {
			fmt.Println("Error reading current reference:", err)
			os.Exit(1)
		}

		// Use flock to prevent other processes from writing to the file
		lockFilePath, err := locks.FilePath(locks.HISTORY_LOCK_FILE)
		if err != nil {
			fmt.Println("Error getting lock file path:", err)
			os.Exit(1)
		}
		lockFile := flock.New(lockFilePath)
		locked, err := lockFile.TryLock()
		if err != nil {
			fmt.Println("Error acquiring file lock:", err)
			os.Exit(1)
		}
		defer lockFile.Unlock()

		if !locked {
			fmt.Println(
				"File is locked, another process is writing to it, wait and try again later.",
			)
			os.Exit(1)
		}

		historyDir, err := shelf.HistoryDir(group)
		if err != nil {
			fmt.Println("Error getting history directory:", err)
			os.Exit(1)
		}
		if err := os.RemoveAll(historyDir); err != nil {
			fmt.Println("Error removing history directory:", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(historyDir, os.ModePerm); err != nil {
			fmt.Println("Error creating history directory:", err)
			os.Exit(1)
		}

		filesIdx := 0
		reference := string(head.Reference)
	filesLoop:
		for {
			file, err := shelf.NewHistoryFile(group, filesIdx)
			if err != nil {
				fmt.Println("Error creating history file:", err)
				os.Exit(1)
			}
			defer file.Close()

			w := zlib.NewWriter(file)
			defer w.Close()

			for i := 0; i < LOGS_PER_PAGE; i++ {
				mark, err := shelf.ReadMark(group, reference)
				if err != nil {
					fmt.Println("Error reading mark:", err)
					os.Exit(1)
				}

				_, err = w.Write([]byte(fmt.Sprintf("%s\n", mark.Hash)))
				if err != nil {
					fmt.Println("Error writing mark hash to file:", err)
					os.Exit(1)
				}

				if mark.Parent == "nil" {
					break filesLoop
				} else if i == LOGS_PER_PAGE-1 {
					_, err = file.Write([]byte("more...\n"))
					if err != nil {
						fmt.Println("Error writing mark hash to file:", err)
						os.Exit(1)
					}
					reference = mark.Parent
				} else {
					reference = mark.Parent
				}
			}
			filesIdx++
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
