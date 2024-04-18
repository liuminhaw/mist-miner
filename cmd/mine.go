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
		for _, plug := range hclConf.Plugs {
			fmt.Printf("Plug Name: %s\n", plug.Name)
			fmt.Printf("Plug Identity: %s\n", plug.Identity)
			if err := run(plug.Name, plug.Identity, logger); err != nil {
				fmt.Printf("Error running plugin: %+v\n", err)
				os.Exit(1)
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

func run(plugName, plugId string, logger hclog.Logger) error {
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
	binaryPath := fmt.Sprintf("%s/%s", pluginsBinDir, plugName)
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

	resources, err := miner.Mine(shared.MinerConfig{Path: configFile})
	if err != nil {
		return err
	}

	labelMap := shelf.IdentifierHashMaps{
		Module:   plugName,
		Identity: plugId,
		Maps:     []shelf.IdentifierHashMap{},
	}
	for _, resource := range resources {
		stuff, err := shelf.NewStuff(plugName, plugId, resource)
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

	labelMap.Sort()
	if err := labelMap.Write(); err != nil {
		return err
	}

	laberMark, err := shelf.NewMark(plugName, plugId, labelMap.Hash)
	if err != nil {
		return err
	}
	err = laberMark.Update()
	if err != nil {
		return err
	}

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
