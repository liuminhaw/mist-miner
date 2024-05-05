/*
Copyright © 2024 Min-Haw, Liu liuminhaw@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/liuminhaw/mist-miner/cmd/mmlog"
)

const (
	defaultConfigFile    = "config.hcl"
	defaultPluginsBinDir = "plugins/bin"
)

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mist-miner",
	Short: "Fetches and stores resources record from cloud services.",
	Long:  `Using customizable plugins to fetch and store resources record from cloud services.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(mmlog.LogCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mist-miner.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// defaultConfigAbs returns the absolute path of the default config file,
// and an error if any occurs.
func defaultConfigAbs() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("default config abs: %w", err)
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, defaultConfigFile), nil
}

// pluginsBinDirAbs returns the absolute path of the plugins bin directory
func pluginsBinDirAbs() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("plugins bin dir abs: %w", err)
	}

	return filepath.Join(filepath.Dir(execPath), defaultPluginsBinDir), nil
}
