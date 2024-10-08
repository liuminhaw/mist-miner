/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/liuminhaw/mist-miner/cmd/mmerr"
	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:          "cat-file <group> <hash>",
	Short:        "Display the content of given hash object",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return mmerr.NewArgsError(
				mmerr.CatFileCmdType,
				fmt.Sprintf("accepts 2 args, received %d", len(args)),
			)
		}
		group := args[0]
		hash := args[1]

		content, err := shelf.NewObjectRecord(group, hash).RecordRead()
		if err != nil {
			return fmt.Errorf("cat-file sub-command failed: %w", err)
		}
		fmt.Println(content)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
