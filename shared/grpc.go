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
	resources, err := m.client.Mine(context.Background(), toProtoMinerConfig(config))
	if err != nil {
		return nil, err
	}

	// Convert proto resources to shared resources
	minerResources := MinerResources{}
	for _, resource := range resources.Resources {
		minerResource := MinerResource{
			Identifier: resource.Identifier,
			Alias:      resource.Alias,
			Properties: []MinerProperty{},
		}
		for _, data := range resource.Properties {
			minerResource.Properties = append(minerResource.Properties, MinerProperty{
				Type: data.Type,
				Label: MinerPropertyLabel{
					Name:   data.Label.Name,
					Unique: data.Label.Unique,
				},
				Content: MinerPropertyContent{
					Format: data.Content.Format,
					Value:  data.Content.Value,
				},
			})
		}
		minerResources = append(minerResources, minerResource)
	}

	return minerResources, nil
}

func toProtoMinerConfig(config MinerConfig) *proto.MinerConfig {
	equipments := []*proto.MinerConfigEquipment{}
	for _, equipment := range config.Equipments {
		equipments = append(equipments, &proto.MinerConfigEquipment{
			Type:       equipment.Type,
			Name:       equipment.Name,
			Attributes: equipment.Attributes,
		})
	}

	return &proto.MinerConfig{
		Auth:       config.Auth,
		Equipments: equipments,
	}
}

// GRPCServer is the server that GRPCClient talks to
type GRPCServer struct {
	// This is the real implementation
	// Impl Greeter
	Impl Miner
}

func (m *GRPCServer) Mine(
	ctx context.Context,
	req *proto.MinerConfig,
) (*proto.MinerResources, error) {
	// func (m *GRPCServer) Mine(ctx context.Context, req *proto.NoParam) (*proto.MinerResources, error) {
	protoResources := []*proto.MinerResource{}

	resources, err := m.Impl.Mine(toSharedMinerConfig(req))
	fmt.Printf("Resources: %+v\n", resources)

	// Convert shared resources to proto resources
	for _, resource := range resources {
		protoResource := proto.MinerResource{
			Identifier: resource.Identifier,
			Alias:      resource.Alias,
			Properties: []*proto.MinerProperty{},
		}
		for _, data := range resource.Properties {
			protoResource.Properties = append(protoResource.Properties, &proto.MinerProperty{
				Type: data.Type,
				Label: &proto.MinerPropertyLabel{
					Name:   data.Label.Name,
					Unique: data.Label.Unique,
				},
				Content: &proto.MinerPropertyContent{
					Format: data.Content.Format,
					Value:  data.Content.Value,
				},
			})
		}
		protoResources = append(protoResources, &protoResource)
	}

	return &proto.MinerResources{
		Resources: protoResources,
	}, err
}

func toSharedMinerConfig(config *proto.MinerConfig) MinerConfig {
	equipments := []MinerConfigEquipment{}
	for _, equipment := range config.Equipments {
		equipments = append(equipments, MinerConfigEquipment{
			Type:       equipment.Type,
			Name:       equipment.Name,
			Attributes: equipment.Attributes,
		})
	}

	return MinerConfig{
		Auth:       config.Auth,
		Equipments: equipments,
	}
}
