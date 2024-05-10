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

			// fmt.Printf("Bucket: %s, region: %s\n", *bucket.Name, bucketRegion)
			bucketResource.Identifier = *bucket.Name

			// Get bucket tags
			taggingProperties, err := getTaggingProperties(client, &bucket)
			if err != nil {
				var configErr *mmS3Error
				if errors.As(err, &configErr) {
					log.Println("No tags found")
				} else {
					log.Printf("Failed to get bucket tags, %v", err)
				}
			} else {
				bucketResource.Properties = append(bucketResource.Properties, taggingProperties...)
			}

			// Get the bucket accelerate configuration
			accelerateProperty, err := getAccelerateProperty(client, &bucket)
			if err != nil {
				var configErr *mmS3Error
				if errors.As(err, &configErr) {
					log.Println("No accelerate configuration found")
				} else {
					log.Printf("Failed to get accelerate configuration, %v", err)
				}
			} else {
				bucketResource.Properties = append(bucketResource.Properties, accelerateProperty)
			}

			// Get the bucket ACL properties
			aclProperties, err := getAclProperties(client, &bucket)
			if err != nil {
				log.Printf("Failed to get ACL properties, %v", err)
			} else {
				bucketResource.Properties = append(bucketResource.Properties, aclProperties...)
			}

			// Get the bucket CORS properties
			corsProperties, err := getCorsProperties(client, &bucket)
			if err != nil {
				var configErr *mmS3Error
				if errors.As(err, &configErr) {
					log.Println("No CORS configuration found")
				} else {
					log.Printf("Failed to get CORS properties, %v", err)
				}
			} else {
				bucketResource.Properties = append(bucketResource.Properties, corsProperties...)
			}

			// Get the bucket encryption properties
			encryptionProperties, err := getEncryptionProperties(client, &bucket)
			if err != nil {
				log.Printf("Failed to get encryption properties, %v", err)
			} else {
				bucketResource.Properties = append(bucketResource.Properties, encryptionProperties...)
			}

			// Get the bucket intelligent tiering properties
			intelligentTieringProperties, err := getIntelligentTieringProperties(client, &bucket)
			if err != nil {
				log.Printf("Failed to get intelligent tiering properties, %v", err)
			} else {
				bucketResource.Properties = append(bucketResource.Properties, intelligentTieringProperties...)
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
