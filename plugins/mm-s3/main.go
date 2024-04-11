package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
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

	// Read the tags configuration
	// tagsMap := readTagsConfig(PLUG_NAME, hclConfig)

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

			// fmt.Printf("Bucket: %s, region: %s\n", *bucket.Name, bucketRegion)
			bucketResource.Identifier = *bucket.Name

			tagsOutput, err := client.GetBucketTagging(
				context.Background(),
				&s3.GetBucketTaggingInput{
					Bucket: bucket.Name,
				},
			)
			if err != nil {
				var apiErr smithy.APIError
				if ok := errors.As(err, &apiErr); ok {
					switch apiErr.ErrorCode() {
					case "NoSuchTagSet":
						log.Println("No tags found")
						// continue
					default:
						log.Printf(
							"Error code: %s, Error message: %s\n",
							apiErr.ErrorCode(),
							apiErr.ErrorMessage(),
						)
						// continue
					}
				} else {
					log.Printf("Failed to get bucket tags, %v", err)
					// continue
				}
			} else {
				for _, tag := range tagsOutput.TagSet {
					log.Printf("Tag name: %s, value: %s\n", *tag.Key, *tag.Value)
					bucketResource.Properties = append(bucketResource.Properties, shared.MinerProperty{
						Type:  "tag",
						Name:  *tag.Key,
						Value: *tag.Value,
					})
				}
			}
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
// if the bucket is in the US East (N. Virginia) region, the region is empty
// so we return us-east-1
func getBucketRegion(client *s3.Client, bucket string) (string, error) {
	result, err := client.GetBucketLocation(context.Background(), &s3.GetBucketLocationInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", fmt.Errorf("getBucketRegion: %w", err)
	}

	region := string(result.LocationConstraint)
	if region == "" {
		region = "us-east-1"
	}

	return region, nil
}
