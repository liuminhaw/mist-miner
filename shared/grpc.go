package shared

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/liuminhaw/mist-miner/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCClient is an implementation of Greeter that talks over RPC
type GRPCClient struct {
	// client proto.GreetServiceClient
	broker *plugin.GRPCBroker
	client proto.MinerServiceClient
}

func (m *GRPCClient) Mine(config MinerConfig, pf PropFormatter) (MinerResources, error) {
	addHelperServer := &GRPCAddHelperServer{Impl: pf}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterPropFormatterServer(s, addHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	resources, err := m.client.Mine(context.Background(), &proto.MinerConfig{
		AddServer: brokerID,
		Path:      config.Path,
	})
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

	s.Stop()

	return minerResources, nil
}

// GRPCServer is the server that GRPCClient talks to
type GRPCServer struct {
	// This is the real implementation
	// Impl Greeter
	Impl   Miner
	broker *plugin.GRPCBroker
}

func (m *GRPCServer) Mine(
	ctx context.Context,
	req *proto.MinerConfig,
) (*proto.MinerResources, error) {
	fmt.Println("GRPCServer Mine run")
	conn, err := m.broker.Dial(req.AddServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	formatter := &GRPCAddHelperClient{proto.NewPropFormatterClient(conn)}

	// func (m *GRPCServer) Mine(ctx context.Context, req *proto.NoParam) (*proto.MinerResources, error) {
	protoResources := []*proto.MinerResource{}

	resources, err := m.Impl.Mine(MinerConfig{Path: req.Path}, formatter)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Resources: %+v\n", resources)
	for _, resource := range resources {
		protoResource := proto.MinerResource{
			Identifier: resource.Identifier,
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
	}, nil
}

// GRPCAddHelperClient is the client API for GRPCAddHelper service
type GRPCAddHelperClient struct{ client proto.PropFormatterClient }

func (m *GRPCAddHelperClient) Format(a *anypb.Any) (MinerProperty, error) {
	resp, err := m.client.Format(context.Background(), a)
	if err != nil {
		hclog.Default().Info("add.Format", "client", "start", "err", err)
		return MinerProperty{}, err
	}

	return MinerProperty{
		Type: resp.Type,
		Label: MinerPropertyLabel{
			Name:   resp.Label.Name,
			Unique: resp.Label.Unique,
		},
		Content: MinerPropertyContent{
			Format: resp.Content.Format,
			Value:  resp.Content.Value,
		},
	}, nil
}

// GRPCAddHelperServer is the server that GRPCClient talks to
type GRPCAddHelperServer struct {
	// This is the real implementation
	Impl PropFormatter
}

func (m *GRPCAddHelperServer) Format(
	ctx context.Context,
	in *anypb.Any,
) (*proto.MinerProperty, error) {
	property, err := m.Impl.Format(in)
	if err != nil {
		return nil, err
	}

	return &proto.MinerProperty{
		Type: property.Type,
		Label: &proto.MinerPropertyLabel{
			Name:   property.Label.Name,
			Unique: property.Label.Unique,
		},
		Content: &proto.MinerPropertyContent{
			Format: property.Content.Format,
			Value:  property.Content.Value,
		},
	}, nil
}
