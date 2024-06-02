package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type loggingProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketLoggingOutput
}

func (lp *loggingProp) fetchConf() error {
	output, err := lp.client.GetBucketLogging(
		context.Background(),
		&s3.GetBucketLoggingInput{
			Bucket: lp.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf loggingProp: %w", err)
	}

	lp.configurations = output
	return nil
}

func (lp *loggingProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := lp.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate loggingProp: %w", err)
	}

	if lp.configurations.LoggingEnabled != nil {
		property := shared.MinerProperty{
			Type: logging,
			Label: shared.MinerPropertyLabel{
				Name:   "Logging",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
		if err := property.FormatContentValue(lp.configurations.LoggingEnabled); err != nil {
			return nil, fmt.Errorf("generate loggingProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
