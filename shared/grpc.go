package shared

import (
	"context"
	"fmt"

	"github.com/liuminhaw/mist-miner/proto"
)

// GRPCClient is an implementation of Greeter that talks over RPC
type GRPCClient struct {
	// client proto.GreetServiceClient
	client proto.MinerServiceClient
}

func (m *GRPCClient) Mine(config MinerConfig) (MinerResources, error) {
	fmt.Printf("GRPCClient Mine: %+v\n", config)
	resources, err := m.client.Mine(context.Background(), &proto.MinerConfig{Path: config.Path})
	if err != nil {
		return nil, err
	}

	// Convert proto resources to shared resources
	minerResources := MinerResources{}
	for _, resource := range resources.Resources {
		minerResource := MinerResource{
			Identifier: resource.Identifier,
			Properties: []MinerProperty{},
		}
		for _, data := range resource.Properties {
			minerResource.Properties = append(minerResource.Properties, MinerProperty{
				Type:  data.Type,
				Name:  data.Name,
				Value: data.Value,
			})
		}
		minerResources = append(minerResources, minerResource)
	}

	return minerResources, nil
}

// GRPCServer is the server that GRPCClient talks to
type GRPCServer struct {
	// This is the real implementation
	// Impl Greeter
	Impl Miner
}

func (m *GRPCServer) Mine(ctx context.Context, req *proto.MinerConfig) (*proto.MinerResources, error) {
	// func (m *GRPCServer) Mine(ctx context.Context, req *proto.NoParam) (*proto.MinerResources, error) {
	protoResources := []*proto.MinerResource{}

	resources, err := m.Impl.Mine(MinerConfig{Path: req.Path})
	fmt.Printf("Resources: %+v\n", resources)
	for _, resource := range resources {
		protoResource := proto.MinerResource{
			Identifier: resource.Identifier,
			Properties: []*proto.MinerProperty{},
		}
		for _, data := range resource.Properties {
			protoResource.Properties = append(protoResource.Properties, &proto.MinerProperty{
				Type:  data.Type,
				Name:  data.Name,
				Value: data.Value,
			})
		}
		protoResources = append(protoResources, &protoResource)
	}

	return &proto.MinerResources{
		Resources: protoResources,
	}, err
}
