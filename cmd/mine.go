/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/liuminhaw/mist-miner/shared"
	"github.com/liuminhaw/mist-miner/shelf"
	"github.com/spf13/cobra"
)

// mineCmd represents the mine command
var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "Mine for cloud services resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Printf("Failed to read config file: %+v\n", err)
			os.Exit(1)
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
				fmt.Printf("Error running plugin: %+v\n", err)
				os.Exit(1)
			}
		}

		for group, label := range gLabels {
			fmt.Printf("Group: %s\n", group)
			if err := label.Update(); err != nil {
				fmt.Printf("Error updating label mark: %+v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Hash: %s\n", label.Hash)
			fmt.Printf("Parent: %s\n", label.Parent)
			for _, mapping := range label.Mappings {
				fmt.Printf("Module: %s, Hash: %s\n", mapping.Module, mapping.Hash)
			}
		}

		// if err := run(); err != nil {
		// 	fmt.Printf("error: %+v\n", err)
		// 	os.Exit(1)
		// }

		os.Exit(0)
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
	// Setup logger
	// logger := hclog.New(&hclog.LoggerOptions{
	// 	Level:      hclog.Debug,
	// 	Output:     os.Stderr,
	// 	JSONFormat: true,
	// })
	pluginsBinDir, err := pluginsBinDirAbs()
	if err != nil {
		return err
	}
	binaryPath := fmt.Sprintf("%s/%s", pluginsBinDir, pMod.name)
	fmt.Printf("Binary Path: %s\n", binaryPath)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		// Cmd:              exec.Command("sh", "-c", os.Getenv("PLUGIN_BINARY")),
		// Cmd:              exec.Command(os.Getenv("PLUGIN_BINARY")),
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
		stuff, err := shelf.NewStuff(pMod.group, resource)
		if err != nil {
			return err
		}

		identifier, err := stuff.ResourceIdentifier()
		if err != nil {
			return err
		}

		labelMap.Maps = append(labelMap.Maps, shelf.IdentifierHashMap{
			Identifier: identifier,
			Hash:       stuff.Hash,
		})

		if err := stuff.Write(); err != nil {
			return err
		}
	}

	// Prevent from writing empty label map
	if len(labelMap.Maps) == 0 {
		fmt.Printf("No resources found in group %s with plugin %s\n", pMod.group, pMod.name)
		return nil
	}

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
		labelMark, err = shelf.NewMark(pMod.name, pMod.group, labelMap.Hash)
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

	// err = labelMark.Update()
	// if err != nil {
	// 	return err
	// }

	// for _, lm := range labelMap.Maps {
	// 	fmt.Printf("Hash: %s, Identifier: %s\n", lm.Hash, lm.Identifier)
	// }

	// b, err := json.Marshal(resources)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Resources: %s\n", string(b))

	return nil
}
