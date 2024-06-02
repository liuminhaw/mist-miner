package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/go-plugin"
	"github.com/liuminhaw/mist-miner/shared"
)

var PLUG_NAME = "mm-s3"

// This is the implementation of Miner
type Miner struct {
	resources shared.MinerResources
}

func (m Miner) Mine(mineConfig shared.MinerConfig) (shared.MinerResources, error) {
	// Quick test for config information
	log.Printf("Config path: %s\n", mineConfig.Path)

	// Read the HCL config file
	hclConfig, err := shared.ReadConfig(mineConfig.Path)
	if err != nil {
		return nil, fmt.Errorf("mine: read config: %w", err)
	}

	resources := shared.MinerResources{}
	for _, plug := range hclConfig.Plugs {
		if plug.Name != PLUG_NAME {
			continue
		}

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithSharedConfigProfile(plug.Profile),
		)

		client := s3.NewFromConfig(cfg)
		bucketsOutput, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
		if err != nil {
			return nil, fmt.Errorf("mine: list buckets: %w", err)
		}

		for _, bucket := range bucketsOutput.Buckets {
			bucketResource := shared.MinerResource{}

			bucketRegion, err := getBucketRegion(client, *bucket.Name)
			if err != nil {
				continue
			}

			cfg, err := config.LoadDefaultConfig(context.Background(),
				config.WithSharedConfigProfile(plug.Profile),
				config.WithRegion(bucketRegion),
			)
			client := s3.NewFromConfig(cfg)

			bucketResource.Identifier = *bucket.Name

			// Add location (region) property
			bucketResource.Properties = append(bucketResource.Properties, shared.MinerProperty{
				Type: location,
				Label: shared.MinerPropertyLabel{
					Name:   "Region",
					Unique: true,
				},
				Content: shared.MinerPropertyContent{
					Format: formatText,
					Value:  bucketRegion,
				},
			})

			for _, propType := range miningProperties {
				propsCrawler, err := New(client, &bucket, propType)
				if err != nil {
					return nil, fmt.Errorf("Failed to create new crawler: %w", err)
				}
				properties, err := propsCrawler.generate()
				if err != nil {
					var configErr *mmS3Error
					if errors.As(err, &configErr) {
						log.Printf("No %s configuration found", propType)
					} else {
						log.Printf("Failed to get %s properties: %v", propType, err)
					}
				} else {
					bucketResource.Properties = append(bucketResource.Properties, properties...)
				}
			}

			bucketResource.Sort()
			resources = append(resources, bucketResource)
		}
	}

	return resources, nil
}

func main() {
	// logger setup for plugin logs
	log.SetOutput(os.Stderr)
	log.Println("Starting miner plugin")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"miner_grpc": &shared.MinerGRPCPlugin{Impl: &Miner{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	},
	)
}

// getBucketRegion returns the region of the bucket
func getBucketRegion(client *s3.Client, bucket string) (string, error) {
	result, err := client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", fmt.Errorf("getBucketRegion: %w", err)
	}

	return *result.BucketRegion, nil
}
