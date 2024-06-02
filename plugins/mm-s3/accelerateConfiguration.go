package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type accelerateProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketAccelerateConfigurationOutput
}

func (a *accelerateProp) fetchConf() error {
	output, err := a.client.GetBucketAccelerateConfiguration(
		context.Background(),
		&s3.GetBucketAccelerateConfigurationInput{
			Bucket: a.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("getAccelerateProperty: %w", err)
	}

	a.configurations = output
	return nil
}

func (a *accelerateProp) generate() ([]shared.MinerProperty, error) {
	if err := a.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate accelerateProp: %w", err)
	}

	if a.configurations.Status == "" {
		return []shared.MinerProperty{}, &mmS3Error{accelerateConfig, noConfig}
	}

    property := shared.MinerProperty{
        Type: accelerateConfig,
        Label: shared.MinerPropertyLabel{
            Name:   "Status",
            Unique: true,
        },
        Content: shared.MinerPropertyContent{   
            Format: shared.FormatText,
        },
    }
    if err := property.FormatContentValue(a.configurations.Status); err != nil {
        return nil, fmt.Errorf("generate accelerateProp: %w", err)
    }

    return []shared.MinerProperty{property}, nil
}
