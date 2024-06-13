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

type awsProfile string

// This is the implementation of Miner
type Miner struct {
	resources shared.MinerResources
}

func (m Miner) Mine(mineConfig shared.MinerConfig) (shared.MinerResources, error) {
	// Quick test for config information
	// log.Printf("Config path: %s\n", mineConfig.Path)
    log.Printf("Plugin name: %s\n", PLUG_NAME)

	// Get authentication profile from config
	profile, err := configAuth(mineConfig)
	if err != nil {
		return nil, fmt.Errorf("mine: %w", err)
	}

	resources := shared.MinerResources{}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(string(profile)),
	)

	client := s3.NewFromConfig(cfg)
	bucketsOutput, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("mine: list buckets: %w", err)
	}

	for _, bucket := range bucketsOutput.Buckets {
		log.Printf("Bucket: %s\n", *bucket.Name)

		bucketResource := shared.MinerResource{}

		bucketRegion, err := getBucketRegion(client, *bucket.Name)
		if err != nil {
			log.Printf("Failed to get bucket region: %v", err)
			continue
		}

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithSharedConfigProfile(string(profile)),
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
				Format: shared.FormatText,
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

func configAuth(mineConfig shared.MinerConfig) (awsProfile, error) {
	if _, ok := mineConfig.Auth["profile"]; !ok {
		return "", fmt.Errorf("configAuth: profile not found")
	}

	return awsProfile(mineConfig.Auth["profile"]), nil
}

// getBucketRegion returns the region of the bucket
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
