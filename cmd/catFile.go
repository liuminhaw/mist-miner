/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file GROUP HASH",
	Short: "Display the content of given hash object",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]
		hash := args[1]

		hashFile, err := shelf.ObjectFile(group, hash)
		if err != nil {
			fmt.Println("Error getting object file:", err)
			os.Exit(1)
		}

		if _, err := os.Stat(hashFile); os.IsNotExist(err) {
			fmt.Printf("Object with hash %s not found\n", hash)
			os.Exit(1)
		}

		f, err := os.Open(hashFile)
		if err != nil {
			fmt.Println("Error opening object file:", err)
			os.Exit(1)
		}
		defer f.Close()

		r, err := zlib.NewReader(f)
		if err != nil {
			fmt.Println("Error creating zlib reader:", err)
			os.Exit(1)
		}
		defer r.Close()

		b, err := io.ReadAll(r)
		if err != nil {
			fmt.Println("Error reading object content:", err)
			os.Exit(1)
		}

		var prettyJson bytes.Buffer
		if err := json.Indent(&prettyJson, b, "", "  "); err != nil {
			fmt.Println(string(b))
		} else {
			fmt.Println(prettyJson.String())
		}
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
