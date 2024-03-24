package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/liuminhaw/mist-miner/shared"
)

var CONFIG_PATH = "config.hcl"

func run(binaryPath string, logger hclog.Logger) error {
	// Setup logger
	// logger := hclog.New(&hclog.LoggerOptions{
	// 	Level:      hclog.Debug,
	// 	Output:     os.Stderr,
	// 	JSONFormat: true,
	// })

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
	// greeter := raw.(shared.Greeter)
	// message, err := greeter.SayHello()
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(message)
	// _, err = greeter.SayHelloServerStream()
	resources, err := miner.Mine(shared.MinerConfig{Path: CONFIG_PATH})
	if err != nil {
		return err
	}

	b, err := json.Marshal(resources)
	if err != nil {
		return err
	}
	fmt.Printf("Resources: %s\n", string(b))

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
        binaryPath := fmt.Sprintf("./plugins/bin/%s", plug.Name)
		fmt.Printf("Plug Name: %s\n", plug.Name)
		fmt.Printf("Plug executable: %s\n", binaryPath)
		if err := run(binaryPath, logger); err != nil {
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
