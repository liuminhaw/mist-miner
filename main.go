package main

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
)

var CONFIG_PATH = "config.hcl"

func run(plugName, plugId string, logger hclog.Logger) error {
	// Setup logger
	// logger := hclog.New(&hclog.LoggerOptions{
	// 	Level:      hclog.Debug,
	// 	Output:     os.Stderr,
	// 	JSONFormat: true,
	// })
	binaryPath := fmt.Sprintf("./plugins/bin/%s", plugName)
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

	resources, err := miner.Mine(shared.MinerConfig{Path: CONFIG_PATH})
	if err != nil {
		return err
	}

	for _, resource := range resources {
		stuff, err := shelf.NewStuff(plugName, plugId, resource)
		if err != nil {
			return err
		}

		if err := stuff.Write(); err != nil {
			return err
		}
	}

	// b, err := json.Marshal(resources)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Resources: %s\n", string(b))

	return nil
}

func main() {
	// Set logger
	// Don't want to see the plugin logs.
	log.SetOutput(io.Discard)
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	// Read the config file
	hclConf, err := shared.ReadConfig(CONFIG_PATH)
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
}
