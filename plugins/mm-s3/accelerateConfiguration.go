package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

// getTaggingProperties retrieves the tagging properties of a bucket.
// Property list:
// - Type: accelerateConfig
// - Name: status
func getAccelerateProperty(client *s3.Client, bucket *types.Bucket) (shared.MinerProperty, error) {
	output, err := client.GetBucketAccelerateConfiguration(
		context.Background(),
		&s3.GetBucketAccelerateConfigurationInput{
			Bucket: bucket.Name,
		},
	)
	if err != nil {
		return shared.MinerProperty{}, fmt.Errorf("getAccelerateProperty: %w", err)
	}

	if output.Status == "" {
		return shared.MinerProperty{}, &mmS3Error{accelerateConfig, noConfig}
	}

	return shared.MinerProperty{
		Type:  accelerateConfig,
		Name:  "Status",
		Value: string(output.Status),
	}, nil
}
