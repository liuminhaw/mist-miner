package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/liuminhaw/mist-miner/proto"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by pluginlugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MINER",
	MagicCookieValue: "mining-elf",
}

// PluginMap is the map of plugins we can dispense
var PluginMap = map[string]plugin.Plugin{
	"miner_grpc": &MinerGRPCPlugin{},
}

// Miner is the interface that we're exposing as a plugin
type Miner interface {
	Mine(MinerConfig) (MinerResources, error)
}

type MinerGRPCPlugin struct {
	plugin.Plugin
	// Miner concreate implementation
	Impl Miner
}

func (p *MinerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterMinerServiceServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *MinerGRPCPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &GRPCClient{client: proto.NewMinerServiceClient(c)}, nil
}
