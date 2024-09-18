/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/liuminhaw/mist-miner/cmd/mmerr"
	"github.com/liuminhaw/mist-miner/locks"
	"github.com/liuminhaw/mist-miner/shared"
	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

// mineCmd represents the mine command
var mineCmd = &cobra.Command{
	Use:          "mine",
	Short:        "Mine for cloud services resources",
	Long:         ``,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return mmerr.NewArgsError(
				mmerr.MineCmdType,
				fmt.Sprintf("accepts no args, received %d", len(args)),
			)
		}

		// Set logger
		// Don't want to see the plugin logs.
		log.SetOutput(io.Discard)
		logger := hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Debug,
			Output:     os.Stderr,
			JSONFormat: true,
		})

		// Read the config file
		hclConf, err := shared.ReadConfig(configFile)
		if err != nil {
			return fmt.Errorf("failed to mine: %w", err)
		}

		// Create objects lock
		objFileLock, err := locks.NewLock("", locks.OBJECTS_LOCKFILE)
		if err != nil {
			return fmt.Errorf("failed to mine: %w", err)
		}
		if err := objFileLock.TryLock(); err != nil {
			if errors.Is(err, locks.ErrIsLocked) {
				return err
			}
			return fmt.Errorf("failed to mine: %w", err)
		}

		// Run plugins
		gLabels := make(groupLabels)
		for _, plug := range hclConf.Plugs {
			fmt.Printf("Plug Name: %s\n", plug.Name)
			fmt.Printf("Plug Group: %s\n", plug.Group)

			err := run(
				pluginModule{name: plug.Name, group: plug.Group, config: plug.GenMinerConfig()},
				&gLabels,
				logger,
			)
			if err != nil {
				return fmt.Errorf("failed to mine: %w", err)
			}
		}

		pointers := []shelf.HistoryPointer{}
		for group, label := range gLabels {
			fmt.Printf("Group: %s\n", group)
			if err := label.Update(); err != nil {
				return fmt.Errorf("failed to mine: %w", err)
			}

			fmt.Printf("Hash: %s\n", label.Hash)
			fmt.Printf("Parent: %s\n", label.Parent)
			for _, mapping := range label.Mappings {
				fmt.Printf("Module: %s, Hash: %s\n", mapping.Module, mapping.Hash)
			}

			// Get history pointers data in each group for later records write
			pointers = append(
				pointers,
				shelf.NewHistoryPointer(group, label.Parent, label.Hash),
			)
		}
		if err := objFileLock.Unlock(); err != nil {
			return fmt.Errorf("failed to mine: %w", err)
		}

		for _, pointer := range pointers {
			// Update history logs record
			if err := shelf.GenerateHistoryRecords(pointer.Group, shelf.SHELF_HISTORY_LOGS_PER_PAGE); err != nil {
				return fmt.Errorf("failed to mine: %w", err)
			}

			// Update history logs pointer
			// fmt.Printf("DEBUG: parent: %s, hash: %s\n", pointer.Parent, label.Hash)
			if err := pointer.WriteNextMap(); err != nil {
				return fmt.Errorf("failed to mine: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mineCmd)

	defaultConf, err := defaultConfigAbs()
	if err != nil {
		fmt.Printf("Error getting default config file: %s\n", err)
		os.Exit(1)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mineCmd.Flags().StringVarP(&configFile, "config", "c", defaultConf, "hcl.conf")
}

type pluginModule struct {
	name   string
	group  string
	config shared.MinerConfig
}

type groupLabels map[string]shelf.LabelMark

func run(pMod pluginModule, gLabel *groupLabels, logger hclog.Logger) error {
	pluginsBinDir, err := pluginsBinDirAbs()
	if err != nil {
		return err
	}
	binaryPath := fmt.Sprintf("%s/%s", pluginsBinDir, pMod.name)
	fmt.Printf("Binary Path: %s\n", binaryPath)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command(binaryPath),
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("miner_grpc")
	if err != nil {
		return err
	}

	// We should have a Greeter now
	miner := raw.(shared.Miner)

	resources, err := miner.Mine(pMod.config)
	if err != nil {
		return err
	}

	labelMap := shelf.IdentifierHashMaps{
		Group: pMod.group,
		Maps:  []shelf.IdentifierHashMap{},
	}
	for _, resource := range resources {
		resource.Sort()

		stuffResource, err := shelf.NewStuff(pMod.group, &resource)
		if err != nil {
			return err
		}
		if err := stuffResource.Write(); err != nil {
			return err
		}

		// TODO: diary resource fetch implementation
		tempDiary := shared.MinerDiary{}
		diaryResource, err := shelf.NewStuff(pMod.group, &tempDiary)
		if err != nil {
			return err
		}
		if err := diaryResource.Write(); err != nil {
			return err
		}

		outline := shelf.NewStuffOutline(pMod.group, stuffResource.Hash, diaryResource.Hash)
		if err := outline.Write(); err != nil {
			return err
		}

		labelMap.Maps = append(labelMap.Maps, shelf.IdentifierHashMap{
			Identifier: resource.Identifier,
			Alias:      resource.Alias,
			Hash:       outline.Hash,
		})

	}

	// Prevent from writing empty label map
	if len(labelMap.Maps) == 0 {
		fmt.Printf("No resources found in group %s with plugin %s\n", pMod.group, pMod.name)
		return nil
	}

	// TODO: Sort should be done within write to avoid forgetting
	labelMap.Sort()
	if err := labelMap.Write(); err != nil {
		return err
	}

	// Check if label mark with plugId (group) exists
	// If not exists, create a new label mark and update
	// If exists, update the existence label mark
	var labelMark *shelf.LabelMark
	if _, ok := (*gLabel)[pMod.group]; !ok {
		(*gLabel)[pMod.group] = shelf.LabelMark{}
		labelMark, err = shelf.NewMark(pMod.group, "mine")
		if err != nil {
			return err
		}
	} else {
		lm := (*gLabel)[pMod.group]
		labelMark = &lm
	}

	// Update labelMark to the groupLabels
	labelMark.AddMapping(pMod.name, labelMap.Hash)
	(*gLabel)[pMod.group] = *labelMark

	return nil
}
