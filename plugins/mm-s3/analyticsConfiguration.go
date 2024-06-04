package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type analyticsProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.ListBucketAnalyticsConfigurationsOutput
	requestToken   string
}

func (ap *analyticsProp) fetchConf() error {
	output, err := ap.client.ListBucketAnalyticsConfigurations(
		context.Background(),
		&s3.ListBucketAnalyticsConfigurationsInput{
			Bucket:            ap.bucket.Name,
			ContinuationToken: &ap.requestToken,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf analytics: %w", err)
	}

	ap.configurations = output
	return nil
}

func (ap *analyticsProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	ap.requestToken = ""
	for {
		if err := ap.fetchConf(); err != nil {
			return nil, fmt.Errorf("generate analytics: %w", err)
		}

		for _, config := range ap.configurations.AnalyticsConfigurationList {
			property := shared.MinerProperty{
				Type: analyticsConfig,
				Label: shared.MinerPropertyLabel{
					Name:   *config.Id,
					Unique: true,
				},
				Content: shared.MinerPropertyContent{
					Format: shared.FormatJson,
				},
			}
			if err := property.FormatContentValue(config); err != nil {
				return nil, fmt.Errorf("generate analytics: %w", err)
			}

			properties = append(properties, property)
		}

		if *ap.configurations.IsTruncated {
			ap.requestToken = *ap.configurations.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
